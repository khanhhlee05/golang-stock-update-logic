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
./stock-lambda# golang-stock-update-logic
