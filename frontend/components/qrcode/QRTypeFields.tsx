"use client";

import { useEffect, useState } from "react";
import type { QRType } from "@/types";

interface QRTypeFieldsProps {
  qrType: QRType;
  onChangeContent: (content: string) => void;
  isProUser: boolean;
}

export function QRTypeFields({ qrType, onChangeContent, isProUser }: QRTypeFieldsProps) {
  // Common states
  const [url, setUrl] = useState("https://example.com");
  const [text, setText] = useState("Hello world");

  // WiFi states
  const [wifiSsid, setWifiSsid] = useState("");
  const [wifiPassword, setWifiPassword] = useState("");
  const [wifiEncryption, setWifiEncryption] = useState<"WPA" | "WEP" | "nopass">("WPA");
  const [wifiHidden, setWifiHidden] = useState(false);

  // vCard states
  const [vFirstName, setVFirstName] = useState("");
  const [vLastName, setVLastName] = useState("");
  const [vPhone, setVPhone] = useState("");
  const [vEmail, setVEmail] = useState("");
  const [vCompany, setVCompany] = useState("");
  const [vTitle, setVTitle] = useState("");

  // Email states
  const [emailTo, setEmailTo] = useState("");
  const [emailSubject, setEmailSubject] = useState("");
  const [emailBody, setEmailBody] = useState("");

  // SMS states
  const [smsPhone, setSmsPhone] = useState("");
  const [smsMessage, setSmsMessage] = useState("");

  // Location states
  const [lat, setLat] = useState("21.028511");
  const [lng, setLng] = useState("105.804817");

  // Social / PDF / Menu states
  const [proUrl, setProUrl] = useState("https://example.com/pro-content");

  useEffect(() => {
    let formatted = "";

    switch (qrType) {
      case "URL":
        formatted = url;
        break;
      case "TEXT":
        formatted = text;
        break;
      case "WIFI":
        formatted = `WIFI:T:${wifiEncryption};S:${wifiSsid};P:${wifiPassword};H:${wifiHidden ? "true" : "false"};;`;
        break;
      case "VCARD":
        formatted = [
          "BEGIN:VCARD",
          "VERSION:3.0",
          `N:${vLastName};${vFirstName}`,
          `FN:${vFirstName} ${vLastName}`.trim(),
          vPhone ? `TEL:${vPhone}` : "",
          vEmail ? `EMAIL:${vEmail}` : "",
          vCompany ? `ORG:${vCompany}` : "",
          vTitle ? `TITLE:${vTitle}` : "",
          "END:VCARD"
        ]
          .filter(Boolean)
          .join("\n");
        break;
      case "EMAIL":
        formatted = `mailto:${emailTo}?subject=${encodeURIComponent(emailSubject)}&body=${encodeURIComponent(emailBody)}`;
        break;
      case "SMS":
        formatted = `smsto:${smsPhone}:${smsMessage}`;
        break;
      case "LOCATION":
        formatted = `geo:${lat},${lng}`;
        break;
      case "SOCIAL":
      case "PDF":
      case "MENU":
        formatted = proUrl;
        break;
    }

    onChangeContent(formatted);
  }, [
    qrType,
    url,
    text,
    wifiSsid,
    wifiPassword,
    wifiEncryption,
    wifiHidden,
    vFirstName,
    vLastName,
    vPhone,
    vEmail,
    vCompany,
    vTitle,
    emailTo,
    emailSubject,
    emailBody,
    smsPhone,
    smsMessage,
    lat,
    lng,
    proUrl,
    onChangeContent
  ]);

  const isProType = qrType === "SOCIAL" || qrType === "PDF" || qrType === "MENU";

  if (isProType && !isProUser) {
    return (
      <div className="md:col-span-2 rounded-md bg-amber-500/10 border border-amber-500/20 p-4 text-amber-800">
        <p className="text-sm font-semibold">Gói Pro yêu cầu</p>
        <p className="mt-1 text-xs">
          Loại QR <span className="font-bold">{qrType}</span> là tính năng dành riêng cho tài khoản Pro. Hãy nâng cấp để sử dụng.
        </p>
      </div>
    );
  }

  switch (qrType) {
    case "URL":
      return (
        <label className="md:col-span-2 text-sm font-medium text-ink">
          Website URL
          <input
            type="url"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            placeholder="https://example.com"
            className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
          />
        </label>
      );

    case "TEXT":
      return (
        <label className="md:col-span-2 text-sm font-medium text-ink">
          Text Content
          <textarea
            value={text}
            onChange={(e) => setText(e.target.value)}
            rows={4}
            placeholder="Nhập văn bản bất kỳ..."
            className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
          />
        </label>
      );

    case "WIFI":
      return (
        <div className="md:col-span-2 grid gap-3 sm:grid-cols-2">
          <label className="text-sm font-medium text-ink">
            Network Name (SSID)
            <input
              value={wifiSsid}
              onChange={(e) => setWifiSsid(e.target.value)}
              placeholder="Tên Wi-Fi"
              className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
            />
          </label>
          <label className="text-sm font-medium text-ink">
            Password
            <input
              type="password"
              value={wifiPassword}
              onChange={(e) => setWifiPassword(e.target.value)}
              placeholder="Mật khẩu Wi-Fi"
              className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
            />
          </label>
          <label className="text-sm font-medium text-ink">
            Encryption Type
            <select
              value={wifiEncryption}
              onChange={(e) => setWifiEncryption(e.target.value as "WPA" | "WEP" | "nopass")}
              className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
            >
              <option value="WPA">WPA / WPA2 / WPA3</option>
              <option value="WEP">WEP</option>
              <option value="nopass">None (Không mật khẩu)</option>
            </select>
          </label>
          <label className="flex items-center gap-2 text-sm font-medium text-ink sm:self-end sm:mb-2">
            <input
              type="checkbox"
              checked={wifiHidden}
              onChange={(e) => setWifiHidden(e.target.checked)}
            />
            Mạng Wi-Fi ẩn (Hidden)
          </label>
        </div>
      );

    case "VCARD":
      return (
        <div className="md:col-span-2 grid gap-3 sm:grid-cols-2">
          <label className="text-sm font-medium text-ink">
            First Name (Tên)
            <input
              value={vFirstName}
              onChange={(e) => setVFirstName(e.target.value)}
              placeholder="Văn A"
              className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
            />
          </label>
          <label className="text-sm font-medium text-ink">
            Last Name (Họ)
            <input
              value={vLastName}
              onChange={(e) => setVLastName(e.target.value)}
              placeholder="Nguyễn"
              className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
            />
          </label>
          <label className="text-sm font-medium text-ink">
            Phone Number
            <input
              value={vPhone}
              onChange={(e) => setVPhone(e.target.value)}
              placeholder="0912345678"
              className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
            />
          </label>
          <label className="text-sm font-medium text-ink">
            Email
            <input
              type="email"
              value={vEmail}
              onChange={(e) => setVEmail(e.target.value)}
              placeholder="contact@example.com"
              className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
            />
          </label>
          <label className="text-sm font-medium text-ink">
            Company / Organization
            <input
              value={vCompany}
              onChange={(e) => setVCompany(e.target.value)}
              placeholder="Công ty ABC"
              className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
            />
          </label>
          <label className="text-sm font-medium text-ink">
            Title / Position
            <input
              value={vTitle}
              onChange={(e) => setVTitle(e.target.value)}
              placeholder="Software Engineer"
              className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
            />
          </label>
        </div>
      );

    case "EMAIL":
      return (
        <div className="md:col-span-2 grid gap-3">
          <label className="text-sm font-medium text-ink">
            Recipient Email
            <input
              type="email"
              value={emailTo}
              onChange={(e) => setEmailTo(e.target.value)}
              placeholder="support@example.com"
              className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
            />
          </label>
          <label className="text-sm font-medium text-ink">
            Subject
            <input
              value={emailSubject}
              onChange={(e) => setEmailSubject(e.target.value)}
              placeholder="Tiêu đề email"
              className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
            />
          </label>
          <label className="text-sm font-medium text-ink">
            Body
            <textarea
              value={emailBody}
              onChange={(e) => setEmailBody(e.target.value)}
              rows={3}
              placeholder="Nội dung email..."
              className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
            />
          </label>
        </div>
      );

    case "SMS":
      return (
        <div className="md:col-span-2 grid gap-3">
          <label className="text-sm font-medium text-ink">
            Phone Number
            <input
              value={smsPhone}
              onChange={(e) => setSmsPhone(e.target.value)}
              placeholder="0912345678"
              className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
            />
          </label>
          <label className="text-sm font-medium text-ink">
            Message Content
            <textarea
              value={smsMessage}
              onChange={(e) => setSmsMessage(e.target.value)}
              rows={3}
              placeholder="Nội dung tin nhắn..."
              className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
            />
          </label>
        </div>
      );

    case "LOCATION":
      return (
        <div className="md:col-span-2 grid gap-3 sm:grid-cols-2">
          <label className="text-sm font-medium text-ink">
            Latitude (Vĩ độ)
            <input
              value={lat}
              onChange={(e) => setLat(e.target.value)}
              placeholder="21.028511"
              className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
            />
          </label>
          <label className="text-sm font-medium text-ink">
            Longitude (Kinh độ)
            <input
              value={lng}
              onChange={(e) => setLng(e.target.value)}
              placeholder="105.804817"
              className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
            />
          </label>
        </div>
      );

    case "SOCIAL":
      return (
        <label className="md:col-span-2 text-sm font-medium text-ink">
          Social Profile URL
          <input
            type="url"
            value={proUrl}
            onChange={(e) => setProUrl(e.target.value)}
            placeholder="https://instagram.com/yourprofile"
            className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
          />
        </label>
      );

    case "PDF":
      return (
        <label className="md:col-span-2 text-sm font-medium text-ink">
          PDF Document URL
          <input
            type="url"
            value={proUrl}
            onChange={(e) => setProUrl(e.target.value)}
            placeholder="https://example.com/document.pdf"
            className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
          />
        </label>
      );

    case "MENU":
      return (
        <label className="md:col-span-2 text-sm font-medium text-ink">
          Menu URL
          <input
            type="url"
            value={proUrl}
            onChange={(e) => setProUrl(e.target.value)}
            placeholder="https://example.com/menu"
            className="focus-ring mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
          />
        </label>
      );

    default:
      return null;
  }
}
