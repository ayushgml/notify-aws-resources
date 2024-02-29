#!/bin/bash

cd /Users/ayushgupta/Downloads/Programming/checkAWSResources # Change this to your path

OUTPUT=$(./monitorAWSResourcesProgram)

echo "$(date) : $OUTPUT" >> /Users/ayushgupta/AWSResourcesLogs.txt  # Change this to your path

osascript -e "display notification \"$OUTPUT\" with title \"Go Program Output\""