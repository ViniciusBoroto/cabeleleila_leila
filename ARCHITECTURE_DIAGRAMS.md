# JWT Authentication Flow Diagrams

## 1. Login Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                          LOGIN FLOW                             │
└─────────────────────────────────────────────────────────────────┘

CLIENT                           API SERVER
   │                                 │
   │  POST /api/auth/login           │
   │  {email, password}              │
   ├────────────────────────────────>│
   │                                 │
   │                          FindByEmail
   │                                 │
   │                          Compare Password
   │                          (bcrypt.CompareHashAndPassword)
   │                                 │
   │                          ✓ Valid? → Generate JWT
   │                                 │
   │  {token, user_info}             │
   │<────────────────────────────────┤
   │                                 │
   ●                                 ●
  (Store Token)              (Token expires in 24h)
```

## 2. Protected Route Access Flow

```
┌──────────────────────────────────────────────────────────────────┐
│                    PROTECTED ROUTE ACCESS                        │
└──────────────────────────────────────────────────────────────────┘

CLIENT                           API SERVER
   │                                 │
   │  GET /api/appointments          │
   │  Authorization: Bearer <token>  │
   ├────────────────────────────────>│
   │                          ┌──────────────────┐
   │                          │ JWTAuthMiddleware│
   │                          └──────────────────┘
   │                                 │
   │                          Extract Token
   │                          Validate Signature
   │                          Check Expiration
   │                                 │
   │                          ✓ Valid? → Continue
   │                                 │
   │                          Set Context:
   │                          - userID
   │                          - email
   │                          - role
   │                          - claims
   │                                 │
   │                          Execute Handler
   │                                 │
   │  [Appointments Data]            │
   │<────────────────────────────────┤
   │                                 │
```

## 3. Admin Route Authorization Flow

```
┌──────────────────────────────────────────────────────────────────┐
│                  ROLE-BASED ACCESS CONTROL                       │
└──────────────────────────────────────────────────────────────────┘

CLIENT                           API SERVER
   │                                 │
   │  GET /api/admin/users           │
   │  Authorization: Bearer <token>  │
   ├────────────────────────────────>│
   │                          ┌──────────────────────┐
   │                          │ RequireRole Middleware│
   │                          └──────────────────────┘
   │                                 │
   │                          Validate Token
   │                          Extract Claims
   │                          Check Role
   │                                 │
   │                    ┌─────────────────────────┐
   │                    │  Role == "admin"?       │
   │                    └─────────────────────────┘
   │                              │
   │              ┌───────────────┴───────────────┐
   │              │                               │
   │          YES │                           NO  │
   │              ▼                               ▼
   │          Continue             403 Forbidden
   │              │                    │
   │         Execute                   │<─── Respond
   │         Handler         
   │              │
   │    [Users Data]
   │              │
   │<─────────────────────────────────┤
   │
```

## 4. Token Structure

```
┌─────────────────────────────────────────────────────────────────┐
│                      JWT TOKEN STRUCTURE                        │
└─────────────────────────────────────────────────────────────────┘

┌──────────────┐
│   HEADER     │
├──────────────┤
│ {            │
│  "alg":      │
│   "HS256",   │
│  "typ":      │
│   "JWT"      │
│ }            │
└──────────────┘
       │
       │ (Base64 Encoded)
       │
       ▼
┌────────────────────────────────────────────┐
│ eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9      │
└────────────────────────────────────────────┘
       │
       ●
┌──────────────────────┐
│    PAYLOAD           │
├──────────────────────┤
│ {                    │
│  "user_id": 1,       │
│  "email":            │
│   "user@test.com",   │
│  "role": "customer", │
│  "exp": 1732118400,  │
│  "iat": 1732032000   │
│ }                    │
└──────────────────────┘
       │
       │ (Base64 Encoded)
       │
       ▼
┌──────────────────────────────────────────────────────┐
│ eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InRlc3RAZXhhbXBsZS... │
└──────────────────────────────────────────────────────┘
       │
       ●
┌──────────────────────────────────────┐
│    SIGNATURE                         │
├──────────────────────────────────────┤
│ HMAC-SHA256(                         │
│   header.payload,                    │
│   JWT_SECRET                         │
│ )                                    │
└──────────────────────────────────────┘
       │
       │ (Base64 Encoded)
       │
       ▼
┌──────────────────────────────────────────────────────┐
│ SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c      │
└──────────────────────────────────────────────────────┘
       │
       │ Combined with dots
       │
       ▼
┌────────────────────────────────────────────────────────────────┐
│ eyJhbGci...9        .        eyJ1c2VyX2l...       .      SflKx...│
│  header              separator        payload              signature   │
└────────────────────────────────────────────────────────────────┘
                    (Final JWT Token)
