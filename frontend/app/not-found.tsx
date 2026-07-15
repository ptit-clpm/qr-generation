import Link from "next/link";
import { PublicHeader } from "@/components/layout/PublicHeader";
import { Button } from "@/components/common/Button";

export default function NotFound() {
  return (
    <div className="min-h-screen bg-panel">
      <PublicHeader />
      <main className="mx-auto max-w-2xl px-4 py-20 text-center">
        <h1 className="text-3xl font-bold text-ink">QR not found</h1>
        <p className="mt-3 text-muted">The QR code is missing, disabled, or deleted.</p>
        <Link href="/" className="mt-6 inline-block"><Button>Back home</Button></Link>
      </main>
    </div>
  );
}
