#!/bin/bash

# Define the output file
OUTPUT_FILE="detailed_commit_log_2024-10-07_to_2024-11-30.txt"

# Specify the date range for the logs
START_DATE="2024-10-07"
END_DATE="2024-11-30"

# Initialize the output file
echo "Detailed Commit Log (From $START_DATE to $END_DATE)" > $OUTPUT_FILE
echo "====================================================" >> $OUTPUT_FILE

# Fetch commit hashes in the specified date range
COMMIT_HASHES=$(git log --since="$START_DATE" --until="$END_DATE" --pretty=format:"%h")

# Loop through each commit hash and append details to the output file
for COMMIT in $COMMIT_HASHES; do
    echo "----------------------------------------------------" >> $OUTPUT_FILE
    git log -1 --pretty=format:"Commit: %h%nAuthor: %an%nDate: %ad%nMessage: %s" --date=short $COMMIT >> $OUTPUT_FILE
    echo -e "\nChanges:" >> $OUTPUT_FILE
    git show --stat $COMMIT >> $OUTPUT_FILE
    echo -e "\nFull Diff:" >> $OUTPUT_FILE
    git show $COMMIT >> $OUTPUT_FILE
    echo "----------------------------------------------------" >> $OUTPUT_FILE
    echo -e "\n" >> $OUTPUT_FILE
done

# Notify the user
echo "Detailed commit log saved to $OUTPUT_FILE"

