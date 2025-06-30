# Balance Service

A Go-based HTTP service for processing balance transactions with PostgreSQL database.

## Features

- Process user balance transactions (win/lose)
- Get current user balance
- Idempotent transaction processing
- Concurrent request handling (20-30 RPS)
- Docker containerization

## API Endpoints

### GET /user/{userId}/balance
Returns the current balance for a user.

**Response:**
```json
{
  "userId": 1,
  "balance": "100.00"
}
```

### POST /user/{userId}/transaction
Processes a balance transaction.

**Headers:**
- `Source-Type`: Must be one of `game`, `server`, or `payment`
- `Content-Type`: `application/json`

**Request Body:**
```json
{
  "state": "win",
  "amount": "10.15",
  "transactionId": "unique-transaction-id"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Transaction processed successfully"
}
```

## Quick Start

### Using Docker Compose (Recommended)

1. Clone the repository
2. Run the application:
```bash
docker-compose up -d
```

The service will be available at `http://localhost:8080`

### Manual Setup

1. Install dependencies:
```bash
go mod download
```

2. Set environment variables:
```bash
export DATABASE_URL="postgres://user:password@localhost:5432/balance?sslmode=disable"
export PORT="8080"
```

3. Run PostgreSQL and execute migrations

4. Start the application:
```bash
go run main.go
```

## Testing

The service comes with 3 predefined users (IDs: 1, 2, 3) with initial balances.

### Test User Balance
```bash
curl http://localhost:8080/user/1/balance
```

### Test Transaction
```bash
curl -X POST http://localhost:8080/user/1/transaction \
  -H "Source-Type: game" \
  -H "Content-Type: application/json" \
  -d '{
    "state": "win",
    "amount": "10.50",
    "transactionId": "test-tx-001"
  }'
```

### Test Win Transaction
```bash
curl -X POST http://localhost:8080/user/1/transaction \
  -H "Source-Type: game" \
  -H "Content-Type: application/json" \
  -d '{
    "state": "win",
    "amount": "25.00",
    "transactionId": "win-tx-001"
  }'
```

### Test Lose Transaction
```bash
curl -X POST http://localhost:8080/user/1/transaction \
  -H "Source-Type: game" \
  -H "Content-Type: application/json" \
  -d '{
    "state": "lose",
    "amount": "5.00",
    "transactionId": "lose-tx-001"
  }'
```

### Test Idempotency
```bash
# This should return an error as the transaction ID already exists
curl -X POST http://localhost:8080/user/1/transaction \
  -H "Source-Type: game" \
  -H "Content-Type: application/json" \
  -d '{
    "state": "win",
    "amount": "10.00",
    "transactionId": "win-tx-001"
  }'
```

## Architecture

- **Config**: Application configuration management
- **Database**: PostgreSQL connection and migration handling
- **Models**: Data structures for users and transactions
- **Repository**: Data access layer with transaction safety
- **Service**: Business logic layer
- **Handlers**: HTTP request handling

## Database Schema

### Users Table
- `id`: Primary key (BIGSERIAL)
- `balance`: User balance (DECIMAL(15,2))
- `created_at`, `updated_at`: Timestamps

### Transactions Table
- `id`: Primary key (BIGSERIAL) 
- `user_id`: Foreign key to users
- `transaction_id`: Unique transaction identifier
- `source_type`: Source of transaction (game/server/payment)
- `state`: win or lose
- `amount`: Transaction amount (DECIMAL(15,2))
- `created_at`: Timestamp

## Performance

The application is designed to handle 20-30 requests per second with:
- Database connection pooling
- Proper indexing on frequently queried columns
- Atomic transactions for balance updates
- Concurrent request handling via Gin framework

## Stopping the Service

```bash
docker-compose down
```

To remove volumes:
```bash
docker-compose down -v
```
