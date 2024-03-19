#!/bin/sh
aws lambda invoke --function-name $1 /dev/null --log-type Tail --query 'LogResult' --output text |  base64 -D
# Async invocation
#aws lambda invoke-async --function-name $1 --invoke-args input.json