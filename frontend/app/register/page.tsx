import Link from "next/link";
import { RegisterForm } from "@/components/forms/AuthForms";
import { PublicHeader } from "@/components/layout/PublicHeader";

export default function RegisterPage() {
  return (
    <div className="min-h-screen bg-panel">
      <PublicHeader />
      <main className="px-4 py-12">
        <RegisterForm />
        <p className="mt-4 text-center text-sm text-muted">
          Already registered? <Link href="/login" className="font-semibold text-teal">Login</Link>
        </p>
      </main>
    </div>
  );
}
