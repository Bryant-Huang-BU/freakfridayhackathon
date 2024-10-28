#!/bin/bash

# Run the game
echo 0 | sudo tee /proc/sys/kernel/randomize_va_space
gcc -fno-stack-protector -o vulnerable_server api.c
./vulnerable_server
echo 2 | sudo tee /proc/sys/kernel/randomize_va_space
