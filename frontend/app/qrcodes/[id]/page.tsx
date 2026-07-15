"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { Save } from "lucide-react";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { QRPreview } from "@/components/qrcode/QRPreview";
import { Button } from "@/components/common/Button";
import { api, messageFromError } from "@/lib/api";
import { backendUrl } from "@/lib/constants";
import type { ApiEnvelope, QRCode } from "@/types";

export default function QRDetailPage() {
  const params = useParams<{ id: string }>();
  const [qr, setQr] = useState<QRCode | null>(null);
  const [message, setMessage] = useState("");

  useEffect(() => {
    api.get<ApiEnvelope<QRCode>>(`/qrcodes/${params.id}`).then((res) => setQr(res.data.data ?? null)).catch((err) => setMessage(messageFromError(err)));
  }, [params.id]);

  async function save() {
    if (!qr) return;
    try {
      const res = await api.put<ApiEnvelope<QRCode>>(`/qrcodes/${qr.id}`, {
        title: qr.title,
        destination_url: qr.destination_url,
        status: qr.status
      });
      setQr(res.data.data ?? qr);
      setMessage("QR updated");
    } catch (err) {
      setMessage(messageFromError(err));
    }
  }

  const preview = qr?.is_dynamic && qr.short_code ? `${backendUrl}/q/${qr.short_code}` : qr?.content ?? "";

  return (
    <DashboardShell>
      <h1 className="text-3xl font-bold text-ink">QR Detail</h1>
      {qr ? (
        <div className="mt-6 grid gap-6 lg:grid-cols-[1fr_360px]">
          <section className="rounded-md border border-slate-200 bg-white p-5 shadow-soft">
            <div className="grid gap-4">
              <label className="text-sm font-medium text-ink">
                Title
                <input value={qr.title} onChange={(e) => setQr({ ...qr, title: e.target.value })} className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2" />
              </label>
              <label className="text-sm font-medium text-ink">
                Status
                <select value={qr.status} onChange={(e) => setQr({ ...qr, status: e.target.value as QRCode["status"] })} className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2">
                  <option value="ACTIVE">ACTIVE</option>
                  <option value="DISABLED">DISABLED</option>
                </select>
              </label>
              {qr.is_dynamic ? (
                <label className="text-sm font-medium text-ink">
                  Destination URL
                  <input value={qr.destination_url ?? ""} onChange={(e) => setQr({ ...qr, destination_url: e.target.value })} className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2" />
                </label>
              ) : (
                <label className="text-sm font-medium text-ink">
                  Static content
                  <textarea value={qr.content} readOnly rows={5} className="mt-1 w-full rounded-md border border-slate-200 bg-panel px-3 py-2 text-muted" />
                </label>
              )}
            </div>
            {message ? <p className="mt-4 rounded-md bg-teal/10 px-3 py-2 text-sm text-teal">{message}</p> : null}
            <Button className="mt-5" onClick={save}><Save className="h-4 w-4" />Save</Button>
          </section>
          <section className="rounded-md border border-slate-200 bg-white p-5 shadow-soft">
            <p className="mb-4 font-semibold text-ink">{qr.qr_type} · {qr.scan_count} scans</p>
            <QRPreview value={preview} foreground={qr.design?.foreground_color} background={qr.design?.background_color} />
          </section>
        </div>
      ) : (
        <p className="mt-6 text-muted">{message || "Loading"}</p>
      )}
    </DashboardShell>
  );
}
