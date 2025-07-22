# agent
Code copied from [How to Build an Agent](https://ampcode.com/how-to-build-an-agent)

Example usage

```
% go run .
Chat with Claude (use 'ctrl-c' to quit)
You: What is this directory?
Claude: To understand what this directory contains, I'll list the files and directories in the current working directory.
tool: list_files({})
Claude: This directory appears to be a Go programming project repository managed with Git. Here's what I can observe:

1. **Core Go files**:
   - `main.go` - The main Go source code file
   - `go.mod` - Go module definition file that manages dependencies
   - `go.sum` - Go module checksums for dependency verification

2. **Documentation**:
   - `README.md` - Project documentation in Markdown format

3. **Git Configuration**:
   - `.git/` directory with all Git version control information
   - `.gitignore` - Specifies files to be ignored by Git

Let me look at the README.md file to get more information about what this project does:
tool: read_file({"path":"README.md"})
tool: read_file({"path":"main.go"})
tool: read_file({"path":"go.mod"})
Claude: Based on the files I've examined, this directory contains a Go programming project that implements an AI agent using Anthropic's Claude API. Here's a detailed overview:

## Project Summary: Claude Agent with File System Tools

This is a Go application that creates an interactive agent powered by Claude (Anthropic's AI model) with the ability to interact with the local file system. The agent acts as a bridge between a user and Claude, allowing Claude to use tools that can read files and list directory contents.

### Key Components:

1. **Main Application (`main.go`)**:
   - Implements a terminal-based chat interface between the user and Claude
   - Defines two tools that Claude can utilize:
     - `read_file`: Reads the contents of a specified file
     - `list_files`: Lists files and directories at a specified path
   - Manages the conversation flow, including handling tool calls

2. **Dependencies** (from `go.mod`):
   - `github.com/anthropics/anthropic-sdk-go` (v1.6.2): The official Anthropic SDK for Go
   - `github.com/invopop/jsonschema` (v0.13.0): Used to generate JSON schemas for tool definitions

3. **Project Structure**:
   - Uses Go modules (Go 1.23.4)
   - Follows a simple architecture with a main package

### How It Works:

1. The application creates a chat loop where:
   - The user enters text input
   - Input is sent to Claude via the Anthropic API
   - Claude responds with text or tool calls
   - If Claude calls a tool, the application executes it and returns the result to Claude
   - The conversation continues with all context maintained

2. The tools enable Claude to:
   - Read files to analyze their contents
   - List directory contents to understand the project structure

According to the README, this code is an implementation example from an article titled "How to Build an Agent" (ampcode.com/how-to-build-an-agent), demonstrating how to create a function-calling agent with Claude.

Would you like me to explain any specific part of the code or provide more information about how the agent works?
You:

```
