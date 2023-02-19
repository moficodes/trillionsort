#!/bin/bash
# Run the commands in the background
for i in {1..10}; do
  ./sortfile -input testdata/input.txt -index $i -output testdata/output_$i.txt
done

# Wait for all background processes to finish
wait

# Output a message indicating that all commands have completed
echo "All commands have completed"