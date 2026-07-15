"use client";

import { useEffect, useState } from "react";
import { BarChart3, CreditCard, QrCode, Zap } from "lucide-react";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { api } from "@/lib/api";
import type { ApiEnvelope, QRCode } from "@/types";

export default function DashboardPage() {
  const [qrs, setQrs] = useState<QRCode[]>([]);
  useEffect(() => {
    api.get<ApiEnvelope<{ items: QRCode[] }>>("/qrcodes").then((res) => setQrs(res.data.data?.items ?? [])).catch(() => setQrs([]));
  }, []);
  const scans = qrs.reduce((sum, qr) => sum + qr.scan_count, 0);

  return (
    <DashboardShell>
      <h1 className="text-3xl font-bold text-ink">Dashboard</h1>
      <div className="mt-6 grid gap-4 md:grid-cols-4">
        {[
          { label: "QR Codes", value: qrs.length, icon: QrCode },
          { label: "Scans", value: scans, icon: BarChart3 },
          { label: "Dynamic", value: qrs.filter((qr) => qr.is_dynamic).length, icon: Zap },
          { label: "Plan", value: "Free/Pro", icon: CreditCard }
        ].map((item) => {
          const Icon = item.icon;
          return (
            <div key={item.label} className="rounded-md border border-slate-200 bg-white p-5 shadow-soft">
              <Icon className="h-5 w-5 text-teal" />
              <p className="mt-4 text-sm text-muted">{item.label}</p>
              <p className="mt-1 text-2xl font-bold text-ink">{item.value}</p>
            </div>
          );
        })}
      </div>
    </DashboardShell>
  );
}
