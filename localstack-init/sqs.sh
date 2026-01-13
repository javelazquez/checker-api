#!/usr/bin/env bash
set -e

echo "🚀 Creating SQS queue: comparison-queue"

# Check if queue already exists
if awslocal sqs get-queue-url --queue-name comparison-queue 2>/dev/null; then
  echo "ℹ️  Queue 'comparison-queue' already exists"
else
  awslocal sqs create-queue \
    --queue-name comparison-queue \
    --attributes VisibilityTimeout=30
  echo "✅ Queue 'comparison-queue' created successfully"
fi

echo "✅ SQS queues available:"
awslocal sqs list-queues
