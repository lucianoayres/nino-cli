#!/bin/bash
echo "Generate a short conventional commit message, 
including prefix, to summarize the changes: 
$(git diff -U0). 
Commit message must use verb on infinive. 
Just output the prefix and message."