# Cloudinary logo upload

QR logo files are uploaded by the authenticated backend to Cloudinary. MySQL stores only the HTTPS URL, public ID, dimensions, format, and byte count; Render's ephemeral filesystem and MySQL are never used for binary storage.

## Cloudinary setup

1. Create a Cloudinary account at [cloudinary.com](https://cloudinary.com/).
2. In the Cloudinary console, copy the Cloud Name, API Key, and API Secret.
3. Put the values only in `backend/.env` locally or in Render's backend environment variables:

```env
CLOUDINARY_ENABLED=true
CLOUDINARY_CLOUD_NAME=your_cloud_name
CLOUDINARY_API_KEY=your_api_key
CLOUDINARY_API_SECRET=your_api_secret
CLOUDINARY_FOLDER=qr-generator/logos
CLOUDINARY_MAX_UPLOAD_BYTES=5242880
```

Do not add a `NEXT_PUBLIC_CLOUDINARY_API_SECRET` variable. The secret is read only by Go. Do not commit `.env` or print credentials in logs.

## Local and Render

Copy `backend/.env.example` to `backend/.env` for local development, fill in the values, and start the backend normally. On Render, add the same variables in the backend service's Environment settings and redeploy. Keep `CLOUDINARY_ENABLED=false` when Cloudinary is not configured; upload requests then return a clear 503 response.

## API flow

1. Create a QR with `POST /api/v1/qrcodes`.
2. Upload `file` with `POST /api/v1/qrcodes/:id/logo` as multipart/form-data.
3. Remove it with `DELETE /api/v1/qrcodes/:id/logo`.

Only an owned QR and an active Pro subscription can upload. Replacing a logo uploads the new asset, saves its metadata, then destroys the old public ID. If database saving fails, the new asset is cleaned up. If Cloudinary deletion fails, the API reports the failure instead of claiming a fully successful deletion.

The download endpoint renders a high-error-correction QR, fetches only HTTPS assets hosted at `res.cloudinary.com`, resizes the logo to at most 18% of the QR width, places it on a white background in the center, and returns PNG. The same logo URL is used by the browser preview.

## Common errors

- `403`: the account is not on an active Pro plan.
- `503`: Cloudinary is disabled or missing backend credentials.
- `400`: missing file, unsupported MIME type, invalid image bytes, excessive size, or dimensions above 4096×4096.
- `502`: Cloudinary upload/delete failed; retry after checking Cloudinary credentials and availability.

Never fetch user-supplied arbitrary URLs for rendering. The renderer validates HTTPS and the Cloudinary hostname, and uses the default TLS verifier.
