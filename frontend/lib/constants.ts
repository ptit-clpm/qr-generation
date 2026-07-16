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

export const bankMap: Record<string, string> = {
  "970422": "MBBank (Ngân hàng TMCP Quân đội)",
  "970415": "VietinBank (Ngân hàng Công thương Việt Nam)",
  "970436": "Vietcombank (Ngân hàng Ngoại thương Việt Nam)",
  "970418": "BIDV (Ngân hàng Đầu tư và Phát triển Việt Nam)",
  "970405": "Agribank (Ngân hàng Nông nghiệp & Phát triển Nông thôn)",
  "970407": "Techcombank (Ngân hàng Kỹ thương Việt Nam)",
  "970416": "ACB (Ngân hàng Á Châu)",
  "970423": "TPBank (Ngân hàng Tiên Phong)",
  "970432": "VPBank (Ngân hàng Thịnh Vượng)",
  "970403": "Sacombank (Ngân hàng Sài Gòn Thương Tín)",
  "970441": "VIB (Ngân hàng Quốc tế)",
  "970437": "HDBank (Ngân hàng Phát triển TP.HCM)",
  "970425": "ABBANK (Ngân hàng An Bình)",
  "970412": "PVcomBank (Ngân hàng Đại chúng Việt Nam)",
  "970427": "VietBank (Ngân hàng Việt Nam Thương Tín)",
  "970443": "SHB (Ngân hàng Sài Gòn - Hà Nội)",
  "970439": "Shinhan Bank (Ngân hàng Shinhan Việt Nam)",
  "970426": "MSB (Ngân hàng Hàng Hải Việt Nam)",
  "970431": "Eximbank (Ngân hàng Xuất Nhập Khẩu Việt Nam)",
  "970421": "VRB (Ngân hàng Liên doanh Việt - Nga)",
  "970438": "BaoVietBank (Ngân hàng Bảo Việt)",
  "970440": "SeABank (Ngân hàng Đông Nam Á)",
  "970429": "SCB (Ngân hàng Sài Gòn)",
  "970433": "VietABank (Ngân hàng Việt Á)",
  "970448": "OCB (Ngân hàng Phương Đông)",
  "970414": "OceanBank (Ngân hàng Đại Dương)",
  "970428": "Nam A Bank (Ngân hàng Nam Á)",
  "970409": "Bac A Bank (Ngân hàng Bắc Á)",
  "970442": "GPBank (Ngân hàng Dầu khí Toàn cầu)",
  "970400": "Saigonbank (Ngân hàng Sài Gòn Công Thương)",
  "970454": "VietCapitalBank (Ngân hàng Bản Việt)",
  "970449": "LienVietPostBank (Ngân hàng Bưu điện Liên Việt)",
  "970452": "Kienlongbank (Ngân hàng Kiên Long)",
  "970419": "NCB (Ngân hàng Quốc Dân)",
};
