# Greenlight API

Greenlight là một dịch vụ RESTful JSON API hoàn chỉnh được phát triển bằng ngôn ngữ **Go (Golang)** và cơ sở dữ liệu **PostgreSQL**, dựa trên kiến trúc và best practices từ cuốn sách *"Let's Go Further!"* của Alex Edwards.

API cung cấp đầy đủ các chức năng quản lý danh mục phim, tìm kiếm, phân trang, lọc dữ liệu nâng cao, cùng hệ thống xác thực người dùng dựa trên token, phân quyền chi tiết, kiểm soát tần suất truy cập (rate limiting), gửi email ngầm và cơ chế tắt server an toàn (graceful shutdown).

---

## 📌 Mục lục
- [Tính năng nổi bật](#-tính-năng-nổi-bật)
- [Cấu trúc thư mục](#-cấu-trúc-thư-mục)
- [Yêu cầu hệ thống](#-yêu-cầu-hệ-thống)
- [Cấu hình cơ sở dữ liệu](#-cấu-hình-cơ-sở-dữ-liệu)
- [Hướng dẫn chạy ứng dụng](#-hướng-dẫn-chạy-ứng-dụng)
- [Danh sách API Endpoints](#-danh-sách-api-endpoints)
- [Tự động hóa với Makefile](#-tự-động-hóa-với-makefile)

---

## 🌟 Tính năng nổi bật

1. **Quản lý danh mục phim (Movies CRUD & Querying)**:
   - Thêm, xem chi tiết, cập nhật và xóa phim.
   - **Cập nhật từng phần (Partial Update / PATCH)**: Sử dụng kiểu con trỏ (`*string`, `*int32`, ...) để phân biệt giá trị mặc định và giá trị không được cung cấp.
   - **Khóa lạc quan (Optimistic Concurrency Control)**: Sử dụng trường `version` ngăn chặn xung đột dữ liệu khi nhiều tiến trình cập nhật cùng thời điểm (trả về mã `409 Conflict`).
   - **Tìm kiếm & Lọc dữ liệu**: Tích hợp PostgreSQL Full-Text Search (`to_tsvector`, `plainto_tsquery`) kết hợp GIN index tăng tốc truy vấn.
   - **Sắp xếp & Phân trang an toàn**: Safelist các trường được phép sắp xếp, phân trang động kèm metadata tính bằng window function `count(*) OVER()`.

2. **Xác thực & Phân quyền (Authentication & Authorization)**:
   - **Đăng ký tài khoản & Hash mật khẩu**: Sử dụng `bcrypt` với cost parameter an toàn.
   - **Xác thực tài khoản (Account Activation)**: Sinh token ngẫu nhiên bảo mật cao (128-bit CSPRNG, base32 không padding, băm SHA-256), gửi email kích hoạt.
   - **Stateful Bearer Token Authentication**: Tra đổi email/mật khẩu lấy bearer token (hết hạn sau 24h).
   - **Phân quyền dựa trên quyền hạn (Permission-based RBAC)**: Hỗ trợ quyền chi tiết (`movies:read`, `movies:write`), tự động gán quyền đọc khi đăng ký tài khoản.

3. **Xử lý tác vụ ngầm & Gửi Email (Background Tasks & Mailer)**:
   - Gửi email kích hoạt và email đặt lại mật khẩu không chặn tiến trình HTTP (chạy trong background goroutine).
   - Nhúng template email HTML và Text trực tiếp vào nhị phân bằng `//go:embed`.
   - Cơ chế tự động thử lại (retry) khi gửi email thất bại.

4. **Độ ổn định & Giám sát (Resilience & Observability)**:
   - **Graceful Shutdown**: Bắt các tín hiệu `SIGINT` (Ctrl+C), `SIGTERM`, dừng nhận request mới và sử dụng `sync.WaitGroup` đợi các tác vụ ngầm hoàn tất trước khi tắt server.
   - **SQL Query Timeout**: Giới hạn thời gian truy vấn cơ sở dữ liệu qua `context.Context` (3 giây) tránh treo kết nối.
   - **Rate Limiting**: Giới hạn tần suất gọi API theo từng địa chỉ IP với `golang.org/x/time/rate` và `tomasen/realip`, tự động dọn dẹp bộ nhớ định kỳ.
   - **Bảo vệ CORS**: Hỗ trợ Simple CORS và Preflight CORS (`OPTIONS`) với cấu hình danh sách domain tin cậy (`-cors-trusted-origins`).
   - **Ghi log có cấu trúc**: JSON Structured Logger tự động định dạng log, phân cấp (`INFO`, `ERROR`, `FATAL`) kèm stack trace cho lỗi runtime.
   - **Hệ thống Metrics**: Báo cáo chỉ số thời gian thực (goroutine count, database pool stats, request/response count theo mã trạng thái HTTP) tại `/debug/vars`.

---

## 📂 Cấu trúc thư mục

```text
greenlight/
├── bin/                          # Thư mục chứa file thực thi sau khi biên dịch
├── cmd/
│   ├── api/                      # Ứng dụng chính Greenlight API
│   │   ├── context.go            # Tiện ích đọc/ghi User vào request context
│   │   ├── errors.go             # Helper phản hồi lỗi định dạng JSON
│   │   ├── healthcheck.go        # Handler kiểm tra trạng thái dịch vụ
│   │   ├── helpers.go            # Helper đọc tham số URL, query string, JSON, tác vụ ngầm
│   │   ├── main.go               # Khởi tạo cờ cấu hình, DB pool, metrics, chạy server
│   │   ├── middleware.go         # Pipeline: CORS, RateLimit, Auth, Panic recovery, Metrics
│   │   ├── movies.go             # Handlers cho nghiệp vụ Movies
│   │   ├── routes.go             # Thiết lập HTTP router và gắn middleware
│   │   ├── server.go             # Khởi chạy server và cơ chế Graceful Shutdown
│   │   ├── tokens.go             # Handlers cấp token xác thực/kích hoạt/đổi mật khẩu
│   │   └── users.go              # Handlers đăng ký, kích hoạt, cập nhật mật khẩu user
│   └── examples/
│       └── cors/                 # Web client kiểm thử CORS (simple & preflight)
├── internal/
│   ├── data/                     # Data Access Layer & Models
│   │   ├── filters.go            # Phân trang, sắp xếp và tính toán metadata
│   │   ├── models.go             # Container tập hợp các Models
│   │   ├── movies.go             # Movie struct & MovieModel CRUD
│   │   ├── permissions.go        # PermissionModel & kiểm tra quyền
│   │   ├── runtime.go            # Custom Runtime type ("<runtime> mins")
│   │   ├── tokens.go             # Token struct & TokenModel
│   │   └── users.go              # User struct & UserModel
│   ├── jsonlog/                  # Ghi log có cấu trúc dạng JSON
│   ├── mailer/                   # Module gửi email và templates HTML/Text nhúng
│   │   ├── templates/            # Các mẫu email (.tmpl)
│   │   └── mailer.go             # Logic kết nối SMTP và gửi thư
│   └── validator/                # Tiện ích kiểm tra tính hợp lệ dữ liệu
├── migrations/                   # Các tệp SQL migration (000001 -> 000006)
├── remote/                       # Cấu hình máy chủ production
│   ├── production/
│   │   ├── api.service           # File unit systemd chạy background service
│   │   └── Caddyfile             # File cấu hình Caddy Reverse Proxy & tự động HTTPS
│   └── setup/
│       └── 01.sh                 # Script tự động thiết lập VPS Ubuntu
├── Makefile                      # Tập hợp lệnh tự động hóa (build, test, migrate, run)
├── .envrc                        # Mẫu cấu hình biến môi trường
└── go.mod / go.sum               # Quản lý phụ thuộc Go modules
```

---

## 🛠 Yêu cầu hệ thống

- **Go**: Phiên bản 1.16 trở lên (khuyến nghị Go 1.20+)
- **PostgreSQL**: Phiên bản 12 trở lên (yêu cầu extension `citext`)
- **golang-migrate**: Công cụ dòng lệnh chạy database migrations
- **Make**: (Tùy chọn) Tiện ích chạy lệnh tự động hóa

---

## 🗄 Cấu hình cơ sở dữ liệu

1. **Khởi tạo Database & User trong PostgreSQL**:
   Đăng nhập vào PostgreSQL (ví dụ qua `psql`) với quyền superuser:
   ```sql
   CREATE DATABASE greenlight;
   \c greenlight
   CREATE EXTENSION IF NOT EXISTS citext;
   CREATE ROLE greenlight WITH LOGIN PASSWORD 'pa55word';
   ```

2. **Chạy các file SQL Migrations**:
   Sử dụng công cụ `migrate`:
   ```bash
   migrate -path ./migrations -database "postgres://greenlight:pa55word@localhost/greenlight?sslmode=disable" up
   ```
   *(Hoặc sử dụng `make db/migrations/up` nếu đã cấu hình `GREENLIGHT_DB_DSN` trong `.envrc`)*

---

## 🚀 Hướng dẫn chạy ứng dụng

### 1. Cấu hình biến môi trường
Thiết lập chuỗi kết nối cơ sở dữ liệu trong file `.envrc` hoặc export trực tiếp:
```bash
export GREENLIGHT_DB_DSN="postgres://greenlight:pa55word@localhost/greenlight?sslmode=disable"
```

### 2. Khởi chạy Server
- **Cách 1: Khởi chạy trực tiếp qua `go run`**
  ```bash
  go run ./cmd/api -db-dsn="postgres://greenlight:pa55word@localhost/greenlight?sslmode=disable"
  ```

- **Cách 2: Khởi chạy bằng `Makefile`**
  ```bash
  make run/api
  ```

Mặc định server sẽ lắng nghe tại cổng `http://localhost:4000`.

### 3. Xem danh sách các cờ cấu hình (Flags)
```bash
go run ./cmd/api -help
```
Các cờ hỗ trợ chính:
- `-port`: Cổng mạng server lắng nghe (mặc định: `4000`).
- `-env`: Môi trường (`development`, `staging`, `production`).
- `-db-dsn`: Chuỗi kết nối PostgreSQL.
- `-db-max-open-conns`, `-db-max-idle-conns`, `-db-max-idle-time`: Cấu hình connection pool.
- `-limiter-enabled`, `-limiter-rps`, `-limiter-burst`: Cấu hình giới hạn tần suất request.
- `-smtp-host`, `-smtp-port`, `-smtp-username`, `-smtp-password`, `-smtp-sender`: Cấu hình máy chủ gửi thư.
- `-cors-trusted-origins`: Danh sách domain tin cậy cho phép truy cập cross-origin.
- `-version`: Hiển thị phiên bản nhị phân và thời gian biên dịch rồi thoát.

---

## 📡 Danh sách API Endpoints

| Phương thức | Đường dẫn URL | Mô tả | Yêu cầu xác thực / quyền |
| :--- | :--- | :--- | :--- |
| `GET` | `/v1/healthcheck` | Kiểm tra trạng thái và phiên bản của API | Không |
| `GET` | `/v1/movies` | Danh sách phim (lọc theo `title`, `genres`, phân trang `page`, `page_size`, sắp xếp `sort`) | User đã kích hoạt + quyền `movies:read` |
| `POST` | `/v1/movies` | Thêm phim mới vào hệ thống | User đã kích hoạt + quyền `movies:write` |
| `GET` | `/v1/movies/:id` | Xem thông tin chi tiết một bộ phim | User đã kích hoạt + quyền `movies:read` |
| `PATCH` | `/v1/movies/:id` | Cập nhật thông tin từng phần của bộ phim | User đã kích hoạt + quyền `movies:write` |
| `DELETE` | `/v1/movies/:id` | Xóa bộ phim khỏi hệ thống | User đã kích hoạt + quyền `movies:write` |
| `POST` | `/v1/users` | Đăng ký tài khoản người dùng mới | Không |
| `PUT` | `/v1/users/activated` | Kích hoạt tài khoản bằng token | Không |
| `PUT` | `/v1/users/password` | Cập nhật mật khẩu mới qua token reset | Không |
| `POST` | `/v1/tokens/authentication` | Tạo Bearer token để đăng nhập | Không |
| `POST` | `/v1/tokens/activation` | Cấp lại token kích hoạt tài khoản | Không |
| `POST` | `/v1/tokens/password-reset` | Yêu cầu gửi token đặt lại mật khẩu | Không |
| `GET` | `/debug/vars` | Xem thống kê số liệu runtime & connection pool | Nội bộ |

### Ví dụ một số yêu cầu mẫu

1. **Kiểm tra Healthcheck**:
   ```bash
   curl -i http://localhost:4000/v1/healthcheck
   ```

2. **Đăng ký tài khoản mới**:
   ```bash
   curl -i -d '{"name":"Alice","email":"alice@example.com","password":"pa55word"}' http://localhost:4000/v1/users
   ```

3. **Kích hoạt tài khoản**:
   ```bash
   curl -i -X PUT -d '{"token":"<ACTIVATION_TOKEN>"}' http://localhost:4000/v1/users/activated
   ```

4. **Đăng nhập lấy Authentication Token**:
   ```bash
   curl -i -d '{"email":"alice@example.com","password":"pa55word"}' http://localhost:4000/v1/tokens/authentication
   ```

5. **Truy vấn danh sách phim (kèm Bearer Token)**:
   ```bash
   curl -i -H "Authorization: Bearer <AUTH_TOKEN>" "http://localhost:4000/v1/movies?title=black&genres=action&page=1&page_size=10&sort=-year"
   ```

---

## 🧰 Tự động hóa với Makefile

Dự án cung cấp sẵn tệp [Makefile](file:///c:/Users/Lenovo%20legion%205/Downloads/OCN_materials/Practice/greenlight/Makefile) hỗ trợ các tác vụ quản trị:

- `make help`: Hiển thị hướng dẫn tất cả các lệnh có sẵn.
- `make run/api`: Khởi chạy ứng dụng với biến môi trường từ `.envrc`.
- `make db/psql`: Mở terminal kết nối trực tiếp đến database PostgreSQL qua `psql`.
- `make db/migrations/new name=<name>`: Tạo cặp file migration mới trong thư mục `migrations/`.
- `make db/migrations/up`: Áp dụng tất cả các migrations mới (có bước xác nhận `Are you sure? [y/N]`).
- `make audit`: Kiểm tra chất lượng mã nguồn (chạy `go mod tidy`, `verify`, `fmt`, `vet`, `staticcheck`, và chạy test có cờ `-race`).
- `make vendor`: Tải và lưu trữ mã nguồn của thư viện bên thứ 3 vào thư mục `vendor/`.
- `make build/api`: Biên dịch ứng dụng thành file nhị phân cho hệ điều hành hiện tại và bản build cross-compile cho Linux (`bin/linux_amd64/api`) tự động nhúng `buildTime` và phiên bản từ Git.
- `make production/deploy/api`: Triển khai nhị phân và chạy migrations tự động lên máy chủ production.
- `make production/configure/api.service`: Cấu hình và kích hoạt systemd service trên server.
- `make production/configure/caddyfile`: Cấu hình Caddy reverse proxy và reload dịch vụ.