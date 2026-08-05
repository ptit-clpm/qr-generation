"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Check, CreditCard, LockKeyhole } from "lucide-react";
import { PublicHeader } from "@/components/layout/PublicHeader";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { Button } from "@/components/common/Button";
import { api, messageFromError } from "@/lib/api";
import { useAuthStore } from "@/stores/auth";
import type { ApiEnvelope, CreatePaymentResponse, Plan, PlanName, Subscription } from "@/types";

export default function PricingPage() {
  const router = useRouter();
  const [plans, setPlans] = useState<Plan[]>([]);
  const [creatingPlanId, setCreatingPlanId] = useState<number | null>(null);
  const user = useAuthStore((state) => state.user);
  const loadMe = useAuthStore((state) => state.loadMe);
  const [currentPlan, setCurrentPlan] = useState<PlanName>("FREE");
  const [error, setError] = useState("");
  const authed = Boolean(user) || currentPlan === "PRO";

  useEffect(() => {
    const hasToken = Boolean(localStorage.getItem("access_token"));
    api.get<ApiEnvelope<Plan[]>>("/plans").then((res) => setPlans(res.data.data ?? [])).catch((err) => setError(messageFromError(err)));
    if (hasToken) {
      loadMe().catch(() => {});
      api.get<ApiEnvelope<Subscription>>("/users/subscription")
        .then((res) => setCurrentPlan(res.data.data?.plan?.name ?? "FREE"))
        .catch(() => setCurrentPlan("FREE"));
    }
  }, [loadMe]);

  async function choose(plan: Plan) {
    setError("");
    if (plan.name === currentPlan) return;
    if (!localStorage.getItem("access_token")) {
      router.push("/login");
      return;
    }
    if (plan.name === "FREE") return;
    setCreatingPlanId(plan.id);
    try {
      const res = await api.post<ApiEnvelope<CreatePaymentResponse>>("/payments/create", {
        plan_id: plan.id,
        payment_method: "SEPAY"
      });
      const transactionCode = res.data.data?.payment.transaction_code;
      if (transactionCode) {
        router.push(`/payments/${transactionCode}`);
      }
    } catch (err) {
      setError(messageFromError(err));
    } finally {
      setCreatingPlanId(null);
    }
  }

  const content = (
    <div>
      <h1 className="text-3xl font-bold text-ink">Pricing</h1>
      {error ? <p className="mt-4 rounded-md bg-coral/10 px-3 py-2 text-sm text-coral">{error}</p> : null}
      <div className="mt-6 grid gap-5 md:grid-cols-2">
        {plans.map((plan) => (
          <div key={plan.id} className="rounded-md border border-slate-200 bg-white p-6 shadow-soft">
            <div className="flex items-start justify-between gap-4">
              <div>
                <h2 className="text-xl font-bold text-ink">{plan.name}</h2>
                <p className="mt-2 text-sm text-muted">{plan.description}</p>
              </div>
              <p className="text-2xl font-bold text-teal">{plan.price.toLocaleString("vi-VN")}đ</p>
            </div>
            <ul className="mt-6 space-y-3 text-sm text-muted">
              {(plan.name === "PRO"
                ? [
                    `Tối đa ${plan.max_qr_codes} mã QR`,
                    "Dynamic QR (áp dụng cho URL)",
                    "Scan Analytics chi tiết (lượt quét, thiết bị, trình duyệt)",
                    "Hỗ trợ tải lên Logo custom",
                    "Danh mục URL nâng cao (Social profile, PDF link, Menu online)",
                    "Tải mã QR định dạng PNG"
                  ]
                : [
                    `Tối đa ${plan.max_qr_codes} mã QR`,
                    "Static QR (mã hóa dữ liệu cố định)",
                    "Các loại QR cơ bản: URL, Text, WiFi, vCard, Email, SMS, Location",
                    "Tải mã QR định dạng PNG",
                    "Thiết kế cơ bản (màu sắc, khung)"
                  ]
              ).map((item) => (
                <li key={item} className="flex items-center gap-2">
                  <Check className="h-4 w-4 text-teal flex-shrink-0" />
                  <span>{item}</span>
                </li>
              ))}
            </ul>
            <Button
              className="mt-6 w-full"
              tone={plan.name === currentPlan || (plan.name === "FREE" && currentPlan === "PRO") ? "secondary" : "primary"}
              onClick={() => choose(plan)}
              disabled={creatingPlanId === plan.id || plan.name === currentPlan || (plan.name === "FREE" && currentPlan === "PRO") || (!user && plan.name === "FREE")}
            >
              {plan.name === currentPlan ? <Check className="h-4 w-4" /> : plan.name === "FREE" && currentPlan === "PRO" ? <LockKeyhole className="h-4 w-4" /> : user ? <CreditCard className="h-4 w-4" /> : <LockKeyhole className="h-4 w-4" />}
              {plan.name === currentPlan
                ? "Gói hiện tại của bạn"
                : creatingPlanId === plan.id
                  ? "Đang tạo thanh toán"
                  : currentPlan === "PRO" && plan.name === "FREE"
                    ? "Không thể hạ xuống Free"
                  : plan.name === "PRO"
                    ? "Nâng cấp lên Pro"
                    : "Đăng nhập để dùng Free"}
            </Button>
          </div>
        ))}
      </div>
    </div>
  );

  return authed ? <DashboardShell>{content}</DashboardShell> : <><PublicHeader /><main className="mx-auto max-w-6xl px-4 py-10">{content}</main></>;
}
