"use client";

import { Bar, BarChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";

export function ScanBarChart({ data, title = "Scans by Date" }: { data: Array<{ label: string; count: number }>; title?: string }) {
  return (
    <div className="rounded-md border border-slate-200 bg-white p-5 shadow-soft">
      <h3 className="mb-4 text-base font-semibold text-ink">{title}</h3>
      <div className="h-72">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart data={data}>
            <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#E2E8F0" />
            <XAxis dataKey="label" stroke="#64748B" fontSize={12} tickLine={false} />
            <YAxis allowDecimals={false} stroke="#64748B" fontSize={12} tickLine={false} />
            <Tooltip
              contentStyle={{ backgroundColor: "#1E293B", borderRadius: "6px", border: "none", color: "#F8FAFC" }}
              itemStyle={{ color: "#38BDF8" }}
            />
            <Bar dataKey="count" fill="#0f9f8f" radius={[4, 4, 0, 0]} name="Scans" />
          </BarChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}

export function CategoryBarChart({
  data,
  title,
  fillColor = "#3B82F6"
}: {
  data: Array<{ label: string; count: number }>;
  title: string;
  fillColor?: string;
}) {
  return (
    <div className="rounded-md border border-slate-200 bg-white p-5 shadow-soft">
      <h3 className="mb-4 text-base font-semibold text-ink">{title}</h3>
      <div className="h-60">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart data={data} layout="vertical">
            <CartesianGrid strokeDasharray="3 3" horizontal={false} stroke="#E2E8F0" />
            <XAxis type="number" allowDecimals={false} stroke="#64748B" fontSize={12} />
            <YAxis dataKey="label" type="category" stroke="#64748B" fontSize={12} width={80} tickLine={false} />
            <Tooltip
              contentStyle={{ backgroundColor: "#1E293B", borderRadius: "6px", border: "none", color: "#F8FAFC" }}
            />
            <Bar dataKey="count" fill={fillColor} radius={[0, 4, 4, 0]} name="Scans" />
          </BarChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}
