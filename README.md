# E-commerce API (Go + Gin)

A clean-architecture backend for an e-commerce platform built with Go and the Gin web framework. This API provides complete functionality for user authentication, product management, order processing, and image uploads with enterprise-grade features like rate limiting, caching, and comprehensive API documentation.

## 🎯 Overview

This e-commerce backend API enables:

- **User Management**: Secure registration and authentication with JWT tokens
- **Product Management**: Full CRUD operations with image uploads (up to 4 images per product)
- **Order Processing**: Transactional order creation with stock validation
- **Admin Features**: Role-based access control with admin promotion capabilities
- **Performance**: In-memory caching for product listings and rate limiting for API protection
- **Documentation**: Interactive Swagger/OpenAPI documentation

## 🛠 Tech Stack

### Core Technologies

- **Go 1.24**: Modern, efficient programming language
- **Gin Web Framework**: High-performance HTTP web framework
- **PostgreSQL**: Robust relational database
- **GORM**: Powerful ORM for database operations

### Security & Authentication

- **JWT (JSON Web Tokens)**: Secure token-based authentication
- **Bcrypt**: Industry-standard password hashing
- **Role-Based Access Control (RBAC)**: Admin and User roles with middleware protection

### Infrastructure & Services

- **Cloudinary**: Cloud-based image storage with signed/unsigned upload support
- **Zap Logger**: Structured, high-performance logging
- **Docker & Docker Compose**: Containerized deployment
- **Swaggo/Swagger**: Auto-generated API documentation

### Performance Features

- **In-Memory Caching**: TTL-based cache for product listings (configurable)
- **Rate Limiting**: Per-IP request limiting to prevent abuse (configurable)
- **Database Transactions**: Atomic operations for order processing

## 📁 Project Structure (Clean Architecture)

This project follows **Clean Architecture** principles with clear separation of concerns:

```
e-commerce/
├── cmd/
│   └── server/              # Application entry point
│       └── main.go         # Main function, config loading, server startup
│
├── config/                 # Configuration management
│   └── config.go          # Config structs and loader (Viper)
│
├── internal/               # Internal application code
│   ├── domain/            # Business entities and interfaces
│   │   ├── user.go        # User entity and Role enum
│   │   ├── product.go     # Product entity
│   │   ├── order.go       # Order and OrderItem entities
│   │   ├── errors.go      # Domain-specific errors
│   │   └── repository/    # Repository interfaces (abstractions)
│   │
│   ├── usecase/           # Business logic layer
│   │   ├── auth/          # Authentication use cases
│   │   ├── product/       # Product management use cases
│   │   └── order/         # Order processing use cases
│   │
│   ├── adapter/           # External adapters (HTTP, Database)
│   │   ├── handler/       # HTTP handlers (Gin)
│   │   ├── middleware/    # Auth, CORS, Rate Limiting
│   │   ├── router/        # Route definitions
│   │   └── repository/
│   │       └── gorm/      # GORM implementations
│   │
│   └── infrastructure/    # Infrastructure setup
│       ├── container/     # Dependency Injection container
│       └── database/      # Database initialization and migrations
│
├── pkg/                   # Reusable packages
│   ├── cache/            # In-memory cache implementation
│   ├── cloudinary/       # Cloudinary client wrapper
│   ├── hash/             # Password hashing utilities
│   ├── jwt/              # JWT token management
│   ├── logger/           # Logger initialization
│   └── response/         # Standardized API response helpers
│
├── docs/                 # Auto-generated Swagger documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── config.yaml           # Application configuration (user-provided)
├── Dockerfile            # Docker image definition
├── docker-compose.yml    # Multi-container orchestration
├── Makefile              # Build and development commands
└── README.md             # This file
```

### Architecture Benefits

- **Separation of Concerns**: Business logic is independent of frameworks
- **Testability**: Easy to mock dependencies and test use cases
- **Maintainability**: Clear boundaries between layers
- **Flexibility**: Easy to swap implementations (e.g., different database)

## 🚀 Quick Start

### Prerequisites

- Go 1.24 or later
- PostgreSQL 12+ (or use Docker)
- Docker and Docker Compose (optional, for containerized setup)
- Cloudinary account (for image uploads, optional)

### 1. Configuration Setup

Create a `config.yaml` file at the project root:

```
In the codebase, rename config.yaml.sample to config.yaml and use it directly. You may modify the credentials as needed. The config.yaml file is excluded from git history for environment security reasons.
```

