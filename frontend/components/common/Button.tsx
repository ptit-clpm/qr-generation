import { clsx } from "clsx";
import type { ButtonHTMLAttributes } from "react";

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  tone?: "primary" | "secondary" | "danger";
};

export function Button({ className, tone = "primary", ...props }: ButtonProps) {
  return (
    <button
      className={clsx(
        "focus-ring inline-flex min-h-10 items-center justify-center gap-2 rounded-md px-4 py-2 text-sm font-semibold transition disabled:cursor-not-allowed disabled:opacity-60",
        tone === "primary" && "bg-teal text-white hover:bg-teal/90",
        tone === "secondary" && "border border-slate-200 bg-white text-ink hover:bg-panel",
        tone === "danger" && "bg-coral text-white hover:bg-coral/90",
        className
      )}
      {...props}
    />
  );
}
