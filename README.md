# StoreX

StoreX is a simple **RESTful Inventory Management System** built with **Go** for tracking and managing inventory for office or school use.

---

## 🚀 Features

- 🗃️ CRUD APIs for inventory items
- 🚦 Structured handlers and middleware
- 📦 Database storage using SQL
- 📋 Basic validation and logging

---

## 🧰 Tech Stack

| Component | Technology |
|-----------|------------|
| Language  | Go         |
| Router    | chi        |
| Validator | go-playground/validator |
| Logger    | logrus     |
| Database  | PostgreSQL / SQL supported by `sqlx` |

---

## 📦 Installation

### 1. Clone the repo

Here's a plain English explanation of every API:

---

### 🔐 Auth

**Register** → Creates your account. You send your name, email, password, phone and role. You get back a token to use for future requests.

**Login** → You send email + password. If correct, you get a token back.

**Logout** → Kills your current session. Your token stops working after this.

---

### 👥 Employee

**Create Employee** *(Admin only)* → Admin creates a new employee account manually instead of them self-registering.

**List Employees with Asset Count** → Shows all employees with how many assets they currently have. You can filter by asset type (like laptop) or status (like assigned). If you don't filter, it shows total assets for everyone.

**Archive Employee** *(Admin only)* → Fires an employee from the system. They don't get deleted permanently, just marked as inactive.

**My Assets** → An employee can see which assets are currently assigned to them. Uses their token to figure out who they are automatically.

---

### 📦 Assets

**Create Asset** *(Admin only)* → Adds a new device to the system. You provide brand, model, serial number, type and depending on the type you also fill in extra details like RAM for laptop or DPI for mouse.

**List All Assets** → Shows all assets in the company. You can search, filter by type/status/owner and paginate. Also returns a dashboard count showing how many are available, assigned, damaged etc.

**Get One Asset** → Fetches full details of a single asset including its type-specific specs like processor, RAM, storage for a laptop.

**Assign Asset** *(Admin only)* → Links an asset to an employee. Asset status changes to assigned. A record is saved in history.

**Return Asset** *(Admin only)* → Takes an asset back from an employee. Asset goes back to available. A return record is saved in history with a note.

**Update Asset** *(Admin only)* → Edit the details of an existing asset like brand, model, warranty dates or laptop specs.

**Delete Asset** *(Admin only)* → Removes an asset from the system. Asset must not be currently assigned to anyone. It's a soft delete so the data is still in the database just hidden.

```sh
git clone https://github.com/wibecoderr/storex.git
cd storex
