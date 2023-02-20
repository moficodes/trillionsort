#!/bin/bash
# Run the commands in the background
echo "rebuild generate binary"

make generate

echo "generate test data"

for i in {0..9}; do
  ./generate -count 100_000 -file testdata/input.txt -fileindex $i
done

# Output a message indicating that all commands have completed
echo "All commands have completed"