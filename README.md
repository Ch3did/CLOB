# CLOB Matching Engine

## Overview

This project implements a simplified **Central Limit Order Book (CLOB)**
with: - Limit orders (BUY/SELL) - Matching engine (price/time
priority) - Trade generation - Balance updates per trade -
Authentication (JWT) - REST API (Gin) - PostgreSQL + GORM - Docker
Compose setup


------------------------------------------------------------------------

## ⚙️ Requirements

-   Go 1.24+
-   Docker + Docker Compose
-   PostgreSQL
-   Make (optional)

------------------------------------------------------------------------

## 🚀 Running with Docker

### 1. Build & start services

    docker-compose up --build

This starts:

-   `postgres` on port **5432**
-   `api` on port **8080**

### 2. API available at:

    http://localhost:8080

------------------------------------------------------------------------

## 🔐 Authentication Flow

### 1. Create account

    POST /accounts/register
    {
      "first_name": "John",
      "last_name": "Doe",
      "email": "john@example.com",
      "password": "123456"
    }

### 2. Login → get JWT

    POST /accounts/login
    {
      "email": "john@example.com",
      "password": "123456"
    }

Response:

    {
      "token": "YOUR.JWT.TOKEN"
    }

### 3. Authenticated requests use:

    Authorization: Bearer <token>

------------------------------------------------------------------------

## 🗄️ Database Models

### Account

-   user auth info

### Order

-   represents intent to BUY/SELL

### Trade

-   represents executed match between two orders

### Balance

-   control account balance


------------------------------------------------------------------------

## 🧩 Application Flow

1. **Register and authenticate**
   - Create an account via the auth endpoint.
   - Log in to obtain a JWT token.
   - All subsequent requests must include `Authorization: Bearer <token>`, which identifies the `account_id`.

2. **Initialize balances**
   - Before placing any order, the account must have sufficient balance:
     - For a **BUY** order on `BASE/QUOTE` (e.g. `BTC/BRL`):
       - You must have enough **QUOTE** (e.g. BRL) to cover `price * quantity`.
     - For a **SELL** order on `BASE/QUOTE`:
       - You must have enough **BASE** (e.g. BTC) to cover the `quantity` being sold.
   - Use the balance endpoints (e.g. credit) to deposit assets into the account.
   - Use the balance listing endpoint to verify current balances.

3. **Place orders**
   - Use the orders endpoint to create **BUY** or **SELL** limit orders with:
     - `instrument` (e.g. `BTC/BRL`)
     - `side` (`BUY` or `SELL`)
     - `price`
     - `quantity`
   - When an order is created:
     - The service validates that the account has enough balance:
       - BUY: checks and reserves/debits the required QUOTE (`price * quantity`).
       - SELL: checks and reserves/debits the required BASE (`quantity`).
     - The order is inserted into the order book.

4. **Matching and trades**
   - After insertion, the matching engine looks for opposite orders:
     - BUY matches with SELL where `sell.price <= buy.price`.
     - SELL matches with BUY where `buy.price >= sell.price`.
     - Matching respects price–time priority.
   - For each match, a `Trade` is created and balances are updated:
     - **Buyer** receives BASE and pays QUOTE.
     - **Seller** delivers BASE and receives QUOTE.
   - Order quantities and statuses are updated (e.g. partial fill, filled).

5. **Inspect state**
   - Use order endpoints to:
     - List your orders.
     - Fetch a single order by ID.
   - Use trade endpoints to:
     - List trades associated with your account.
     - Fetch a single trade by ID.
   - Use balance endpoints to:
     - View final balances after trades and cancellations.

6. **Cancel open orders**
   - You can cancel an open order owned by your account.
   - Only the remaining (unfilled) quantity is affected:
     - The reserved balance for that remaining quantity is released back to the account.
   - Filled or fully matched orders cannot be canceled.


------------------------------------------------------------------------

## 🧪 Testing with Postman

A Postman collection exists in:

    docs/postman.json

Import and test all routes easily.


------------------------------------------------------------------------

## 📜 License

MIT

