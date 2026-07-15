"use client";

import { useEffect } from "react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { BarChart3, CreditCard, Folder, LayoutDashboard, LogOut, QrCode, Settings } from "lucide-react";
import { clsx } from "clsx";
import { useAuthStore } from "@/stores/auth";
import { Button } from "@/components/common/Button";

const nav = [
  { href: "/dashboard", label: "Dashboard", icon: LayoutDashboard },
  { href: "/qrcodes", label: "QR Codes", icon: QrCode },
  { href: "/folders", label: "Folders", icon: Folder },
  { href: "/analytics", label: "Analytics", icon: BarChart3 },
  { href: "/account", label: "Account", icon: Settings },
  { href: "/pricing", label: "Plan", icon: CreditCard }
];

export function DashboardShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();
  const user = useAuthStore((state) => state.user);
  const loading = useAuthStore((state) => state.loading);
  const loadMe = useAuthStore((state) => state.loadMe);
  const logout = useAuthStore((state) => state.logout);
  const isAdmin = useAuthStore((state) => state.isAdmin());

  useEffect(() => {
    const token = window.localStorage.getItem("access_token");
    if (!token) {
      router.replace("/login");
      return;
    }
    if (!user) {
      loadMe().catch(() => router.replace("/login"));
    }
  }, [loadMe, router, user]);

  if (loading || !user) {
    return (
      <div className="grid min-h-screen place-items-center bg-panel px-4">
        <p className="text-sm font-medium text-muted">Loading...</p>
      </div>
    );
  }

  if (isAdmin) {
    return (
      <div className="grid min-h-screen place-items-center bg-panel px-4">
        <div className="w-full max-w-md rounded-md border border-slate-200 bg-white p-6 text-center shadow-soft">
          <div className="mx-auto grid h-12 w-12 place-items-center rounded-md bg-ink text-white">
            <QrCode className="h-6 w-6" />
          </div>
          <h1 className="mt-5 text-2xl font-bold text-ink">Trang này dành cho người dùng</h1>
          <p className="mt-3 text-sm text-muted">Tài khoản admin cần đăng nhập vào trang quản trị riêng.</p>
          <div className="mt-6 flex flex-col gap-3 sm:flex-row">
            <a className="focus-ring inline-flex min-h-10 flex-1 items-center justify-center rounded-md bg-teal px-4 py-2 text-sm font-semibold text-white hover:bg-teal/90" href={process.env.NEXT_PUBLIC_ADMIN_URL ?? "http://localhost:5173"}>
              Mở trang admin
            </a>
            <Button className="flex-1" tone="secondary" onClick={logout}>Đăng xuất</Button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-panel">
      <aside className="fixed left-0 top-0 hidden h-full w-64 border-r border-slate-200 bg-white p-4 lg:block">
        <Link href="/" className="mb-8 flex items-center gap-2 font-bold text-ink">
          <span className="grid h-9 w-9 place-items-center rounded-md bg-ink text-white">
            <QrCode className="h-5 w-5" />
          </span>
          QR Studio
        </Link>
        <nav className="space-y-1">
          {nav.map((item) => {
            const Icon = item.icon;
            return (
              <Link
                key={item.href}
                href={item.href}
                className={clsx(
                  "flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium",
                  pathname === item.href ? "bg-teal text-white" : "text-muted hover:bg-panel hover:text-ink"
                )}
              >
                <Icon className="h-4 w-4" />
                {item.label}
              </Link>
            );
          })}
        </nav>
        <button onClick={logout} className="absolute bottom-4 left-4 flex items-center gap-2 rounded-md px-3 py-2 text-sm font-medium text-muted hover:bg-panel">
          <LogOut className="h-4 w-4" />
          Logout
        </button>
      </aside>
      <main className="lg:pl-64">
        <div className="mx-auto max-w-6xl px-4 py-6">{children}</div>
      </main>
    </div>
  );
}
