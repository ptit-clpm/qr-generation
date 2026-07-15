"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { Download, Search, Trash2 } from "lucide-react";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { CreateQRForm } from "@/components/qrcode/CreateQRForm";
import { EmptyState } from "@/components/common/State";
import { Button } from "@/components/common/Button";
import { api, messageFromError } from "@/lib/api";
import type { ApiEnvelope, QRCode } from "@/types";

export default function QRCodesPage() {
  const [items, setItems] = useState<QRCode[]>([]);
  const [q, setQ] = useState("");
  const [error, setError] = useState("");

  async function load() {
    try {
      const res = await api.get<ApiEnvelope<{ items: QRCode[] }>>("/qrcodes", { params: { q } });
      setItems(res.data.data?.items ?? []);
    } catch (err) {
      setError(messageFromError(err));
    }
  }

  useEffect(() => { load(); }, []);

  async function remove(id: number) {
    await api.delete(`/qrcodes/${id}`);
    load();
  }

  async function download(id: number) {
    const res = await api.get(`/qrcodes/${id}/download`, { responseType: "blob" });
    const url = URL.createObjectURL(res.data);
    const link = document.createElement("a");
    link.href = url;
    link.download = "qr-code.png";
    link.click();
    URL.revokeObjectURL(url);
  }

  return (
    <DashboardShell>
      <div className="flex flex-wrap items-center justify-between gap-3">
        <h1 className="text-3xl font-bold text-ink">QR Codes</h1>
        <div className="flex gap-2">
          <input value={q} onChange={(e) => setQ(e.target.value)} className="focus-ring rounded-md border border-slate-200 px-3 py-2 text-sm" placeholder="Search" />
          <Button tone="secondary" onClick={load}><Search className="h-4 w-4" /></Button>
        </div>
      </div>
      <div className="mt-6">
        <CreateQRForm />
      </div>
      <section className="mt-8">
        <h2 className="text-xl font-bold text-ink">Saved QR</h2>
        {error ? <p className="mt-3 rounded-md bg-coral/10 px-3 py-2 text-sm text-coral">{error}</p> : null}
        {items.length === 0 ? <div className="mt-4"><EmptyState title="No QR codes found" description="Saved QR codes will appear here." /></div> : null}
        <div className="mt-4 grid gap-4 md:grid-cols-2">
          {items.map((qr) => (
            <div key={qr.id} className="rounded-md border border-slate-200 bg-white p-4 shadow-soft">
              <div className="flex items-start justify-between gap-3">
                <div>
                  <Link href={`/qrcodes/${qr.id}`} className="font-semibold text-ink hover:text-teal">{qr.title}</Link>
                  <p className="mt-1 text-sm text-muted">{qr.qr_type} · {qr.status} · {qr.scan_count} scans</p>
                </div>
                <span className="rounded-md bg-panel px-2 py-1 text-xs font-semibold text-muted">{qr.is_dynamic ? "Dynamic" : "Static"}</span>
              </div>
              <div className="mt-4 flex gap-2">
                <Button tone="secondary" onClick={() => download(qr.id)}><Download className="h-4 w-4" /></Button>
                <Button tone="danger" onClick={() => remove(qr.id)}><Trash2 className="h-4 w-4" /></Button>
              </div>
            </div>
          ))}
        </div>
      </section>
    </DashboardShell>
  );
}
