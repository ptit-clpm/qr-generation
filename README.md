# QR Generator - QR Studio

Website tao va quan ly ma QR theo SRS: Next.js frontend, Golang REST API, MySQL, JWT access/refresh token, Free/Pro plan, Dynamic QR, analytics, Sepay payment webhook va admin.

## Chay bang Docker

```bash
docker compose up -d --build
```

- Frontend: http://localhost:3000
- Backend: http://localhost:8080
- Health check: http://localhost:8080/health

Tai khoan admin seed:

```txt
Email: admin@qr.local
Password: Admin@123456
```

## Chay rieng tung phan

Backend:

```bash
cd backend
cp .env.example .env
go mod download
go run ./cmd/server
```

Frontend:

```bash
cd frontend
cp .env.example .env.local
npm install
npm run dev
```

## Sepay payment

Payment Pro dung Sepay bank transfer. Backend tao payment `PENDING` voi noi dung chuyen khoan `QRPRO-*`; webhook hop le toi `POST /api/v1/payments/sepay/webhook` moi chuyen payment sang `SUCCESS` va kich hoat/gia han Pro.

Cau hinh Sepay trong `backend/.env`:

```env
SEPAY_ENABLED=true
SEPAY_WEBHOOK_SECRET=
SEPAY_BANK_ACCOUNT=
SEPAY_BANK_NAME=
SEPAY_ACCOUNT_NAME=
SEPAY_PAYMENT_PREFIX=QRPRO
SEPAY_RETURN_URL=http://localhost:3000/pricing
SEPAY_API_BASE_URL=
SEPAY_API_TOKEN=
```

Route `/payments/mock-success` chi duoc dang ky khi `APP_ENV=development` va `SEPAY_ENABLED=false`; frontend khong su dung route mock nay.

## Ghi chu trien khai

- SRS duoc doc tu `docs/SRS_QR_Generator.md`.
- Database dung MySQL theo yeu cau cong nghe chinh. SRS cung cho phep PostgreSQL/MySQL, va backend dung GORM nen co the doi dialect sau.
- MVP dung AutoMigrate + seed de chay nhanh trong moi truong hoc tap/demo. Production nen chuyen sang migration versioned.

Chi tiet kien truc nam trong [docs/architecture-plan.md](docs/architecture-plan.md).
