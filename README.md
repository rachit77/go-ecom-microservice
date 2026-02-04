# Go E-Commerce Microservice

A modern, scalable microservice-based e-commerce platform built with Go, gRPC, and GraphQL. The system follows a microservice architecture pattern with a GraphQL API gateway for unified client access.

## Project Overview

This project demonstrates a production-ready microservice architecture with the following characteristics:

- **Microservice Architecture**: Independent services for account management, product catalog, and order processing
- **gRPC Communication**: High-performance inter-service communication using Protocol Buffers
- **GraphQL API Gateway**: Unified GraphQL API for seamless client integration
- **Database per Service**: Each service maintains its own database (separation of concerns)
- **Docker Containerization**: Easy deployment and orchestration using Docker Compose
- **PostgreSQL & Elasticsearch**: Relational and search capabilities for different use cases

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Client Applications                       │
└──────────────────────────┬──────────────────────────────────┘
                           │
                    HTTP/GraphQL
                           │
        ┌──────────────────▼──────────────────┐
        │    GraphQL API Gateway (Port 8000)   │
        └──────────────────┬──────────────────┘
                           │
         ┌─────────────────┼─────────────────┐
         │                 │                 │
      gRPC              gRPC              gRPC
         │                 │                 │
    ┌────▼───┐       ┌────▼────┐      ┌────▼──┐
    │ Account │       │ Catalog  │      │ Order │
    │ Service │       │ Service  │      │Service│
    │ :8080   │       │ :8080    │      │ :8080 │
    └────┬───┘       └────┬────┘      └────┬──┘
         │                │                │
    ┌────▼──────┐    ┌────▼──────┐   ┌────▼────┐
    │ PostgreSQL │    │ Elasticsearch  │ PostgreSQL
    │ (Account)  │    │ (Catalog)      │ (Order)
    └────────────┘    └───────────┘   └─────────┘
```

## Services

### 1. Account Service
- **Purpose**: Manage user accounts and authentication
- **Port**: 8080 (gRPC)
- **Database**: PostgreSQL
- **Key Features**:
  - Create and retrieve accounts
  - Link accounts to orders

### 2. Catalog Service
- **Purpose**: Manage product inventory and search
- **Port**: 8080 (gRPC)
- **Database**: Elasticsearch
- **Key Features**:
  - Create products with details (name, description, price)
  - Full-text search capabilities
  - Product pagination

### 3. Order Service
- **Purpose**: Handle order creation and management
- **Port**: 8080 (gRPC)
- **Database**: PostgreSQL
- **Key Features**:
  - Create orders with multiple products
  - Link orders to accounts
  - Calculate total order prices
  - Retrieve order history per account

### 4. GraphQL Gateway
- **Purpose**: Unified API for client applications
- **Port**: 8000 (HTTP/GraphQL)
- **Features**:
  - Aggregates data from all services
  - Provides a clean GraphQL interface
  - GraphQL Playground for testing

## Prerequisites

- **Go**: Version 1.19 or higher
- **Docker**: For containerized deployment
- **Docker Compose**: For orchestration
- **PostgreSQL**: For Account and Order services (included in Docker Compose)
- **Elasticsearch**: For Catalog service (included in Docker Compose)

## Getting Started

### Option 1: Local Development (Direct Go Execution)

1. **Clone the repository**
   ```bash
   git clone https://github.com/rachit77/go-ecom-microservice.git
   cd go-ecom-microservice
   ```

2. **Set up environment variables**
   Create a `.env` file or export the following:
   ```bash
   export ACCOUNT_SERVICE_URL=localhost:8081
   export CATALOG_SERVICE_URL=localhost:8082
   export ORDER_SERVICE_URL=localhost:8083
   ```

3. **Start individual services** (in separate terminals)
   
   **Account Service:**
   ```bash
   cd account/cmd/account
   go run main.go
   ```

   **Catalog Service:**
   ```bash
   cd catalog/cmd/catalog
   go run main.go
   ```

   **Order Service:**
   ```bash
   cd order/cmd/order
   go run main.go
   ```

   **GraphQL Gateway:**
   ```bash
   cd graphql
   go run .
   ```

### Option 2: Docker Compose (Recommended for Production)

1. **Clone the repository**
   ```bash
   git clone https://github.com/rachit77/go-ecom-microservice.git
   cd go-ecom-microservice
   ```

2. **Start all services**
   ```bash
   docker-compose up --build
   ```

   This will:
   - Build all service images
   - Start all microservices
   - Initialize databases
   - Expose the GraphQL gateway at `http://localhost:8000`

3. **Stop services**
   ```bash
   docker-compose down
   ```

## Accessing the Application

### GraphQL Playground
- **URL**: `http://localhost:8000/playground`
- Interactive IDE for testing GraphQL queries and mutations

