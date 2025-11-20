
# JWT Authentication System Documentation

## Overview

The application now features a complete JWT-based authentication system with role-based access control (RBAC). Users can be either `admin` or `customer`, with different permissions for each role.

## Architecture

### Models

#### User Model (`internal/models/user.go`)
The core user entity with the following fields:
- `ID`: Unique identifier
- `Email`: User's email (unique)
- `Password`: Bcrypt-hashed password (never serialized)
- `Role`: User role (`admin` or `customer`)
- `Name`: User's full name
- `Phone`: User's phone number
- `IsActive`: Account status flag
- `CreatedAt`, `UpdatedAt`: Timestamps

#### Customer Model (Updated)
Now references a User through the `UserID` foreign key relationship.

### Services

#### AuthService (`internal/service/auth_service.go`)
Handles JWT token generation and validation:
- `GenerateToken(user User) (string, error)`: Creates JWT token from user data
- `ValidateToken(token string) (*CustomClaims, error)`: Validates token and returns claims
- `ValidateTokenWithRole(token string, roles... UserRole) (*CustomClaims, error)`: Validates token with role-based access

**JWT Claims Structure:**
```go
type CustomClaims struct {
    UserID     uint             // User ID
    Email      string           // User email
    Role       UserRole         // User role (admin/customer)
    CustomerID *uint            // Optional customer ID
    jwt.RegisteredClaims         // Standard JWT claims (expiration, etc.)
}
```

Token expiration: 24 hours

### Repositories

#### UserRepository Interface (`internal/repository/user_repository.go`)
Methods:
- `Create(user User) (User, error)`
- `FindByID(id uint) (User, error)`
- `FindByEmail(email string) (User, error)`
- `Update(user User) error`
- `Delete(id uint) error`
- `FindAll() ([]User, error)`
- `FindByRole(role UserRole) ([]User, error)`

#### SQL User Repository (`internal/repository/sql_user_repository.go`)
GORM-based implementation of UserRepository using SQLite.

### Handlers

#### Authentication Handler (`internal/handlers/login_handler.go`)
**Endpoints:** 
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - Public customer registration

**Login Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Register Request:**
```json
{
  "email": "newcustomer@example.com",
  "password": "password123",
  "name": "Customer Name",
  "phone": "123456789"
}
```

**Response (Success - 200/201)**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe",
    "role": "customer"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid email format or password too short
- `401 Unauthorized`: Invalid credentials or inactive user (login only)
- `409 Conflict`: Email already registered (register only)

#### User Management Handler (`internal/handlers/user_handler.go`)
Admin-only endpoints for user management:
- `GET /api/admin/users` - List all users
- `POST /api/admin/users` - Create new user
- `GET /api/admin/users/{id}` - Get specific user
- `PUT /api/admin/users/{id}` - Update user
- `DELETE /api/admin/users/{id}` - Delete user

### Middleware

#### JWT Authentication Middleware (`internal/handlers/jwt_middleware.go`)

**JWTAuthMiddleware**
Validates JWT token and extracts claims. All authenticated users can pass through.

Usage:
```go
protected := r.Group("/api")
protected.Use(handlers.JWTAuthMiddleware(authSvc))
```

**RequireRole Middleware**
Role-based access control. Only users with specified roles can pass.

Usage:
```go
protected.GET("/admin/users", 
    handlers.RequireRole(authSvc, models.RoleAdmin), 
    handlers.GetAllUsers(userRepo))
```

Context values set by middleware:
- `userID`: User ID (uint)
- `email`: User email (string)
- `role`: User role (UserRole)
- `claims`: Full CustomClaims object

## Usage Flow

### 1. Customer Registration (Public)
New customers can self-register without admin intervention:

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newcustomer@example.com",
    "password": "password123",
    "name": "Jane Doe",
    "phone": "123456789"
  }'
