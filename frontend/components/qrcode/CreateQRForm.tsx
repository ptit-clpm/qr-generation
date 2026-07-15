"use client";

import { useMemo, useState } from "react";
import { Download, Save } from "lucide-react";
import { api, messageFromError } from "@/lib/api";
import { backendUrl, qrTypes } from "@/lib/constants";
import { Button } from "@/components/common/Button";
import { QRPreview } from "@/components/qrcode/QRPreview";
import type { ApiEnvelope, QRCode, QRType } from "@/types";

export function CreateQRForm() {
  const [title, setTitle] = useState("Campaign QR");
  const [qrType, setQrType] = useState<QRType>("URL");
  const [content, setContent] = useState("https://example.com");
  const [isDynamic, setDynamic] = useState(false);
  const [destinationUrl, setDestinationUrl] = useState("https://example.com");
  const [foreground, setForeground] = useState("#111827");
  const [background, setBackground] = useState("#FFFFFF");
  const [created, setCreated] = useState<QRCode | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const previewValue = useMemo(() => {
    if (created?.is_dynamic && created.short_code) return `${backendUrl}/q/${created.short_code}`;
    return isDynamic ? `${backendUrl}/q/preview` : content;
  }, [created, content, isDynamic]);

  async function submit() {
    setLoading(true);
    setError("");
    try {
      const res = await api.post<ApiEnvelope<QRCode>>("/qrcodes", {
        title,
        qr_type: qrType,
        content,
        is_dynamic: isDynamic,
        destination_url: destinationUrl,
        design: { foreground_color: foreground, background_color: background, size: 512, error_correction_level: "M" }
      });
      setCreated(res.data.data ?? null);
    } catch (err) {
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

  return (
    <div className="grid gap-6 lg:grid-cols-[1fr_420px]">
      <section className="rounded-md border border-slate-200 bg-white p-5 shadow-soft">
        <div className="grid gap-4 md:grid-cols-2">
          <label className="text-sm font-medium text-ink">
            Title
            <input value={title} onChange={(e) => setTitle(e.target.value)} className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2" />
          </label>
          <label className="text-sm font-medium text-ink">
            Type
            <select value={qrType} onChange={(e) => setQrType(e.target.value as QRType)} className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2">
              {qrTypes.map((type) => <option key={type.value} value={type.value}>{type.label}{type.pro ? " Pro" : ""}</option>)}
            </select>
          </label>
          <label className="md:col-span-2 text-sm font-medium text-ink">
            Content
            <textarea value={content} onChange={(e) => setContent(e.target.value)} rows={5} className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2" />
          </label>
          <label className="flex items-center gap-2 text-sm font-medium text-ink">
            <input type="checkbox" checked={isDynamic} onChange={(e) => setDynamic(e.target.checked)} />
            Dynamic QR
          </label>
          {isDynamic ? (
            <label className="md:col-span-2 text-sm font-medium text-ink">
              Destination URL
              <input value={destinationUrl} onChange={(e) => setDestinationUrl(e.target.value)} className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2" />
            </label>
          ) : null}
          <label className="text-sm font-medium text-ink">
            Foreground
            <input type="color" value={foreground} onChange={(e) => setForeground(e.target.value)} className="mt-1 h-10 w-full rounded-md border border-slate-200" />
          </label>
          <label className="text-sm font-medium text-ink">
            Background
            <input type="color" value={background} onChange={(e) => setBackground(e.target.value)} className="mt-1 h-10 w-full rounded-md border border-slate-200" />
          </label>
        </div>
        {error ? <p className="mt-4 rounded-md bg-coral/10 px-3 py-2 text-sm text-coral">{error}</p> : null}
        <div className="mt-5 flex flex-wrap gap-3">
          <Button onClick={submit} disabled={loading}><Save className="h-4 w-4" />{loading ? "Saving" : "Save QR"}</Button>
          {created ? <Button tone="secondary" onClick={download}><Download className="h-4 w-4" />Download</Button> : null}
        </div>
      </section>
      <section className="rounded-md border border-slate-200 bg-white p-5 shadow-soft">
        <h2 className="text-lg font-semibold text-ink">Preview</h2>
        <div className="mt-4 flex justify-center">
          <QRPreview value={previewValue} foreground={foreground} background={background} />
        </div>
      </section>
    </div>
  );
}
