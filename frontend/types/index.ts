export type RoleName = "USER" | "ADMIN";
export type PlanName = "FREE" | "PRO";
export type PaymentStatus = "PENDING" | "SUCCESS" | "FAILED" | "CANCELLED" | "REFUNDED";
export type PaymentMethod = "VNPAY" | "MOMO" | "ZALOPAY" | "PAYPAL" | "STRIPE" | "BANK_TRANSFER" | "SEPAY";
export type QRType = "URL" | "TEXT" | "WIFI" | "VCARD" | "EMAIL" | "SMS" | "LOCATION" | "SOCIAL" | "PDF" | "MENU";
export type QRStatus = "ACTIVE" | "DISABLED" | "DELETED";

export interface Role {
  id: number;
  name: RoleName;
  description?: string;
}

export interface User {
  id: number;
  full_name: string;
  email: string;
  phone_number?: string;
  avatar_url?: string;
  status: "ACTIVE" | "LOCKED" | "DELETED";
  roles?: Role[];
}

export interface Plan {
  id: number;
  name: PlanName;
  price: number;
  duration_days: number;
  max_qr_codes: number;
  allow_dynamic_qr: boolean;
  allow_logo: boolean;
  allow_analytics: boolean;
  allow_svg_pdf_export: boolean;
  description: string;
}

export interface Subscription {
  id: number;
  user_id: number;
  plan_id: number;
  plan?: Plan;
  start_date: string;
  end_date: string;
  status: "PENDING" | "ACTIVE" | "EXPIRED" | "CANCELLED";
  auto_renew: boolean;
}

export interface QRDesign {
  foreground_color: string;
  background_color: string;
  logo_url?: string;
  size: number;
  error_correction_level: "L" | "M" | "Q" | "H";
}

export interface QRCode {
  id: number;
  title: string;
  qr_type: QRType;
  content: string;
  short_code?: string;
  is_dynamic: boolean;
  destination_url?: string;
  scan_count: number;
  status: QRStatus;
  design?: QRDesign;
  created_at: string;
}

export interface Payment {
  id: number;
  user_id: number;
  subscription_id?: number;
  amount: number;
  currency: string;
  payment_method: PaymentMethod;
  transaction_code: string;
  provider?: string;
  provider_ref?: string;
  status: PaymentStatus;
  paid_at?: string;
  created_at: string;
}

export interface PaymentInstructions {
  provider: string;
  bank_account: string;
  bank_name: string;
  account_name: string;
  amount: number;
  currency: string;
  transaction_code: string;
  transfer_content: string;
  return_url: string;
  enabled: boolean;
}

export interface CreatePaymentResponse {
  payment: Payment;
  instructions: PaymentInstructions;
}

export interface ApiEnvelope<T> {
  success: boolean;
  message: string;
  data?: T;
  errors?: unknown;
}