```

Response includes JWT token (201 Created).

**Key Points:**
- Anyone can register
- New users are automatically assigned `customer` role
- Email must be unique (409 Conflict if already exists)
- Password minimum 6 characters

### 2. User Login
Users (both customers and admins) login with credentials:

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

Response includes JWT token.

### 3. Accessing Protected Routes
Include token in Authorization header:

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

Response includes JWT token.

### 3. Accessing Protected Routes
Include token in Authorization header:

```bash
curl -X GET http://localhost:8080/api/appointments \
  -H "Authorization: Bearer <your-jwt-token>"
```

### 4. Role-Based Access
Admin routes only allow users with `admin` role:

```bash
curl -X GET http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer <admin-token>"
```

If a customer tries to access: `403 Forbidden`

## Admin User Creation

Since the registration endpoint only creates `customer` accounts, admin users must be created through direct database access or an admin-only endpoint. For development/setup:

```go
// Create admin user (typically done once during setup)
admin := models.User{
    Email:    "admin@example.com",
    Password: hashPassword("adminPassword123"),
    Name:     "Admin User",
    Role:     models.RoleAdmin,
    IsActive: true,
}
userRepo.Create(admin)
```

Or via admin API (requires existing admin token):

```bash
curl -X POST http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newadmin@example.com",
    "password": "adminPass123",
    "name": "New Admin",
    "phone": "987654321",
    "role": "admin"
  }'
```

## Testing

### Unit Tests
Comprehensive test coverage using `testify/assert` and `gomock`:

**Auth Service Tests** (88% coverage):
- Token generation and validation
- Role-based token validation
- Invalid token handling
- Wrong secret detection

**JWT Middleware Tests** (27.9% coverage):
- Valid token verification
- Missing/invalid headers
- Role-based access control
- Unauthorized role rejection

**Auth Handler Tests**:
- Successful login
- Invalid credentials
- User account status validation
- Password verification

**User Repository Tests**:
- Create, read, update, delete operations
- Find by email and role
- Behavioral tests with mock repository

**Run Tests:**
```bash
# All tests
go test -v ./internal/service ./internal/repository ./internal/handlers -cover

# Specific test
go test -v ./internal/service -run TestGenerateToken

# With coverage report
go test ./internal/service -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Environment Setup

### .env Configuration
```
DB_PATH=app.db
JWT_SECRET=your-strong-secret-key-here
```

**Important:** In production, use a strong, randomly-generated JWT_SECRET with at least 32 characters.

### Database Initialization
Database schema is automatically migrated on application startup via `AutoMigrate()`.

Tables created:
- `users`: User credentials and profile
- `customers`: Customer-specific information
- `services`: Service offerings
- `appointments`: Appointment records

## Security Considerations

1. **Password Hashing**: Uses `golang.org/x/crypto/bcrypt` with default cost
2. **JWT Secret**: Must be strong and kept secure (minimum 32 characters recommended)
3. **Token Expiration**: Tokens expire after 24 hours
4. **HTTPS**: Use HTTPS in production (not enforced by application)
5. **Active Status**: Users can be marked inactive; inactive users cannot login
6. **Password Never Serialized**: Passwords are never included in API responses

## Adding New Roles

To add new roles, extend the `UserRole` type in `models/user.go`:

```go
type UserRole string

const (
    RoleAdmin     UserRole = "admin"
    RoleCustomer  UserRole = "customer"
    RoleStaff     UserRole = "staff"      // New role
)
```

Then use in role checks:
```go
handlers.RequireRole(authSvc, models.RoleStaff)
```

## Future Enhancements

- [ ] Email verification on registration
- [ ] Password reset functionality
- [ ] Refresh token implementation
- [ ] Multi-factor authentication (MFA)
- [ ] Rate limiting on login attempts
- [ ] Audit logging for authentication events
- [ ] OAuth2/OpenID Connect integration
- [ ] API key authentication for service-to-service calls
