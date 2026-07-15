import Link from "next/link";
import { LoginForm } from "@/components/forms/AuthForms";
import { PublicHeader } from "@/components/layout/PublicHeader";

export default function LoginPage() {
  return (
    <div className="min-h-screen bg-panel">
      <PublicHeader />
      <main className="px-4 py-12">
        <LoginForm />
        <p className="mt-4 text-center text-sm text-muted">
          No account? <Link href="/register" className="font-semibold text-teal">Register</Link>
        </p>
      </main>
    </div>
  );
}
