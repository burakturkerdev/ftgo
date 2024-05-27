# FtGo
FtGo is server and client application designed for secure file transfer. It utilizes a custom protocol that ensures end-to-end encryption for data transmission. This repository contains the file transfer protocol, server application, and client application.

## Features
- **Custom Protocol**: FtGo uses a custom file transfer protocol designed for secure and efficient data transmission.
- **End-to-End Encryption**: All data transferred between the server and client is encrypted, ensuring privacy and security.
- **Server Application**: The server application allows users to host their own file transfer servers.
- **Client Application**: The client application enables users to connect to servers and transfer files securely.

## Client Commands
```bash
ftgo server add <server-name> <server-address> # Adds server to list for using again
ftgo server list # Lists all servers saved to client
ftgo server rm <server-name> # Removes server from client
ftgo package new <package-name> # Creates new package
ftgo package add <package-name> <file/directory> # Adds file or directory to package
ftgo package push <package-name> <server-name>/<server-address> # Pushs package to server
ftgo push <file/directory> <server-name>/<server-address> # Pushs file or directory to server
ftgo server connect <server-name>/<server-address> # Connects to server and lists all files and directories
ftgo cd <directory> # Cd into directory in server
ftgo pull <file/directory> <path> # Pulls directory or file from server to path (if path blank it will pull to default dir)
ftgo dir set <path> # Sets the default directory for pulling
ftgo dir get # Gets the default directory for pulling
```
## Server Commands
```bash
ftgosv serve # Starts server daemon for serving
ftgosv status # Lists all server information
ftgosv port add <port> # Adds port for listening
ftgosv port rm <port> # Removes port
ftgosv port list # Lists all ports
ftgosv dir set <path> # Sets the serving directory
ftgosv dir get # Gets the serving directory
ftgosv perm write set <perm> # Sets the perm for write operations
ftgosv perm write get # Gets the perm for write operations
ftgosv perm read set <perm> # Sets the perm for read operations
ftgosv perm read get # Gets the perm for write operations
ftgosv perm list <perm> # Lists all usable perms
ftgosv perm ip add <ip> # Adds ip for ip based perm
ftgosv perm ip rm <ip> # Removes ip from allowed ip list
ftgosv perm ip list # Lists allowed ip's
ftgosv perm password set <password> # Sets password for password authentication perm
```
## Protocol Api
The protocol API facilitates bidirectional communication between client and server, both using the same API. To establish a connection, the `CreateConnection` function requires a `net.Conn` interface. For client-to-server connection, we use `net.Dial`, while for server-side, we use `net.Listen`. Once connected, both can utilize the same methods.

However, there are some fundamental rules:

1. Each message between client and server is preceded by 4 bytes. The first 3 bytes are zeroes, and the 4th byte represents the message itself. You can find message definitions in `common/messages.go`. These messages are integers ranging from 0 to 127.
   
2. It's mandatory to add a message if the protocol function doesn't handle it for us.

3. After establishing a connection, we use the `Read` method to wait for data. Following this, we must utilize either `GetMessage` or `IgnoreMessage` method to extract the message from the buffer. Then, we can retrieve the desired data using methods like `GetString`, `GetJson`, etc.

4. For sending data, we have two types of methods for each data type (For example: JSON or raw string):
   - One doesn't require a message as a parameter and sends the data along with a blank message.
   - The other type wants a message parameter in the function body.

By adhering to these guidelines, we ensure smooth communication between the client and server using the protocol.

```go
// Creating connection object for using protocol
conn := common.CreateConnection(n)

// We are sending "LIST DIRECTORIES/FILES inside /hello" message to server
conn.SendMessageWithData(common.CListDirs, "/hello")

// Creating a message for holding the response
var m common.Message

// Firstly, we are reading the response from the server. After that, we are extracting the message to our message holder.
conn.Read().GetMessage(&m)

// Creating a file info slice
var infos []common.FileInfo

if m == common.Blank {
    // We don't need to authenticate because the message is blank, so let's directly extract JSON.
    conn.GetJson(&infos)

    // We can use our file info JSON.
    // ...
} else if m == common.SAuthenticate {
    // Server wants authentication so we are sending our password
    conn.SendString("testpassword") // This method adds a blank message for us; we don't need to pass any message.

    // Reading the response from the server
    conn.Read().GetMessage(&m)

    if m == common.Blank {
        // Message is blank; we successfully authenticated, let's get our file infos.
        conn.GetJson(&infos)

        // We can use our file info JSON.
        // ...
    } else { // We can look for other message cases as well. But we are skipping in tutorial
        // Message is not blank; we couldn't authenticate.
        // ...
    }
}
```

## Disclaimer
FTGO is an open-source project designed for fast and secure file transfer. It utilizes a custom security layer with TCP rather than TLS for encryption, providing end-to-end encryption for file transfers. However, it's important to note that while FTGO strives to ensure security, as an open-source project, there's no absolute guarantee of security. While FTGO endeavors to provide a secure file transfer solution, users should exercise caution when transferring sensitive or important files. If security is a primary concern, it's advisable to explore alternative solutions or consult with security experts.

## License

This project is licensed under the [GNU General Public License v3.0](LICENSE), which means that everyone is free to view, modify, distribute, and use the software for non-commercial purposes. Any derivative works must also be licensed under the GNU GPL. For more details, please refer to the [LICENSE](LICENSE) file.
