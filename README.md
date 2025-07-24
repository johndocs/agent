# agent
Code copied from [How to Build an Agent](https://ampcode.com/how-to-build-an-agent)

Example usage

```
% go run .                                                                                                peter@Peters-MacBook-Pro agent % go run .
Chat with Claude (use 'ctrl-c' to quit)
You:  What is this directory?
Claude: To understand what this directory contains, I'll list the files and directories in the current location.
tool: list_files({})
Claude: This directory appears to be a Git repository for a Go programming project. Here's what I can determine:

1. It's a Go project with:
   - `main.go` - The main Go source code file
   - `go.mod` - Go module definition file
   - `go.sum` - Go module checksum file

2. It includes Git version control files:
   - `.git/` directory with all Git-related data
   - `.gitignore` - Specifies files that should be ignored by Git

3. Documentation:
   - `README.md` - Project documentation in Markdown format

To get more specific information about the project, I could read the README.md or main.go file. Would you like me to examine either of those files to provide more details about what this project does?
You: yes
Claude: I'll examine both the README.md and main.go files to give you a better understanding of this project.

First, let's look at the README.md:
tool: read_file({"path":"README.md"})
Claude: Now, let's examine the main.go file to understand the code implementation:
tool: read_file({"path":"main.go"})
Claude: Let's also check the go.mod file to understand the dependencies:
tool: read_file({"path":"go.mod"})
Claude: Based on the files I've examined, here's a comprehensive description of this directory:

## Project Overview: Claude-Powered File System Agent

This is a Go application that creates an interactive terminal-based agent powered by Anthropic's Claude AI model (specifically Claude 3.7 Sonnet). The agent allows users to chat with Claude while giving Claude the ability to interact with the local file system through two defined tools.

### Key Components:

1. **Main Application (main.go)**:
   - Implements a terminal-based chat interface to interact with Claude
   - Defines two file system tools that Claude can use during the conversation:
     - `read_file`: Allows Claude to read the contents of specified files
     - `list_files`: Enables Claude to list files and directories at a given path
   - Manages the conversation flow between users and Claude
   - Handles tool invocations and passing the results back to Claude

2. **Dependencies (go.mod)**:
   - Primary dependencies:
     - `github.com/anthropics/anthropic-sdk-go` (v1.6.2): Official Anthropic SDK for Go
     - `github.com/invopop/jsonschema` (v0.13.0): Used to generate JSON schemas for tool definitions

3. **Documentation (README.md)**:
   - Indicates this is an example implementation from a guide titled "How to Build an Agent"
   - Includes an example of the agent in action, demonstrating Claude examining the code itself

### How It Works:

1. The application starts a chat session with Claude 3.7 Sonnet
2. Users type messages in the terminal
3. Claude responds with either:
   - Text responses (displayed directly to the user)
   - Tool invocation requests (when it needs to read files or list directories)
4. When Claude requests to use a tool, the application:
   - Executes the requested file system operation
   - Returns the results back to Claude
   - Claude processes the information and continues the conversation

### Technical Implementation Details:

- The `Agent` struct manages the conversation state and handles interactions with Claude
- The `ToolDefinition` struct defines the tools that Claude can use
- JSON schema generation is used to properly format tool definitions for Claude
- Signal handling is implemented for graceful shutdown on Ctrl+C

This is a practical example of building an AI agent that can interact with a file system while maintaining a conversational interface with the user, demonstrating how AI models can be given access to specific tools through a well-defined interface.
You:


```
