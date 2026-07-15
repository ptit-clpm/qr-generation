import Link from "next/link";
import { QrCode } from "lucide-react";
import { Button } from "@/components/common/Button";

export function PublicHeader() {
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
          <Link href="/qrcodes">Create QR</Link>
          <Link href="/pricing">Pricing</Link>
          <Link href="/dashboard">Dashboard</Link>
        </nav>
        <div className="flex items-center gap-2">
          <Link href="/login">
            <Button tone="secondary">Login</Button>
          </Link>
          <Link href="/register">
            <Button>Register</Button>
          </Link>
        </div>
      </div>
    </header>
  );
}
