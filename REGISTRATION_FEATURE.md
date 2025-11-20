# Customer Self-Registration Feature

## What Was Added

A public customer registration endpoint that allows potential clients to create their own accounts without requiring admin intervention.

## New Endpoint

### POST `/api/auth/register` (Public - No Auth Required)

**Purpose**: Allow new customers to self-register

**Request Body**:
```json
{
  "email": "newcustomer@example.com",
  "password": "password123",
  "name": "Customer Name",
  "phone": "123456789"
}
```

**Response (201 Created)**:
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

## Features

✅ **Public Registration**: Anyone can create a customer account  
✅ **Automatic JWT Token**: User receives token immediately after registration  
✅ **Email Uniqueness**: Prevents duplicate email addresses (409 Conflict)  
✅ **Password Hashing**: Uses bcrypt for security  
✅ **Customer Role**: All new registrations are automatically assigned `customer` role  
✅ **Input Validation**: 
   - Email must be valid format
   - Password minimum 6 characters
   - Name and phone required

## Error Responses

| Status | Condition |
|--------|-----------|
| 400 | Invalid email format, password < 6 chars, or missing fields |
| 409 | Email already registered |
| 500 | Server error |

## Implementation Details

### Request Structure (RegisterRequest)
```go
type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    Name     string `json:"name" binding:"required"`
    Phone    string `json:"phone" binding:"required"`
}
```

### Response Structure (RegisterResponse)
```go
type RegisterResponse struct {
    Token string   `json:"token"`
    User  UserInfo `json:"user"`
}
```

### Handler Logic
1. Extract and validate registration data
2. Check if email already exists (prevent duplicates)
3. Hash password using bcrypt
4. Create User with `customer` role
5. Generate JWT token
6. Return token + user info

## Tests Added

6 new test cases for registration:

1. **TestRegister_Success** - Happy path registration
2. **TestRegister_EmailAlreadyExists** - Duplicate email handling (409)
3. **TestRegister_InvalidEmail** - Invalid email format (400)
4. **TestRegister_MissingFields** - Missing required fields (400)
5. **TestRegister_ShortPassword** - Password too short (400)
6. **TestRegister_AlwaysCreatesCustomer** - Verifies role is always `customer`

**Test Results**: ✅ All 6 tests passing

## Usage Example

```bash
# Register as new customer
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "jane@example.com",
    "password": "myPassword123",
    "name": "Jane Doe",
    "phone": "555-1234"
  }'

# Response
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "jane@example.com",
    "name": "Jane Doe",
    "role": "customer"
  }
}

# Use token to access protected routes
curl -X GET http://localhost:8080/api/appointments \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## Flow Comparison

### Before (Admin-Only Registration)
```
1. Admin creates user account
2. Customer receives credentials
3. Customer logs in
4. Customer accesses features
```

### After (Self-Registration)
```
1. Customer registers themselves
2. Customer immediately gets token
3. Customer accesses features
4. Admin can still create additional admin accounts
```

## Database Impact

No new tables or schema changes needed. Uses existing `users` table with new rows:
- `role` = "customer"
- `is_active` = true

## Files Modified

1. **internal/handlers/login_handler.go**
   - Added `RegisterRequest` struct
   - Added `RegisterResponse` struct
   - Added `Register()` handler method

2. **main.go**
   - Added route: `public.POST("/api/auth/register", authHandler.Register)`

3. **internal/handlers/auth_handler_test.go**
   - Added 6 registration test cases

4. **QUICK_REFERENCE.md**
   - Added registration documentation
   - Added example usage

5. **AUTH_DOCUMENTATION.md**
   - Updated usage flow section
   - Added admin creation guidance

## Security Considerations

✅ Passwords are hashed with bcrypt (not stored in plaintext)  
✅ Email is unique per user  
✅ No password in API responses  
✅ Input validation on all fields  
✅ Tokens are signed with HMAC-SHA256  
✅ New customers get `customer` role (cannot escalate to admin)  

## Backward Compatibility

✅ Existing login endpoint unchanged  
✅ Existing authentication unchanged  
✅ Admin user creation via API still works  
✅ All existing tests pass  

## Total Test Coverage

- Before: 44 tests
- After: 50 tests (added 6)
- All passing ✅

## Summary

Now your hair salon application allows:
1. **Customers** can self-register and immediately start booking appointments
2. **Admins** can still create admin accounts and manage users
3. **Zero friction** for new customers wanting to use the service

The implementation is production-ready with full test coverage and comprehensive error handling.