```

## 5. User Model Relationships

```
┌──────────────────────────────────────────────┐
│              USER TABLE                      │
├──────────────────────────────────────────────┤
│ id (PK)                                      │
│ email (UNIQUE)                               │
│ password (bcrypt hash)                       │
│ role (admin|customer)                        │
│ name                                         │
│ phone                                        │
│ is_active (bool)                             │
│ created_at, updated_at                       │
└──────────────┬───────────────────────────────┘
               │
               │ (1:1 relationship)
               │ FK: user_id
               │
               ▼
┌──────────────────────────────────────────────┐
│           CUSTOMER TABLE                     │
├──────────────────────────────────────────────┤
│ id (PK)                                      │
│ user_id (FK) ──────────┐                     │
│ is_active              │ References User
│ created_at, updated_at │                     │
└───────────────────────┬────────────────────┘
                        │
                        ▼
               (Customer-specific data)
```

## 6. Authorization Middleware Chain

```
┌─────────────────────────────────────────────────────────────┐
│               MIDDLEWARE CHAIN                              │
└─────────────────────────────────────────────────────────────┘

Request
   │
   ▼
┌──────────────────────────────────┐
│  JWTAuthMiddleware               │
│ ├─ Extract Bearer token          │
│ ├─ Validate signature            │
│ ├─ Check expiration              │
│ ├─ Set context: userID, email,   │
│ │  role, claims                  │
│ └─ Continue if valid             │
└──────────────────────────────────┘
   │
   ├─ Invalid? ──────────────────┐
   │                             │
   │ Valid                       ▼
   ▼                      401 Unauthorized
   │
   ▼ (Optional) RequireRole Middleware
┌──────────────────────────────────┐
│  RequireRole(role1, role2...)    │
│ ├─ Get role from context         │
│ ├─ Check if role in allowed list │
│ └─ Continue if authorized        │
└──────────────────────────────────┘
   │
   ├─ Not authorized? ───────────┐
   │                             │
   │ Authorized                  ▼
   ▼                      403 Forbidden
   │
   ▼
┌──────────────────────────────────┐
│  Handler Function                │
│ ├─ Access context values         │
│ ├─ Execute business logic        │
│ └─ Return response               │
└──────────────────────────────────┘
   │
   ▼
Response
```

## 7. Password Hashing & Verification

```
┌──────────────────────────────────────────────────────────┐
│         PASSWORD HASHING & VERIFICATION                 │
└──────────────────────────────────────────────────────────┘

REGISTRATION                         LOGIN
───────────                          ─────

User enters                          User enters
password                             password
   │                                    │
   ▼                                    ▼
┌──────────────────┐             ┌──────────────────┐
│ bcrypt.Generate  │             │ Find user by     │
│ FromPassword()   │             │ email            │
└──────────────────┘             └──────────────────┘
   │                                    │
   ▼                                    ▼
Hashed password                    ┌──────────────────────┐
(e.g. $2a$10$...)                  │ bcrypt.Compare       │
   │                                │ HashAndPassword      │
   ▼                                │ (storedHash,         │
Store in DB                         │  enteredPassword)    │
(Never store plain)                 └──────────────────────┘
                                       │
                                       ├─ Match ────────────┐
                                       │                    │
                                       │ No match    ▼
                                       │          Reject
                                       ▼ Match
                                    Generate
                                    JWT Token
                                       │
                                       ▼
                                    Return to
                                    Client
```

## 8. Complete Request Lifecycle

```
┌──────────────────────────────────────────────────────────────┐
│            COMPLETE REQUEST LIFECYCLE                        │
└──────────────────────────────────────────────────────────────┘

┌─────────────┐
│   CLIENT    │
└──────┬──────┘
       │
       │ HTTP Request
       │ (with Bearer token)
       │
       ▼
┌──────────────────────────────────┐
│  Gin Router                      │
│  (Route matching)                │
└──────────────┬───────────────────┘
               │
               ▼
      ┌────────────────────────────┐
      │ Middleware 1               │
      │ JWTAuthMiddleware          │
      └────────┬───────────────────┘
               │ (Valid token)
               │
               ▼
      ┌────────────────────────────┐
      │ Middleware 2               │
      │ RequireRole (optional)     │
      └────────┬───────────────────┘
               │ (Role authorized)
               │
               ▼
      ┌────────────────────────────┐
      │ Handler Function           │
      │                            │
      │ ctx.Get("userID")          │
      │ ctx.Get("role")            │
      │ ctx.Get("claims")          │
      │                            │
      │ Business Logic             │
      └────────┬───────────────────┘
               │
               ▼
      ┌────────────────────────────┐
      │ Database Query             │
      │ (using userID)             │
      └────────┬───────────────────┘
               │
               ▼
      ┌────────────────────────────┐
      │ Response Generation        │
      │ c.JSON(200, data)          │
      └────────┬───────────────────┘
               │
               ▼
       │ HTTP Response
       │ (JSON)
       │
       ▼
┌─────────────┐
│   CLIENT    │
│  (Receives  │
│   Response) │
└─────────────┘
```

---

These diagrams illustrate the complete authentication flow and architecture of your JWT-based system.
