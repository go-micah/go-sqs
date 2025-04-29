# Go SQS Demo

A simple command-line application demonstrating how to work with Amazon SQS using the AWS SDK for Go v2.

## Prerequisites

- Go 1.21 or later
- AWS CLI configured with appropriate credentials
- AWS SAM CLI (for deploying the SQS queue)

## Setup

1. Clone this repository:
   ```
   git clone https://github.com/go-micah/go-sqs.git
   cd go-sqs
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Deploy the SQS queue using AWS SAM:
   ```
   sam deploy --guided
   ```
   
   This will create the SQS queue defined in `template.yaml`.

4. Note the queue URL from the SAM deployment outputs.

## Usage

### List all SQS queues

```
go run main.go list --region us-east-1
```

### Send a message to the queue

```
go run main.go send --queue https://sqs.us-east-1.amazonaws.com/123456789012/go-sqs-demo-queue --message "Hello, SQS!"
```

### Receive messages from the queue

```
go run main.go receive --queue https://sqs.us-east-1.amazonaws.com/123456789012/go-sqs-demo-queue
```

## Building the application

To build the application into a binary:

```
go build -o go-sqs
```

Then you can run it directly:

```
./go-sqs list
```

## AWS SAM Template

The `template.yaml` file defines an Amazon SQS queue with the following properties:

- Queue Name: go-sqs-demo-queue
- Visibility Timeout: 30 seconds
- Message Retention Period: 4 days

## Clean Up

To delete the SQS queue and associated resources:

```
sam delete
```
