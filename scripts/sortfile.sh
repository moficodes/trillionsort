#!/bin/bash
# Run the commands in the background
echo "rebuild sortfile binary"

make sortfile
for i in {0..9}; do
  # generate string in the form of 0001 from 1
  # https://stackoverflow.com/questions/11266400/bash-scripting-formatting-integer-output-with-leading-zeros
  ./sortfile -input testdata/split.txt -index $i -output testdata/sorted.txt
done

# Output a message indicating that all commands have completed
echo "All commands have completed"