### GraphQL API Endpoint
- **URL**: `http://localhost:8000/graphql`
- Send POST requests with GraphQL queries/mutations

## GraphQL API Reference

### Queries

#### Get All Accounts
```graphql
query {
  accounts {
    id
    name
    orders {
      id
      createdAt
      totalPrice
    }
  }
}
```

#### Get Account by ID
```graphql
query {
  accounts(id: "account-id") {
    id
    name
    orders {
      id
      createdAt
      totalPrice
      products {
        id
        name
        price
        quantity
      }
    }
  }
}
```

#### Get All Products
```graphql
query {
  products {
    id
    name
    description
    price
  }
}
```

#### Search Products
```graphql
query {
  products(query: "laptop") {
    id
    name
    description
    price
  }
}
```

### Mutations

#### Create Account
```graphql
mutation {
  createAccount(account: { name: "John Doe" }) {
    id
    name
  }
}
```

#### Create Product
```graphql
mutation {
  createProduct(product: {
    name: "Laptop"
    description: "High-performance laptop"
    price: 999.99
  }) {
    id
    name
    price
  }
}
```

#### Create Order
```graphql
mutation {
  createOrder(order: {
    accountId: "account-id"
    products: [
      { id: "product-id-1", quantity: 2 },
      { id: "product-id-2", quantity: 1 }
    ]
  }) {
    id
    createdAt
    totalPrice
    products {
      id
      name
      price
      quantity
    }
  }
}
```

## Project Structure

```
go-ecom-microservice/
├── account/              # Account microservice
│   ├── cmd/account/      # Service entry point
│   ├── pb/               # Protocol Buffer generated code
│   ├── account.proto     # Service definition
│   ├── client.go         # gRPC client
│   ├── server.go         # gRPC server
│   ├── service.go        # Business logic
│   ├── repository.go     # Data access layer
│   ├── app.dockerfile    # Docker image for service
│   ├── db.dockerfile     # Docker image for database
│   └── up.sql            # Database initialization script
├── catalog/              # Catalog microservice (similar structure)
├── order/                # Order microservice (similar structure)
├── graphql/              # GraphQL API Gateway
│   ├── main.go           # Gateway entry point
│   ├── graph.go          # Server initialization
│   ├── schema.graphql    # GraphQL schema definition
│   ├── query_resolver.go # Query resolver implementation
│   ├── mutation_resolver.go # Mutation resolver implementation
│   ├── models.go         # Data models
│   └── app.dockerfile    # Docker image
├── docker-compose.yml    # Service orchestration
├── go.mod                # Go module definition
└── README.md             # This file
```

## Database Schemas

### Account Database (PostgreSQL)
- **Table**: `accounts`
  - `id` (UUID, Primary Key)
  - `name` (VARCHAR)
  - `created_at` (TIMESTAMP)

### Catalog Database (Elasticsearch)
- **Index**: Products
  - `id` (Keyword)
  - `name` (Text)
  - `description` (Text)
  - `price` (Float)

### Order Database (PostgreSQL)
- **Table**: `orders`
  - `id` (UUID, Primary Key)
  - `account_id` (UUID, Foreign Key)
  - `total_price` (DECIMAL)
  - `created_at` (TIMESTAMP)

- **Table**: `order_products`
  - `id` (UUID, Primary Key)
  - `order_id` (UUID, Foreign Key)
  - `product_id` (UUID)
  - `quantity` (INT)

## Development Workflow

1. **Make changes to a service**
   ```bash
   # Edit service code
   vim account/service.go
   ```

2. **If using local development**
   - Kill the running service (Ctrl+C)
   - Recompile: `go run .` in the service directory

3. **If using Docker**
   - Rebuild images: `docker-compose up --build`

4. **Test changes**
   - Use GraphQL Playground to test queries/mutations

## Troubleshooting

### gRPC Connection Errors
- Ensure all services are running on the correct ports
- Check `docker-compose.yml` for service URLs
- Verify environment variables are set correctly

### Database Connection Errors
- Verify database services are running
- Check database credentials in `docker-compose.yml`
- Ensure databases are initialized with SQL scripts

### GraphQL Playground Not Loading
- Confirm the GraphQL gateway is running on port 8000
- Check browser console for errors
- Verify CORS settings if accessing from different origin

## Future Enhancements

- [ ] Authentication and authorization layer
- [ ] Order status tracking and notifications
- [ ] Product reviews and ratings
- [ ] Payment gateway integration
- [ ] Advanced search and filtering
- [ ] API rate limiting
- [ ] Distributed tracing with Jaeger
- [ ] Metrics collection with Prometheus

## Contributing

1. Create a feature branch
2. Make your changes
3. Test thoroughly
4. Submit a pull request

## License

MIT License - see LICENSE file for details

