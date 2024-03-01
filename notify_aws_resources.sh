#!/bin/bash

DIRECTORY="/Users/ayushgupta/Downloads/Programming/checkAWSResources" # Change this to your path
cd "$DIRECTORY" || { echo "Directory not found"; exit 1; }

OUTPUT=$(./monitorAWSResourcesProgram)

LOG_FILE="/Users/ayushgupta/AWSResourcesLogs.txt"  # Change this to your path
echo "$(date) : $OUTPUT" >> $LOG_FILE
echo "-----------------------------------" >> $LOG_FILE

osascript -e "display notification \"$OUTPUT\" with title \"Go Program Output\""