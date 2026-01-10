<a id="readme-top"></a>

<!-- PROJECT SHIELDS -->
[![Go Version][go-shield]][go-url]
[![Nuxt][nuxt-shield]][nuxt-url]
[![PostgreSQL][postgres-shield]][postgres-url]
[![Docker][docker-shield]][docker-url]

<!-- PROJECT HEADER -->
<br />
<div align="center">
  <h1 align="center">Expense Management System</h1>

  <p align="center">
    Full-stack expense management system for Indonesian fintech platform
    <br />
    Automated approval workflow with background payment processing
    <br />
    <br />
    <a href="#getting-started"><strong>Quick Start Guide »</strong></a>
    <br />
    <br />
    <a href="#api-documentation">API Docs</a>
    ·
    <a href="docs/openapi.yaml">Swagger/OpenAPI</a>
    ·
    <a href="#assumptions-made">Assumptions</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
        <li><a href="#key-features">Key Features</a></li>
        <li><a href="#business-rules">Business Rules</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#api-documentation">API Documentation</a></li>
    <li><a href="#architecture">Architecture</a></li>
    <li><a href="#database-schema">Database Schema</a></li>
    <li><a href="#design-decisions">Design Decisions & Trade-offs</a></li>
    <li><a href="#security">Security Considerations</a></li>
    <li><a href="#assumptions-made">Assumptions Made</a></li>
    <li><a href="#roadmap">Future Improvements</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->

## About The Project

This is a comprehensive expense management system designed for Indonesian fintech platforms, enabling employees to submit expenses in Rupiah (IDR) and managers to approve them with automated payment processing.

**Why This System?**

- **Automated Approvals**: Expenses below IDR 1,000,000 are automatically approved, reducing manager workload
- **Real-time Processing**: Background workers handle payment integration asynchronously for better performance
- **Audit Trail**: Complete tracking of all expense status changes for compliance
- **Clean Architecture**: Separation of concerns makes the codebase maintainable and testable

The system implements Indonesian-specific currency handling (no decimal points for IDR), proper thousand separator formatting (Rp 1.500.000), and integrates with mock payment processors following local banking standards.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Built With

**Backend:**
[![Go][Go.dev]][Go-url]
[![PostgreSQL][PostgreSQL.org]][PostgreSQL-url]
[![JWT][JWT.io]][JWT-url]

**Frontend:**
[![Nuxt][Nuxt.com]][Nuxt-url]
[![Vue][Vue.js]][Vue-url]
[![TailwindCSS][Tailwind.com]][Tailwind-url]
[![TypeScript][TypeScript.org]][TypeScript-url]

**Infrastructure:**
[![Docker][Docker.com]][Docker-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Key Features

**Core Functionality:**
- JWT-based authentication with session persistence
- Expense submission with IDR currency validation
- Auto-approval for expenses below IDR 1,000,000
- Manager approval workflow with notes
- Background payment processing with idempotency
- Status filtering (pending, approved, rejected, auto-approved)
- Responsive design optimized for mobile and desktop

**Technical Features:**
- Clean Architecture with clear layer separation
- Rate limiting on API endpoints
- Audit trail for all expense changes
- Proper IDR formatting (Rp 1.500.000)
- Session restoration across page refreshes
- Containerized setup with Docker Compose

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Business Rules

```go
const (
    MinExpenseAmount  = 10_000      // IDR 10,000
    MaxExpenseAmount  = 50_000_000  // IDR 50,000,000
    ApprovalThreshold = 1_000_000   // IDR 1,000,000
)
```

- **Currency**: All amounts in Indonesian Rupiah (IDR) only
- **Amount Validation**: Must be between IDR 10,000 and IDR 50,000,000
- **Auto-Approval**: Expenses < IDR 1,000,000 bypass manual approval
- **Access Control**: Employees see only their expenses; managers see all
- **Payment Processing**: Approved expenses trigger background payment jobs
- **Idempotency**: Payment processor handles duplicate requests via external_id

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->

## Getting Started

Follow these steps to set up the project locally.

### Prerequisites

Ensure you have the following installed:

- **Docker & Docker Compose**
  ```sh
  docker --version
  docker compose --version
  ```
- **Go 1.22+** (for local development)
  ```sh
  go version
  ```
- **Node.js 18+** (for local development)
  ```sh
  node --version
  ```

### Installation

**1. Clone the repository**

```sh
git clone <repository-url>
cd Expense-Management-System
```

**2. Start services with Docker (Recommended)**

```sh
# Make scripts executable
chmod +x start.sh cleanup.sh

# Start all services
./start.sh

# Or manually:
docker compose up -d
```

**3. Access the application**

Services will be available at:
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **API Documentation (Swagger)**: http://localhost:8080/docs
- **PostgreSQL**: localhost:5432

**4. Login with test credentials**

Employee Account:
```
Email: employee1@example.com
Password: password123
```

Manager Account:
```
Email: manager@example.com
Password: password123
```

**5. Stop and cleanup**

```sh
./cleanup.sh
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Local Development Setup

**Backend:**

```sh
cd backend

# Install dependencies
go mod download

# Set environment variables
cp .env.example .env
# Edit .env with your configuration

# Run migrations
cd migrations
./migrate.sh

# Run server
cd ..
go run cmd/api/main.go
```

**Frontend:**

```sh
cd frontend

# Install dependencies
npm install

# Set environment variables
cp .env.example .env
# Edit .env with your configuration

# Run development server
npm run dev
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- USAGE -->

## Usage

### Submit Expense (Auto-Approved)

1. Login as employee
2. Navigate to Dashboard
3. Click "Add Expense"
4. Enter amount below IDR 1,000,000
5. Status will be automatically set to "Approved"

### Submit Expense (Requires Approval)

1. Login as employee
2. Enter amount ≥ IDR 1,000,000
3. Status will be "Pending"
4. Manager will receive notification (audit log)

### Approve/Reject Expense

1. Login as manager
2. Navigate to "Approvals" page
3. Click on pending expense
4. Add approval notes
5. Click "Approve" or "Reject"

### Filter Expenses

Use filter buttons on Dashboard:
- **All**: Show all expenses
- **Pending**: Only pending approval
- **Approved**: Approved by manager or auto-approved
- **Rejected**: Rejected by manager

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- API DOCUMENTATION -->

## API Documentation

### Base URL
```
http://localhost:8080/api
```

### Authentication

**Login**
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "employee1@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "employee1@example.com",
    "name": "Employee One",
    "role": "employee"
  }
}
```

### Expense Operations

**Submit Expense**
```http
POST /api/expenses
Authorization: Bearer <token>
Content-Type: application/json