```yaml
app:
  name: ecommerce-api
  environment: development

server:
  port: 8080

database:
  host: localhost # Use 'host.docker.internal' for Docker on Mac
  port: 5432 # Use 5433 for Docker
  user: postgres
  password: postgres
  name: commerce
  sslmode: disable

jwt:
  secret: your-secret-key-change-in-production
  issuer: ecommerce-api
  access_token_ttl: 30m
  refresh_token_ttl: 168h

cloudinary:
  cloud_name: your-cloud-name
  api_key: your-api-key
  api_secret: your-api-secret
  upload_preset: your-preset # For unsigned uploads (optional)
  folder: ecommerce

rate_limit:
  enabled: true
  limit: 100 # Requests per window
  window: 1m # Time window

cache:
  enabled: true
  product_list_ttl: 1m # Cache TTL for product listings
  max_product_entries: 1000

admin_seed:
  enabled: true
  email: admin@example.com
  username: admin
  password: Admin#1234
```

**Note**: All configuration values can be overridden using environment variables (e.g., `DATABASE_HOST`, `JWT_SECRET`). Use underscores instead of dots (e.g., `DATABASE_HOST` for `database.host`).

### 2. Local Development

#### Option A: Run Locally (requires local PostgreSQL)

```bash
# Install dependencies
go mod download

# Run the application
make run
# or
go run cmd/server/main.go
```

#### Option B: Run with Docker (recommended)

```bash
# Start PostgreSQL and the application
make up

# Or run in detached mode with logs
make up-detach-logs

# Stop containers
make down
```

The API will be available at `http://localhost:8080`

