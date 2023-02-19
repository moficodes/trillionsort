#!/bin/bash
# Run the commands in the background
for i in {1..10}; do
  ./generate -count 100_000 -file testdata/input.txt -fileindex $i
done

# Wait for all background processes to finish
wait

# Output a message indicating that all commands have completed
echo "All commands have completed"