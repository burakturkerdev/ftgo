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

## Disclaimer
FTGO is an open-source project designed for fast and secure file transfer. It utilizes a custom security layer with TCP rather than TLS for encryption, providing end-to-end encryption for file transfers. However, it's important to note that while FTGO strives to ensure security, as an open-source project, there's no absolute guarantee of security. While FTGO endeavors to provide a secure file transfer solution, users should exercise caution when transferring sensitive or important files. If security is a primary concern, it's advisable to explore alternative solutions or consult with security experts.

## License

This project is licensed under the [GNU General Public License v3.0](LICENSE), which means that everyone is free to view, modify, distribute, and use the software for non-commercial purposes. Any derivative works must also be licensed under the GNU GPL. For more details, please refer to the [LICENSE](LICENSE) file.
