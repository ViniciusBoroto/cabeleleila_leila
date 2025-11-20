# Quick Reference: JWT Authentication

## Setup

1. **Environment Variables** (create `.env`):
```
DB_PATH=app.db
JWT_SECRET=your-32-character-random-secret-key
```

2. **Run Application**:
```bash
go run main.go
```

## Registration & Login

### Register New Customer
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newcustomer@example.com",
    "password": "password123",
    "name": "Customer Name",
    "phone": "123456789"
  }'
```

Response (201):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "newcustomer@example.com",
    "name": "Customer Name",
    "role": "customer"
  }
}
```

### Login Existing User
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

## API Endpoints

### Public Endpoints

#### POST `/api/auth/login`
Login with email and password.

**Request**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200)**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe",
    "role": "customer"
  }
}
```

**Errors**:
- `400`: Invalid email format or password < 6 chars
- `401`: Wrong credentials or inactive user

#### POST `/api/auth/register`
Create a new customer account (public registration).

**Request**:
```json
{
  "email": "newcustomer@example.com",
  "password": "password123",
  "name": "Customer Name",
  "phone": "123456789"
}
```

**Response (201)**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 2,
    "email": "newcustomer@example.com",
    "name": "Customer Name",
    "role": "customer"
  }
}
```

**Errors**:
- `400`: Invalid email, password < 6 chars, or missing fields
- `409`: Email already registered

### Protected Endpoints

All protected endpoints require:
```
Authorization: Bearer <jwt_token>
```

#### GET `/api/appointments`
List appointments (all authenticated users).

#### POST `/api/appointments`
Create appointment (all authenticated users).

### Admin Endpoints

Require `Authorization: Bearer <admin-token>` + admin role

#### GET `/api/admin/users`
List all users.

#### POST `/api/admin/users`
Create new user.

**Request**:
```json
{
  "email": "newuser@example.com",
  "password": "password123",
  "name": "New User",
  "phone": "123456789",
  "role": "customer"
}
```

#### GET `/api/admin/users/{id}`
Get specific user.

#### PUT `/api/admin/users/{id}`
Update user.

#### DELETE `/api/admin/users/{id}`
Delete user.

## Roles

### Customer (`customer`)
- Can login
- Can view/manage own appointments

### Admin (`admin`)
- Can login
- Can manage all users
- Can view all appointments

## JWT Token Structure

Token payload (claims):
```json
{
  "user_id": 1,
  "email": "user@example.com",
  "role": "customer",
  "exp": 1732118400,
  "iat": 1732032000
}
```

Token expires in: **24 hours**

## Creating Users via API

### Create Customer User
```bash
curl -X POST http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "password": "password123",
    "name": "Customer Name",
    "phone": "123456789",
    "role": "customer"
  }'
```

### Create Admin User
```bash
curl -X POST http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123",
    "name": "Admin Name",
    "phone": "987654321",
    "role": "admin"
  }'
```

## Testing

### Run All Tests
```bash
go test -v ./internal/service ./internal/repository ./internal/handlers -cover
```

### Run Specific Test Suite
```bash
# Auth service tests
go test -v ./internal/service/auth_service_test.go ./internal/service/auth_service.go

# Middleware tests
go test -v ./internal/handlers -run Middleware

# Handler tests
go test -v ./internal/handlers -run Login
```

### Coverage Report
```bash
go test -coverprofile=coverage.out ./internal/...
go tool cover -html=coverage.out
```

## Common Issues & Solutions

### Issue: "JWT_SECRET environment variable not set"
**Solution**: Create `.env` file with `JWT_SECRET=your-secret`

### Issue: "user not found" on login
**Solution**: User doesn't exist. Create user via admin endpoint first.

### Issue: "invalid token" on protected route
**Solution**: Token is invalid or expired. Login again to get new token.

### Issue: "user role not authorized"
**Solution**: Your user role doesn't have permission. Use admin token for admin routes.

### Issue: Database locked
**Solution**: Close other connections to `app.db` file.

## Password Requirements

- Minimum 6 characters
- Uses bcrypt hashing
- Never sent back in API responses

## Security Checklist

- [ ] Change `JWT_SECRET` in production (use 32+ character random string)
- [ ] Use HTTPS in production
- [ ] Keep `app.db` file secure
- [ ] Regularly rotate JWT_SECRET
- [ ] Monitor failed login attempts
- [ ] Deactivate unused accounts (`IsActive: false`)
- [ ] Use strong passwords for admin accounts

## Accessing Context Values in Handlers

Inside any protected route handler, access:

```go
func MyHandler(c *gin.Context) {
    // Get values from context
    userID, _ := c.Get("userID")           // uint
    email, _ := c.Get("email")             // string
    role, _ := c.Get("role")               // UserRole
    claims, _ := c.Get("claims")           // *CustomClaims
    
    // Use in handler logic
    log.Printf("User %d (%s) with role %s", userID, email, role)
}
```

## Token Debugging

Decode JWT at [jwt.io](https://jwt.io) to inspect token contents (for development only).

Example token:
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJyb2xlIjoiY3VzdG9tZXIiLCJleHAiOjE3MzIxMTg0MDB9.
abc123...
```

Three parts separated by `.`:
1. **Header**: Algorithm info
2. **Payload**: Claims (user data)
3. **Signature**: HMAC signature

## File Structure
```
├── internal/
│   ├── models/
│   │   ├── user.go          # User model + roles
│   │   ├── customer.go      # Customer model (updated)
│   │   └── ...
│   ├── service/
│   │   ├── auth_service.go  # JWT generation/validation
│   │   └── auth_service_test.go
│   ├── repository/
│   │   ├── user_repository.go       # Interface
│   │   ├── sql_user_repository.go   # GORM impl
│   │   └── sql_user_repository_test.go
│   ├── handlers/
│   │   ├── login_handler.go         # Login endpoint
│   │   ├── auth_handler_test.go
│   │   ├── jwt_middleware.go        # Auth + role middleware
│   │   ├── jwt_middleware_test.go
│   │   └── user_handler.go          # Admin user management
│   └── mocks/
│       └── mock_user_repository.go  # For testing
├── main.go                          # Application entry point
├── .env.example                     # Environment template
├── go.mod                           # Dependencies
├── AUTH_DOCUMENTATION.md            # Full guide
└── IMPLEMENTATION_SUMMARY.md        # What was done
```

---

For detailed information, see `AUTH_DOCUMENTATION.md`
