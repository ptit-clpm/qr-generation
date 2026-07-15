import type { QRType } from "@/types";

export const qrTypes: Array<{ value: QRType; label: string; pro?: boolean }> = [
  { value: "URL", label: "URL" },
  { value: "TEXT", label: "Text" },
  { value: "WIFI", label: "WiFi" },
  { value: "VCARD", label: "vCard" },
  { value: "EMAIL", label: "Email" },
  { value: "SMS", label: "SMS" },
  { value: "LOCATION", label: "Location" },
  { value: "SOCIAL", label: "Social", pro: true },
  { value: "PDF", label: "PDF", pro: true },
  { value: "MENU", label: "Menu", pro: true }
];

export const backendUrl = process.env.NEXT_PUBLIC_BACKEND_URL ?? "http://localhost:8080";