{
  "amount_idr": 750000,
  "description": "Office supplies",
  "receipt_url": "https://example.com/receipt.jpg"
}
```

**List Expenses**
```http
GET /api/expenses?status=pending&page=1&limit=10
Authorization: Bearer <token>
```

**Get Expense Details**
```http
GET /api/expenses/{id}
Authorization: Bearer <token>
```

### Manager Actions

**Approve Expense**
```http
PUT /api/expenses/{id}/approve
Authorization: Bearer <token>
Content-Type: application/json

{
  "notes": "Approved for Q1 budget"
}
```

**Reject Expense**
```http
PUT /api/expenses/{id}/reject
Authorization: Bearer <token>
Content-Type: application/json

{
  "notes": "Receipt not clear"
}
```

### Health Check

```http
GET /api/health
```

**Interactive API Documentation:**

Visit http://localhost:8080/docs for full interactive Swagger UI with:
- Try out API endpoints directly from browser
- Request/response examples
- Schema definitions
- Authentication flows


<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Architecture


This project follows **Clean Architecture** principles with clear separation of concerns:

```
backend/
├── cmd/api/              # Application entry point
├── internal/
│   ├── domain/          # Business entities & interfaces
│   ├── usecase/         # Business logic layer
│   ├── repository/      # Data access layer
│   ├── handler/         # HTTP handlers (API layer)
│   ├── middleware/      # Auth, logging, rate limiting
│   └── worker/          # Background payment processing
├── pkg/                 # Shared packages
│   ├── config/         # Configuration management
│   ├── database/       # Database connection
│   └── logger/         # Logging utilities
└── migrations/          # Database migrations

frontend/
├── pages/              # Vue pages (login, dashboard, approvals)
├── stores/             # Pinia state management
├── middleware/         # Route protection
├── plugins/            # Session persistence
├── composables/        # Reusable logic
└── types/              # TypeScript types
```

**Layer Responsibilities:**

- **Domain**: Core business entities and repository interfaces
- **Use Case**: Business logic implementation (validation, approval rules)
- **Repository**: Database operations and data access
- **Handler**: HTTP request/response handling
- **Middleware**: Cross-cutting concerns (auth, logging, rate limiting)
- **Worker**: Background job processing

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('employee', 'manager')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Expenses Table
```sql
CREATE TABLE expenses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    amount_idr INTEGER NOT NULL,
    description TEXT NOT NULL,
    receipt_url TEXT,
    status VARCHAR(50) NOT NULL,
    auto_approved BOOLEAN DEFAULT FALSE,
    payment_id VARCHAR(255),
    payment_external_id VARCHAR(255),
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Approvals Table
```sql
CREATE TABLE approvals (
    id SERIAL PRIMARY KEY,
    expense_id INTEGER REFERENCES expenses(id),
    approver_id INTEGER REFERENCES users(id),
    status VARCHAR(50) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Audit Logs Table
```sql
CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    expense_id INTEGER REFERENCES expenses(id),
    user_id INTEGER REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    old_status VARCHAR(50),
    new_status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Database Seeding:**

The system includes seed data with:
- 2 Employee accounts
- 1 Manager account
- Sample expenses in all statuses (pending, approved, rejected, auto-approved)
- Mix of amounts to test approval threshold

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Design Decisions

### 1. Currency as Integer

**Decision**: Store IDR amounts as integers (no decimals)

**Rationale**: 
- Avoids floating-point precision errors
- IDR currency doesn't use decimal points
- Matches Indonesian banking standards

**Trade-off**: Must format for display (1500000 → Rp 1.500.000)

