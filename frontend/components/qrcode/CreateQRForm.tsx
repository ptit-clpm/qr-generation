"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { Download, Save, Trash2, Upload } from "lucide-react";
import { api, messageFromError } from "@/lib/api";
import { backendUrl, qrTypes } from "@/lib/constants";
import { Button } from "@/components/common/Button";
import { QRPreview } from "@/components/qrcode/QRPreview";
import { QRTypeFields } from "@/components/qrcode/QRTypeFields";
import { useAuthStore } from "@/stores/auth";
import type { ApiEnvelope, QRCode, QRType, Subscription } from "@/types";

interface FolderItem {
  id: number;
  name: string;
}

interface CreateQRFormProps {
  onSaved?: (qr: QRCode) => void;
}

export function CreateQRForm({ onSaved }: CreateQRFormProps) {
  const { user } = useAuthStore();
  const [hasProSubscription, setHasProSubscription] = useState(false);
  const isProUser = user?.roles?.some((r) => r.name === "ADMIN") || hasProSubscription;

  const [title, setTitle] = useState("Campaign QR");
  const [qrType, setQrType] = useState<QRType>("URL");
  const [content, setContent] = useState("https://example.com");
  const [isDynamic, setDynamic] = useState(false);
  const [destinationUrl, setDestinationUrl] = useState("https://example.com");
  const [folderId, setFolderId] = useState<number | null>(null);
  const [folders, setFolders] = useState<FolderItem[]>([]);
  const [foreground, setForeground] = useState("#111827");
  const [background, setBackground] = useState("#FFFFFF");
  const [created, setCreated] = useState<QRCode | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [logoFile, setLogoFile] = useState<File | null>(null);
  const [logoPreview, setLogoPreview] = useState("");
  const [logoUploading, setLogoUploading] = useState(false);

  useEffect(() => {
    if (!user) {
      setHasProSubscription(false);
      return;
    }

    let cancelled = false;

    api
      .get<ApiEnvelope<Subscription>>("/users/subscription")
      .then((res) => {
        if (cancelled) return;
        setHasProSubscription(res.data.data?.plan?.name === "PRO");
      })
      .catch(() => {
        if (cancelled) return;
        setHasProSubscription(false);
      });

    return () => {
      cancelled = true;
    };
  }, [user]);

  // Fetch folders if user is logged in
  useEffect(() => {
    if (user) {
      api
        .get<ApiEnvelope<FolderItem[]>>("/folders")
        .then((res) => setFolders(res.data.data ?? []))
        .catch(() => {});
    }
  }, [user]);

  const handleContentChange = useCallback((newContent: string) => {
    setContent(newContent);
  }, []);

  const previewValue = useMemo(() => {
    if (created?.is_dynamic && created.short_code) return `${backendUrl}/q/${created.short_code}`;
    return qrType === "URL" && isDynamic ? `${backendUrl}/q/preview` : content;
  }, [created, content, isDynamic, qrType]);

  async function submit() {
    setLoading(true);
    setError("");
    try {
      const isDynamicUrl = qrType === "URL" && isDynamic;
      const res = await api.post<ApiEnvelope<QRCode>>("/qrcodes", {
        title,
        qr_type: qrType,
        content,
        is_dynamic: isDynamicUrl,
        destination_url: isDynamicUrl ? destinationUrl : undefined,
        folder_id: folderId,
        design: { foreground_color: foreground, background_color: background, size: 512, error_correction_level: "M" }
      });
      const newQr = res.data.data ?? null;
      let savedQr = newQr;
      if (newQr && logoFile) {
        setLogoUploading(true);
        const form = new FormData();
        form.append("file", logoFile);
        const logoRes = await api.post<ApiEnvelope<QRCode["design"]>>(`/qrcodes/${newQr.id}/logo`, form);
        savedQr = { ...newQr, design: logoRes.data.data ?? newQr.design };
        setLogoFile(null);
        setLogoPreview("");
        setLogoUploading(false);
      }
      setCreated(savedQr);
      if (savedQr) {
        onSaved?.(savedQr);
      }
    } catch (err) {
      setLogoUploading(false);
      setError(messageFromError(err));
    } finally {
      setLoading(false);
    }
  }

  async function download() {
    if (!created) return;
    const res = await api.get(`/qrcodes/${created.id}/download`, { responseType: "blob" });
    const url = URL.createObjectURL(res.data);
    const link = document.createElement("a");
    link.href = url;
    link.download = "qr-code.png";
    link.click();
    URL.revokeObjectURL(url);
  }

  async function deleteLogo() {
    if (!created?.design?.logo_url) return;
    try {
      const res = await api.delete<ApiEnvelope<QRCode["design"]>>(`/qrcodes/${created.id}/logo`);
      setCreated({ ...created, design: res.data.data ?? { ...created.design, logo_url: "" } });
    } catch (err) {
      setError(messageFromError(err));
    }
  }

  const isLocalhostBackend = backendUrl.includes("localhost") || backendUrl.includes("127.0.0.1");

  return (
    <div className="grid gap-6 lg:grid-cols-[1fr_420px]">
      <section className="rounded-md border border-slate-200 bg-white p-5 shadow-soft">
        <div className="grid gap-4 md:grid-cols-2">
          <label className="text-sm font-medium text-ink">
            Title
            <input value={title} onChange={(e) => setTitle(e.target.value)} className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2" />
          </label>

          <div className="md:col-span-2 rounded-md border border-slate-200 p-4">
            <div className="flex items-center justify-between gap-3">
              <div>
                <p className="text-sm font-semibold text-ink">Logo / Icon</p>
                <p className="mt-1 text-xs text-muted">PNG, JPEG hoặc WebP · tối đa 5 MB · chỉ tài khoản Pro</p>
              </div>
              {!isProUser ? <span className="rounded bg-amber-500/10 px-2 py-1 text-xs text-amber-800">Cần nâng cấp Pro</span> : null}
            </div>
            {isProUser ? (
              <div className="mt-3 flex flex-wrap items-center gap-3">
                <label className="inline-flex cursor-pointer items-center gap-2 rounded-md border border-slate-300 px-3 py-2 text-sm">
                  <Upload className="h-4 w-4" /> Chọn logo
                  <input type="file" accept="image/png,image/jpeg,image/webp" className="hidden" onChange={(e) => {
                    const file = e.target.files?.[0];
                    if (!file) return;
                    if (!["image/png", "image/jpeg", "image/webp"].includes(file.type)) { setError("Logo phải là PNG, JPEG hoặc WebP"); return; }
                    if (file.size > 5 * 1024 * 1024) { setError("Logo không được vượt quá 5 MB"); return; }
                    setLogoFile(file);
                    setLogoPreview(URL.createObjectURL(file));
                    setError("");
                  }} />
                </label>
                {logoFile ? <span className="text-xs text-muted">{logoFile.name} · {(logoFile.size / 1024).toFixed(0)} KB</span> : null}
                {logoFile ? <button type="button" className="text-coral" onClick={() => { setLogoFile(null); setLogoPreview(""); }}><Trash2 className="h-4 w-4" /></button> : null}
                {logoPreview ? <img src={logoPreview} alt="Logo preview" className="h-12 w-12 rounded border object-contain" /> : null}
              </div>
            ) : <p className="mt-3 text-xs text-muted">Logo bị khóa trên gói Free.</p>}
          </div>
          <label className="text-sm font-medium text-ink">
            Type
            <select
              value={qrType}
              onChange={(e) => {
                const nextType = e.target.value as QRType;
                setQrType(nextType);
                if (nextType !== "URL") {
                  setDynamic(false);
                  setDestinationUrl("");
                }
              }}
              className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
            >
              {qrTypes.map((type) => (
                <option key={type.value} value={type.value}>
                  {type.label}
                  {type.pro ? " (Pro)" : ""}
                </option>
              ))}
            </select>
          </label>

          {/* Dynamic field layout according to qrType */}
          <QRTypeFields qrType={qrType} onChangeContent={handleContentChange} isProUser={isProUser} />

          {/* Folder selection if logged in */}
          {user ? (
            <label className="md:col-span-2 text-sm font-medium text-ink">
              Folder (Thư mục)
              <select
                value={folderId ?? ""}
                onChange={(e) => setFolderId(e.target.value ? Number(e.target.value) : null)}
                className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
              >
                <option value="">-- Uncategorized (Không chọn thư mục) --</option>
                {folders.map((f) => (
                  <option key={f.id} value={f.id}>
                    {f.name}
                  </option>
                ))}
              </select>
            </label>
          ) : null}

          {qrType === "URL" ? (
            <>
              <label className="flex items-center gap-2 text-sm font-medium text-ink md:col-span-2">
                <input type="checkbox" checked={isDynamic} onChange={(e) => setDynamic(e.target.checked)} />
                Dynamic QR (Cho phép đếm scan và đổi URL đích về sau)
              </label>

              {isDynamic ? (
                <label className="md:col-span-2 text-sm font-medium text-ink">
                  Destination URL (URL đích khi quét mã Dynamic)
                  <input
                    value={destinationUrl}
                    onChange={(e) => setDestinationUrl(e.target.value)}
                    placeholder="https://example.com"
                    className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
                  />
                </label>
              ) : null}
            </>
          ) : null}

          <label className="text-sm font-medium text-ink">
            Foreground Color
            <input type="color" value={foreground} onChange={(e) => setForeground(e.target.value)} className="mt-1 h-10 w-full rounded-md border border-slate-200" />
          </label>
          <label className="text-sm font-medium text-ink">
            Background Color
            <input type="color" value={background} onChange={(e) => setBackground(e.target.value)} className="mt-1 h-10 w-full rounded-md border border-slate-200" />
          </label>
        </div>

        {error ? <p className="mt-4 rounded-md bg-coral/10 px-3 py-2 text-sm text-coral">{error}</p> : null}
        <div className="mt-5 flex flex-wrap gap-3">
          <Button onClick={submit} disabled={loading || logoUploading}>
            <Save className="h-4 w-4" />
            {logoUploading ? "Uploading logo" : loading ? "Saving" : "Save QR"}
          </Button>
          {created ? (
            <Button tone="secondary" onClick={download}>
              <Download className="h-4 w-4" />
              Download PNG
            </Button>
          ) : null}
        </div>
      </section>

      <section className="rounded-md border border-slate-200 bg-white p-5 shadow-soft">
        <h2 className="text-lg font-semibold text-ink">Preview</h2>
        <div className="mt-4 flex justify-center">
          <QRPreview value={previewValue} foreground={foreground} background={background} logoUrl={logoPreview || created?.design?.logo_url} />
        </div>
        {created?.design?.logo_url ? <button type="button" onClick={deleteLogo} className="mt-3 inline-flex items-center gap-2 text-sm text-coral"><Trash2 className="h-4 w-4" /> Xóa logo đã lưu</button> : null}

        {isDynamic && isLocalhostBackend ? (
          <div className="mt-4 rounded-md bg-sky-500/10 border border-sky-500/20 p-3 text-xs text-sky-800">
            <p className="font-semibold">💡 Lưu ý test quét bằng điện thoại:</p>
            <p className="mt-1">
              Mã Dynamic QR mã hóa đường dẫn <code className="font-mono">{backendUrl}/q/...</code>.
              Vì URL hiện tại là <code className="font-mono">localhost</code>, điện thoại ngoài sẽ không mở được trừ khi bạn cấu hình <code className="font-mono">NEXT_PUBLIC_BACKEND_URL</code> thành IP LAN (ví dụ <code className="font-mono">http://192.168.x.x:8080</code>) hoặc Ngrok tunnel.
            </p>
          </div>
        ) : null}
      </section>
    </div>
  );
}
