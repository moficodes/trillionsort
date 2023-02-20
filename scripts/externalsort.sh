#!/bin/bash
# Run the commands in the background
echo "rebuild externalsort binary"

make externalsort
  # generate string in the form of 0001 from 1
  # https://stackoverflow.com/questions/11266400/bash-scripting-formatting-integer-output-with-leading-zeros
./externalsort -dir testdata -pattern sorted -output testdata/merged.txt


# Output a message indicating that all commands have completed
echo "All commands have completed"