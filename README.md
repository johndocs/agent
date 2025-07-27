# agent
**Code copied from [How to Build an Agent](https://ampcode.com/how-to-build-an-agent)**

## Summary
This Go project demonstrates how to build a command-line agent that uses the Claude 3.7 Sonnet model
and the official Anthropic Go SDK to interact with a local file system.

The agent uses an `Agent` struct to manage the conversation and tool execution.

Key implementation details for Go developers include:

* **Tool Definition**: The agent defines two tools for Claude to use: `list_files` and `read_file`.
 It uses the `github.com/invopop/jsonschema` library to generate the necessary JSON schemas for
 these tool definitions.
* **API Integration**: The project uses `github.com/anthropics/anthropic-sdk-go` (v1.6.2) to
communicate with the Claude API.
* **Execution Loop**: The core logic handles the back-and-forth between the user and the model. When
Claude requests a tool, the Go application executes the file system operation, sends the result back
to the model, and then Claude uses that information to formulate its final response to the user.
* **User Interface**: The application provides a simple, color-coded terminal interface for the chat.
It can accept initial prompts as command-line arguments or interactively.


## Example Usage and Variation Between Runs

This shows the variation in output between runs. The log files are in the [logs](logs) directory.
[diff.txt](logs/diff.txt) shows the diff between the 2 runs.

```
% find logs -type f -exec wc -l {} +
     207 logs.directory_contents/log.e.2.txt
      76 logs.directory_contents/log.1.txt
     181 logs.directory_contents/log.e.1.txt
     123 logs.directory_contents/log.2.txt
     191 logs.directory_contents/diff.txt
     778 total
```

### Command (See [run_directory_contents.sh.sh](run_directory_contents.sh.sh))
```
% go run . -d "What is this directory?" "Please examine the contents of all files" "Provide a short overview of this project for the README.txt ## Summary"
```

### Output first run

    Chat with Claude (use 'ctrl-c' to quit)
    [94mYou[0m: [93mClaude[0m: I'll help you understand what's in this directory by listing the files and directories present.
    [92mtool[0m: list_files({})
    ListFiles result ".": ["README.md","go.mod","go.sum","main.go","prompts.txt","utils.go"]
    [93mClaude[0m: Based on the files listed, this appears to be a Go language project directory. It contains:

    1. `README.md` - A markdown file that typically contains project documentation
    2. `go.mod` and `go.sum` - Go module files that manage dependencies
    3. `main.go` - The main Go source code file
    4. `utils.go` - A utility Go source code file
    5. `prompts.txt` - A text file that likely contains prompts or instructions

    Would you like me to examine any of these files in more detail to better understand the purpose of this project?
    [94mYou[0m: [93mClaude[0m: I'll examine the contents of all the files in this directory to give you a comprehensive understanding of the project.
    [92mtool[0m: read_file({"path":"README.md"})
    [92mtool[0m: read_file({"path":"main.go"})
    [92mtool[0m: read_file({"path":"go.mod"})
    [92mtool[0m: read_file({"path":"go.sum"})
    [92mtool[0m: read_file({"path":"prompts.txt"})
    [92mtool[0m: read_file({"path":"utils.go"})
    [93mClaude[0m: Based on examining all the files, I can provide you with a comprehensive overview of this directory:

    ## Agent Project Overview

    This directory contains a Go application that implements an AI assistant powered by Anthropic's Claude 3.7 Sonnet model. The agent is designed to interact with a local file system through defined tools while maintaining a conversational interface with users.

    ### Core Components

    1. **main.go**: The primary application file containing:
    - The `main` function that initializes the agent
    - The `Agent` struct and its methods for managing conversation flow
    - Tool definitions and implementations for file system interactions
    - Functions for generating JSON schemas for tool definitions

    2. **utils.go**: Contains utility functions for:
    - Context handling with signal cancellation for graceful shutdown
    - Debug logging and content description
    - JSON key counting for monitoring API responses
    - Prompt saving functionality

    3. **Go Module Files**:
    - **go.mod**: Defines the module and its dependencies
    - **go.sum**: Contains checksums for dependency verification

    ### Dependencies

    The project relies on two main external libraries:
    - `github.com/anthropics/anthropic-sdk-go v1.6.2`: The official Anthropic SDK for Go to interact with Claude API
    - `github.com/invopop/jsonschema v0.13.0`: Used to generate JSON schemas for tool definitions

    ### Tool Functionality

    The agent implements two file system tools:

    1. **read_file**: Allows Claude to read the contents of a specified file
    - Input: A relative file path
    - Output: The contents of the file as a string

    2. **list_files**: Enables Claude to list files and directories at a given path
    - Input: An optional path (defaults to current directory)
    - Output: A JSON array of file paths

    ### Application Flow

    1. The application starts and initializes the Anthropic client
    2. It sets up a user input mechanism (either from command-line arguments or interactive input)
    3. The agent runs in a loop where:
    - It reads user input
    - Sends the message to Claude
    - Processes Claude's response:
        - Shows text responses directly to the user
        - Executes tool invocations when Claude requests to use a tool
        - Returns tool results back to Claude
    4. The loop continues until the user terminates the program (e.g., with Ctrl+C)

    ### Additional Features

    - **Graceful Shutdown**: The application handles termination signals for clean exits
    - **Debug Mode**: Optional debug logging activated with the `-d` flag
    - **Prompt Saving**: User inputs are saved to `prompts.txt` for reference
    - **Colorized Output**: Uses terminal color codes to distinguish between user, Claude, and tool interactions

    ### File Structure Summary

    - `main.go`: Main application code
    - `utils.go`: Utility functions for the application
    - `go.mod` & `go.sum`: Go module definition and dependency checksums
    - `README.md`: Project documentation showing this is from a tutorial on "How to Build an Agent"
    - `prompts.txt`: Record of user inputs from previous runs

    This project serves as a practical example of building an AI agent with tool-use capabilities, showcasing how to integrate Claude with external systems (in this case, the file system) while maintaining a natural conversation flow.
    [94mYou[0m: [93mClaude[0m: Here's a short overview for a README.txt file:

    ## Summary

    This Go project demonstrates how to build a command-line agent that uses the Claude 3.7 Sonnet model and the official Anthropic Go SDK to interact with a local file system.

    The agent uses a main `Agent` struct to manage the conversation and tool execution. It provides Claude with two file system tools: `list_files` and `read_file`, allowing the AI to explore and read files while maintaining a natural conversation with the user.

    Key implementation details include:
    - Tool definition using JSON schemas
    - API integration with Anthropic's Claude
    - Conversation flow management
    - Graceful handling of interruptions and signals
    - Simple, color-coded terminal interface

    The application can accept initial prompts as command-line arguments or operate in an interactive mode, making it flexible for various use cases. It demonstrates a pattern for creating AI agents with access to external tools in a controlled, well-defined manner.
    [94mYou[0m:


### Output second run

    Chat with Claude (use 'ctrl-c' to quit)
    [94mYou[0m: [93mClaude[0m: To understand what this directory contains, I'll list the files and directories at the current path.
    [92mtool[0m: list_files({})
    ListFiles result ".": ["README.md","go.mod","go.sum","main.go","prompts.txt","utils.go"]
    [93mClaude[0m: This appears to be a Go programming project directory. It contains:

    1. `README.md` - Likely contains information about the project
    2. `go.mod` and `go.sum` - Go module definition files that manage dependencies
    3. `main.go` - The main Go source file for the application
    4. `prompts.txt` - Possibly contains text prompts used by the application
    5. `utils.go` - A Go file likely containing utility functions

    Would you like me to examine any of these files in more detail to better understand the project?
    [94mYou[0m: [93mClaude[0m: I'll examine the contents of all files in this directory to give you a comprehensive overview.

    Let's start with the README.md:
    [92mtool[0m: read_file({"path":"README.md"})
    [93mClaude[0m: Now let's look at main.go:
    [92mtool[0m: read_file({"path":"main.go"})
    [93mClaude[0m: Let's look at go.mod:
    [92mtool[0m: read_file({"path":"go.mod"})
    [93mClaude[0m: Let's look at go.sum:
    [92mtool[0m: read_file({"path":"go.sum"})
    [93mClaude[0m: Let's look at prompts.txt:
    [92mtool[0m: read_file({"path":"prompts.txt"})
    [93mClaude[0m: Let's check utils.go:
    [92mtool[0m: read_file({"path":"utils.go"})
    [93mClaude[0m: Based on examining all the files in this directory, here's a comprehensive overview:

    # Project Overview: AI Agent with File System Tools

    This is a Go project that demonstrates how to build a command-line AI agent powered by Anthropic's Claude 3.7 Sonnet model. The agent can interact with the local file system through defined tools while maintaining a conversational interface with the user.

    ## Key Components

    ### 1. Main Application (main.go)
    - The core of the application that:
    - Sets up a terminal-based chat interface with Claude
    - Defines two file system tools: `read_file` and `list_files`
    - Handles the conversation flow between users and Claude
    - Processes tool invocations and returns results to Claude
    - Implements command-line argument processing for batch prompts

    ### 2. Agent Structure
    - The application uses an `Agent` struct to manage:
    - Claude client connection
    - User message handling
    - Tool definitions and execution
    - Conversation state management

    ### 3. File System Tools
    - Two tools are implemented:
    - `read_file`: Reads the contents of a specified file
    - `list_files`: Lists files and directories at a given path
    - Each tool has a defined JSON schema using the `github.com/invopop/jsonschema` library

    ### 4. Utility Functions (utils.go)
    - Signal handling for graceful shutdown
    - JSON key counting and content description for debugging
    - Prompt saving to maintain a history of user inputs

    ### 5. Dependencies
    - Primary dependencies:
    - `github.com/anthropics/anthropic-sdk-go v1.6.2`: Official Anthropic SDK for Go
    - `github.com/invopop/jsonschema v0.13.0`: Used for JSON schema generation

    ## How It Works

    1. The application starts a chat session with Claude 3.7 Sonnet
    2. Users type messages in the terminal (or prompts are provided as command-line arguments)
    3. The agent:
    - Sends user messages to Claude
    - Handles tool invocation requests from Claude
    - Executes file system operations when requested
    - Returns tool results back to Claude
    - Presents Claude's final responses to the user
    4. Prompts are saved to a file (`prompts.txt`) for future reference
    5. The application gracefully handles interruptions with Ctrl+C

    ## Technical Implementation Details

    - Tool definitions use JSON schema generation to properly format for Claude's API
    - The app implements a structured conversation flow with state management
    - Debug mode is available for detailed logging of API interactions
    - Signal handling ensures proper cleanup when the application is terminated
    - The code demonstrates proper use of Go concurrency patterns with context management

    ## Usage

    The application can be used in two modes:
    1. Interactive mode: `go run .`
    2. Batch mode with predefined prompts: `go run . "What is this directory?" "Please examine the contents of all files"`

    The application is self-referential - it can examine its own code and structure using the file system tools it has access to, as demonstrated in the README examples and the existing prompts.txt file.

    This project serves as a practical demonstration of how to build an AI agent with tool-use capabilities, focusing specifically on file system interactions through a conversational interface with Claude.
    [94mYou[0m: [93mClaude[0m: Here's a short overview for the README.txt ## Summary section:

    ## Summary

    This Go project demonstrates how to build a command-line agent that uses the Claude 3.7 Sonnet model and the official Anthropic Go SDK to interact with a local file system.

    The agent uses a main `Agent` struct to manage the conversation and tool execution.

    Key implementation details for Go developers include:

    * **Tool Definition**: The agent defines two tools for Claude to use: `list_files` and `read_file`. It uses the `github.com/invopop/jsonschema` library to generate the necessary JSON schemas for these tool definitions.
    * **API Integration**: The project uses `github.com/anthropics/anthropic-sdk-go` (v1.6.2) to communicate with the Claude API.
    * **Execution Loop**: The core logic handles the back-and-forth between the user and the model. When Claude requests a tool, the Go application executes the file system operation, sends the result back to the model, and then Claude uses that information to formulate its final response to the user.
    * **User Interface**: The application provides a simple, color-coded terminal interface for the chat. It can accept initial prompts as command-line arguments or interactively.
    [94mYou[0m:


## Summary comparison

### First run

    ## Summary

    This Go project demonstrates how to build a command-line agent that uses the Claude 3.7 Sonnet model and the official Anthropic Go SDK to interact with a local file system.

    The agent uses a main `Agent` struct to manage the conversation and tool execution. It provides Claude with two file system tools: `list_files` and `read_file`, allowing the AI to explore and read files while maintaining a natural conversation with the user.

    Key implementation details include:
    - Tool definition using JSON schemas
    - API integration with Anthropic's Claude
    - Conversation flow management
    - Graceful handling of interruptions and signals
    - Simple, color-coded terminal interface

    The application can accept initial prompts as command-line arguments or operate in an interactive mode, making it flexible for various use cases. It demonstrates a pattern for creating AI agents with access to external tools in a controlled, well-defined manner.

## Second run

    ## Summary

    This Go project demonstrates how to build a command-line agent that uses the Claude 3.7 Sonnet model and the official Anthropic Go SDK to interact with a local file system.

    The agent uses a main `Agent` struct to manage the conversation and tool execution.

    Key implementation details for Go developers include:

    * **Tool Definition**: The agent defines two tools for Claude to use: `list_files` and `read_file`. It uses the `github.com/invopop/jsonschema` library to generate the necessary JSON schemas for these tool definitions.
    * **API Integration**: The project uses `github.com/anthropics/anthropic-sdk-go` (v1.6.2) to communicate with the Claude API.
    * **Execution Loop**: The core logic handles the back-and-forth between the user and the model. When Claude requests a tool, the Go application executes the file system operation, sends the result back to the model, and then Claude uses that information to formulate its final response to the user.
    * **User Interface**: The application provides a simple, color-coded terminal interface for the chat. It can accept initial prompts as command-line arguments or interactively.
