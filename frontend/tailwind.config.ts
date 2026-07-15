import type { Config } from "tailwindcss";

const config: Config = {
  content: ["./app/**/*.{ts,tsx}", "./components/**/*.{ts,tsx}", "./lib/**/*.{ts,tsx}", "./stores/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        ink: "#18212f",
        muted: "#657084",
        panel: "#f7f9fb",
        teal: "#0f9f8f",
        coral: "#e4564f",
        amber: "#d89b25"
      },
      boxShadow: {
        soft: "0 16px 44px rgba(24, 33, 47, 0.08)"
      }
    }
  },
  plugins: []
};

export default config;
