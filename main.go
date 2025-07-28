package main

// Run main Agent that interacts with the Claude model using anthropic-sdk-go.

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/invopop/jsonschema"
)

const (
	maxPromptTokens = 200_000 // Maximum allowed tokens in a single request
	maxTotalTokens  = 70_000  // Stay safely below the 80k tokens/minute rate limit
)

var (
	debug bool
)

func main() {
	var (
		outputDir string
		logNum    int
	)
	// Define command line flags
	flag.StringVar(&outputDir, "o", "logs", "Log output directory")
	flag.IntVar(&logNum, "l", 0, "Log file number (for debugging)")
	flag.BoolVar(&debug, "d", false, "enable debug mode")
	flag.Parse()

	if err := makePrintfFunctions(outputDir, logNum); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logging: %v\n", err)
		os.Exit(1)
	}

	// oprintf("***stdout: Starting Agent with output directory: %s\n", outputDir)
	// eprintf("***stderr: Debug mode: %v\n", debug)
	// oprintf("args: %d %q\n", len(flag.Args()), flag.Args())
	// os.Exit(0) // Exit early for debugging purposes

	client := anthropic.NewClient()

	// Get remaining arguments as prompts
	prompts := flag.Args()

	var getUserMessage func() (string, bool)

	promptIndex := 0
	scanner := bufio.NewScanner(os.Stdin)
	getUserMessage = func() (string, bool) {
		if promptIndex < len(prompts) {
			prompt := prompts[promptIndex]
			eprintf("Prompt %d of %d: %s\n", promptIndex+1, len(prompts), prompt)
			promptIndex++
			return prompt, true
		}
		// After all prompts are used, switch to interactive mode.
		if promptIndex == len(prompts) && len(prompts) > 0 {
			eprintf("--- All prompts used. Switching to interactive mode. ---\n")
			// Increment promptIndex to prevent this message from showing again.
			promptIndex++
		}

		// Interactive mode: read from stdin.
		if scanner.Scan() {
			return scanner.Text(), true
		}
		return "", false
	}

	getUserMessage = prependSystemPrompt(getUserMessage, []string{
		"If you create a new file, put it in the ./bin directory.",
		"If you create a new data file, such as JSON file, put it in the ./data directory.",
		"If you create a program, check that it works by running it first.",
	})
	getUserMessage = wrapAndSavePrompts(getUserMessage, "prompts.txt")

	tools := []ToolDefinition{readFileDefinition, listFilesDefinition, editFileDefinition,
		execProgramDefinition}
	agent := newAgent(&client, getUserMessage, tools)

	if err := agent.Run(context.TODO()); err != nil {
		eprintf(fmt.Sprintf("Error: %s\n", err.Error()))
	}
}

// Agent represents a chat agent that interacts with the Claude model.
type Agent struct {
	client         *anthropic.Client
	getUserMessage func() (string, bool)
	tools          []ToolDefinition
}

// newAgent creates a new Agent instance with the provided client, user message function, and tools
func newAgent(
	client *anthropic.Client,
	getUserMessage func() (string, bool),
	tools []ToolDefinition,
) *Agent {
	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
		tools:          tools,
	}
}

// Run starts the chat loop with the user, allowing them to interact with the Claude model.
// It reads user input, sends messages to the model, and handles tool invocations.
// The loop continues until the user decides to quit (e.g. by pressing ctrl-c).
func (a *Agent) Run(ctx context.Context) error {
	var conversation []anthropic.MessageParam

	fmt.Println("Chat with Claude (use 'ctrl-c' to quit)")

	// Create a context that will be canceled on SIGINT (Ctrl+C).
	ctx = withSignalCancellation(ctx, syscall.SIGINT, syscall.SIGTERM)

	readUserInput := true
	for {
		if readUserInput {
			oprintf("\u001b[94mYou\u001b[0m: ")
			userInput, ok := a.getUserMessage()
			if !ok {
				break
			}
			eprintf("*User input: %s\n", userInput)

			userMessage := anthropic.NewUserMessage(anthropic.NewTextBlock(userInput))
			conversation = append(conversation, userMessage)
		}

		message, err := a.runInference(ctx, conversation)
		if err != nil {
			return fmt.Errorf("failed to run inference: %w", err)
		}

		// Only add the assistant's message to the conversation if it's not empty.
		if len(message.Content) > 0 {
			conversation = append(conversation, message.ToParam())
		} else {
			// If the message is empty, there's nothing to process, so we should
			// probably ask for user input again.
			eprintf("--- Assistant returned empty message, asking for user input. ---\n")
			readUserInput = true
			continue
		}

		var toolResults []anthropic.ContentBlockParamUnion
		for _, content := range message.Content {
			switch content.Type {
			case "text":
				oprintf("\u001b[93mClaude\u001b[0m: %s\n", content.Text)
			case "tool_use":
				result := a.executeTool(content.ID, content.Name, content.Input)
				toolResults = append(toolResults, result)
			}
		}
		if len(toolResults) == 0 {
			readUserInput = true
			continue
		}
		readUserInput = false
		conversation = append(conversation, anthropic.NewUserMessage(toolResults...))
	}

	return nil
}

