# Architecture And Implementation Plan

## 1. Tóm tắt yêu cầu từ SRS

QR Generator - QR Studio là website client-server cho phép tạo, tùy chỉnh, lưu trữ, quản lý và thống kê mã QR. Hệ thống hỗ trợ QR tĩnh và QR động, hai gói Free/Pro, phân quyền USER/ADMIN, payment gateway cho nâng cấp Pro, và redirect service `/q/{short_code}` cho Dynamic QR.

Nhóm chức năng chính:
- Auth và account: đăng ký, đăng nhập, đăng xuất, đổi mật khẩu, cập nhật hồ sơ, khóa tài khoản.
- QR generation: URL, TEXT, WIFI là Must; VCARD là Should; EMAIL/SMS/LOCATION là Should; SOCIAL/PDF/MENU là Could/Pro.
- QR design: màu sắc, kích thước, template; logo, eye style, dot style nâng cao cho Pro.
- QR management: lưu, danh sách, tìm kiếm/lọc, chi tiết, xóa mềm, download, duplicate.
- Dynamic QR: Pro-only, short code, sửa URL đích, redirect, ghi nhận scan.
- Analytics: tổng lượt quét, theo ngày, thiết bị, trình duyệt, hệ điều hành, vị trí tương đối.
- Subscription/payment: Free mặc định, Pro qua payment mock trong MVP, SUCCESS mới kích hoạt.
- Admin: dashboard, quản lý user, QR, plan, payment, template, log.

## 2. Assumptions

- Project nguồn ở ổ `D:` hiện chỉ có `docs`, nên source code được triển khai trong workspace Codex tại `qr-generator/`.
- SRS cho phép MySQL hoặc PostgreSQL, prompt yêu cầu MySQL, vì vậy MVP dùng MySQL.
- AutoMigrate được dùng cho MVP để chạy nhanh; thư mục `migrations/` được giữ để bổ sung migration versioned sau.
- Refresh token là JWT có thời hạn dài. Revocation store có thể bổ sung khi cần quản lý phiên theo thiết bị.
- Static QR không cho sửa `content` sau khi tạo theo BR-08; Dynamic QR cho sửa `destination_url`.
- IP scan được cắt/ẩn danh nhẹ ở service layer nếu cần mở rộng; MVP lưu IP từ request để phục vụ demo analytics theo SRS.

## 3. Kiến trúc đề xuất

```txt
Next.js App Router frontend
        |
        | REST /api/v1, Bearer access token
        v
Golang Gin API
        |
        | GORM
        v
MySQL
```

Backend chia lớp:
- `handler`: parse request, bind validation, trả response chuẩn.
- `service`: business rules Free/Pro/Admin, QR dynamic, payment success, analytics.
- `repository`: truy cập GORM.
- `middleware`: JWT, Admin, CORS, logging.
- `models/shared`: DB schema, enum, response, pagination.

Frontend chia lớp:
- App Router pages cho public/user/admin.
- API client wrapper trong `lib/api.ts`.
- Auth state bằng Zustand.
- Form validation bằng Zod.
- UI components cho QR form/preview/design, analytics chart, layout, empty/error/loading state.

## 4. Database design

Các bảng theo SRS: `users`, `roles`, `user_roles`, `plans`, `subscriptions`, `payments`, `folders`, `qr_codes`, `qr_designs`, `qr_templates`, `qr_scans`, `system_logs`.

Các enum bắt buộc được định nghĩa trong backend:
- `user_status`: ACTIVE, LOCKED, DELETED
- `role_name`: USER, ADMIN
- `plan_name`: FREE, PRO
- `plan_status`: ACTIVE, INACTIVE, DELETED
- `subscription_status`: ACTIVE, EXPIRED, CANCELLED, PENDING
- `payment_status`: PENDING, SUCCESS, FAILED, CANCELLED, REFUNDED
- `payment_method`: VNPAY, MOMO, ZALOPAY, PAYPAL, STRIPE, BANK_TRANSFER
- `qr_type`: URL, TEXT, WIFI, VCARD, EMAIL, SMS, LOCATION, SOCIAL, PDF, MENU
- `qr_status`: ACTIVE, DISABLED, DELETED
- `template_status`: ACTIVE, HIDDEN, DELETED
- `error_correction_level`: L, M, Q, H
- `log_level`: INFO, WARNING, ERROR, SECURITY

## 5. API design

Base path: `/api/v1`

Auth:
- `POST /auth/register`
- `POST /auth/login`
- `POST /auth/refresh`
- `POST /auth/logout`
- `GET /auth/me`
- `POST /auth/change-password`

User:
- `GET /users/profile`
- `PUT /users/profile`
- `GET /users/subscription`
- `GET /users/payments`

QR:
- `POST /qrcodes`
- `GET /qrcodes`
- `GET /qrcodes/:id`
- `PUT /qrcodes/:id`
- `DELETE /qrcodes/:id`
- `POST /qrcodes/:id/duplicate`
- `GET /qrcodes/:id/download`
- `GET /qrcodes/:id/design`
- `PUT /qrcodes/:id/design`
- `GET /q/:shortCode`

Analytics:
- `GET /qrcodes/:id/analytics/summary`
- `GET /qrcodes/:id/analytics/by-date`
- `GET /qrcodes/:id/analytics/by-device`
- `GET /qrcodes/:id/analytics/by-browser`
- `GET /qrcodes/:id/analytics/by-location`

Folders, plans, payments, admin APIs được triển khai theo SRS và prompt.

## 6. Implementation plan

Phase 1: Setup monorepo, Docker Compose, env, README, docs.

Phase 2: Backend foundation: config, DB, models, seed roles/plans/admin, response helpers, JWT, middleware.

Phase 3: Auth/user/profile: register, login, refresh, me, change password, profile, subscription, payments.

Phase 4: QR core: create/list/detail/update/delete/duplicate/download PNG, validation, Free limits.

Phase 5: Dynamic QR: Pro check, short code, redirect, scan log, scan count.

Phase 6: Subscription/payment Sepay: create pending Pro payment with Sepay transfer content, accept authenticated Sepay webhook, validate amount/status, activate or extend Pro only for SUCCESS. Local mock success is dev-only and disabled from the normal MVP flow.

Phase 7: Analytics/admin: summary/grouped stats and admin management endpoints.

Phase 8: Frontend: public pages, auth, dashboard, QR CRUD UI, pricing/payment, analytics charts, admin pages, responsive polish.
