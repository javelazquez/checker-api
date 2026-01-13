#!/usr/bin/env bash
set -e

echo "🚀 Creating DynamoDB table: comparisons"

# Check if table already exists
if awslocal dynamodb describe-table --table-name comparisons 2>/dev/null; then
  echo "ℹ️  Table 'comparisons' already exists"
else
  awslocal dynamodb create-table \
    --table-name comparisons \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST
  echo "✅ Table 'comparisons' created successfully"
fi

echo "✅ DynamoDB tables available:"
awslocal dynamodb list-tables