// executeTool executes a tool with the given name and input.
func (a *Agent) executeTool(id, name string, input json.RawMessage,
) anthropic.ContentBlockParamUnion {
	tool, found := a.getTool(name)
	if !found {
		return anthropic.NewToolResultBlock(id, "tool not found", true)
	}

	oprintf("\u001b[92mtool\u001b[0m: %s(%s)\n", name, input)
	response, err := tool.Function(input)
	if err != nil {
		return anthropic.NewToolResultBlock(id, err.Error(), true)
	}
	return anthropic.NewToolResultBlock(id, response, false)
}

// getTool retrieves a tool by its name from the agent's list of tools.
func (a *Agent) getTool(name string) (ToolDefinition, bool) {
	for _, tool := range a.tools {
		if tool.Name == name {
			return tool, true
		}
	}
	return ToolDefinition{}, false
}

// runInference sends the conversation to the Claude model and returns the response message.
// It includes the defined tools in the request so that Claude can invoke them if needed.
func (a *Agent) runInference(ctx context.Context, conversation []anthropic.MessageParam,
) (*anthropic.Message, error) {
	// Loop to shorten the conversation if it's too long, while preserving the first message (system prompt).
	for estimateTokenCount(conversation) > maxTotalTokens && len(conversation) > 2 {
		// Remove the second element (the oldest message after the system prompt)
		conversation = append(conversation[:1], conversation[2:]...)
		eprintf("--- Conversation history too long, removing oldest message. ---\n")
	}

	var anthropicTools []anthropic.ToolUnionParam
	for _, tool := range a.tools {
		anthropicTools = append(anthropicTools, anthropic.ToolUnionParam{
			OfTool: &anthropic.ToolParam{
				Name:        tool.Name,
				Description: anthropic.String(tool.Description),
				InputSchema: tool.InputSchema,
			},
		})
	}

	// Estimate and log the token count before making the API call.
	estimatedTokens := estimateTokenCount(conversation)

	if estimatedTokens > maxPromptTokens {
		panic(fmt.Errorf("prompt is too long: %d tokens > %d maximum",
			estimatedTokens, maxPromptTokens))
	}

	startTime := time.Now()
	message, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_7SonnetLatest,
		MaxTokens: int64(1024),
		Messages:  conversation,
		Tools:     anthropicTools,
	})
	duration := time.Since(startTime)

	description := ""
	if debug && len(conversation) > 0 {
		lastMessage := conversation[len(conversation)-1]
		contentDesc := describeContents(lastMessage.Content)
		description = fmt.Sprintf("last message: role=%s, content=[%s]",
			lastMessage.Role, contentDesc)
	}
	eprintf("anthropic call: %4.1f seconds, %d messages, ~%d tokens\n",
		duration.Seconds(), len(conversation), estimatedTokens)
	eprintf("  key counts (top): %v\n", keyCountsTop)
	eprintf("  key counts (all): %v\n", keyCountsAll)
	eprintf("  %s\n", description)

	return message, err
}

// ToolDefinition defines a tool that can be invoked by the Claude model.
type ToolDefinition struct {
	Name        string                         `json:"name"`
	Description string                         `json:"description"`
	InputSchema anthropic.ToolInputSchemaParam `json:"input_schema"`
	Function    func(input json.RawMessage) (string, error)
}

// readFileDefinition defines a tool for reading the contents of a file at a specified relative path.
var readFileDefinition = ToolDefinition{
	Name:        "read_file",
	Description: "Read the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names.",
	InputSchema: readFileInputSchema,
	Function:    readFile,
}

var readFileInputSchema = generateSchema[ReadFileInput]()

type ReadFileInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of a file in the working directory."`
}

