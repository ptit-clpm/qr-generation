"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { Activity, Calendar, Compass, Crown, Globe, Smartphone } from "lucide-react";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { CategoryBarChart, ScanBarChart } from "@/components/analytics/ScanCharts";
import { EmptyState, ErrorState, LoadingState } from "@/components/common/State";
import { Button } from "@/components/common/Button";
import { api, messageFromError } from "@/lib/api";
import type { ApiEnvelope, QRCode } from "@/types";

interface SummaryData {
  qr_id: number;
  scan_count: number;
  first_scan?: string | null;
  last_scan?: string | null;
  top_device?: string;
  top_browser?: string;
}

interface StatRow {
  label: string;
  count: number;
}

type DateRangeOption = 7 | 30 | 90;

export default function AnalyticsPage() {
  const [qrs, setQrs] = useState<QRCode[]>([]);
  const [selectedQrId, setSelectedQrId] = useState<number | null>(null);
  const [dateRangeDays, setDateRangeDays] = useState<DateRangeOption>(30);
  const [summary, setSummary] = useState<SummaryData | null>(null);
  const [byDateRows, setByDateRows] = useState<StatRow[]>([]);
  const [byDeviceRows, setByDeviceRows] = useState<StatRow[]>([]);
  const [byBrowserRows, setByBrowserRows] = useState<StatRow[]>([]);
  const [byLocationRows, setByLocationRows] = useState<StatRow[]>([]);
  
  const [initialLoading, setInitialLoading] = useState(true);
  const [analyticsLoading, setAnalyticsLoading] = useState(false);
  const [errorStatus, setErrorStatus] = useState<number | null>(null);
  const [errorMessage, setErrorMessage] = useState("");

  // Load user's QR codes and filter for Dynamic URL QR codes
  useEffect(() => {
    setInitialLoading(true);
    setErrorStatus(null);
    setErrorMessage("");

    api
      .get<ApiEnvelope<{ items: QRCode[] }>>("/qrcodes", { params: { limit: 100 } })
      .then((res) => {
        const items = res.data.data?.items ?? [];
        const dynamicUrlQrs = items.filter((qr) => qr.is_dynamic && qr.qr_type === "URL");
        setQrs(dynamicUrlQrs);
        if (dynamicUrlQrs.length > 0) {
          setSelectedQrId(dynamicUrlQrs[0].id);
        }
      })
      .catch((err) => {
        const msg = messageFromError(err);
        setErrorMessage(msg);
        if (msg.includes("403") || msg.toLowerCase().includes("pro")) {
          setErrorStatus(403);
        } else {
          setErrorStatus(500);
        }
      })
      .finally(() => {
        setInitialLoading(false);
      });
  }, []);

  // Fetch analytics data when selected QR code or date range changes
  useEffect(() => {
    if (!selectedQrId) return;

    setAnalyticsLoading(true);
    setErrorStatus(null);
    setErrorMessage("");

    const now = new Date();
    const fromDateObj = new Date();
    fromDateObj.setDate(now.getDate() - dateRangeDays);

    const fromDateStr = fromDateObj.toISOString().split("T")[0];
    const toDateStr = now.toISOString().split("T")[0];

    Promise.all([
      api.get<ApiEnvelope<SummaryData>>(`/qrcodes/${selectedQrId}/analytics/summary`),
      api.get<ApiEnvelope<StatRow[]>>(`/qrcodes/${selectedQrId}/analytics/by-date`, { params: { from: fromDateStr, to: toDateStr } }),
      api.get<ApiEnvelope<StatRow[]>>(`/qrcodes/${selectedQrId}/analytics/by-device`),
      api.get<ApiEnvelope<StatRow[]>>(`/qrcodes/${selectedQrId}/analytics/by-browser`),
      api.get<ApiEnvelope<StatRow[]>>(`/qrcodes/${selectedQrId}/analytics/by-location`)
    ])
      .then(([summaryRes, byDateRes, byDeviceRes, byBrowserRes, byLocationRes]) => {
        setSummary(summaryRes.data.data ?? null);
        setByDateRows(byDateRes.data.data ?? []);
        setByDeviceRows(byDeviceRes.data.data ?? []);
        setByBrowserRows(byBrowserRes.data.data ?? []);
        setByLocationRows(byLocationRes.data.data ?? []);
      })
      .catch((err) => {
        const status = err?.response?.status;
        const msg = messageFromError(err);
        setErrorStatus(status || 500);
        setErrorMessage(msg);
      })
      .finally(() => {
        setAnalyticsLoading(false);
      });
  }, [selectedQrId, dateRangeDays]);
  const formatDate = (isoStr?: string | null) => {
    if (!isoStr) return "N/A";
    const d = new Date(isoStr);
    return isNaN(d.getTime()) ? "N/A" : d.toLocaleDateString("vi-VN", { day: "2-digit", month: "2-digit", year: "numeric", hour: "2-digit", minute: "2-digit" });
  };

  if (initialLoading) {
    return (
      <DashboardShell>
        <LoadingState label="Loading analytics dashboard..." />
      </DashboardShell>
    );
  }

  if (errorStatus === 403) {
    return (
      <DashboardShell>
        <div className="mx-auto max-w-2xl rounded-lg border border-amber-200 bg-amber-50/50 p-8 text-center shadow-soft">
          <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-amber-100 text-amber-600">
            <Crown className="h-6 w-6" />
          </div>
          <h2 className="mt-4 text-xl font-bold text-ink">Yêu cầu gói Pro</h2>
          <p className="mt-2 text-sm text-muted">
            Tính năng thống kê và phân tích QR code (Scan Analytics) chỉ dành riêng cho tài khoản Pro còn thời hạn.
          </p>
          <Link href="/pricing" className="mt-6 inline-block">
            <Button>
              <Crown className="h-4 w-4 mr-2" />
              Nâng cấp lên Pro ngay
            </Button>
          </Link>
        </div>
      </DashboardShell>
    );
  }

  return (
    <DashboardShell>
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-3xl font-bold text-ink">Scan Analytics</h1>
          <p className="text-sm text-muted">Thống kê chi tiết lượt quét các mã Dynamic QR</p>
        </div>

        {qrs.length > 0 ? (
          <div className="flex flex-wrap items-center gap-3">
            {/* QR Selector */}
            <select
              value={selectedQrId ?? ""}
              onChange={(e) => setSelectedQrId(Number(e.target.value))}
              className="focus-ring rounded-md border border-slate-200 bg-white px-3 py-2 text-sm font-medium text-ink shadow-sm"
            >
              {qrs.map((qr) => (
                <option key={qr.id} value={qr.id}>
                  {qr.title} ({qr.status})
                </option>
              ))}
            </select>

            {/* Date Range Selector */}
            <div className="inline-flex rounded-md border border-slate-200 bg-white p-1 shadow-sm">
              {([7, 30, 90] as DateRangeOption[]).map((days) => (
                <button
                  key={days}
                  onClick={() => setDateRangeDays(days)}
                  className={`rounded px-3 py-1 text-xs font-medium transition ${
                    dateRangeDays === days ? "bg-teal text-white shadow-sm" : "text-muted hover:text-ink"
                  }`}
                >
                  {days} ngày
                </button>
              ))}
            </div>
          </div>
        ) : null}
      </div>

      {qrs.length === 0 ? (
        <div className="mt-8">
          <EmptyState
            title="Chưa có Dynamic QR nào"
            description="Chỉ mã URL Dynamic QR mới ghi nhận scan analytics. Hãy tạo mã Dynamic URL QR để bắt đầu theo dõi."
          />
        </div>
      ) : errorStatus && errorStatus !== 403 ? (
        <div className="mt-8">
          <ErrorState message={errorMessage || "Không thể tải dữ liệu analytics"} />
        </div>
      ) : analyticsLoading ? (
        <div className="mt-8">
          <LoadingState label="Đang tải dữ liệu thống kê..." />
        </div>
      ) : (
        <div className="mt-8 space-y-6">
          {/* Summary Cards */}
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-5">
            <div className="rounded-md border border-slate-200 bg-white p-4 shadow-soft">
              <div className="flex items-center gap-2 text-xs font-medium text-muted">
                <Activity className="h-4 w-4 text-teal" />
                <span>Total scans</span>
              </div>
              <p className="mt-2 text-2xl font-bold text-ink">{summary?.scan_count ?? 0}</p>
              <p className="mt-1 text-[11px] text-muted">Tổng số lượt redirect</p>
            </div>

            <div className="rounded-md border border-slate-200 bg-white p-4 shadow-soft">
              <div className="flex items-center gap-2 text-xs font-medium text-muted">
                <Calendar className="h-4 w-4 text-sky-500" />
                <span>First scan</span>
              </div>
              <p className="mt-2 text-sm font-bold text-ink">{formatDate(summary?.first_scan)}</p>
              <p className="mt-1 text-[11px] text-muted">Lần quét đầu tiên</p>
            </div>

            <div className="rounded-md border border-slate-200 bg-white p-4 shadow-soft">
              <div className="flex items-center gap-2 text-xs font-medium text-muted">
                <Calendar className="h-4 w-4 text-emerald-500" />
                <span>Last scan</span>
              </div>
              <p className="mt-2 text-sm font-bold text-ink">{formatDate(summary?.last_scan)}</p>
              <p className="mt-1 text-[11px] text-muted">Lần quét gần nhất</p>
            </div>

            <div className="rounded-md border border-slate-200 bg-white p-4 shadow-soft">
              <div className="flex items-center gap-2 text-xs font-medium text-muted">
                <Smartphone className="h-4 w-4 text-indigo-500" />
                <span>Top device</span>
              </div>
              <p className="mt-2 text-base font-bold text-ink">{summary?.top_device || "N/A"}</p>
              <p className="mt-1 text-[11px] text-muted">Thiết bị phổ biến nhất</p>
            </div>

            <div className="rounded-md border border-slate-200 bg-white p-4 shadow-soft">
              <div className="flex items-center gap-2 text-xs font-medium text-muted">
                <Compass className="h-4 w-4 text-amber-500" />
                <span>Top browser</span>
              </div>
              <p className="mt-2 text-base font-bold text-ink">{summary?.top_browser || "N/A"}</p>
              <p className="mt-1 text-[11px] text-muted">Trình duyệt phổ biến nhất</p>
            </div>
          </div>

          {/* If scan_count is 0 */}
          {summary?.scan_count === 0 ? (
            <EmptyState
              title="Mã Dynamic QR này chưa có lượt quét nào"
              description="Hãy chia sẻ hoặc quét mã để ghi nhận lượt redirect thực tế."
            />
          ) : (
            <>
              {/* Scans by Date Chart */}
              <ScanBarChart data={byDateRows} title={`Lượt quét theo ngày (${dateRangeDays} ngày gần nhất)`} />

              {/* Scans by Device & Browser */}
              <div className="grid gap-6 md:grid-cols-2">
                <CategoryBarChart data={byDeviceRows} title="Theo thiết bị (Device)" fillColor="#0f9f8f" />
                <CategoryBarChart data={byBrowserRows} title="Theo trình duyệt (Browser)" fillColor="#6366F1" />
              </div>

              {/* Scans by Location */}
              <div className="rounded-md border border-slate-200 bg-white p-5 shadow-soft">
                <h3 className="mb-4 text-base font-semibold text-ink flex items-center gap-2">
                  <Globe className="h-4 w-4 text-teal" />
                  Vị trí địa lý (Location)
                </h3>
                {byLocationRows.length > 0 ? (
                  <div className="divide-y divide-slate-100">
                    {byLocationRows.map((row) => (
                      <div key={row.label} className="flex items-center justify-between py-2 text-sm">
                        <span className="font-medium text-ink">{row.label}</span>
                        <span className="font-bold text-teal">{row.count} scans</span>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="py-6 text-center text-sm text-muted">
                    Location data unavailable
                  </div>
                )}
              </div>
            </>
          )}
        </div>
      )}
    </DashboardShell>
  );
}
