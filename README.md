# Personal Finance Tracker API

## Models (4 total)

### 1. User

```json
{
  "id": "uuid",
  "email": "string",
  "name": "string",
  "currency": "string",
  "created_at": "datetime"
}
```

### 2. Account

```json
{
  "id": "uuid",
  "user_id": "uuid",
  "name": "string", // "Chase Checking", "Cash Wallet", "Credit Card"
  "type": "enum", // bank, cash, credit_card
  "balance": "decimal",
  "is_active": "boolean"
}
```

### 3. Category

```json
{
  "id": "uuid",
  "name": "string", // "Food", "Transport", "Salary"
  "type": "enum", // income, expense
  "color": "string"
}
```

### 4. Transaction

```json
{
  "id": "uuid",
  "account_id": "uuid",
  "category_id": "uuid",
  "amount": "decimal", // positive for income, negative for expenses
  "description": "string",
  "date": "date",
  "created_at": "datetime"
}
```

## REST Endpoints (16 total)

### Users (2)

- `GET /users/me` - Get current user
- `PUT /users/me` - Update user profile

### Accounts (4)

- `GET /accounts` - Get all user accounts (bank + cash + credit cards)
- `POST /accounts` - Create new account
- `PUT /accounts/{id}` - Update account
- `DELETE /accounts/{id}` - Delete account

### Categories (4)

- `GET /categories` - Get all categories
- `POST /categories` - Create category
- `PUT /categories/{id}` - Update category
- `DELETE /categories/{id}` - Delete category

### Transactions (6)

- `GET /transactions` - Get transactions (with date/account/category filters)
- `GET /transactions/{id}` - Get specific transaction
- `POST /transactions` - Add transaction (works for bank, cash, or credit card)
- `PUT /transactions/{id}` - Update transaction
- `DELETE /transactions/{id}` - Delete transaction
- `GET /transactions/summary` - Get spending summary by category

## How It Works

**Account Types:**

- **Bank**: Checking/savings accounts
- **Cash**: Physical cash wallet
- **Credit Card**: Credit card accounts

**Transaction Examples:**

- Cash coffee purchase: POST to `/transactions` with cash account_id, amount: -4.50
- Salary deposit: POST with bank account_id, amount: +3000.00
- Credit card purchase: POST with credit_card account_id, amount: -89.99

**Key Features:**

- All account types work the same way - just different `type` field
- Account balances auto-update when transactions are added/modified
- Simple filtering on transactions endpoint
- Transaction summary for spending analytics

This covers all CRUD operations while keeping the business logic straightforward!