### 2. Clean Architecture

**Decision**: Strict layer separation (domain → usecase → repository → handler)

**Rationale**:
- Testability through interface-based design
- Easy to swap implementations (e.g., database, payment provider)
- Clear separation of business logic from infrastructure

**Trade-off**: More boilerplate for simple CRUD operations

### 3. Background Payment Processing

**Decision**: Use Go channels + worker pool for async payment processing

**Rationale**:
- Non-blocking API responses improve UX
- Better performance under load
- Handles payment processor failures gracefully

**Trade-off**: Eventual consistency (status updates asynchronously)

### 4. Session Persistence with localStorage

**Decision**: Store JWT in localStorage with client-side plugin

**Rationale**:
- Better user experience (no re-login on refresh)
- Simple implementation for test assignment
- Matches common SPA patterns

**Trade-off**: Vulnerable to XSS attacks (acceptable for this scope)

### 5. Mock File Upload

**Decision**: Use placeholder URLs instead of real file storage

**Rationale**:
- Assignment requirement: "mock the receipt upload"
- Avoids S3/GCS complexity for test project
- Focuses evaluation on core business logic

**Trade-off**: Not production-ready

### 6. No User Registration

**Decision**: Only login endpoint, no registration

**Rationale**:
- The system represents an internal corporate tool where user accounts are provisioned centrally.
- This keeps the implementation focused on core expense and approval workflows.

**Trade-off**: Cannot create new users via UI

<p align="right">(<a href="#readme-top">back to top</a>)</p>


## Roadmap

### Phase 1: Core Functionality (Completed)
- [x] Authentication system with JWT
- [x] Expense submission with validation
- [x] Auto-approval logic
- [x] Manager approval workflow
- [x] Background payment processing
- [x] Audit trail implementation

### Phase 2: Infrastructure (Completed)
- [x] Docker Compose setup
- [x] Database migrations
- [x] Seed data for testing
- [x] Rate limiting
- [x] CORS configuration

### Phase 3: Documentation (Completed)
- [x] Comprehensive README
- [x] Swagger/OpenAPI specification
- [x] API examples with IDR
- [x] Architecture documentation

### Future Enhancements

**Immediate Improvements:**
- Comprehensive unit test coverage
- Integration tests for API endpoints
- E2E tests for critical user flows
- More descriptive error messages
- Structured logging with levels

**Feature Extensions:**
- Real email notifications via SMTP/SendGrid
- S3/GCS integration for receipt storage
- PDF/Excel export for expense reports
- Analytics dashboard with charts
- Comment threads on approvals
- Expense categories and tags
- Bulk approval operations

**Infrastructure:**
- CI/CD pipeline (GitHub Actions)
- Monitoring with Prometheus + Grafana
- Distributed tracing (OpenTelemetry)
- Redis caching for sessions
- Load balancing with Nginx
- Kubernetes deployment manifests

**Security:**
- OAuth2/OIDC integration
- Two-factor authentication
- IP tracking in audit logs
- Field-level encryption for sensitive data
- Automated security scanning

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->

<!-- Shields -->
[go-shield]: https://img.shields.io/badge/Go-1.22-00ADD8?style=for-the-badge&logo=go&logoColor=white
[go-url]: https://go.dev
[nuxt-shield]: https://img.shields.io/badge/Nuxt-3-00DC82?style=for-the-badge&logo=nuxt.js&logoColor=white
[nuxt-url]: https://nuxt.com
[postgres-shield]: https://img.shields.io/badge/PostgreSQL-15-336791?style=for-the-badge&logo=postgresql&logoColor=white
[postgres-url]: https://www.postgresql.org
[docker-shield]: https://img.shields.io/badge/Docker-Compose-2496ED?style=for-the-badge&logo=docker&logoColor=white
[docker-url]: https://www.docker.com

<!-- Technology Badges -->
[Go.dev]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[Go-url]: https://go.dev
[PostgreSQL.org]: https://img.shields.io/badge/PostgreSQL-336791?style=for-the-badge&logo=postgresql&logoColor=white
[PostgreSQL-url]: https://www.postgresql.org
[JWT.io]: https://img.shields.io/badge/JWT-000000?style=for-the-badge&logo=jsonwebtokens&logoColor=white
[JWT-url]: https://jwt.io
[Nuxt.com]: https://img.shields.io/badge/Nuxt.js-00DC82?style=for-the-badge&logo=nuxt.js&logoColor=white
[Nuxt-url]: https://nuxt.com
[Vue.js]: https://img.shields.io/badge/Vue.js-4FC08D?style=for-the-badge&logo=vue.js&logoColor=white
[Vue-url]: https://vuejs.org
[Tailwind.com]: https://img.shields.io/badge/Tailwind_CSS-38B2AC?style=for-the-badge&logo=tailwind-css&logoColor=white
[Tailwind-url]: https://tailwindcss.com
[TypeScript.org]: https://img.shields.io/badge/TypeScript-3178C6?style=for-the-badge&logo=typescript&logoColor=white
[TypeScript-url]: https://www.typescriptlang.org
[Docker.com]: https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white
[Docker-url]: https://www.docker.com