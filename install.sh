#!/bin/bash

# Clone the repository
echo "Installation starting."

git clone https://github.com/burakturkerdev/ftgo

# Build and move client
cd ftgo/src/client
go build -o ftgo
sudo mv ftgo /usr/bin

# Build and move daemon
cd ../../src/server/daemon
go build -o ftgodaemon
sudo mv ftgodaemon /usr/bin

# Build and move server CLI
cd ../../src/server/cli
go build -o ftgosv
sudo mv ftgosv /usr/bin

echo "Installation completed."