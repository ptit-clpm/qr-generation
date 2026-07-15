import nextPlugin from "@next/eslint-plugin-next";
import nextTypescript from "eslint-config-next/typescript";

export default [
  {
    ignores: [".next/**", "node_modules/**", "next-env.d.ts", "tsconfig.tsbuildinfo"]
  },
  ...nextTypescript,
  {
    files: ["**/*.{js,jsx,ts,tsx}"],
    plugins: {
      "@next/next": nextPlugin
    },
    rules: {
      ...nextPlugin.configs.recommended.rules,
      ...nextPlugin.configs["core-web-vitals"].rules
    }
  }
];
