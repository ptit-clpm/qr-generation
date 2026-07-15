"use client";

import { QRCodeCanvas } from "qrcode.react";

export function QRPreview({ value, foreground = "#111827", background = "#FFFFFF" }: { value: string; foreground?: string; background?: string }) {
  return (
    <div className="flex aspect-square w-full max-w-sm items-center justify-center rounded-md border border-slate-200 bg-white p-6">
      <QRCodeCanvas value={value || "QR Studio"} size={260} fgColor={foreground} bgColor={background} includeMargin />
    </div>
  );
}
