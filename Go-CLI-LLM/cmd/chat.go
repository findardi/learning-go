/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"go-cli-llm/tools"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/henomis/lingoose/llm/groq"
	"github.com/henomis/lingoose/thread"
	"github.com/spf13/cobra"
)

// Helper function to convert string to *string
func stringPtr(s string) *string {
	return &s
}

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "LLM",
	Long:  `Simple Implementing LLM using langchain`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

		go func() {
			<-sigChan
			fmt.Println("\nInterrupt signal. Bye...")
			os.Exit(0)
		}()

		// Initialize web search tool
		searchTool := tools.NewWebSearchTool()

		// Initialize LLM with tools
		llm := groq.New().
			WithModel("llama-3.1-8b-instant").
			WithMaxTokens(2048).
			WithTemperature(0.7).
			WithTools(searchTool).
			WithToolChoice(stringPtr("auto"))

		fmt.Print("Input initial system prompt (or press Enter for default): ")
		initial, _ := reader.ReadString('\n')
		initial = strings.TrimSpace(initial)

		if initial == "" {
			initial = "You are a helpful AI assistant with web search capabilities. When users ask about current events, recent information, or anything you're unsure about, use the web_search tool to find accurate information."
		}

		myThread := thread.New().AddMessage(
			thread.NewSystemMessage().AddContent(
				thread.NewTextContent(initial),
			),
		)

		fmt.Println("🤖 Chat started! Type 'exit', 'quit', or 'q' to exit.")
		fmt.Println("📡 Web search is enabled - ask me anything!")

		for {
			fmt.Print("You: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				continue
			}

			input = strings.TrimSpace(input)

			switch strings.ToLower(input) {
			case "exit", "quit", "q":
				fmt.Println("👋 Bye...")
				os.Exit(0)
			case "":
				continue
			default:
				// Add user message to thread
				myThread.AddMessage(
					thread.NewUserMessage().AddContent(
						thread.NewTextContent(input),
					),
				)

				// Tool execution loop
				maxIterations := 5 // Prevent infinite loops
				iteration := 0

				for iteration < maxIterations {
					iteration++

					// Generate response
					err := llm.Generate(context.Background(), myThread)
					if err != nil {
						fmt.Printf("❌ Error: %v\n", err)
						break
					}

					lastMsg := myThread.LastMessage()
					if lastMsg == nil {
						fmt.Println("❌ No response from LLM")
						break
					}

					// Check for tool calls
					hasToolCalls := false
					var toolCallsToExecute []thread.Content

					for _, content := range lastMsg.Contents {
						if content.Type == thread.ContentTypeToolCall {
							hasToolCalls = true
							toolCallsToExecute = append(toolCallsToExecute, *content)
						}
					}

					// If no tool calls, print response and exit loop
					if !hasToolCalls {
						fmt.Print("AI: ")
						for _, content := range lastMsg.Contents {
							if content.Type == thread.ContentTypeText {
								fmt.Println(content.Data)
							}
						}
						break
					}

					// Execute tool calls
					fmt.Println("🔍 Searching the web...")
					for _, toolCall := range toolCallsToExecute {
						toolCallDataList := toolCall.AsToolCallData()
						if toolCallDataList == nil {
							continue
						}

						for _, toolCallData := range toolCallDataList {
							toolName := toolCallData.Name
							toolCallID := toolCallData.ID

							// Execute the tool
							result, err := searchTool.Execute(toolCallData.Arguments)
							if err != nil {
								result = fmt.Sprintf("Error executing search: %v", err)
							}

							// Add tool result to thread
							myThread.AddMessage(
								thread.NewToolMessage().
									AddContent(
										thread.NewToolResponseContent(thread.ToolResponseData{
											ID:     toolCallID,
											Name:   toolName,
											Result: result,
										}),
									),
							)
						}
					}
				}

				if iteration >= maxIterations {
					fmt.Println("⚠️  Maximum tool iterations reached")
				}

				fmt.Println() // Empty line for readability
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
