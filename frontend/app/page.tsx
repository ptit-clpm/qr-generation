import Link from "next/link";
import { BarChart3, Palette, QrCode, Zap } from "lucide-react";
import { PublicHeader } from "@/components/layout/PublicHeader";
import { Button } from "@/components/common/Button";
import { CreateQRForm } from "@/components/qrcode/CreateQRForm";

export default function HomePage() {
  return (
    <div className="min-h-screen bg-white">
      <PublicHeader />
      <main>
        <section className="border-b border-slate-200">
          <div className="mx-auto grid max-w-6xl gap-10 px-4 py-16 lg:grid-cols-[0.9fr_1.1fr] lg:items-center">
            <div>
              <p className="font-semibold text-teal">QR Generator - QR Studio</p>
              <h1 className="mt-3 text-4xl font-bold tracking-normal text-ink md:text-6xl">QR Studio</h1>
              <p className="mt-5 max-w-xl text-lg leading-8 text-muted">
                Tạo, quản lý và đo lường hiệu quả mã QR tĩnh hoặc động cho cá nhân, cửa hàng và chiến dịch marketing.
              </p>
              <div className="mt-7 flex flex-wrap gap-3">
                <Link href="/qrcodes">
                  <Button><QrCode className="h-4 w-4" />Create QR</Button>
                </Link>
                <Link href="/pricing">
                  <Button tone="secondary">Pricing</Button>
                </Link>
              </div>
            </div>
            <div className="rounded-md border border-slate-200 bg-panel p-4 shadow-soft">
              <div className="grid gap-3 sm:grid-cols-2">
                {[
                  ["URL", "https://qr-studio.local"],
                  ["WiFi", "SSID, WPA, password"],
                  ["vCard", "Name, phone, email"],
                  ["Dynamic", "/q/shortCode"]
                ].map(([label, value]) => (
                  <div key={label} className="rounded-md bg-white p-4">
                    <p className="text-sm font-semibold text-ink">{label}</p>
                    <p className="mt-2 text-sm text-muted">{value}</p>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </section>
        <section className="mx-auto max-w-6xl px-4 py-12">
          <div className="grid gap-4 md:grid-cols-3">
            {[
              { icon: Zap, title: "Dynamic QR", text: "Đổi URL đích sau khi in mã." },
              { icon: BarChart3, title: "Analytics", text: "Thống kê lượt quét theo ngày và thiết bị." },
              { icon: Palette, title: "Design", text: "Màu sắc, kích thước, logo và template." }
            ].map((item) => {
              const Icon = item.icon;
              return (
                <div key={item.title} className="rounded-md border border-slate-200 bg-white p-5">
                  <Icon className="h-5 w-5 text-teal" />
                  <h2 className="mt-4 font-semibold text-ink">{item.title}</h2>
                  <p className="mt-2 text-sm leading-6 text-muted">{item.text}</p>
                </div>
              );
            })}
          </div>
        </section>
        <section className="border-t border-slate-200 bg-panel px-4 py-12">
          <div className="mx-auto max-w-6xl">
            <h2 className="mb-5 text-2xl font-bold text-ink">Create QR</h2>
            <CreateQRForm />
          </div>
        </section>
      </main>
    </div>
  );
}
