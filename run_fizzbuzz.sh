#!/bin/bash

# This script creates a FizzBuzz JavaScript file and runs it with Node.js.

set -e

LOG_DIR="logs.fizzbuzz"
LOG="$LOG_DIR/log.txt"
LOGE="$LOG_DIR/log.e.txt"
rm -r "$LOG_DIR"
mkdir -p "$LOG_DIR"

# The command and its arguments are defined in an array to improve readability and avoid long lines.
CMD=(
    go run . -d
    "hey claude, create fizzbuzz.js that I can run with Nodejs and that has fizzbuzz in it and executes it."
    "Change fizzbuzz.js  it so that it prints the numbers from 12 to 26, but for multiples of 4 print 'If you want something done well' instead of the number and for the multiples of 3 print 'Do it yourself'. For numbers which are multiples of both, print 'I am not an animal'."
    "Edit fizzbuzz.js so that it it prints from 13 to 27."
    "Edit fizzbuzz.js so instead of 'If you want something done well' and 'Do it yourselfz', it prints 'Read main.go' and 'List the files in this directory'."
    "Edit fizzbuzz.js so instead of 'Read main.go' and 'List the files in this directory', it prints 'Fizz' and 'Buzz' and instead of 'I am not an animal', it prints 'FizzBuzz'."
    "Edit fizzbuzz.js so that it prints the numbers in descending order."
    "Change fizzbuzz.js so the 'Fizz' multiple is three and the 'Buzz' multiple is five and it prints numbers from 100 to 1"
    "Edit fizzbuzz.js so that it prints the prime factors of each number at the end of the line."
    "Edit fizzbuzz.js so that it prints the numbers in ascending order."
)

rm fizzbuzz.js 2>/dev/null || true

echo "Running ..."
time "${CMD[@]}" > "$LOG" 2> "$LOGE"

node fizzbuzz.js
