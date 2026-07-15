"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { Button } from "@/components/common/Button";
import { api, messageFromError } from "@/lib/api";
import { loginSchema, registerSchema } from "@/lib/validators";
import { useAuthStore } from "@/stores/auth";
import type { ApiEnvelope, User } from "@/types";

type AuthPayload = { access_token: string; refresh_token: string; user: User };

export function LoginForm() {
  const router = useRouter();
  const setUser = useAuthStore((state) => state.setUser);
  const [error, setError] = useState("");
  const form = useForm<z.infer<typeof loginSchema>>({ resolver: zodResolver(loginSchema), defaultValues: { email: "", password: "" } });

  async function onSubmit(values: z.infer<typeof loginSchema>) {
    setError("");
    try {
      const res = await api.post<ApiEnvelope<AuthPayload>>("/auth/login", values);
      const data = res.data.data;
      if (data) {
        if (data.user.roles?.some((role) => role.name === "ADMIN")) {
          localStorage.removeItem("access_token");
          localStorage.removeItem("refresh_token");
          setUser(null);
          setError("Tài khoản admin cần đăng nhập vào trang quản trị riêng.");
          return;
        }
        localStorage.setItem("access_token", data.access_token);
        localStorage.setItem("refresh_token", data.refresh_token);
        setUser(data.user);
        router.push("/dashboard");
      }
    } catch (err) {
      setError(messageFromError(err));
    }
  }

  return <AuthFormShell title="Login" error={error} loading={form.formState.isSubmitting} onSubmit={form.handleSubmit(onSubmit)} fields={[
    { label: "Email", type: "email", props: form.register("email") },
    { label: "Password", type: "password", props: form.register("password") }
  ]} />;
}

export function RegisterForm() {
  const router = useRouter();
  const setUser = useAuthStore((state) => state.setUser);
  const [error, setError] = useState("");
  const form = useForm<z.infer<typeof registerSchema>>({
    resolver: zodResolver(registerSchema),
    defaultValues: { full_name: "", email: "", phone_number: "", password: "", confirm_password: "" }
  });

  async function onSubmit(values: z.infer<typeof registerSchema>) {
    setError("");
    try {
      const res = await api.post<ApiEnvelope<AuthPayload>>("/auth/register", values);
      const data = res.data.data;
      if (data) {
        localStorage.setItem("access_token", data.access_token);
        localStorage.setItem("refresh_token", data.refresh_token);
        setUser(data.user);
        router.push("/dashboard");
      }
    } catch (err) {
      setError(messageFromError(err));
    }
  }

  return <AuthFormShell title="Register" error={error} loading={form.formState.isSubmitting} onSubmit={form.handleSubmit(onSubmit)} fields={[
    { label: "Full name", type: "text", props: form.register("full_name") },
    { label: "Email", type: "email", props: form.register("email") },
    { label: "Phone", type: "text", props: form.register("phone_number") },
    { label: "Password", type: "password", props: form.register("password") },
    { label: "Confirm password", type: "password", props: form.register("confirm_password") }
  ]} />;
}

function AuthFormShell({
  title,
  fields,
  error,
  loading,
  onSubmit
}: {
  title: string;
  fields: Array<{ label: string; type: string; props: object }>;
  error: string;
  loading: boolean;
  onSubmit: () => void;
}) {
  return (
    <form onSubmit={onSubmit} className="mx-auto w-full max-w-md rounded-md border border-slate-200 bg-white p-6 shadow-soft">
      <h1 className="text-2xl font-bold text-ink">{title}</h1>
      <div className="mt-6 space-y-4">
        {fields.map((field) => (
          <label key={field.label} className="block text-sm font-medium text-ink">
            {field.label}
            <input className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2" type={field.type} {...field.props} />
          </label>
        ))}
      </div>
      {error ? <p className="mt-4 rounded-md bg-coral/10 px-3 py-2 text-sm text-coral">{error}</p> : null}
      <Button className="mt-6 w-full" disabled={loading}>{loading ? "Submitting" : title}</Button>
    </form>
  );
}
