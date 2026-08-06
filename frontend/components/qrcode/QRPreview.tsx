"use client";

import { QRCodeCanvas } from "qrcode.react";

export function QRPreview({ value, foreground = "#111827", background = "#FFFFFF", logoUrl }: { value: string; foreground?: string; background?: string; logoUrl?: string }) {
  return (
    <div className="flex aspect-square w-full max-w-sm items-center justify-center rounded-md border border-slate-200 bg-white p-6">
      <QRCodeCanvas value={value || "QR Studio"} size={260} fgColor={foreground} bgColor={background} level={logoUrl ? "H" : "M"} includeMargin imageSettings={logoUrl ? { src: logoUrl, height: 46, width: 46, excavate: true } : undefined} />
    </div>
  );
}
