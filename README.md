# agent
Code copied from [How to Build an Agent](https://ampcode.com/how-to-build-an-agent)

## Example usage

This shows the variation in output between runs

```
% go run . "What is this directory?" "Please examine the contents of all files"
```

### Output first run

```
Chat with Claude (use 'ctrl-c' to quit)
[94mYou[0m: [93mClaude[0m: I'll help you understand what this directory contains by listing the files and directories present in it.
[92mtool[0m: list_files({})
[93mClaude[0m: Based on the file listing, this appears to be a Go programming project that's managed with Git version control. Here's what I can tell you about this directory:

1. **Go Application Files**:
   - `main.go` - The main Go source code file
   - `go.mod` - Go module definition file that defines dependencies
   - `go.sum` - Checksum file for Go dependencies

2. **Documentation**:
   - `README.md` - Project documentation in Markdown format
   - `prompts.txt` - Possibly contains prompts or instructions

3. **Version Control**:
   - `.git/` directory - Contains all Git version control information
   - `.gitignore` - Specifies files that Git should ignore

4. **Git Branches**:
   - The repository has at least three branches: `main`, `start01`, and `start02`

To get more specific information about the project, I could examine the README.md or main.go file. Would you like me to look at either of those files to provide more details about the project's purpose?
[94mYou[0m: [93mClaude[0m: I'll examine the contents of all the non-Git files to give you a comprehensive overview of this project. Let me start by looking at each file:
[92mtool[0m: read_file({"path":"README.md"})
[92mtool[0m: read_file({"path":"main.go"})
[92mtool[0m: read_file({"path":"go.mod"})
[92mtool[0m: read_file({"path":"go.sum"})
[92mtool[0m: read_file({"path":".gitignore"})
[92mtool[0m: read_file({"path":"prompts.txt"})
[93mClaude[0m: Based on my examination of all the files in this directory, here's a comprehensive overview:

## Project Summary

This is a Go application that implements an AI agent powered by Anthropic's Claude 3.7 Sonnet model. The agent can interact with your local file system through defined tools, allowing Claude to read files and list directories while maintaining a chat conversation with the user.

### Key Components

1. **Main Application (`main.go`)**:
   - This is the core of the application, which:
     - Sets up a terminal-based chat interface with Claude
     - Defines and implements two file system tools:
       - `read_file`: Allows Claude to read file contents
       - `list_files`: Allows Claude to list files in a directory
     - Handles the conversation flow between users and Claude
     - Manages tool invocations and returns results to Claude
     - Includes graceful handling of interruptions (ctrl-c)
     - Saves user prompts to a file

2. **Dependencies (from `go.mod` and `go.sum`)**:
   - Primary dependencies:
     - `github.com/anthropics/anthropic-sdk-go` v1.6.2: Official Anthropic SDK for Go
     - `github.com/invopop/jsonschema` v0.13.0: Used for generating JSON schemas for tool definitions
   - Go version: 1.23.4

3. **Documentation**:
   - `README.md`: Brief project overview explaining that this code is from a tutorial titled "How to Build an Agent" with an example of usage
   - Shows a sample conversation demonstrating Claude examining the project directory itself

4. **Other Files**:
   - `.gitignore`: Standard Go project gitignore file, excluding binaries, test outputs, etc.
   - `prompts.txt`: Contains a record of user inputs/prompts from previous runs of the application

### Technical Implementation Details

- The application uses the Anthropic SDK to communicate with Claude 3.7 Sonnet
- It implements a structured way to define tools that Claude can use (through the `ToolDefinition` struct)
- Uses JSON schema generation to properly format tool definitions for Claude
- Implements an interactive terminal chat interface with colorized output
- Has context management for graceful shutdown on interrupt signals
- Maintains conversation history for context in the ongoing dialogue

### How It Works

1. The application starts a chat session with Claude 3.7 Sonnet
2. Users type messages in the terminal
3. Claude responds with either:
   - Regular text responses (displayed to the user)
   - Tool invocation requests (when it needs to access files or list directories)
4. When Claude requests to use a tool, the application:
   - Executes the requested file operation
   - Returns the results back to Claude
   - Claude processes the information and continues the conversation

This project serves as a practical example of building an AI agent with tool-use capabilities, specifically focused on file system interactions. It demonstrates how Large Language Models like Claude can be extended with custom tools to perform useful tasks while maintaining a conversational interface.
[94mYou[0m:
```


### Output second run

```
Chat with Claude (use 'ctrl-c' to quit)
[94mYou[0m: [93mClaude[0m: I'll help you explore this directory. Let's first list the files and directories in the current location to understand what's available.
[92mtool[0m: list_files({})
[93mClaude[0m: This appears to be a Go language project directory. Here's what I can tell from the file listing:

1. It's a Git repository (indicated by the `.git/` directory)
2. It contains:
   - `main.go` - the main Go source code file
   - `go.mod` and `go.sum` - Go module dependency files
   - `README.md` - project documentation
   - `prompts.txt` - possibly contains prompts or instructions
   - `.gitignore` - specifies files to be ignored by Git

Let's look at the README.md file to get more information about this project:
[92mtool[0m: read_file({"path":"README.md"})
[93mClaude[0m: Now, let's look at the main Go code file to understand what this project does:
[92mtool[0m: read_file({"path":"main.go"})
[93mClaude[0m: Let's also check the Go module dependencies:
[92mtool[0m: read_file({"path":"go.mod"})
[93mClaude[0m: Based on examining the files, this directory contains:

## Agent Project - Claude-Powered File System Tool

This is a Go application that implements an agent powered by Anthropic's Claude 3.7 Sonnet AI model. The agent allows users to have a text-based chat with Claude while giving Claude the ability to interact with the local file system through defined tools.

### Key Components:

1. **Main Application (main.go)**:
   - Creates a terminal-based chat interface
   - Defines two file system tools that Claude can use:
     - `read_file`: Reads the contents of specified files
     - `list_files`: Lists files and directories at a given path
   - Handles the conversation flow and tool execution

2. **Dependencies (go.mod)**:
   - Uses `github.com/anthropics/anthropic-sdk-go` to interact with Claude API
   - Uses `github.com/invopop/jsonschema` for JSON schema generation

3. **Project Structure**:
   - `main.go`: Contains all the code for the agent
   - `go.mod` and `go.sum`: Go module definition and dependency checksums
   - `README.md`: Project documentation
   - `prompts.txt`: Saves user prompts/inputs

### How It Works:

1. The program establishes a chat session with Claude 3.7 Sonnet
2. Users type messages in the terminal
3. Claude responds, and can request to use tools (read files or list directories)
4. The app executes the requested tool operation and returns results to Claude
5. Claude then incorporates this information into its responses

As shown in the README example, the agent is capable of analyzing its own codebase by using the file system tools it was given, demonstrating a practical example of an AI assistant with access to controlled external capabilities.

The code was copied from a tutorial called "How to Build an Agent" as mentioned in the README.md file.
[94mYou[0m: [93mClaude[0m: I'll examine the contents of all files in this directory to give you a comprehensive overview of the project. Let's start by listing the files again to make sure we review everything important.
[92mtool[0m: list_files({})
[93mClaude[0m: Let's examine the main files of interest one by one. I'll skip the `.git` directory contents as those are standard Git version control files. Let's look at the key project files:

1. Let's check the `.gitignore` file:
[92mtool[0m: read_file({"path":".gitignore"})
[93mClaude[0m: 2. We've already seen the README.md, but let's check the `prompts.txt` file:
[92mtool[0m: read_file({"path":"prompts.txt"})
[93mClaude[0m: 3. Let's check the `go.sum` file which contains checksums for all the dependencies:
[92mtool[0m: read_file({"path":"go.sum"})
[93mClaude[0m: Now that I've examined all the key files in the directory, I can provide a comprehensive overview of this project:

## Agent Project - Claude-Powered File System Tool

### Overview
This is a Go application that implements a terminal-based chat agent powered by Anthropic's Claude 3.7 Sonnet AI model. The agent allows users to chat with Claude while giving it the ability to interact with the local file system through specifically defined tools.

### Project Files
1. **main.go**: The main Go source file containing all the code for the agent, including:
   - Terminal-based chat interface implementation
   - Tool definitions (read_file and list_files)
   - Conversation flow management
   - Claude API integration

2. **go.mod & go.sum**: Go module definition and dependency verification files
   - Primary dependencies:
     - `github.com/anthropics/anthropic-sdk-go v1.6.2`: The official Anthropic SDK for Go
     - `github.com/invopop/jsonschema v0.13.0`: Used for JSON schema generation for tool definitions

3. **README.md**: Project documentation showing this is a code example from a tutorial titled "How to Build an Agent"

4. **prompts.txt**: A file that stores user prompts entered during agent interactions (contains repeated queries like "What is this directory?" and "Please examine the contents of all files")

5. **.gitignore**: Standard Go project gitignore file specifying which files should be excluded from version control

### Key Features
1. **Terminal-based Chat Interface**:
   - Allows users to have a natural language conversation with Claude
   - Handles user input via terminal and presents Claude's responses

2. **File System Tool Integration**:
   - `read_file`: Allows Claude to read the contents of specified files
   - `list_files`: Enables Claude to list files and directories at a given path

3. **Tool Execution Framework**:
   - Defines tools with standardized input schemas using JSON schema
   - Manages tool invocation requests from Claude
   - Executes tool operations and returns results back to Claude

4. **Conversation Management**:
   - Maintains conversation history between user and Claude
   - Handles tool results as part of the conversation flow
   - Formats messages appropriately for Claude's API

5. **Graceful Shutdown**:
   - Implements signal handling for proper termination on Ctrl+C

6. **Prompt Saving**:
   - Automatically saves user prompts to prompts.txt for future reference

### Implementation Details
The core architecture revolves around the `Agent` struct, which manages:
- The Claude client connection
- User message retrieval
- Tool definitions and execution
- Conversation flow

The application demonstrates a pattern for creating AI agents with access to external tools by:
1. Defining tool schemas and functions
2. Handling the conversational state
3. Processing tool requests from the AI
4. Returning tool results back to the AI

### Usage
As shown in the README example, the agent can be started with:
```
go run .
```

The agent is self-referential - it can examine its own code and structure by using the file system tools it was given access to, as demonstrated in the example interaction in the README.

The `prompts.txt` file shows that the agent has been used multiple times to examine the directory contents and file contents, which suggests it's being used or tested for its file exploration capabilities.
[94mYou[0m:
```
