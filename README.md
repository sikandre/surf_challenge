# SURF Code Challenge

### Table of Contents

i. [Overview](#overview)  
ii. [Tech Stack](#tech-stack)  
iii. [Project Structure](#project-structure)  
iv. [Endpoints](#endpoints)  
v. [Commands](#commands)  
vi. [Notes & Assumptions](#notes--assumptions)


<a name="Overview"></a>
## Overview

This project is a back-end coding challenge implementation in **Go**.  
It simulates a simple web server that loads data from JSON files and exposes a set of REST API endpoints.

The two main entities are:

- **User** – basic information about a product’s user
- **Action** – a record of a user’s activity inside the product (e.g. refer another user, view a conversation, add to CRM)

The challenge goals were to:

1. Fetch a user by ID
2. Return the total number of actions performed by a user
3. Calculate the probability distribution of next actions given a starting action type
4. Compute the *Referral Index* for all users, i.e. the number of direct and indirect users invited by a given user

The server reads data from the provided `users.json` and `actions.json` files into memory at startup — no database is required.

<a name="tech-stack"></a>
## Tech Stack
- Go 1.24
- Chi (HTTP router)
- Zap (structured logging)
- GoMock / Testify (testing & mocking)
- golangci-lint (linting)

<a name="endpoints"></a>

<a name="project-structure"></a>
## Project Structure
```
.
├── README.md
├── cmd
│   └── main.go
├── go.mod
├── go.sum
└── internal
    ├── action
    │   ├── domain
    │   │   └── domain.go
    │   ├── mapper
    │   │   └── mapper.go
    │   ├── service.go
    │   ├── service_mock.go
    │   ├── service_test.go
    │   └── storage
    │       ├── db
    │       │   └── actions.json
    │       ├── entity
    │       │   └── entity.go
    │       ├── repository.go
    │       └── repository_mock.go
    ├── api
    │   ├── action
    │   │   ├── dto
    │   │   │   └── response.go
    │   │   ├── handler.go
    │   │   ├── handler_test.go
    │   │   └── mapper
    │   │       └── mapper.go
    │   ├── apierror
    │   │   └── error.go
    │   ├── router
    │   │   └── router.go
    │   └── user
    │       ├── dto
    │       │   └── response.go
    │       ├── handler.go
    │       ├── handler_test.go
    │       └── mapper
    │           └── mapper.go
    ├── container
    │   └── container.go
    ├── converter
    │   └── utils.go
    └── user
        ├── domain
        │   └── domain.go
        ├── mapper
        │   └── mapper.go
        ├── service.go
        ├── service_mock.go
        ├── service_test.go
        └── storage
            ├── db
            │   └── users.json
            ├── entity
            │   └── entity.go
            ├── repository.go
            └── repository_mock.go
```

<a name="endpoints"></a>
## Endpoints

Base URL: `/api/v1`

---

### 1) List users (paginated)
**GET** `/users`

**Query params**
- `userId` _(optional, int)_ — if provided, returns only that user (still wrapped in the list response)
- `page` _(optional, int, default: 1)_
- `pageSize` _(optional, int, default: 10)_

**Response 200**
```json
{
  "users": [
    {
      "id": "1",
      "name": "John Doe",
      "createdAt": "2023-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "totalItems": 1,
    "totalPages": 1,
    "page": 1,
    "pageSize": 10
  }
}
```

**Errors**
- `400` invalid query parameters
- `500` internal error

---

### 2) Get user by ID
**GET** `/users/{userId}`

**Path params**
- `userId` _(required, int)_

**Response 200**
```json
{
  "id": "1234",
  "name": "John Doe",
  "createdAt": "2022-04-14T11:12:22.758Z"
}
```

**Errors**
- `400` invalid `userId`
- `404` user not found
- `500` internal error

---

### 3) Get total number of actions for a user
**GET** `/users/{userId}/actions/count`

**Path params**
- `userId` _(required, int)_

**Response 200**
```json
{
  "count": 100
}
```

**Errors**
- `400` invalid `userId`
- `404` user not found
- `500` internal error

---

### 4) Get next-action probability breakdown
**GET** `/actions/next-probability`

**Query params**
- `next` _(required, string)_ — the current action type (e.g., `ADD_TO_CRM`, `REFER_USER`, `VIEW_CONVERSATION`)

**Response 200**
```json
{
  "ADD_TO_CRM": 0.70,
  "REFER_USER": 0.20,
  "VIEW_CONVERSATION": 0.10
}
```
> Values are probabilities in the range `[0,1]` formatted as float with two decimal places ordered by most probable.

**Errors**
- `400` missing/invalid `next`
- `500` internal error

---

### 5) Get Referral Index for all users
**GET** `/actions/referrals`

Computes, for each user, the number of **unique** users they referred directly or indirectly (a user can be invited only once).

**Response 200**
```json
{
  "1": 3,
  "2": 0,
  "3": 7
}
```

**Errors**
- `500` internal error

---

### Error format
When returned as JSON, errors follow:
```json
{
  "message": "Resource not found",
  "code": 404
}
```

---

### Quick cURL examples

```bash
# List users (page 2, size 5)
curl "http://localhost:3000/api/v1/users?page=2&pageSize=5"

# Get user by ID
curl "http://localhost:3000/api/v1/users/1"

# Actions count for a user
curl "http://localhost:3000/api/v1/users/1/actions/count"

# Next-action probabilities (after VIEW_CONVERSATION)
curl "http://localhost:3000/api/v1/actions/next-probability?next=EDIT_CONTACT"

# Referral index (all users)
curl "http://localhost:3000/api/v1/actions/referrals"
```


<a name="commands"></a>
## Commands
> Requirements: Go 1.21+ (or the version in `go.mod`), `golangci-lint` (optional), `mockgen` (optional).

### Run the server
```bash
  go run ./cmd
# SStarting server on :3000
```

**Build**
```bash
  go build -o bin/surf-challenge ./cmd
./bin/surf-challenge
```
### Linting
(optional, if you use golangci-lint)
```bash
  golangci-lint run ./...
```

### Generating mocks
```bash
  go generate ./...
```

### Testing
```bash
  go test -v ./...
```

<a name="notes--assumptions"></a>
## Notes & Assumptions

- **In-memory data**:  
  Users and Actions are loaded from the provided JSON files (`users.json` and `actions.json`) at startup.  
  There is no database or persistence layer.

- **Pagination**:
    - Default `page=1` and `pageSize=10`
    - Query parameters are validated, invalid values return an API error.

- **Case-insensitive queries**:  
  Action types are compared using case-insensitive matching when computing probabilities.

- **Next action probability**:
    - For each user, actions are ordered by `createdAt`.
    - We compute the probability distribution of the *immediately following* actions after a given action type.
    - Returned values are normalized decimals (e.g., `0.70`) instead of percentages.

- **Referral index**:
    - A user can only be invited **once**.
    - Computed using DFS traversal of the referral graph.
    - Complexity: **O(N + M)** where `N` is number of users and `M` number of referral edges.
    - Cycles are prevented by tracking visited nodes.

- **Error handling**:
    - All errors are mapped to a consistent JSON structure `{ "message": "...", "code": ... }`.
    - `404` for not found, `400` for invalid input, `500` for unexpected errors.

- **Testing**:
    - Handlers tested with `httptest` and mocked services.
    - Services tested with mocked repositories.

- **Assumption**:
    - Input data (`users.json`, `actions.json`) is well-formed and valid.
    - IDs are unique and consistent across files.
    - Dates are in ISO-8601 / RFC3339 format.  