// readFile reads the contents of a file at the specified relative path and returns it as a string.
func readFile(input json.RawMessage) (string, error) {
	var readFileInput ReadFileInput
	if err := json.Unmarshal(input, &readFileInput); err != nil {
		panic(err)
	}

	content, err := os.ReadFile(readFileInput.Path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

var listFilesDefinition = ToolDefinition{
	Name:        "list_files",
	Description: "List files and directories at a given path. If no path is provided, lists files in the current directory.",
	InputSchema: listFilesInputSchema,
	Function:    listFiles,
}

var listFilesInputSchema = generateSchema[ListFilesInput]()

type ListFilesInput struct {
	Path string `json:"path,omitempty" jsonschema_description:"Optional relative path to list files from. Defaults to current directory if not provided."`
}

// listFiles lists all files and directories in the specified path, returning a JSON array of file paths.
func listFiles(input json.RawMessage) (string, error) {
	var listFilesInput ListFilesInput
	if err := json.Unmarshal(input, &listFilesInput); err != nil {
		panic(err)
	}

	dir := "."
	if listFilesInput.Path != "" {
		dir = listFilesInput.Path
	}

	var files []string
	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		basePath := filepath.Base(path)
		if strings.HasPrefix(basePath, ".") || strings.HasPrefix(relPath, ".") {
			// Skip hidden files and directories
			return nil
		}
		if basePath == "README.md" ||
			basePath == "prompts.txt" ||
			basePath == "fizzbuzz.js" ||
			strings.HasPrefix(relPath, "log.") {
			// Skip files that record output of previous runs
			return nil
		}
		if info.IsDir() {
			files = append(files, relPath+"/")
		} else {
			files = append(files, relPath)
		}
		return nil
	}
	if err := filepath.Walk(dir, walk); err != nil {
		return "", err
	}

	result, err := json.Marshal(files)
	if err != nil {
		return "", err
	}

	if debug {
		oprintf("ListFiles result %q: %s\n", dir, string(result))
	}

	return string(result), nil
}

var editFileDefinition = ToolDefinition{
	Name: "edit_file",
	Description: `Make edits to a text file.

Replaces 'old_str' with 'new_str' in the given file. 'old_str' and 'new_str' MUST be different from each other.

If the file specified with path doesn't exist, it will be created.
`,
	InputSchema: EditFileInputSchema,
	Function:    editFile,
}

type EditFileInput struct {
	Path   string `json:"path" jsonschema_description:"The path to the file"`
	OldStr string `json:"old_str" jsonschema_description:"Text to search for - must match exactly and must only have one match exactly"`
	NewStr string `json:"new_str" jsonschema_description:"Text to replace old_str with"`
}

var EditFileInputSchema = generateSchema[EditFileInput]()

func editFile(input json.RawMessage) (string, error) {
	var editFileInput EditFileInput
	if err := json.Unmarshal(input, &editFileInput); err != nil {
		return "", err
	}

	if editFileInput.Path == "" || editFileInput.OldStr == editFileInput.NewStr {
		return "", fmt.Errorf("invalid input parameters")
	}

	content, err := os.ReadFile(editFileInput.Path)
	if err != nil {
		if os.IsNotExist(err) && editFileInput.OldStr == "" {
			return createNewFile(editFileInput.Path, editFileInput.NewStr)
		}
		return "", err
	}

	oldContent := string(content)
	newContent := strings.Replace(oldContent, editFileInput.OldStr, editFileInput.NewStr, -1)

	if oldContent == newContent && editFileInput.OldStr != "" {
		return "", fmt.Errorf("old_str not found in file")
	}

	if err := os.WriteFile(editFileInput.Path, []byte(newContent), 0644); err != nil {
		return "", err
	}

	return "OK", nil
}

var execProgramDefinition = ToolDefinition{
	Name: "exec_program",
	Description: `Executes a program, binary, or script with the given arguments.
WARNING: This tool allows the execution of arbitrary code, which can be dangerous.
The command output (stdout) will be returned. If the command fails, stderr will be included in the error message.`,
	InputSchema: execProgramInputSchema,
	Function:    execProgram,
}

var execProgramInputSchema = generateSchema[ExecProgramInput]()

type ExecProgramInput struct {
	Command string   `json:"command" jsonschema_description:"The program/binary/script to execute."`
	Args    []string `json:"args,omitempty" jsonschema_description:"A list of arguments to pass to the command."`
}

// execProgram executes a command with the provided arguments and returns its output.
// If the command fails, it returns an error with the stderr output.
func execProgram(input json.RawMessage) (string, error) {
	var execInput ExecProgramInput
	if err := json.Unmarshal(input, &execInput); err != nil {
		panic(err)
	}

	if execInput.Command == "" {
		return "", fmt.Errorf("command cannot be empty")
	}

	cmd := exec.Command(execInput.Command, execInput.Args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %w. Output: %s", err, string(output))
	}

	return string(output), nil
}

// generateSchema generates a JSON schema for type `T` using the jsonschema package.
func generateSchema[T any]() anthropic.ToolInputSchemaParam {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T

	schema := reflector.Reflect(v)

	return anthropic.ToolInputSchemaParam{
		Properties: schema.Properties,
	}
}

// estimateTokenCount provides a rough estimation of the number of tokens in a conversation.
// It uses a simple character count, assuming ~4 characters per token.
func estimateTokenCount(messages []anthropic.MessageParam) int {
	var totalChars int
	for _, message := range messages {
		for _, content := range message.Content {
			// This is a simplified way to get the character count.
			// We marshal each content block to JSON and count the characters.
			if data, err := json.Marshal(content); err == nil {
				totalChars += len(data)
			}
		}
	}
	// A common heuristic is 4 characters per token.
	return totalChars / 4
}
