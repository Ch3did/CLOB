# CLOB Matching Engine -- README

## Overview

This project implements a simplified **Central Limit Order Book (CLOB)**
with: - Limit orders (BUY/SELL) - Matching engine (price/time
priority) - Trade generation - Balance updates per trade -
Authentication (JWT) - REST API (Gin) - PostgreSQL + GORM - Docker
Compose setup


------------------------------------------------------------------------

## ⚙️ Requirements

-   Go 1.22+
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

-   `(account_id, asset, amount)`

------------------------------------------------------------------------

## 🧪 Testing with Postman

A Postman collection exists in:

    docs/postman.json

Import and test all routes easily.

------------------------------------------------------------------------

## 📜 License

MIT

