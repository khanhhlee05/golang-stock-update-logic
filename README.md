# Go Stock Update Lambda

A Lambda-based service that updates user portfolios with real-time stock prices.

## Overview

This service fetches current stock prices from the Finnhub API for all unique stocks held by users in the database, then updates each user's portfolio with current stock values and calculates account balances.

## Features

- Fetches unique stock symbols across all user portfolios
- Retrieves real-time stock prices from Finnhub API
- Concurrently updates user portfolios and financial data
- Tracks daily portfolio value history
- Performance measurement for execution time analysis

## Architecture

The project follows a modular structure:
- `cmd/main`: Contains the entry point for the lambda function
- `internal/db`: MongoDB connection and configuration
- `internal/handlers`: Business logic for fetching stocks and updating portfolios
- `internal/models`: Data structures for MongoDB documents

## Technology Stack

- Go 1.24.2
- MongoDB Atlas for data storage
- Finnhub API for stock price data

## Prerequisites

- Go 1.24.2 or higher
- MongoDB Atlas account
- Finnhub API key

## Configuration

The application connects to MongoDB using the connection string defined in `internal/db/mongodb.go`. Update this with your MongoDB credentials.

The Finnhub API key is defined in `internal/handlers/handlers.go`. Replace the example key with your own.

## Usage

Build and run the application:

```bash
cd cmd/main
go build -o stock-lambda .
./stock-lambda
```

# 🐳 AWS Lambda Go Function Deployment (Manual - Windows Console)

This guide walks you through manually building and deploying a Go-based AWS Lambda function using the AWS Console from a Windows machine. It uses the **provided.al2023** custom runtime and targets the **arm64** architecture.

---

## 📦 Prerequisites

- [Go installed](https://go.dev/dl)
- [AWS CLI installed](https://aws.amazon.com/cli/)
- AWS Account with Lambda permissions
- IAM Role with `AWSLambdaBasicExecutionRole`

---

## 🚀 Deployment Steps

### 1. Write the Lambda Handler

Create `main.go`:

```go
package main

import (
    "context"
    "github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context) (string, error) {
    return "Hello from Lambda!", nil
}

func main() {
    lambda.Start(handler)
}
```

### 2. Build for Linux ARM64
Open PowerShell:

```powershell
$env:GOOS = "linux"
$env:GOARCH = "arm64"
go build -o bootstrap main.go
```
✅ This creates a `bootstrap` binary required by AWS Lambda custom runtime.

### 3. Zip the Executable
```powershell
Compress-Archive -Path bootstrap -DestinationPath myFunction.zip
```
📁 Your folder should now contain:

* `main.go`
* `bootstrap`
* `myFunction.zip`

### 4. Deploy to AWS Lambda Console

1.  **Go to the AWS Lambda Console** and click **Create function**.
2.  Select **Author from scratch**.
3.  Configure the basic settings:
    *   **Function name**: Enter a name, e.g., `myFunction`.
    *   **Runtime**: Select `Provide your own bootstrap on Amazon Linux 2023`.
    *   **Architecture**: Select `arm64`.
    *   **Permissions**: Choose an existing execution role with `AWSLambdaBasicExecutionRole` permissions, or create a new one.
4.  Click **Create function**.
5.  Once the function is created, you'll be on its configuration page.
    *   Locate the **Runtime settings** panel (usually visible on the **Code** tab or under the **Configuration** tab, then **Function overview**). Click **Edit**.
    *   In the **Handler** field, enter `main.handler`.
    *   Click **Save**.
6.  Upload the deployment package:
    *   In the **Code source** section, select **Upload from**.
    *   Choose **.zip file**.
    *   Click **Upload** and select the `myFunction.zip` file you created earlier.
    *   Click **Save** (or **Deploy** if the button is labeled as such after uploading).

### 5. Update Timeout
1. Go to **Configuration** > **General configuration**
2. Click **Edit**
3. Set timeout to 10 seconds or more
4. Click **Save**

### 6. Test the Function
1. Go to the **Test** tab
2. Click **Test**

You should see:
```json
"Hello from Lambda!"
```
✅ Done!
You've successfully deployed a Go Lambda manually using AWS Console on Windows.
