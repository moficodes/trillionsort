#!/bin/bash
# Run the commands in the background
echo "rebuild joinfile binary"

make joinfile
echo "join all files starting with input in testdata directory"
./joinfile -dir testdata -pattern input -output testdata/joined.txt
# Output a message indicating that all commands have completed
echo "All commands have completed"