### 3. Verify Installation

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Expected response: {"success":true,"message":"ok","data":null}
```

## 📚 Available Commands (Makefile)

| Command               | Description                                        |
| --------------------- | -------------------------------------------------- | --- |
| `make run`            | Run the application locally                        |
| `make build`          | Build the application binary to `bin/ecommerce`    |
| `make test`           | Run tests with handler/endpoint tests              |     |
| `make fmt`            | Format all Go code                                 |
| `make swagger`        | Generate/update Swagger API documentation          |
| `make docker-build`   | Build Docker image                                 |
| `make up`             | Start Docker containers (PostgreSQL + API)         |
| `make up-detach-logs` | Start Docker containers in detached mode with logs |
| `make down`           | Stop and remove Docker containers                  |

## 🔐 Authentication & Authorization

### User Registration

- **Endpoint**: `POST /api/v1/auth/register`
- **Access**: Public
- **Role Assignment**: All new users are automatically assigned the "user" role
- **Password Requirements**: Strong password validation (minimum length, special characters, etc.)

### User Login

- **Endpoint**: `POST /api/v1/auth/login`
- **Access**: Public
- **Response**: Returns JWT token with user information
- **Token Usage**: Include in `Authorization: Bearer <token>` header for protected endpoints

### Role-Based Access

- **User Role**: Can place orders and view own orders
- **Admin Role**: Full access including product management, image uploads, and user promotion
- **Protected Endpoints**: Require valid JWT token and appropriate role

## 📡 API Endpoints

### Base URL

All API endpoints are prefixed with `/api/v1`

### Health Check

- **GET** `/api/v1/health` - Check API health status (public)

### Authentication Endpoints

#### Register User

- **POST** `/api/v1/auth/register`
- **Access**: Public
- **Request Body**:
  ```json
  {
    "username": "john_doe",
    "email": "john@example.com",
    "password": "Strong#Pass123"
  }
  ```
- **Success Response** (201):
  ```json
  {
    "success": true,
    "message": "user registered successfully",
    "data": {
      "userId": "uuid",
      "username": "john_doe",
      "email": "john@example.com",
      "role": "user"
    }
  }
  ```

#### Login

- **POST** `/api/v1/auth/login`
- **Access**: Public
- **Request Body**:
  ```json
  {
    "email": "john@example.com",
    "password": "Strong#Pass123"
  }
  ```
- **Success Response** (200):
  ```json
  {
    "success": true,
    "message": "login successful",
    "data": {
      "token": "jwt-token-here",
      "expiresAt": "2024-01-01T12:00:00Z",
      "userId": "uuid",
      "username": "john_doe",
      "email": "john@example.com",
      "role": "user"
    }
  }
  ```

### Product Endpoints

#### List Products (Public)

- **GET** `/api/v1/products`
- **Access**: Public
- **Query Parameters**:
  - `search` (optional): Search products by name
  - `page` (optional, default: 1): Page number
  - `limit` (optional, default: 10): Items per page
- **Features**:
  - Pagination support
  - Search functionality
  - Includes product images in response
  - Cached responses (configurable TTL)
- **Success Response** (200):
  ```json
  {
    "success": true,
    "message": "products retrieved",
    "data": [
      {
        "id": "uuid",
        "name": "Product Name",
        "description": "Product description",
        "price": 99.99,
        "stock": 50,
        "category": "Electronics",
        "images": [
          {
            "id": "uuid",
            "url": "https://cloudinary.com/image.jpg",
            "productId": "uuid"
          }
        ],
        "createdAt": "2024-01-01T12:00:00Z",
        "updatedAt": "2024-01-01T12:00:00Z"
      }
    ],
    "currentPage": 1,
    "pageSize": 10,
    "totalPages": 5,
    "totalProducts": 50
  }
  ```

#### Get Product Details (Public)

- **GET** `/api/v1/products/:id`
- **Access**: Public
- **Success Response** (200): Single product object with images
- **Error Response** (404): Product not found

#### Create Product (Admin Only)

- **POST** `/api/v1/products`
- **Access**: Admin (requires JWT token)
- **Request Body**:
  ```json
  {
    "name": "New Product",
    "description": "Product description",
    "price": 99.99,
    "stock": 100,
    "category": "Electronics"
  }
  ```
- **Success Response** (201): Created product object

#### Update Product (Admin Only)

- **PUT** `/api/v1/products/:id`
- **Access**: Admin (requires JWT token)
- **Request Body**: Partial update (only include fields to update)
- **Success Response** (200): Updated product object
- **Error Response** (404): Product not found

#### Delete Product (Admin Only)

- **DELETE** `/api/v1/products/:id`
- **Access**: Admin (requires JWT token)
- **Business Rule**: Cannot delete products with pending orders
- **Success Response** (200): Success message
- **Error Responses**:
  - 400: Product has pending orders
  - 404: Product not found

#### Upload Product Images (Admin Only)

- **POST** `/api/v1/products/:id/images`
- **Access**: Admin (requires JWT token)
- **Content-Type**: `multipart/form-data`
- **Form Field**: `files` (1-4 image files)
- **Limits**: Maximum 4 images per product (total, not per request)
- **Upload Method**: Uses signed uploads if Cloudinary API key/secret are configured, otherwise falls back to unsigned
- **Success Response** (201):
  ```json
  {
    "success": true,
    "message": "images uploaded",
    "data": [
      {
        "id": "uuid",
        "url": "https://cloudinary.com/image.jpg",
        "productId": "uuid"
      }
    ]
  }
  ```

### Order Endpoints

#### Create Order (User/Admin)

- **POST** `/api/v1/orders`
- **Access**: Authenticated users (requires JWT token)
- **Request Body**:
  ```json
  {
    "description": "Order description",
    "items": [
      {
        "productId": "uuid",
        "quantity": 2
      }
    ]
  }
  ```
- **Features**:
  - Transactional stock validation
  - Automatic stock deduction
  - Prevents overselling
- **Success Response** (201): Created order with items
- **Error Responses**:
  - 400: Insufficient stock or invalid product
  - 404: Product not found

#### List My Orders (User/Admin)

- **GET** `/api/v1/orders`
- **Access**: Authenticated users (requires JWT token)
- **Features**: Returns only orders belonging to the authenticated user
- **Success Response** (200): Array of order objects with items

### Admin Endpoints

#### Promote User to Admin

- **POST** `/api/v1/admin/users/:id/admin`
- **Access**: Admin only (requires JWT token with admin role)
- **Path Parameter**: `id` - User UUID to promote
- **Success Response** (200): Success message
- **Error Response** (404): User not found

## 🧪 Testing

### Running Tests

```bash
# Run endpoint tests
make test

