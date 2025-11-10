# 🧭 Stocky
![Go Version](https://img.shields.io/badge/Go-1.22+-blue)
![Gin](https://img.shields.io/badge/Framework-Gin-green)
![Database](https://img.shields.io/badge/PostgreSQL-15-blue)
![License](https://img.shields.io/badge/License-MIT-yellow)


> 🪙 **Stocky** is a backend simulation of a stock-based reward system that allows users to earn fractional Indian stock shares through onboarding and referrals, with live INR valuation updates and complete transaction traceability.


### Tech Stack
**Golang**, **Gin**, **PostgreSQL**, **Redis**, **Logrus**

---

## 📘 Overview
**Stocky** is a backend service that simulates a stock reward system where users earn fractional Indian stock shares (such as *RELIANCE*, *TCS*, or *HDFC*) through onboarding and referrals.  
The system tracks user rewards, maintains a ledger for every transaction, and recalculates INR portfolio values dynamically using simulated stock prices.

---

## ⚙️ Features
- RESTful APIs for rewards, portfolios, and user statistics  
- Idempotency handling to prevent duplicate rewards  
- Structured JSON logging using **Logrus**  
- Hourly simulated price updates  
- Persistent PostgreSQL-backed data model  
- Ledger-based transaction tracking for every reward  

---

## 🔄 Example Workflow

1. **User registers** → `/register`  
   → A new user is created and rewarded initial stock shares.  
2. **User refers a friend** → `/refer`  
   → Both users receive stock rewards.  
3. **System updates prices hourly** → `/price/updater` (background task)  
   → All holdings' INR values are recalculated.  
4. **User views dashboard data**  
   - `/today-stocks/:userId` → Today's rewards  
   - `/portfolio/:userId` → Total holdings with INR value  
   - `/historical-inr/:userId` → Daily reward history  
   - `/stats/:userId` → Summary of portfolio + daily rewards
     
---

## 🚀 API Endpoints

| Endpoint | Method | Description |
|-----------|---------|-------------|
| `/register` | **POST** | Register a new user and reward them |
| `/reward` | **POST** | Add a stock reward for a user |
| `/today-stocks/:userId` | **GET** | Fetch today’s rewarded stocks |
| `/historical-inr/:userId` | **GET** | Fetch INR value of past rewards |
| `/stats/:userId` | **GET** | Get today’s rewards + total INR portfolio |
| `/portfolio/:userId` | **GET** | Get current holdings grouped by stock |
| `/refer` | **POST** | Refer another user and reward both |

---

## 🧩 Sample Payloads

### **POST /reward**
#### Request
```json
{
  "user_id": 1,
  "symbol": "RELIANCE",
  "quantity": 2.5
}

{
  "id": 10,
  "user_id": 1,
  "stock_symbol": "RELIANCE",
  "quantity": 2.5,
  "rewarded_at": "2025-11-09T12:50:00Z"
}

```

---

## 🗃️ Database Schema

### **users**

| Column | Type | Description |
|--------|------|-------------|
| id | integer (PK) | User ID |
| name | varchar | User name |

---

### **rewards**

| Column | Type | Description |
|--------|------|-------------|
| id | integer (PK) | Reward ID |
| user_id | integer (FK → users.id) | Rewarded user |
| stock_symbol | varchar(20) | Stock symbol |
| quantity | numeric(18,6) | Quantity rewarded |
| rewarded_at | timestamp | Timestamp of reward |

---

### **ledger_entries**

| Column | Type | Description |
|--------|------|-------------|
| id | integer (PK) | Ledger entry ID |
| reward_id | integer (FK → rewards.id) | Linked reward |
| stock_symbol | varchar(20) | Stock symbol |
| stock_units | numeric(18,6) | Number of units |
| cash_outflow | numeric(18,4) | Cash equivalent of reward |
| brokerage_fee | numeric(18,4) | Brokerage charge |
| stt | numeric(18,4) | Securities transaction tax |
| gst | numeric(18,4) | GST on brokerage |
| created_at | timestamp | Creation time |

---

## 🧩 Brief Explanation of the Code

The project follows a modular clean architecture with clear separation between routes, services, and database layers.

- cmd/server → Application entrypoint that initializes environment variables, database, Redis, logger, and HTTP routes using Gin.

- internal/db → Handles PostgreSQL connection setup and schema initialization.

- internal/reward → Core business logic for stock rewards, ledger tracking, and user statistics.
Inserts reward events in the rewards table.
Automatically logs corresponding company expenses in ledger_entries.

- internal/price → Simulated stock price service that generates random stock prices for real-time INR valuation.

- internal/users → Manages user onboarding and registration.

- internal/referrals → Implements referral flow rewarding both inviter and invitee.

- pkg/logger → Configures Logrus for structured JSON logging across all services.

When a reward is created:

- A row is inserted into the rewards table with user ID, stock, quantity, and timestamp.
- A corresponding entry is added into the ledger_entries table capturing the brokerage, STT, GST, and total cash outflow.
- Stock price lookups are handled by the price service to compute INR values for stats and portfolio APIs.
- This modular design ensures scalability, clear maintainability, and accurate financial tracking through double-entry bookkeeping.

---

## 🧠 Edge Cases Handled

- **Duplicate reward prevention** — via request-level idempotency handling  
- **Stale price recovery** — retries with cached/fallback values  
- **Rounding precision** — enforced using `NUMERIC(18,4)` and controlled math rounding  
- **Hourly updates** — price refresh ensures accurate INR valuations  
- **Safe database writes** — transactional inserts for rewards and ledger entries  

---

## ⚡ Scalability

- Stateless REST API — horizontally scalable with load balancers  
- Redis caching for frequently accessed data (e.g., idempotency keys)  
- PostgreSQL ensures data integrity for rewards and ledger relationships  
- Background goroutines handle price updates asynchronously  
- Centralized structured logging (Logrus) for observability and monitoring  

--- 

## ⚙️ Setup Instructions

### 1. **Clone the repo**
```bash
git clone https://github.com/angad363/stocky-assignment.git
cd stocky-assignment
```
### 2. Create PostgreSQL database
```bash
CREATE DATABASE assignment;
```
### 3. Set up environment variables

Create a .env file in the project root:
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=assignment
REDIS_ADDR=localhost:6379
```
### 4. Run the server
```bash
go run ./cmd/server
```
### 5. Test APIs

Import the postman_collection.json file (included in the repo) into Postman to test all endpoints.

---

## 📈 Example Logs
```
INFO[2025-11-09 12:51:14] Starting Stocky Server initialization...
✅ Connected to PostgreSQL successfully!
✅ Connected to Redis successfully!
💹 Price updater started
🛣 Registering routes...
📡 All API routes registered
✅ Routes registered successfully
INFO[2025-11-09 12:51:14] Starting HTTP server port=8080
```
---

## 👨‍💻 Author

**Angad Anil Gosain**  
📧 angadgosain@gmail.com](mailto:angadgosain@gmail.com)  
🔗 [GitHub](https://github.com/angad363)
