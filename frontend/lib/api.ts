import axios from "axios";
import type { ApiEnvelope } from "@/types";

const baseURL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080/api/v1";

export const api = axios.create({
  baseURL,
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
    // Don't intercept refresh token calls or retried requests
    if (
      error.response?.status === 401 &&
      !original?._retry &&
      !original?.url?.includes("/auth/refresh") &&
      typeof window !== "undefined"
    ) {
      original._retry = true;
      const refreshToken = window.localStorage.getItem("refresh_token");
      if (refreshToken) {
        try {
          // Use standard axios call to bypass interceptor loop
          const res = await axios.post<ApiEnvelope<{ access_token: string; refresh_token: string }>>(
            `${baseURL}/auth/refresh`,
            { refresh_token: refreshToken }
          );
          if (res.data.data) {
            window.localStorage.setItem("access_token", res.data.data.access_token);
            window.localStorage.setItem("refresh_token", res.data.data.refresh_token);
            original.headers.Authorization = `Bearer ${res.data.data.access_token}`;
            return api(original);
          }
        } catch {
          // Refresh failed - purge invalid tokens
          window.localStorage.removeItem("access_token");
          window.localStorage.removeItem("refresh_token");
        }
      } else {
        window.localStorage.removeItem("access_token");
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
