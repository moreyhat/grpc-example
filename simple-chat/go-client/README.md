# gRPC chat app client
Golang based client for gRPC based simple chat app.

# Commands references
The following sub commands are available.

| subcommand | explain |
| --- | --- |
|list-messages| List all messages |
|post-message| Post messages. Message to post is provided with '-m' option|

# Examples
- Post message  
  
  ```shell
  go run main.go post-message -m "Hello, World"
  ```

- List messages  
  
  ```shell
  go run main.go list-messages
  ```