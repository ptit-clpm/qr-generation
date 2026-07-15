"use client";

import { create } from "zustand";
import { api } from "@/lib/api";
import type { ApiEnvelope, RoleName, User } from "@/types";

interface AuthState {
  user: User | null;
  loading: boolean;
  setUser: (user: User | null) => void;
  hasRole: (role: RoleName) => boolean;
  isAdmin: () => boolean;
  loadMe: () => Promise<void>;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  loading: false,
  setUser: (user) => set({ user }),
  hasRole: (role) => get().user?.roles?.some((userRole) => userRole.name === role) ?? false,
  isAdmin: () => get().hasRole("ADMIN"),
  loadMe: async () => {
    set({ loading: true });
    try {
      const res = await api.get<ApiEnvelope<User>>("/auth/me");
      set({ user: res.data.data ?? null });
    } catch (error) {
      if (typeof window !== "undefined") {
        window.localStorage.removeItem("access_token");
        window.localStorage.removeItem("refresh_token");
      }
      set({ user: null });
      throw error;
    } finally {
      set({ loading: false });
    }
  },
  logout: () => {
    window.localStorage.removeItem("access_token");
    window.localStorage.removeItem("refresh_token");
    set({ user: null });
  }
}));
