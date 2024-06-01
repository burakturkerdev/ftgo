#!/bin/bash

# Clone the repository
echo "Installation starting."

# Build and move client
cd /src/client
go build -o ftgo
sudo mv ftgo /usr/bin

# Build and move daemon
cd ../server/daemon
go build -o ftgodaemon
sudo mv ftgodaemon /usr/bin

# Build and move server CLI
cd ../cli
go build -o ftgosv
sudo mv ftgosv /usr/bin

echo "Installation completed."
