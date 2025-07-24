package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/invopop/jsonschema"
)

func main() {
	// Define command line flags
	var debug bool
	flag.BoolVar(&debug, "d", false, "enable debug mode")
	flag.Parse()

	client := anthropic.NewClient()

	// Get remaining arguments as prompts
	prompts := flag.Args()

	var getUserMessage func() (string, bool)

	if len(prompts) > 0 {
		// If prompts are provided as arguments, use them sequentially
		promptIndex := 0
		getUserMessage = func() (string, bool) {
			if promptIndex < len(prompts) {
				prompt := prompts[promptIndex]
				fmt.Fprintf(os.Stderr, "Prompt %d of %d: %s\n", promptIndex+1, len(prompts), prompt)
				promptIndex++
				return prompt, true
			}
			fmt.Fprintf(os.Stderr, "--- All prompts used.\n")
			return "", false
		}
	} else {
		// No prompts provided, use interactive mode from the start
		scanner := bufio.NewScanner(os.Stdin)
		getUserMessage = func() (string, bool) {
			if !scanner.Scan() {
				return "", false
			}
			return scanner.Text(), true
		}
	}

	getUserMessage = wrapAndSavePrompts(getUserMessage, "prompts.txt")

	tools := []ToolDefinition{ReadFileDefinition, ListFilesDefinition}
	agent := NewAgent(&client, getUserMessage, tools, debug)

	if err := agent.Run(context.TODO()); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}

// Agent represents a chat agent that interacts with the Claude model.
type Agent struct {
	client         *anthropic.Client
	getUserMessage func() (string, bool)
	tools          []ToolDefinition
	debug          bool
}

// NewAgent creates a new Agent instance with the provided client, user message function, and tools
func NewAgent(
	client *anthropic.Client,
	getUserMessage func() (string, bool),
	tools []ToolDefinition,
	debug bool,
) *Agent {
	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
		tools:          tools,
		debug:          debug,
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
			fmt.Print("\u001b[94mYou\u001b[0m: ")
			userInput, ok := a.getUserMessage()
			if !ok {
				break
			}

			userMessage := anthropic.NewUserMessage(anthropic.NewTextBlock(userInput))
			conversation = append(conversation, userMessage)
		}

		message, err := a.runInference(ctx, conversation)
		if err != nil {
			return err
		}
		conversation = append(conversation, message.ToParam())

		var toolResults []anthropic.ContentBlockParamUnion
		for _, content := range message.Content {
			switch content.Type {
			case "text":
				fmt.Printf("\u001b[93mClaude\u001b[0m: %s\n", content.Text)
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

	fmt.Printf("\u001b[92mtool\u001b[0m: %s(%s)\n", name, input)
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

	startTime := time.Now()
	message, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_7SonnetLatest,
		MaxTokens: int64(1024),
		Messages:  conversation,
		Tools:     anthropicTools,
	})
	duration := time.Since(startTime)
	fmt.Fprintf(os.Stderr, "anthropic call took %4.1f seconds\n", duration.Seconds())

	return message, err
}

// ToolDefinition defines a tool that can be invoked by the Claude model.
type ToolDefinition struct {
	Name        string                         `json:"name"`
	Description string                         `json:"description"`
	InputSchema anthropic.ToolInputSchemaParam `json:"input_schema"`
	Function    func(input json.RawMessage) (string, error)
}

// ReadFileDefinition defines a tool for reading the contents of a file at a specified relative path.
var ReadFileDefinition = ToolDefinition{
	Name:        "read_file",
	Description: "Read the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names.",
	InputSchema: ReadFileInputSchema,
	Function:    ReadFile,
}

var ReadFileInputSchema = GenerateSchema[ReadFileInput]()

type ReadFileInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of a file in the working directory."`
}

// ReadFile reads the contents of a file at the specified relative path and returns it as a string.
func ReadFile(input json.RawMessage) (string, error) {
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

var ListFilesDefinition = ToolDefinition{
	Name:        "list_files",
	Description: "List files and directories at a given path. If no path is provided, lists files in the current directory.",
	InputSchema: ListFilesInputSchema,
	Function:    ListFiles,
}

var ListFilesInputSchema = GenerateSchema[ListFilesInput]()

type ListFilesInput struct {
	Path string `json:"path,omitempty" jsonschema_description:"Optional relative path to list files from. Defaults to current directory if not provided."`
}

// ListFiles lists all files and directories in the specified path, returning a JSON array of file paths.
func ListFiles(input json.RawMessage) (string, error) {
	var listFilesInput ListFilesInput
	if err := json.Unmarshal(input, &listFilesInput); err != nil {
		panic(err)
	}

	dir := "."
	if listFilesInput.Path != "" {
		dir = listFilesInput.Path
	}

	var files []string
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		if relPath != "." {
			if info.IsDir() {
				files = append(files, relPath+"/")
			} else {
				files = append(files, relPath)
			}
		}
		return nil
	}); err != nil {
		return "", err
	}

	result, err := json.Marshal(files)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// GenerateSchema generates a JSON schema for type `T` using the jsonschema package.
func GenerateSchema[T any]() anthropic.ToolInputSchemaParam {
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

// withSignalCancellation creates a context that is canceled when one of the specified signals is received.
// It's a common pattern for graceful shutdown.
func withSignalCancellation(parent context.Context, sigs ...os.Signal) context.Context {
	// Create a new context that we can cancel manually.
	// The parent context is passed so that cancellation propagates downwards.
	ctx, cancel := context.WithCancel(parent)

	// Create a channel to receive OS signals.
	// A buffer of 1 is recommended so the signal package doesn't block.
	sigChan := make(chan os.Signal, 1)

	// Register the given signals to be sent to the channel.
	// If no signals are provided, we'll default to SIGINT and SIGTERM.
	if len(sigs) == 0 {
		sigs = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}
	signal.Notify(sigChan, sigs...)

	// Start a new goroutine.
	// This goroutine will block until a signal is received.
	go func() {
		// Wait for a signal.
		sig := <-sigChan
		fmt.Printf("\n[Signal Handler] Received signal: %s. Cancelling context & shutting down...\n",
			sig)

		// Once a signal is received, call the cancel function.
		// This will cause the context's Done() channel to be closed.
		cancel()

		// It's good practice to clean up the signal notification.
		signal.Stop(sigChan)

		os.Exit(0)
	}()

	return ctx
}

// wrapAndSavePrompts creates a wrapper around a getUserMessage function
// to save the prompts to a file.
func wrapAndSavePrompts(innerGetUserMessage func() (string, bool), filePath string) func() (string, bool) {
	return func() (string, bool) {
		prompt, ok := innerGetUserMessage()
		if ok && prompt != "" {
			f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening prompts file: %v\n", err)
			} else {
				defer f.Close()
				if _, err := f.WriteString(prompt + "\n"); err != nil {
					fmt.Fprintf(os.Stderr, "Error writing to prompts file: %v\n", err)
				}
			}
		}
		return prompt, ok
	}
}
