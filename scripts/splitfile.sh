#!/bin/bash
# Run the commands in the background
echo "rebuild filesplit binary"

make filesplit
echo "join all files starting with input in testdata directory"
./filesplit -count 10 -input testdata/joined.txt -output testdata/split.txt
# Output a message indicating that all commands have completed
echo "All commands have completed"