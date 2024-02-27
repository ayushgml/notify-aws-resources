#!/bin/bash

cd /Users/ayushgupta/Downloads/Programming/checkAWSResources # Change this to your path

OUTPUT=$(./monitorAWSResourcesProgram)

osascript -e "display notification \"$OUTPUT\" with title \"Go Program Output\""