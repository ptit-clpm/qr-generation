"use client";

import { useEffect, useState } from "react";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { ScanBarChart } from "@/components/analytics/ScanCharts";
import { EmptyState } from "@/components/common/State";
import { api, messageFromError } from "@/lib/api";
import type { ApiEnvelope, QRCode } from "@/types";

interface ChartRow {
  label: string;
  count: number;
}

export default function AnalyticsPage() {
  const [qrs, setQrs] = useState<QRCode[]>([]);
  const [qrId, setQrId] = useState<number | null>(null);
  const [rows, setRows] = useState<ChartRow[]>([]);
  const [error, setError] = useState("");

  useEffect(() => {
    api.get<ApiEnvelope<{ items: QRCode[] }>>("/qrcodes").then((res) => {
      const dynamic = (res.data.data?.items ?? []).filter((qr) => qr.is_dynamic);
      setQrs(dynamic);
      setQrId(dynamic[0]?.id ?? null);
    }).catch((err) => setError(messageFromError(err)));
  }, []);

  useEffect(() => {
    if (!qrId) return;
    api.get<ApiEnvelope<ChartRow[]>>(`/qrcodes/${qrId}/analytics/by-date`).then((res) => setRows(res.data.data ?? [])).catch((err) => setError(messageFromError(err)));
  }, [qrId]);

  return (
    <DashboardShell>
      <h1 className="text-3xl font-bold text-ink">Analytics</h1>
      <div className="mt-6 max-w-sm">
        <select value={qrId ?? ""} onChange={(e) => setQrId(Number(e.target.value))} className="focus-ring w-full rounded-md border border-slate-200 bg-white px-3 py-2">
          {qrs.map((qr) => <option key={qr.id} value={qr.id}>{qr.title}</option>)}
        </select>
      </div>
      {error ? <p className="mt-4 rounded-md bg-coral/10 px-3 py-2 text-sm text-coral">{error}</p> : null}
      <div className="mt-6">
        {rows.length ? <ScanBarChart data={rows} /> : <EmptyState title="No analytics data" description="Dynamic QR scans will populate the chart." />}
      </div>
    </DashboardShell>
  );
}
