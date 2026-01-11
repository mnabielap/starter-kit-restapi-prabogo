# 🚀 Prabogo: Hexagonal Go REST API Starter Kit

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org/doc/devel/release.html)
[![Architecture](https://img.shields.io/badge/Architecture-Hexagonal-purple)](https://alistair.cockburn.us/hexagonal-architecture/)
[![Framework](https://img.shields.io/badge/Framework-Fiber-black)](https://gofiber.io/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**Prabogo** is a robust, production-ready Go framework designed to simplify project development using **Hexagonal Architecture** (Ports and Adapters). It provides an interactive command interface, built-in AI assistance instructions, and a suite of pre-configured tools for modern web development.

---

## ✨ Key Features

- **🏗 Hexagonal Architecture**: Clean separation between Domain, Ports, and Adapters.
- **⚡ High Performance**: Built on top of **Fiber** (Fastest Go HTTP engine).
- **🔐 Advanced Authentication**:
  - JWT Access & Refresh Token rotation.
  - Role-Based Access Control (**RBAC**) (Admin vs User).
  - Password Reset & Email Verification flows.
- **💾 Database & SQL**:
  - **PostgreSQL** integration.
  - **Goqu** for type-safe, fluent SQL query building (No heavy ORM).
  - Auto-migrations via **Goose**.
- **🐇 Event Driven**: RabbitMQ integration for asynchronous messaging.
- **⚡ Caching**: Redis integration for high-speed data access.
- **🛠 Code Generation**: Powerful `Makefile` to generate Models, Adapters, and Ports instantly.
- **🧪 Automated Testing**: Python-based API test suite included (No Postman needed!).

---

## 📂 Project Structure

```text
starter-kit-restapi-prabogo/
├── cmd/                    # Application entry point
├── internal/
│   ├── domain/             # Core Business Logic (User, Auth, Client)
│   ├── model/              # Data Structures / Entities
│   ├── port/               # Interfaces (Inbound/Outbound)
│   ├── adapter/            # Implementations (Fiber, Postgres, Redis, etc.)
│   └── migration/          # Database migration scripts
├── utils/                  # Shared utilities (JWT, Password, Logger)
├── api_tests/              # Python automated test scripts
├── docker-compose.yml      # Docker orchestration
├── Makefile                # Command runner and code generator
└── README.md               # Documentation
```

---

## ⚙️ Configuration

Before running the application (Local or Docker), you **must** configure your environment variables.

1. Copy the example file:
   ```bash
   cp .env.example .env
   ```

2. Edit `.env` and fill in your details. **Crucial variables for Auth:**

   ```properties
   # Server
   SERVER_PORT=8000
   APP_MODE=debug

   # Security
   INTERNAL_KEY=your_secure_internal_key
   JWT_SECRET=change_this_to_a_super_secure_secret
   JWT_ACCESS_EXPIRATION_MINUTES=30
   JWT_REFRESH_EXPIRATION_DAYS=30

   # Database
   DATABASE_HOST=localhost   # Use 'postgres' if running inside Docker
   DATABASE_USER=prabogo
   DATABASE_PASSWORD=changeme
   DATABASE_NAME=prabogo

   # SMTP (Email)
   SMTP_HOST=smtp.example.com
   SMTP_PORT=587
   SMTP_USERNAME=user@example.com
   SMTP_PASSWORD=password
   ```

---

## 🐳 Running with Docker (Recommended)

This is the easiest way to get started. You do not need Go or PostgreSQL installed on your machine.

### 1. Load Environment Variables
Docker Compose automatically loads variables from your `.env` file. **Ensure your `.env` file exists** (see Configuration section).

> **Note:** Inside `.env`, set `DATABASE_HOST=postgres`, `CACHE_HOST=redis`, and `MESSAGE_HOST=rabbitmq` so the container can find the services.

### 2. Start Services
Build and start the containers in the background:
```bash
docker-compose up -d --build
```

### 🛠 Docker Management Commands

Here is a cheat sheet for managing your Docker environment:

| Action | Command |
| :--- | :--- |
| **👀 View Logs** | `docker-compose logs -f app` |
| **🛑 Stop Containers** | `docker-compose stop` |
| **▶️ Start Containers** | `docker-compose start` |
| **♻️ Restart App** | `docker-compose restart app` |
| **🗑 Remove Containers** | `docker-compose down` (Stops and removes containers & networks) |
| **📦 List Volumes** | `docker volume ls` |
| **⚠️ Delete Volumes** | `docker-compose down -v` <br>*(WARNING: Permanently deletes database data!)* |

---

## 🧪 API Testing (Python Suite)

Forget manual Postman collections! This project comes with a comprehensive **Python Test Suite** located in `api_tests/`.

### Prerequisites
*   Python 3.x
*   `requests` library (`pip install requests`)

### How to Run Tests
Navigate to the `api_tests` folder and run the scripts sequentially. They automatically handle token management (saving tokens to `secrets.json`).

**1. Authentication Tests**
```bash
# Register a new user
python api_tests/A1.auth_register.py

# Login (Saves tokens)
python api_tests/A2.auth_login.py

# Refresh Token
python api_tests/A3.auth_refresh.py
```

---

## 🛠 Makefile & Code Generation

Prabogo features an interactive CLI for generating boilerplate code, keeping your Hexagonal Architecture clean.

### Interactive Mode
Simply run:
```bash
make run
```
*Select an option from the menu to generate models, adapters, or migrations.*

### Direct Commands
*   **Generate Model:** `make model VAL=product`
*   **Generate Migration:** `make migration-postgres VAL=add_product_table`
*   **Generate Handler (Fiber):** `make inbound-http-fiber VAL=product`
*   **Generate Repository (Postgres):** `make outbound-database-postgres VAL=product`

---

## 🛡 Security Features

*   **Role-Based Access Control (RBAC):** Middleware ensures only users with `role: admin` can perform sensitive operations.
*   **Argon2/Bcrypt:** Password hashing implementation (via `utils/password`).
*   **JWT Security:** Short-lived Access Tokens and long-lived Refresh Tokens.
*   **Input Validation:** Strict struct validation on all incoming requests.

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<p align="center">
  Built with ❤️ by <b>Moch Dieqy Dzulqaidar</b>
</p>