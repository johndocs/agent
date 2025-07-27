package main

// Utilities for the main Agent that interacts with the Claude model using anthropic-sdk-go.

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/anthropics/anthropic-sdk-go"
)

// describeContents generates a string description of the content blocks in the conversation.
func describeContents(content []anthropic.ContentBlockParamUnion) string {
	if len(content) == 0 {
		return "no content"
	}
	descriptions := make([]string, len(content))
	for i, c := range content {
		descriptions[i] = describeContent(c)
	}
	return strings.Join(descriptions, ", ")
}

// describeContent generates a string representation of a single content block,
func describeContent(content anthropic.ContentBlockParamUnion) string {
	data, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		panic(fmt.Errorf("failed to marshal content: %w", err))
	}
	countJSONKeys(data)
	return string(data)
}

// countJSONKeys counts the keys in the JSON data and updates the global key counts.
func countJSONKeys(data []byte) {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		panic(fmt.Errorf("failed to unmarshal content data for key counting: %w", err))
	}

	countsTop := countKeys(m, false)
	countsAll := countKeys(m, true)
	addKeyCounts(keyCountsTop, countsTop)
	addKeyCounts(keyCountsAll, countsAll)
}

var (
	keyCountsTop = make(map[string]int) // Counts of keys in the top-level JSON object
	keyCountsAll = make(map[string]int) // Counts of keys in all JSON objects
)

// countKeys counts the keys in the given map and returns a map of key counts.
// If recursive is true, it will also count keys in nested "content" lists.
func countKeys(m map[string]any, recursive bool) map[string]int {
	counts := make(map[string]int, len(m))
	for k := range m {
		counts[k]++
		if recursive && k == "content" {
			nestedList, ok := m[k].([]any)
			if !ok {
				panic(fmt.Errorf("expected content to be a list of objects, got %T", m[k]))
			}
			for _, nestedAny := range nestedList {
				nested, ok := nestedAny.(map[string]any)
				if !ok {
					panic(fmt.Errorf("expected content to be a list of objects, got %T", nestedAny))
				}
				c := countKeys(nested, true)
				addKeyCounts(counts, c)
			}
		}
	}
	return counts
}

// addKeyCounts adds the counts from the source map to the destination map.
func addKeyCounts(m, o map[string]int) {
	for k, v := range o {
		m[k] += v
	}
}

// withSignalCancellation creates a context that is canceled when one of the specified signals is
// received. It's a common pattern for graceful shutdown.
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

// wrapAndSavePrompts creates a wrapper around a getUserMessage function to save the prompts to a file.
func wrapAndSavePrompts(innerGetUserMessage func() (string, bool), filePath string,
) func() (string, bool) {
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		panic(fmt.Errorf("error removing existing prompts file: %w", err))
	}
	return func() (string, bool) {
		prompt, ok := innerGetUserMessage()
		if ok && prompt != "" {
			f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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
