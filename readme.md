# AIO Bot
Currently only able to do Solana Pump Fun coins.
Planning to expand to EVM, and then into Web2. Personal Bot for now.

## Running Tests
To run tests for backend, run the following command

```bash
  cd backend/
  go test ./...
```

To view tests coverage

```bash
  cd backend
  go test -cover ./...        
```

To run the program:
```bash
  cd backend
  go build ./cmd/server
  ./server
```
in the backend directory
