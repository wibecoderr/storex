# Storex — Asset Management System

A REST API for managing company assets and employees. Built with Go.

---

## Tech Stack

- **Go** with [chi](https://github.com/go-chi/chi) router
- **PostgreSQL** with [sqlx](https://github.com/jmoiron/sqlx)
- **JWT** authentication with session management
- **golang-migrate** for database migrations
- **bcrypt** for password hashing

---

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL

### Environment Variables

```env
DB_HOST=localhost
DB_PORT=5433
DB_DATABASE=postgres
DB_USER=local
DB_PASSWORD=local
JWT_SECRET=your_secret_key
```

### Run

```bash
go run main.go
```

Server starts on `:8080`. Migrations run automatically on startup.

---

## How It Works

### 🔐 Auth

**Register** `POST /register`
Creates your account. You send your name, email, password, phone and role. You get back a token to use for future requests. Only non-admin roles can self-register — admin accounts must be created by an existing admin.

**Login** `POST /login`
You send email and password. If correct, you get a JWT token back to use in future requests.

**Logout** `POST /logout`
Kills your current session. Your token stops working immediately after this.

---

### 👥 Employees

**Create Employee** `POST /register/employee` *(Admin only)*
Admin creates a new employee account manually instead of them self-registering. Can assign any role including admin.

**List Employees** `GET /employees`
Shows all employees with how many assets they currently have. You can filter by asset type (e.g. `laptop`) or status (e.g. `assigned`). If you don't filter, it shows total asset count for everyone.

```
GET /employees?type=laptop&status=assigned
```

**Archive Employee** `DELETE /user/{id}` *(Admin only)*
Removes an employee from the system. They are not permanently deleted — just marked as inactive. Employee must have no assets currently assigned to them before archiving.

**My Assets** `GET /employees/asset`
An employee can see which assets are currently assigned to them. Uses their token to figure out who they are automatically — no need to pass an ID.

---

### 📦 Assets

**Create Asset** `POST /assets` *(Admin only)*
Adds a new device to the system. You provide brand, model, serial number, type and owner. Depending on the type you also fill in extra details:

| Type | Extra Fields Required |
|------|-----------------------|
| `laptop` | processor, ram, storage, os, charger |
| `mouse` | dpi, is_wireless |
| `keyboard` | layout |
| `mobile` | os, ram, storage, charger |
| `hardware` | storage |

**List All Assets** `GET /assets`
Shows all assets in the company. Also returns a dashboard summary showing how many are available, assigned, in service, waiting for repair or damaged. Supports filtering and pagination:

```
GET /assets?type=laptop&status=available&owner=company&search=dell&page=1&limit=10
```

**Get One Asset** `GET /assets/{id}`
Fetches full details of a single asset including its type-specific specs such as processor, RAM and storage for a laptop.

**Assign Asset** `POST /assets/assign` *(Admin only)*
Links an asset to an employee. The asset status changes to `assigned` and a record is saved in history. Asset must not already be assigned to someone else.

**Return Asset** `POST /assets/return/{id}` *(Admin only)*
Takes an asset back from an employee. Asset status goes back to `available`. A return record is saved in history along with a note explaining why it was returned.

**Update Asset** `PUT /assets/update/{id}` *(Admin only)*
Edit the details of an existing asset such as brand, model, warranty dates or type-specific specs like laptop RAM or mouse DPI.

**Delete Asset** `DELETE /assets/{id}` *(Admin only)*
Removes an asset from the system. The asset must not be currently assigned to anyone. This is a soft delete — the data stays in the database but is hidden from all listings.

---

## Authentication

All protected routes require a Bearer token in the Authorization header:

```
Authorization: Bearer <token>
```

Tokens expire after 10 minutes. Sessions are tracked in the database and invalidated immediately on logout.

---

## Roles

| Role | Can Self-Register | Admin Access |
|------|------------------|--------------|
| `admin` | ❌ Must be created by admin | ✅ Full access |
| `manager` | ✅ | ❌ |
| `employee` | ✅ | ❌ |
| `intern` | ✅ | ❌ |
| `freelancer` | ✅ | ❌ |

---

## Asset Status

| Status | Meaning |
|--------|---------|
| `available` | Not assigned to anyone, ready to use |
| `assigned` | Currently with an employee |
| `in_service` | Being serviced or maintained |
| `for_repair` | Waiting to be repaired |
| `damaged` | Damaged, not usable |

---

## Project Structure

```
.
├── main.go
├── server/
│   └── server.go            # Route definitions
├── handler/
│   ├── user.go              # Auth & employee handlers
│   └── assets.go            # Asset handlers
├── middleware/
│   └── auth.go              # JWT & role middleware
├── database/
│   ├── db.go                # DB connection & transactions
│   ├── dbhelper/
│   │   ├── employee.go      # Employee queries
│   │   └── assest.go        # Asset queries
│   └── migrations/
│       └── 000001_initial_up.sql
└── model/
    └── model.go             # Structs & types
```
