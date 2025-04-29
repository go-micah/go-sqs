package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/spf13/cobra"
)

var (
	queueURL string
	message  string
	region   string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "go-sqs",
		Short: "A simple CLI for Amazon SQS operations",
		Long:  `A command line interface to demonstrate Amazon SQS operations using AWS SDK for Go v2.`,
	}

	// Send message command
	var sendCmd = &cobra.Command{
		Use:   "send",
		Short: "Send a message to SQS queue",
		Run: func(cmd *cobra.Command, args []string) {
			if queueURL == "" || message == "" {
				log.Fatal("Queue URL and message are required")
			}
			sendMessage(queueURL, message)
		},
	}

	// Receive message command
	var receiveCmd = &cobra.Command{
		Use:   "receive",
		Short: "Receive messages from SQS queue",
		Run: func(cmd *cobra.Command, args []string) {
			if queueURL == "" {
				log.Fatal("Queue URL is required")
			}
			receiveMessages(queueURL)
		},
	}

	// List queues command
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all SQS queues",
		Run: func(cmd *cobra.Command, args []string) {
			listQueues()
		},
	}

	// Add flags
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "us-east-1", "AWS region")
	sendCmd.Flags().StringVarP(&queueURL, "queue", "q", "", "SQS Queue URL")
	sendCmd.Flags().StringVarP(&message, "message", "m", "", "Message to send")
	receiveCmd.Flags().StringVarP(&queueURL, "queue", "q", "", "SQS Queue URL")

	// Add commands to root
	rootCmd.AddCommand(sendCmd, receiveCmd, listCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// getSQSClient creates and returns an SQS client
func getSQSClient() *sqs.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), 
		config.WithRegion(region),
	)
	if err != nil {
		log.Fatalf("Unable to load SDK config: %v", err)
	}

	return sqs.NewFromConfig(cfg)
}

// sendMessage sends a message to the specified SQS queue
func sendMessage(queueURL, messageBody string) {
	client := getSQSClient()
	
	input := &sqs.SendMessageInput{
		QueueUrl:    &queueURL,
		MessageBody: &messageBody,
	}
	
	result, err := client.SendMessage(context.TODO(), input)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
	
	fmt.Printf("Message sent successfully! Message ID: %s\n", *result.MessageId)
}

// receiveMessages receives and processes messages from the specified SQS queue
func receiveMessages(queueURL string) {
	client := getSQSClient()
	
	input := &sqs.ReceiveMessageInput{
		QueueUrl:            &queueURL,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     20, // Long polling
	}
	
	result, err := client.ReceiveMessage(context.TODO(), input)
	if err != nil {
		log.Fatalf("Failed to receive messages: %v", err)
	}
	
	if len(result.Messages) == 0 {
		fmt.Println("No messages received")
		return
	}
	
	fmt.Printf("Received %d messages:\n", len(result.Messages))
	
	for i, msg := range result.Messages {
		fmt.Printf("%d. Message ID: %s\n   Body: %s\n", i+1, *msg.MessageId, *msg.Body)
		
		// Delete the message after processing
		deleteInput := &sqs.DeleteMessageInput{
			QueueUrl:      &queueURL,
			ReceiptHandle: msg.ReceiptHandle,
		}
		
		_, err := client.DeleteMessage(context.TODO(), deleteInput)
		if err != nil {
			log.Printf("Failed to delete message %s: %v", *msg.MessageId, err)
		} else {
			fmt.Printf("   Message deleted successfully\n")
		}
	}
}

// listQueues lists all SQS queues in the account
func listQueues() {
	client := getSQSClient()
	
	result, err := client.ListQueues(context.TODO(), &sqs.ListQueuesInput{})
	if err != nil {
		log.Fatalf("Failed to list queues: %v", err)
	}
	
	if len(result.QueueUrls) == 0 {
		fmt.Println("No queues found")
		return
	}
	
	fmt.Println("Available SQS Queues:")
	for i, url := range result.QueueUrls {
		fmt.Printf("%d. %s\n", i+1, url)
	}
}