# Run tests for a specific package
go test ./internal/adapter/handler/... -v
```

### Test Coverage

The project includes unit tests for all API endpoints:

- ✅ **Auth Handler Tests**: Register and Login endpoints
- ✅ **Product Handler Tests**: List products endpoint
- ✅ **Order Handler Tests**: Create and List orders endpoints
- ✅ **Admin Handler Tests**: Promote user to admin endpoint

Tests use mocks (testify/mock) to isolate handler logic and verify:

- HTTP status codes
- Response structure
- Error handling
- Service method calls

### Writing New Tests

Tests are located alongside handlers in `internal/adapter/handler/*_test.go`. Follow the existing pattern:

1. Create mock services implementing the use case interfaces
2. Set up Gin test context
3. Call handler methods
4. Assert responses using `testify/assert`

## 📖 API Documentation (Swagger)

### Accessing Swagger UI

Once the server is running, access the interactive API documentation at:

```
http://localhost:8080/swagger/index.html
```

### Features

- **Interactive Testing**: Try API endpoints directly from the browser
- **Authentication**: Click "Authorize" to add your JWT token
- **Request/Response Examples**: See example payloads and responses
- **Endpoint Details**: View all parameters, requirements, and status codes

### Generating Documentation

```bash
# Generate/update Swagger docs
make swagger
```

This command:

1. Installs/updates the `swag` tool
2. Scans code for Swagger annotations
3. Generates `docs/swagger.json` and `docs/swagger.yaml`

### Swagger Annotations

All endpoints are annotated with:

- `@Summary`: Brief endpoint description
- `@Description`: Detailed explanation
- `@Tags`: Grouping (Auth, Products, Orders, Admin, Health)
- `@Param`: Request parameters
- `@Success`: Success response structure
- `@Failure`: Error response structure
- `@Security`: Authentication requirements
- `@Router`: Route path and method

## 🔧 Configuration Details

### Database Configuration

- **Host**: Database server address
  - Local: `localhost`
  - Docker (Mac): `host.docker.internal`
- **Port**: Database port (default: 5432, Docker: 5433)
- **SSL Mode**: `disable` for local development, `require` for production

### JWT Configuration

- **Secret**: Strong secret key (change in production!)
- **Access Token TTL**: Default 30 minutes
- **Refresh Token TTL**: Default 7 days

### Cloudinary Configuration

- **Cloud Name**: Your Cloudinary cloud name
- **API Key/Secret**: For signed uploads (recommended)
- **Upload Preset**: For unsigned uploads (optional)
- **Folder**: Organize images in a specific folder

### Rate Limiting

- **Enabled**: Toggle rate limiting on/off
- **Limit**: Maximum requests per window (default: 100)
- **Window**: Time window (default: 1 minute)
- **Note**: Swagger UI routes are excluded from rate limiting

### Caching

- **Enabled**: Toggle caching on/off
- **Product List TTL**: Cache expiration time (default: 1 minute)
- **Max Entries**: Maximum cached entries (default: 1000)
- **Scope**: Only product listing endpoint is cached

### Admin Seeding

- **Enabled**: Automatically create admin user on startup
- **Idempotent**: Won't create duplicate admins (checks email)
- **Use Case**: Simplifies initial setup for development/testing

## 🐳 Docker Setup

### Docker Compose Services

1. **PostgreSQL Database** (`db`)

   - Image: `postgres:15`
   - Port: `5433:5432` (host:container)
   - Database: `commerce`
   - Persistent volume: `db_data`

2. **E-commerce API** (`app`)
   - Built from `Dockerfile`
   - Port: `8080:8080`
   - Depends on: `db`
   - Mounts: `config.yaml` (read-only)

### Docker Commands

```bash
# Build and start all services
make up

# Start in background with logs
make up-detach-logs

# Stop and remove containers
make down

# View logs
docker-compose logs -f app

# Execute commands in container
docker exec -it ecommerce-api sh
```

### Dockerfile Details

- **Base Image**: `golang:1.24-alpine` (build stage)
- **Final Image**: `alpine:3.18` (minimal runtime)
- **Security**: Runs as non-root user `appuser`
- **Optimization**: Multi-stage build for smaller image size

## 🔒 Security Features

### Password Security

- **Hashing**: Bcrypt with automatic salt generation
- **Validation**: Strong password requirements enforced
- **Storage**: Passwords never stored in plain text

### Authentication Security

- **JWT Tokens**: Secure, stateless authentication
- **Token Expiration**: Configurable TTL for access tokens
- **Role-Based Access**: Middleware enforces role requirements

### API Security

- **Rate Limiting**: Prevents abuse and DDoS attacks
- **CORS**: Configurable cross-origin resource sharing
- **Input Validation**: All inputs validated before processing
- **Error Handling**: No sensitive information leaked in errors

### Database Security

- **Parameterized Queries**: GORM prevents SQL injection
- **Transactions**: Atomic operations prevent data corruption
- **Connection Security**: SSL mode configurable

## 📊 Performance Features

### Caching

- **Product Listings**: Results cached by search query and page
- **TTL-Based Expiration**: Automatic cache invalidation
- **Memory Efficient**: Configurable maximum entries

### Rate Limiting

- **Per-IP Tracking**: Tracks requests by client IP
- **Sliding Window**: Cleans old requests automatically
- **Configurable**: Adjust limits per environment

### Database Optimization

- **Indexes**: GORM auto-creates indexes on foreign keys
- **Eager Loading**: Product images loaded efficiently
- **Connection Pooling**: GORM manages connection pool

## 🚨 Error Handling

### Standard Error Response Format

```json
{
  "success": false,
  "message": "Error message",
  "errors": ["Detailed error 1", "Detailed error 2"]
}
```

### Common HTTP Status Codes

- **200 OK**: Successful GET/PUT request
- **201 Created**: Successful POST request
- **400 Bad Request**: Validation errors, invalid input
- **401 Unauthorized**: Missing or invalid authentication
- **403 Forbidden**: Insufficient permissions (wrong role)
- **404 Not Found**: Resource not found
- **429 Too Many Requests**: Rate limit exceeded
- **500 Internal Server Error**: Server-side errors

### Domain-Specific Errors

The application uses custom domain errors for better error handling:

- `ErrEmailAlreadyExists`: Email already registered
- `ErrUsernameAlreadyExists`: Username already taken
- `ErrInvalidCredentials`: Invalid login credentials
- `ErrProductNotFound`: Product doesn't exist
- `ErrInsufficientStock`: Not enough stock for order
- `ErrProductHasPendingOrders`: Cannot delete product with orders
- `ErrUserNotFound`: User doesn't exist

## 🔄 Business Rules

### User Registration

- All new users are assigned "user" role (admin cannot be created via registration)
- Username and email must be unique
- Strong password validation enforced

### Product Management

- Products can only be deleted if they have no pending orders
- Product images limited to 4 per product (total, not per upload)
- Stock is validated and decremented transactionally during order creation

### Order Processing

- Orders are created within database transactions
- Stock is checked and decremented atomically
- Users can only view their own orders
- Orders cannot be created for out-of-stock products

### Admin Operations

- Only existing admins can promote other users to admin
- Admin promotion is idempotent (safe to call multiple times)
- Admin user is seeded automatically on startup (if configured)

## 🛠 Development Guidelines

### Code Style

- Follow Go standard formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions focused and small

### Testing

- Write tests for all handlers
- Use mocks for external dependencies
- Test both success and error cases
- Maintain test coverage

### Error Handling

- Use domain-specific errors
- Return appropriate HTTP status codes
- Log errors with context
- Never expose sensitive information

### Database

- Use transactions for multi-step operations
- Validate inputs before database operations
- Use GORM's built-in features (migrations, relationships)
- Handle connection errors gracefully

## 📝 Environment Variables

All configuration can be overridden using environment variables. Use underscores instead of dots:

| Config Path          | Environment Variable | Example         |
| -------------------- | -------------------- | --------------- |
| `database.host`      | `DATABASE_HOST`      | `localhost`     |
| `database.port`      | `DATABASE_PORT`      | `5432`          |
| `jwt.secret`         | `JWT_SECRET`         | `my-secret-key` |
| `cloudinary.api_key` | `CLOUDINARY_API_KEY` | `123456789`     |
| `rate_limit.limit`   | `RATE_LIMIT_LIMIT`   | `100`           |

## 🐛 Troubleshooting

### Common Issues

#### Database Connection Failed

- **Check**: Database is running and accessible
- **Docker**: Ensure `host.docker.internal` is used on Mac
- **Port**: Verify port mapping in `docker-compose.yml`

#### Swagger UI Shows 429 Error

- **Solution**: Swagger routes are excluded from rate limiting (already fixed)
- **Check**: Rate limiter configuration in `config.yaml`

#### Cloudinary Upload Fails

- **Check**: API credentials in `config.yaml`
- **Network**: Verify DNS resolution (Docker uses Google DNS)
- **Logs**: Check application logs for detailed error messages

#### Tests Fail

- **Check**: All dependencies installed (`go mod download`)
- **Verify**: Test database is accessible (if integration tests)

### Getting Help

1. Check application logs: `docker-compose logs -f app`
2. Verify configuration: Review `config.yaml`
3. Test endpoints: Use Swagger UI or `curl`
4. Check database: Connect directly to verify data

## 📄 License

**Built with ❤️ using Go and Clean Architecture principles**

                                         ### Thank you
