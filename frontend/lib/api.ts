import axios from "axios";
import type { ApiEnvelope } from "@/types";

export const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080/api/v1",
  headers: { "Content-Type": "application/json" }
});

api.interceptors.request.use((config) => {
  if (typeof window !== "undefined") {
    const token = window.localStorage.getItem("access_token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const original = error.config;
    if (error.response?.status === 401 && !original?._retry && typeof window !== "undefined") {
      original._retry = true;
      const refreshToken = window.localStorage.getItem("refresh_token");
      if (refreshToken) {
        const res = await api.post<ApiEnvelope<{ access_token: string; refresh_token: string }>>("/auth/refresh", {
          refresh_token: refreshToken
        });
        if (res.data.data) {
          window.localStorage.setItem("access_token", res.data.data.access_token);
          window.localStorage.setItem("refresh_token", res.data.data.refresh_token);
          original.headers.Authorization = `Bearer ${res.data.data.access_token}`;
          return api(original);
        }
      }
    }
    return Promise.reject(error);
  }
);

export function messageFromError(error: unknown) {
  if (axios.isAxiosError(error)) {
    return (error.response?.data as { message?: string })?.message ?? error.message;
  }
  return "Unexpected error";
}
