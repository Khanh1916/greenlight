# Greenlight API

<div align="center">

![Go](https://img.shields.io/badge/Go-1.24%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13%2B-336791?style=for-the-badge&logo=postgresql&logoColor=white)
![HTTPRouter](https://img.shields.io/badge/HTTPRouter-v1.3-blue?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Status](https://img.shields.io/badge/Status-Completed-success?style=for-the-badge)

<p align="center">
  A production-ready, feature-rich RESTful JSON API built with Go and PostgreSQL.<br>
  Dịch vụ RESTful JSON API chuẩn production được xây dựng bằng Go và PostgreSQL.
</p>

[**English**](#-english) &nbsp;|&nbsp; [**Tiếng Việt**](#-tiếng-việt)

</div>

---

<a name="english"></a>
## 🇬🇧 English

### 📖 About the Project
**Greenlight** is a full-featured RESTful JSON API designed and implemented in **Go** backed by **PostgreSQL**. Following the architecture and patterns from Alex Edwards' *"Let's Go Further!"*, Greenlight demonstrates production-grade Go design patterns: idiomatic package structuring, strict request/response handling, custom JSON encoding/decoding, robust data validation, stateful token authentication, role-based authorization, rate limiting, background workers, and graceful shutdown.

---

### ✨ Key Features

- **Movies Management (CRUD & Querying)**:
  - Create, fetch, update, and delete movies with relational array handling for genres.
  - **Partial Updates (`PATCH`)**: Differentiates between omitted JSON fields and zero-valued fields using pointers.
  - **Optimistic Concurrency Control**: Prevents race conditions during simultaneous updates using version tokens (`409 Conflict`).
  - **Full-Text Search & Reductive Filtering**: Powered by PostgreSQL `to_tsvector` / `plainto_tsquery` with GIN indexing.
  - **Pagination & Sorting**: Dynamic sorting with safelisting, offset/limit pagination, and metadata calculation via `count(*) OVER()`.

- **Authentication & Authorization**:
  - **Account Registration & Password Hashing**: Secure password hashing using `bcrypt` (cost: 12).
  - **Account Activation**: Cryptographically secure 128-bit CSPRNG tokens encoded in base32 (unpadded) and hashed with SHA-256.
  - **Stateful Bearer Token Authentication**: Exchange user credentials for temporary bearer tokens (24h lifespan).
  - **Permission-Based RBAC**: Fine-grained endpoint permissions (`movies:read`, `movies:write`) mapped through join tables.

- **Background Tasks & Mailer**:
  - Decoupled asynchronous email delivery (welcome emails, activation tokens, password reset) via background goroutines.
  - Templates compiled into the binary using Go's `//go:embed` filesystem.
  - Built-in retry mechanism with exponential backoff for SMTP delivery.

- **Resilience, Observability & Security**:
  - **Graceful Shutdown**: Intercepts `SIGINT` and `SIGTERM` signals; allows up to 5 seconds for in-flight requests and background goroutines (`sync.WaitGroup`) to complete cleanly.
  - **IP-Based Rate Limiting**: Token-bucket algorithm (`x/time/rate`) with reverse-proxy support (`tomasen/realip`) and background memory cleanup.
  - **CORS Handling**: Supports both simple and preflight (`OPTIONS`) requests with customizable trusted origin safelisting.
  - **SQL Query Timeouts**: Hard 3-second context cancellation on all database transactions.
  - **Structured JSON Logging**: Leveled logger (`INFO`, `ERROR`, `FATAL`) formatted in JSON with stack traces on runtime errors.
  - **Metrics**: Real-time stats exposed at `/debug/vars` (goroutines, database connection pool, request count by HTTP status code).

---

### 📂 Directory Structure

```text
greenlight/
├── bin/                          # Compiled application binaries
├── cmd/
│   ├── api/                      # Application core (handlers, routes, server)
│   │   ├── context.go            # Request context user helpers
│   │   ├── errors.go             # JSON error response helpers
│   │   ├── healthcheck.go        # Health check handler
│   │   ├── helpers.go            # URL/query string parsing & background task helpers
│   │   ├── main.go               # Config flags, DB pool setup, server bootstrap
│   │   ├── middleware.go         # CORS, rate limiting, auth, panic recovery, metrics
│   │   ├── movies.go             # Movie domain handlers
│   │   ├── routes.go             # Routing table with attached middleware
│   │   ├── server.go             # Server initialization & graceful shutdown
│   │   ├── tokens.go             # Token issuance & password reset handlers
│   │   └── users.go              # User registration & activation handlers
│   └── examples/
│       └── cors/                 # Web clients to verify simple and preflight CORS
├── internal/
│   ├── data/                     # Database models & data validation
│   │   ├── filters.go            # Query string filtering, sorting, pagination
│   │   ├── models.go             # Models container (Movies, Users, Tokens, Permissions)
│   │   ├── movies.go             # Movie model implementation
│   │   ├── permissions.go        # Permission model implementation
│   │   ├── runtime.go            # Custom Runtime type ("<runtime> mins")
│   │   ├── tokens.go             # Token generator & model
│   │   └── users.go              # User model with bcrypt hashing
│   ├── jsonlog/                  # Leveled structured JSON logging
│   ├── mailer/                   # SMTP mailer with embedded HTML/Text templates
│   └── validator/                # Validation helper utilities
├── migrations/                   # Sequential SQL migration files (000001 -> 000006)
├── remote/                       # Production setup scripts, systemd unit, Caddyfile
├── Makefile                      # Build, audit, run, and migration recipes
├── .envrc                        # Environment variable configuration
└── go.mod / go.sum               # Go module dependencies
```

---

### 🛠 Prerequisites

- **Go**: Version 1.16 or newer (tested on Go 1.24/1.26)
- **PostgreSQL**: Version 12 or newer (with `citext` extension)
- **golang-migrate**: Database migration CLI tool
- **GNU Make**: Automation runner

---

### 🚀 Getting Started

#### 1. Setup Database
Connect to PostgreSQL as a superuser and create the project database and user:
```sql
CREATE DATABASE greenlight;
\c greenlight
CREATE EXTENSION IF NOT EXISTS citext;
CREATE ROLE greenlight WITH LOGIN PASSWORD 'pa55word';
```

#### 2. Execute SQL Migrations
Apply migrations using `golang-migrate`:
```bash
migrate -path ./migrations -database "postgres://greenlight:pa55word@localhost/greenlight?sslmode=disable" up
```
*(Or run `make db/migrations/up` if configured in `.envrc`)*

#### 3. Run Application
- **Using `go run`:**
  ```bash
  go run ./cmd/api -db-dsn="postgres://greenlight:pa55word@localhost/greenlight?sslmode=disable"
  ```
- **Using `make`:**
  ```bash
  make run/api
  ```

Server will start listening on `http://localhost:4000`.

---

### 📡 API Endpoints

| Method | URL Pattern | Description | Access / Permissions |
| :--- | :--- | :--- | :--- |
| `GET` | `/v1/healthcheck` | View API health and version info | Public |
| `GET` | `/v1/movies` | List movies (filters, search, pagination) | Activated User (`movies:read`) |
| `POST` | `/v1/movies` | Create a new movie | Activated User (`movies:write`) |
| `GET` | `/v1/movies/:id` | View specific movie details | Activated User (`movies:read`) |
| `PATCH` | `/v1/movies/:id` | Partially update a movie | Activated User (`movies:write`) |
| `DELETE` | `/v1/movies/:id` | Delete a movie | Activated User (`movies:write`) |
| `POST` | `/v1/users` | Register a new user | Public |
| `PUT` | `/v1/users/activated` | Activate account via token | Public |
| `PUT` | `/v1/users/password` | Update password via reset token | Public |
| `POST` | `/v1/tokens/authentication` | Generate 24-hour authentication token | Public |
| `POST` | `/v1/tokens/activation` | Resend activation token | Public |
| `POST` | `/v1/tokens/password-reset` | Request password reset token | Public |
| `GET` | `/debug/vars` | Inspect memory, goroutines, and request metrics | Internal / Localhost |

---

<br>

---

<a name="tiếng-việt"></a>
## 🇻🇳 Tiếng Việt

### 📖 Giới thiệu dự án
**Greenlight** là ứng dụng RESTful JSON API hoàn chỉnh được phát triển bằng ngôn ngữ **Go (Golang)** kết hợp cơ sở dữ liệu **PostgreSQL**, tuân thủ các quy chuẩn thiết kế từ cuốn sách *"Let's Go Further!"* của tác giả Alex Edwards.

Dự án cung cấp giải pháp mẫu mực cho hệ thống backend microservice/monolith: quản lý danh mục phim, tìm kiếm toàn văn, phân trang, lọc dữ liệu, hệ thống xác thực người dùng dựa trên token, phân quyền chi tiết (RBAC), kiểm soát tần suất truy cập (rate limiting), gửi email ngầm và cơ chế tắt server mượt mà (graceful shutdown).

---

### ✨ Các tính năng chính

- **Quản lý danh mục phim (Movies CRUD & Querying)**:
  - Thêm, xem chi tiết, cập nhật và xóa phim; lưu trữ mảng thể loại bằng PostgreSQL `text[]`.
  - **Cập nhật từng phần (Partial Update / `PATCH`)**: Dùng trường con trỏ (`*string`, `*int32`, ...) để cập nhật chính xác các trường được gửi lên, không ghi đè dữ liệu còn lại.
  - **Khóa lạc quan (Optimistic Concurrency Control)**: Sử dụng trường `version` ngăn chặn xung đột ghi đồng thời (trả về `409 Conflict` nếu dữ liệu đã bị sửa đổi).
  - **Tìm kiếm toàn văn & Bộ lọc**: Tích hợp PostgreSQL Full-Text Search (`to_tsvector`, `plainto_tsquery`) kết hợp GIN index tăng tốc tối đa.
  - **Phân trang & Sắp xếp an toàn**: Whitelist tham số sắp xếp, tính toán metadata phân trang qua window function `count(*) OVER()`.

- **Xác thực & Phân quyền (Authentication & Authorization)**:
  - **Đăng ký người dùng & Hash mật khẩu**: Sử dụng `bcrypt` với cost parameter an toàn.
  - **Kích hoạt tài khoản (Account Activation)**: Sinh token ngẫu nhiên bảo mật cao (128-bit CSPRNG, base32 không padding, băm SHA-256), gửi email kích hoạt.
  - **Xác thực Token có trạng thái (Stateful Bearer Token)**: Tra đổi email/mật khẩu lấy bearer token (hết hạn sau 24h).
  - **Phân quyền chi tiết (Permission-based RBAC)**: Quản lý quyền theo bảng liên kết (`movies:read`, `movies:write`), tự động cấp quyền đọc khi đăng ký.

- **Xử lý tác vụ ngầm & Gửi Email (Background Tasks & Mailer)**:
  - Gửi email kích hoạt, đổi mật khẩu không chặn tiến trình HTTP (chạy trong background goroutine).
  - Nhúng trực tiếp template email HTML và Text vào file nhị phân bằng `//go:embed`.
  - Cơ chế tự động thử lại (retry) tối đa 3 lần khi gửi email gặp sự cố kết nối.

- **Độ ổn định, Giám sát & Bảo mật**:
  - **Graceful Shutdown**: Bắt tín hiệu `SIGINT` (Ctrl+C) và `SIGTERM`, dừng nhận request mới và sử dụng `sync.WaitGroup` đợi các tác vụ ngầm hoàn tất trước khi thoát hẳn.
  - **Giới hạn tốc độ (Rate Limiting)**: Token-bucket algorithm (`x/time/rate`) theo từng IP client (`tomasen/realip`), tự dọn dẹp bộ nhớ định kỳ.
  - **Hỗ trợ CORS**: Hỗ trợ Simple CORS và Preflight (`OPTIONS`) với danh sách domain tin cậy được cấu hình linh hoạt.
  - **SQL Query Timeout**: Giới hạn thời gian truy vấn 3 giây qua `context.Context` tránh nghẽn kết nối database.
  - **Ghi log có cấu trúc**: JSON Structured Logger đa cấp độ (`INFO`, `ERROR`, `FATAL`) kèm stack trace cho lỗi runtime.
  - **Hệ thống Metrics**: Báo cáo thống kê hiệu năng thời gian thực tại `/debug/vars` (goroutines, database connection pool, số lượng request/response theo mã trạng thái HTTP).

---

### 🛠 Yêu cầu hệ thống

- **Go**: Phiên bản 1.16 trở lên (khuyến nghị Go 1.20+)
- **PostgreSQL**: Phiên bản 12 trở lên (yêu cầu extension `citext`)
- **golang-migrate**: Công cụ dòng lệnh chạy migrations
- **GNU Make**: Tiện ích chạy các lệnh tự động hóa

---

### 🚀 Hướng dẫn cài đặt & Khởi chạy

#### 1. Khởi tạo cơ sở dữ liệu
Đăng nhập PostgreSQL bằng quyền superuser và thực thi:
```sql
CREATE DATABASE greenlight;
\c greenlight
CREATE EXTENSION IF NOT EXISTS citext;
CREATE ROLE greenlight WITH LOGIN PASSWORD 'pa55word';
```

#### 2. Áp dụng SQL Migrations
Chạy migrations từ thư mục [migrations/](file:///c:/Users/Lenovo%20legion%205/Downloads/OCN_materials/Practice/greenlight/migrations):
```bash
migrate -path ./migrations -database "postgres://greenlight:pa55word@localhost/greenlight?sslmode=disable" up
```
*(Hoặc chạy lệnh `make db/migrations/up`)*

#### 3. Cấu hình & Chạy ứng dụng
Cấu hình biến môi trường kết nối trong file [.envrc](file:///c:/Users/Lenovo%20legion%205/Downloads/OCN_materials/Practice/greenlight/.envrc) hoặc export trực tiếp:
```bash
export GREENLIGHT_DB_DSN="postgres://greenlight:pa55word@localhost/greenlight?sslmode=disable"
```

Khởi chạy server:
- **Qua lệnh `go run`:**
  ```bash
  go run ./cmd/api -db-dsn="postgres://greenlight:pa55word@localhost/greenlight?sslmode=disable"
  ```
- **Qua lệnh `make`:**
  ```bash
  make run/api
  ```

API sẽ hoạt động tại địa chỉ: `http://localhost:4000`.

---

### 📡 Danh sách API Endpoints

| Phương thức | Đường dẫn URL | Mô tả | Quyền hạn yêu cầu |
| :--- | :--- | :--- | :--- |
| `GET` | `/v1/healthcheck` | Kiểm tra trạng thái và phiên bản API | Công khai |
| `GET` | `/v1/movies` | Danh sách phim (tìm kiếm, phân trang, lọc) | Đã kích hoạt + `movies:read` |
| `POST` | `/v1/movies` | Thêm phim mới | Đã kích hoạt + `movies:write` |
| `GET` | `/v1/movies/:id` | Xem chi tiết phim | Đã kích hoạt + `movies:read` |
| `PATCH` | `/v1/movies/:id` | Cập nhật thông tin từng phần của phim | Đã kích hoạt + `movies:write` |
| `DELETE` | `/v1/movies/:id` | Xóa phim | Đã kích hoạt + `movies:write` |
| `POST` | `/v1/users` | Đăng ký tài khoản người dùng mới | Công khai |
| `PUT` | `/v1/users/activated` | Kích hoạt tài khoản bằng token | Công khai |
| `PUT` | `/v1/users/password` | Cập nhật mật khẩu bằng token reset | Công khai |
| `POST` | `/v1/tokens/authentication` | Tạo Bearer token để đăng nhập (24h) | Công khai |
| `POST` | `/v1/tokens/activation` | Cấp lại token kích hoạt tài khoản | Công khai |
| `POST` | `/v1/tokens/password-reset` | Yêu cầu gửi email đặt lại mật khẩu | Công khai |
| `GET` | `/debug/vars` | Xem thống kê số liệu runtime & connection pool | Nội bộ |

---

### 💻 Một số lệnh gọi API mẫu

1. **Kiểm tra trạng thái (Healthcheck)**:
   ```bash
   curl -i http://localhost:4000/v1/healthcheck
   ```

2. **Đăng ký tài khoản**:
   ```bash
   curl -i -d '{"name":"Alice","email":"alice@example.com","password":"pa55word"}' http://localhost:4000/v1/users
   ```

3. **Kích hoạt tài khoản bằng token**:
   ```bash
   curl -i -X PUT -d '{"token":"<ACTIVATION_TOKEN>"}' http://localhost:4000/v1/users/activated
   ```

4. **Đăng nhập lấy Bearer token**:
   ```bash
   curl -i -d '{"email":"alice@example.com","password":"pa55word"}' http://localhost:4000/v1/tokens/authentication
   ```

5. **Truy vấn danh sách phim (kèm token và bộ lọc)**:
   ```bash
   curl -i -H "Authorization: Bearer <AUTH_TOKEN>" "http://localhost:4000/v1/movies?title=black&genres=action&page=1&page_size=10&sort=-year"
   ```

---

### 🧰 Quản trị tự động với Makefile

Tệp [Makefile](file:///c:/Users/Lenovo%20legion%205/Downloads/OCN_materials/Practice/greenlight/Makefile) cung cấp đầy đủ các tác vụ tự động hóa:

- `make help`: Hiển thị danh mục các lệnh hỗ trợ kèm hướng dẫn.
- `make run/api`: Khởi chạy ứng dụng với cấu hình lấy từ `.envrc`.
- `make db/psql`: Mở phiên làm việc trực tiếp với PostgreSQL qua `psql`.
- `make db/migrations/new name=<name>`: Tạo cặp file migration up/down mới.
- `make db/migrations/up`: Áp dụng tất cả migrations mới vào cơ sở dữ liệu.
- `make audit`: Kiểm tra chất lượng code (`go mod tidy`, `verify`, `fmt`, `vet`, `staticcheck`, và chạy test với `-race`).
- `make vendor`: Vendor các dependency của bên thứ 3 vào thư mục `vendor/`.
- `make build/api`: Biên dịch mã nguồn thành file thực thi kèm nạp động `buildTime` và phiên bản Git commit.
- `make production/deploy/api`: Tự động deploy nhị phân và chạy migrations lên server Ubuntu.
- `make production/configure/api.service`: Cấu hình file systemd service trên server.
- `make production/configure/caddyfile`: Cấu hình Caddy Reverse Proxy & HTTPS tự động.