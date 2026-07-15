"use client";

import { useEffect, useState } from "react";
import { Save } from "lucide-react";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { Button } from "@/components/common/Button";
import { api, messageFromError } from "@/lib/api";
import type { ApiEnvelope, User } from "@/types";

export default function AccountPage() {
  const [user, setUser] = useState<User | null>(null);
  const [message, setMessage] = useState("");

  useEffect(() => {
    api.get<ApiEnvelope<User>>("/users/profile").then((res) => setUser(res.data.data ?? null)).catch((err) => setMessage(messageFromError(err)));
  }, []);

  async function save() {
    if (!user) return;
    try {
      const res = await api.put<ApiEnvelope<User>>("/users/profile", {
        full_name: user.full_name,
        phone_number: user.phone_number,
        avatar_url: user.avatar_url
      });
      setUser(res.data.data ?? user);
      setMessage("Profile updated");
    } catch (err) {
      setMessage(messageFromError(err));
    }
  }

  return (
    <DashboardShell>
      <h1 className="text-3xl font-bold text-ink">Account</h1>
      <section className="mt-6 max-w-2xl rounded-md border border-slate-200 bg-white p-5 shadow-soft">
        <div className="grid gap-4">
          <label className="text-sm font-medium text-ink">
            Full name
            <input value={user?.full_name ?? ""} onChange={(e) => setUser((old) => old ? { ...old, full_name: e.target.value } : old)} className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2" />
          </label>
          <label className="text-sm font-medium text-ink">
            Email
            <input value={user?.email ?? ""} readOnly className="mt-1 w-full rounded-md border border-slate-200 bg-panel px-3 py-2 text-muted" />
          </label>
          <label className="text-sm font-medium text-ink">
            Phone
            <input value={user?.phone_number ?? ""} onChange={(e) => setUser((old) => old ? { ...old, phone_number: e.target.value } : old)} className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2" />
          </label>
        </div>
        {message ? <p className="mt-4 rounded-md bg-teal/10 px-3 py-2 text-sm text-teal">{message}</p> : null}
        <Button className="mt-5" onClick={save}><Save className="h-4 w-4" />Save</Button>
      </section>
    </DashboardShell>
  );
}
