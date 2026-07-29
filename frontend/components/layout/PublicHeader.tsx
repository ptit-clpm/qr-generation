"use client";

import { useEffect } from "react";
import Link from "next/link";
import { LogOut, QrCode, User as UserIcon } from "lucide-react";
import { Button } from "@/components/common/Button";
import { useAuthStore } from "@/stores/auth";

export function PublicHeader() {
  const { user, loadMe, logout } = useAuthStore();

  useEffect(() => {
    if (typeof window !== "undefined") {
      const token = window.localStorage.getItem("access_token");
      if (token && !user) {
        loadMe().catch(() => {});
      }
    }
  }, [user, loadMe]);

  return (
    <header className="sticky top-0 z-20 border-b border-slate-200 bg-white/95 backdrop-blur">
      <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-4">
        <Link href="/" className="flex items-center gap-2 font-bold text-ink">
          <span className="grid h-9 w-9 place-items-center rounded-md bg-ink text-white">
            <QrCode className="h-5 w-5" />
          </span>
          QR Studio
        </Link>
        <nav className="hidden items-center gap-6 text-sm font-medium text-muted md:flex">
          <Link href="/qrcodes" className="hover:text-ink">Create QR</Link>
          <Link href="/pricing" className="hover:text-ink">Pricing</Link>
          <Link href="/dashboard" className="hover:text-ink">Dashboard</Link>
        </nav>
        <div className="flex items-center gap-2">
          {user ? (
            <div className="flex items-center gap-3">
              <Link
                href="/account"
                className="flex items-center gap-2 rounded-md border border-slate-200 bg-panel px-3 py-1.5 text-sm font-semibold text-ink hover:bg-slate-100"
              >
                <UserIcon className="h-4 w-4 text-teal" />
                <span className="max-w-[120px] truncate">{user.full_name || user.email}</span>
              </Link>
              <Button tone="secondary" onClick={logout} title="Logout">
                <LogOut className="h-4 w-4" />
                <span className="hidden sm:inline">Logout</span>
              </Button>
            </div>
          ) : (
            <>
              <Link href="/login">
                <Button tone="secondary">Login</Button>
              </Link>
              <Link href="/register">
                <Button>Register</Button>
              </Link>
            </>
          )}
        </div>
      </div>
    </header>
  );
}

