# SSE Example

This is a simple Server-Sent Events (SSE) implementation in Go with both server and client components.

## Requirements

- Go 1.13 or higher

## Features

- SSE endpoint at GET /event
- Sends 10 data messages and then ends the connection
- Each data message contains "#i" where i is the message number
- Every other message triggers an error with payload "error #i" where i is the error count
- After 10 data messages, the SSE response ends

## Project Structure

- `server` - Contains the SSE server implementation
- `client` - Contains a Go client for testing the SSE server
- `client-web` - Web-based client for testing the SSE server

## Running the Server

```bash
cd server
go run main.go
```

The server will start on port 8080.

## Running the Go Client

In a separate terminal:

```bash
cd client
go run main.go
```

The client will connect to the server, receive and display all events, and show a summary when the connection ends.

## Running the Web Client

Open the `client-web/index.html` in a web browser. The web client will connect to the SSE server and display all events in the browser console.
