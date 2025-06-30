#!/bin/bash

echo "Testing Balance Service API"
echo "==========================="

BASE_URL="http://localhost:8080"

echo "1. Getting initial balance for user 1..."
curl -s "$BASE_URL/user/1/balance" | jq .
echo

echo "2. Processing a win transaction..."
curl -s -X POST "$BASE_URL/user/1/transaction" \
  -H "Source-Type: game" \
  -H "Content-Type: application/json" \
  -d '{
    "state": "win",
    "amount": "25.75",
    "transactionId": "test-win-001"
  }' | jq .
echo

echo "3. Getting balance after win..."
curl -s "$BASE_URL/user/1/balance" | jq .
echo

echo "4. Processing a lose transaction..."
curl -s -X POST "$BASE_URL/user/1/transaction" \
  -H "Source-Type: game" \
  -H "Content-Type: application/json" \
  -d '{
    "state": "lose",
    "amount": "10.25",
    "transactionId": "test-lose-001"
  }' | jq .
echo

echo "5. Getting final balance..."
curl -s "$BASE_URL/user/1/balance" | jq .
echo

echo "6. Testing idempotency (should fail)..."
curl -s -X POST "$BASE_URL/user/1/transaction" \
  -H "Source-Type: game" \
  -H "Content-Type: application/json" \
  -d '{
    "state": "win",
    "amount": "25.75",
    "transactionId": "test-win-001"
  }' | jq .
echo

echo "7. Testing insufficient balance (should fail)..."
curl -s -X POST "$BASE_URL/user/1/transaction" \
  -H "Source-Type: game" \
  -H "Content-Type: application/json" \
  -d '{
    "state": "lose",
    "amount": "1000.00",
    "transactionId": "test-fail-001"
  }' | jq .
echo

echo "Tests completed!"
