"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { ArrowLeft, CheckCircle2, RefreshCw } from "lucide-react";
import Link from "next/link";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { Button } from "@/components/common/Button";
import { api, messageFromError } from "@/lib/api";
import { bankMap } from "@/lib/constants";
import type { ApiEnvelope, CreatePaymentResponse } from "@/types";

export default function PaymentPage() {
  const params = useParams<{ transactionCode: string }>();
  const router = useRouter();
  const [payment, setPayment] = useState<CreatePaymentResponse | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [timeLeft, setTimeLeft] = useState<number | null>(null);

  // Stable refs — never trigger re-renders or effect re-runs
  const transactionCode = params.transactionCode;
  const pollingIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const countdownIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const isCancellingRef = useRef(false);
  const routerRef = useRef(router);
  routerRef.current = router;

  // ─── Stop polling helper ──────────────────────────────────────────────────
  const stopPolling = useCallback(() => {
    if (pollingIntervalRef.current !== null) {
      clearInterval(pollingIntervalRef.current);
      pollingIntervalRef.current = null;
    }
  }, []);

  // ─── Stop countdown helper ────────────────────────────────────────────────
  const stopCountdown = useCallback(() => {
    if (countdownIntervalRef.current !== null) {
      clearInterval(countdownIntervalRef.current);
      countdownIntervalRef.current = null;
    }
  }, []);

  // ─── Handle terminal status (SUCCESS / CANCELLED / FAILED) ───────────────
  const handleTerminalStatus = useCallback(
    (status: string) => {
      stopPolling();
      stopCountdown();
      if (status === "SUCCESS") {
        routerRef.current.push("/pricing");
      }
    },
    [stopPolling, stopCountdown]
  );

  // ─── Load payment once on mount ───────────────────────────────────────────
  const loadPayment = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const res = await api.get<ApiEnvelope<CreatePaymentResponse>>(
        `/payments/${transactionCode}`
      );
      const data = res.data.data ?? null;
      setPayment(data);
      if (data && data.payment.status !== "PENDING") {
        handleTerminalStatus(data.payment.status);
      }
    } catch (err) {
      setError(messageFromError(err));
    } finally {
      setLoading(false);
    }
  }, [transactionCode, handleTerminalStatus]);

  // ─── Expire / cancel ─────────────────────────────────────────────────────
  const handleExpire = useCallback(async () => {
    if (isCancellingRef.current) return;
    isCancellingRef.current = true;
    stopPolling();
    stopCountdown();
    try {
      const res = await api.post<ApiEnvelope<CreatePaymentResponse>>(
        `/payments/${transactionCode}/cancel`
      );
      const data = res.data.data ?? null;
      if (data) setPayment(data);
    } catch (err) {
      console.error("Failed to cancel expired payment", err);
    } finally {
      isCancellingRef.current = false;
    }
  }, [transactionCode, stopPolling, stopCountdown]);

  // ─── Mount: load payment ──────────────────────────────────────────────────
  useEffect(() => {
    loadPayment();
  }, [loadPayment]);

  // ─── Start polling once payment is loaded and PENDING ────────────────────
  // Deps: only transaction_code — so this never reruns when status changes inside setPayment.
  useEffect(() => {
    if (!payment) return;
    if (payment.payment.status !== "PENDING") return;
    if (pollingIntervalRef.current !== null) return; // already running

    pollingIntervalRef.current = setInterval(async () => {
      try {
        const res = await api.get<ApiEnvelope<CreatePaymentResponse>>(
          `/payments/${transactionCode}`
        );
        const newPayment = res.data.data ?? null;
        if (!newPayment) return;

        setPayment(newPayment);

        if (newPayment.payment.status !== "PENDING") {
          // Stop interval first to avoid race with cleanup
          stopPolling();
          stopCountdown();
          if (newPayment.payment.status === "SUCCESS") {
            routerRef.current.push("/pricing");
          }
        }
      } catch {
        // ignore transient polling errors
      }
    }, 5000);

    return () => stopPolling();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [payment?.payment.transaction_code]); // only re-run when a new payment record is loaded

  // ─── Start countdown once payment is loaded and PENDING ──────────────────
  // Same dep strategy: keyed on transaction_code, not on status.
  useEffect(() => {
    if (!payment) return;
    if (payment.payment.status !== "PENDING") {
      setTimeLeft(0);
      return;
    }
    if (countdownIntervalRef.current !== null) return; // already running

    const createdAt = new Date(payment.payment.created_at).getTime();

    const calcRemaining = () => {
      const diff = Math.floor((Date.now() - createdAt) / 1000);
      const remaining = 600 - diff;
      return remaining > 0 ? remaining : 0;
    };

    const initial = calcRemaining();
    setTimeLeft(initial);

    if (initial <= 0) {
      handleExpire();
      return;
    }

    countdownIntervalRef.current = setInterval(() => {
      const remaining = calcRemaining();
      setTimeLeft(remaining);
      if (remaining <= 0) {
        handleExpire();
      }
    }, 1000);

    return () => stopCountdown();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [payment?.payment.transaction_code]); // only re-run when a new payment record is loaded

  // ─── Status badge ─────────────────────────────────────────────────────────
  const getStatusBadge = (status: string) => {
    switch (status) {
      case "SUCCESS":
        return (
          <div className="mt-5 flex items-center justify-center gap-2 rounded-md bg-teal/10 px-3 py-2 text-sm font-semibold text-teal">
            <CheckCircle2 className="h-4 w-4" />
            Trạng thái: Đã thanh toán (Thành công)
          </div>
        );
      case "CANCELLED":
        return (
          <div className="mt-5 flex items-center justify-center gap-2 rounded-md bg-coral/10 px-3 py-2 text-sm font-semibold text-coral">
            Trạng thái: Đã hủy (Hết hạn)
          </div>
        );
      case "FAILED":
        return (
          <div className="mt-5 flex items-center justify-center gap-2 rounded-md bg-coral/10 px-3 py-2 text-sm font-semibold text-coral">
            Trạng thái: Giao dịch thất bại
          </div>
        );
      case "PENDING":
      default:
        return (
          <div className="mt-5 flex items-center justify-center gap-2 rounded-md bg-amber-500/10 px-3 py-2 text-sm font-semibold text-amber-700 animate-pulse">
            Trạng thái: Chờ thanh toán...
          </div>
        );
    }
  };

  return (
    <DashboardShell>
      <div className="mb-5">
        <Link href="/pricing" className="inline-flex items-center gap-2 text-sm font-semibold text-muted hover:text-ink">
          <ArrowLeft className="h-4 w-4" />
          Pricing
        </Link>
      </div>
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold text-ink">Thanh toán Pro</h1>
          <p className="mt-2 text-sm text-muted">Quét QR hoặc chuyển khoản đúng nội dung để hệ thống tự kích hoạt gói.</p>
        </div>
        <Button tone="secondary" onClick={loadPayment} disabled={loading}>
          <RefreshCw className="h-4 w-4" />
          {loading ? "Đang tải" : "Cập nhật"}
        </Button>
      </div>

      {error ? <p className="mt-4 rounded-md bg-coral/10 px-3 py-2 text-sm text-coral">{error}</p> : null}

      {payment ? (
        <section className="mt-6 grid gap-6 lg:grid-cols-[360px_1fr]">
          <div className="rounded-md border border-slate-200 bg-white p-6 shadow-soft">
            <div className="mx-auto flex aspect-square max-w-xs items-center justify-center rounded-md border border-slate-200 bg-white p-4 relative overflow-hidden">
              {payment.payment.status === "PENDING" ? (
                <img
                  src={payment.instructions.qr_image_url}
                  alt="Mã VietQR"
                  className="w-[240px] h-[240px] object-contain"
                />
              ) : (
                <div className="flex flex-col items-center justify-center text-center p-4">
                  <span className="text-4xl mb-2">⚠️</span>
                  <p className="text-sm font-semibold text-muted">Mã QR đã hết hiệu lực</p>
                </div>
              )}
            </div>
            {payment.payment.status === "PENDING" && timeLeft !== null && timeLeft > 0 ? (
              <div className="mt-4 rounded-md bg-amber-500/10 border border-amber-500/20 p-3 text-center">
                <p className="text-sm text-amber-800 font-medium">
                  Thời gian thanh toán còn lại:{" "}
                  <span className="font-mono text-base font-bold">
                    {Math.floor(timeLeft / 60)}:
                    {String(timeLeft % 60).padStart(2, "0")}
                  </span>
                </p>
              </div>
            ) : null}
            {getStatusBadge(payment.payment.status)}
            {!payment.instructions.enabled ? (
              <p className="mt-3 rounded-md bg-coral/10 px-3 py-2 text-sm text-coral">Sepay đang tắt ở môi trường backend này.</p>
            ) : null}
          </div>

          <div className="rounded-md border border-slate-200 bg-white p-6 shadow-soft">
            <h2 className="text-xl font-bold text-ink">Thông tin chuyển khoản</h2>
            <dl className="mt-5 grid gap-4 text-sm md:grid-cols-2">
              <div><dt className="text-muted">Ngân hàng</dt><dd className="mt-1 font-semibold text-ink">{bankMap[payment.instructions.bank_code] || payment.instructions.bank_code || "Chưa cấu hình"}</dd></div>
              <div><dt className="text-muted">Số tài khoản</dt><dd className="mt-1 font-semibold text-ink">{payment.instructions.account_no || "Chưa cấu hình"}</dd></div>
              <div><dt className="text-muted">Chủ tài khoản</dt><dd className="mt-1 font-semibold text-ink">{payment.instructions.account_name || "Chưa cấu hình"}</dd></div>
              <div><dt className="text-muted">Số tiền</dt><dd className="mt-1 font-semibold text-ink">{payment.instructions.amount.toLocaleString("vi-VN")} {payment.instructions.currency}</dd></div>
              <div className="md:col-span-2">
                <dt className="text-muted">Nội dung chuyển khoản</dt>
                <dd className="mt-1 break-all rounded-md bg-panel px-3 py-2 font-mono text-base font-semibold text-ink">{payment.instructions.transfer_content}</dd>
              </div>
              <div className="md:col-span-2">
                <dt className="text-muted">Mã giao dịch</dt>
                <dd className="mt-1 break-all font-mono font-semibold text-ink">{payment.payment.transaction_code}</dd>
              </div>
            </dl>
          </div>
        </section>
      ) : null}
    </DashboardShell>
  );
}
