# payment-service API

Servicio Backend de test para gestionar transferencias entre cuentas

## Technologies

- **Go** (Gin, Gorm, JWT, etc.)
- **PostgreSQL** (Sqlite para pruebas)
- **Docker**

## Structure

```
payment-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── constant/
│   ├── controller/
│   ├── middleware/
│   ├── model/
│   ├── router/
│   ├── scheduler/
│   └── service/
├── config/
└── db/
```

## Running the Application

Create a `.env` file in the root directory with the following content:

```
DB_URL=postgres://user:password@localhost:5432/dbname
JWT_SECRET=your_jwt_secret
PORT=8080
APP_ENV=development
```

Run with docker:

```bash
docker-compose up --build
```

## Running Tests

Run the tests with:

```bash
go test ./...
```

## Authentication

The API uses JWT for authentication. You can obtain a token by sending a POST request to `/token/:account_id` with an existing account ID.

## API Endpoints

### Accounts

- **GET** `/accounts/:account_id/balance` - Get account balance

### Transfers

- **POST** `/transfers` - Create a transfer between accounts
- **POST** `/transfers/:transfer_id/webhook` - Webhook for transfer status updates

### Transfer Scheduler

Checks for pending transfers every 5 minutes and expires them if not completed.

## Functional Requirements

1. **Create Transfers with Pending Status**
2. **Expire Pending Transfers After 5 Minutes**
3. **Get Account Balance**
4. **Handle Transfer Webhooks**
5. **JWT Authentication**
6. **Unit Tests for Services**

## Author

Marcos Villanueva (tttato_)
