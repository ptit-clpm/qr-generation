"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { useParams } from "next/navigation";
import { ArrowLeft, CheckCircle2, RefreshCw } from "lucide-react";
import Link from "next/link";
import { QRCodeCanvas } from "qrcode.react";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { Button } from "@/components/common/Button";
import { api, messageFromError } from "@/lib/api";
import type { ApiEnvelope, CreatePaymentResponse } from "@/types";

export default function PaymentPage() {
  const params = useParams<{ transactionCode: string }>();
  const [payment, setPayment] = useState<CreatePaymentResponse | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const qrValue = useMemo(() => {
    if (!payment) return "";
    const info = payment.instructions;
    return [
      "SEPAY BANK TRANSFER",
      `Bank: ${info.bank_name}`,
      `Account: ${info.bank_account}`,
      `Account name: ${info.account_name}`,
      `Amount: ${info.amount} ${info.currency}`,
      `Content: ${info.transfer_content}`
    ].join("\n");
  }, [payment]);

  const transactionCode = params.transactionCode;

  const loadPayment = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const res = await api.get<ApiEnvelope<CreatePaymentResponse>>(`/payments/${transactionCode}`);
      setPayment(res.data.data ?? null);
    } catch (err) {
      setError(messageFromError(err));
    } finally {
      setLoading(false);
    }
  }, [transactionCode]);

  useEffect(() => {
    loadPayment();
  }, [loadPayment]);

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
            <div className="mx-auto flex aspect-square max-w-xs items-center justify-center rounded-md border border-slate-200 bg-white p-4">
              <QRCodeCanvas value={qrValue} size={240} includeMargin />
            </div>
            <div className="mt-5 flex items-center justify-center gap-2 rounded-md bg-teal/10 px-3 py-2 text-sm font-semibold text-teal">
              <CheckCircle2 className="h-4 w-4" />
              Trạng thái: {payment.payment.status}
            </div>
            {!payment.instructions.enabled ? (
              <p className="mt-3 rounded-md bg-coral/10 px-3 py-2 text-sm text-coral">Sepay đang tắt ở môi trường backend này.</p>
            ) : null}
          </div>

          <div className="rounded-md border border-slate-200 bg-white p-6 shadow-soft">
            <h2 className="text-xl font-bold text-ink">Thông tin chuyển khoản</h2>
            <dl className="mt-5 grid gap-4 text-sm md:grid-cols-2">
              <div><dt className="text-muted">Ngân hàng</dt><dd className="mt-1 font-semibold text-ink">{payment.instructions.bank_name || "Chưa cấu hình"}</dd></div>
              <div><dt className="text-muted">Số tài khoản</dt><dd className="mt-1 font-semibold text-ink">{payment.instructions.bank_account || "Chưa cấu hình"}</dd></div>
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
