import { z } from "zod";

export const loginSchema = z.object({
  email: z.string().email(),
  password: z.string().min(1)
});

export const registerSchema = z.object({
  full_name: z.string().min(2).max(150),
  email: z.string().email(),
  phone_number: z.string().optional(),
  password: z.string().min(8),
  confirm_password: z.string().min(8)
}).refine((value) => value.password === value.confirm_password, {
  path: ["confirm_password"],
  message: "Passwords do not match"
});

export const qrSchema = z.object({
  title: z.string().min(2).max(150),
  qr_type: z.string(),
  content: z.string().min(1),
  is_dynamic: z.boolean().default(false),
  destination_url: z.string().optional(),
  foreground_color: z.string().default("#111827"),
  background_color: z.string().default("#FFFFFF")
});
