import { useEffect, useMemo, useState, type FormEvent, type ReactNode } from "react";
import {
  BarChart3,
  CreditCard,
  LayoutDashboard,
  ListChecks,
  LogOut,
  QrCode,
  RefreshCw,
  Search,
  Shield,
  Users
} from "lucide-react";

type RoleName = "USER" | "ADMIN";
type UserStatus = "ACTIVE" | "LOCKED" | "DELETED";
type QRStatus = "ACTIVE" | "DISABLED" | "DELETED";
type PaymentStatus = "PENDING" | "SUCCESS" | "FAILED" | "CANCELLED" | "REFUNDED";
type PlanStatus = "ACTIVE" | "INACTIVE" | "DELETED";

interface ApiEnvelope<T> {
  success: boolean;
  message: string;
  data?: T;
}

interface Role {
  id: number;
  name: RoleName;
}

interface User {
  id: number;
  full_name: string;
  email: string;
  phone_number?: string;
  status: UserStatus;
  roles?: Role[];
  created_at?: string;
}

interface AdminDashboard {
  users: number;
  qrcodes: number;
  scans: number;
  successful_payments: number;
  revenue: number;
}

interface QRCodeItem {
  id: number;
  user_id: number;
  title: string;
  qr_type: string;
  content: string;
  short_code?: string;
  is_dynamic: boolean;
  scan_count: number;
  status: QRStatus;
  created_at?: string;
}

interface Payment {
  id: number;
  user_id: number;
  amount: number;
  currency: string;
  payment_method: string;
  transaction_code: string;
  status: PaymentStatus;
  paid_at?: string;
  created_at?: string;
}

interface Plan {
  id: number;
  name: "FREE" | "PRO";
  price: number;
  duration_days: number;
  max_qr_codes: number;
  allow_dynamic_qr: boolean;
  allow_logo: boolean;
  allow_analytics: boolean;
  allow_svg_pdf_export: boolean;
  description: string;
  status: PlanStatus;
}

interface SystemLog {
  id: number;
  action: string;
  entity_type?: string;
  level: string;
  message: string;
  ip_address?: string;
  created_at?: string;
}

interface AuthPayload {
  access_token: string;
  refresh_token: string;
  user: User;
}

type TabKey = "dashboard" | "users" | "qrcodes" | "payments" | "plans" | "logs";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080/api/v1";

const tabs: Array<{ key: TabKey; label: string; icon: typeof LayoutDashboard }> = [
  { key: "dashboard", label: "Tổng quan", icon: LayoutDashboard },
  { key: "users", label: "Người dùng", icon: Users },
  { key: "qrcodes", label: "QR toàn hệ thống", icon: QrCode },
  { key: "payments", label: "Doanh thu", icon: CreditCard },
  { key: "plans", label: "Gói dịch vụ", icon: ListChecks },
  { key: "logs", label: "Nhật ký", icon: BarChart3 }
];

function App() {
  const [token, setToken] = useState(() => localStorage.getItem("admin_access_token") ?? "");
  const [activeTab, setActiveTab] = useState<TabKey>("dashboard");
  const [dashboard, setDashboard] = useState<AdminDashboard | null>(null);
  const [users, setUsers] = useState<User[]>([]);
  const [qrcodes, setQrcodes] = useState<QRCodeItem[]>([]);
  const [payments, setPayments] = useState<Payment[]>([]);
  const [plans, setPlans] = useState<Plan[]>([]);
  const [logs, setLogs] = useState<SystemLog[]>([]);
  const [query, setQuery] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const api = useMemo(() => createApi(token, () => {
    localStorage.removeItem("admin_access_token");
    localStorage.removeItem("admin_refresh_token");
    setToken("");
  }), [token]);

  async function loadAll() {
    if (!token) return;
    setLoading(true);
    setError("");
    try {
      const [dashboardData, usersData, qrsData, paymentsData, plansData, logsData] = await Promise.all([
        api.get<AdminDashboard>("/admin/dashboard"),
        api.get<User[]>("/admin/users"),
        api.get<QRCodeItem[]>("/admin/qrcodes"),
        api.get<Payment[]>("/admin/payments"),
        api.get<Plan[]>("/admin/plans"),
        api.get<SystemLog[]>("/admin/logs")
      ]);
      setDashboard(dashboardData);
      setUsers(usersData);
      setQrcodes(qrsData);
      setPayments(paymentsData);
      setPlans(plansData);
      setLogs(logsData);
    } catch (err) {
      setError(messageFromError(err));
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadAll();
  }, [token]);

  function logout() {
    localStorage.removeItem("admin_access_token");
    localStorage.removeItem("admin_refresh_token");
    setToken("");
  }

  async function updateUserStatus(user: User, status: UserStatus) {
    await api.put(`/admin/users/${user.id}/status`, { status });
    setUsers((items) => items.map((item) => item.id === user.id ? { ...item, status } : item));
  }

  async function updateQRStatus(qr: QRCodeItem, status: QRStatus) {
    await api.put(`/admin/qrcodes/${qr.id}/status`, { status });
    setQrcodes((items) => items.map((item) => item.id === qr.id ? { ...item, status } : item));
  }

  if (!token) {
    return <LoginScreen onLogin={setToken} />;
  }

  const visibleUsers = filterItems(users, query, (item) => [item.full_name, item.email, item.status]);
  const visibleQRCodes = filterItems(qrcodes, query, (item) => [item.title, item.qr_type, item.status, String(item.user_id)]);
  const visiblePayments = filterItems(payments, query, (item) => [item.transaction_code, item.status, String(item.user_id)]);
  const visibleLogs = filterItems(logs, query, (item) => [item.action, item.entity_type ?? "", item.level, item.message]);

  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="brand">
          <span className="brand-mark"><Shield size={20} /></span>
          <span>QR Admin</span>
        </div>
        <nav className="nav-list">
          {tabs.map((tab) => {
            const Icon = tab.icon;
            return (
              <button key={tab.key} className={activeTab === tab.key ? "nav-item active" : "nav-item"} onClick={() => setActiveTab(tab.key)}>
                <Icon size={18} />
                {tab.label}
              </button>
            );
          })}
        </nav>
        <button className="nav-item logout" onClick={logout}>
          <LogOut size={18} />
          Đăng xuất
        </button>
      </aside>

      <main className="main">
        <header className="topbar">
          <div>
            <p className="eyebrow">Quản trị toàn cục</p>
            <h1>{tabs.find((tab) => tab.key === activeTab)?.label}</h1>
          </div>
          <div className="topbar-actions">
            <label className="search">
              <Search size={17} />
              <input value={query} onChange={(event) => setQuery(event.target.value)} placeholder="Tìm kiếm" />
            </label>
            <button className="icon-button" onClick={loadAll} disabled={loading} title="Tải lại dữ liệu">
              <RefreshCw size={18} />
            </button>
          </div>
        </header>

        {error ? <p className="alert">{error}</p> : null}

        {activeTab === "dashboard" ? <DashboardView dashboard={dashboard} payments={payments} qrcodes={qrcodes} /> : null}
        {activeTab === "users" ? <UsersView users={visibleUsers} onStatusChange={updateUserStatus} /> : null}
        {activeTab === "qrcodes" ? <QRCodesView qrcodes={visibleQRCodes} onStatusChange={updateQRStatus} /> : null}
        {activeTab === "payments" ? <PaymentsView payments={visiblePayments} /> : null}
        {activeTab === "plans" ? <PlansView plans={plans} /> : null}
        {activeTab === "logs" ? <LogsView logs={visibleLogs} /> : null}
      </main>
    </div>
  );
}

function LoginScreen({ onLogin }: { onLogin: (token: string) => void }) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function submit(event: FormEvent) {
    event.preventDefault();
    setLoading(true);
    setError("");
    try {
      const res = await request<AuthPayload>("/auth/login", {
        method: "POST",
        body: JSON.stringify({ email, password })
      });
      const isAdmin = res.user.roles?.some((role) => role.name === "ADMIN");
      if (!isAdmin) {
        setError("Chỉ tài khoản admin được đăng nhập trang quản trị.");
        return;
      }
      localStorage.setItem("admin_access_token", res.access_token);
      localStorage.setItem("admin_refresh_token", res.refresh_token);
      onLogin(res.access_token);
    } catch (err) {
      setError(messageFromError(err));
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="login-page">
      <form className="login-panel" onSubmit={submit}>
        <div className="brand login-brand">
          <span className="brand-mark"><Shield size={20} /></span>
          <span>QR Admin</span>
        </div>
        <h1>Đăng nhập quản trị</h1>
        <label>
          Email
          <input type="email" value={email} onChange={(event) => setEmail(event.target.value)} required />
        </label>
        <label>
          Mật khẩu
          <input type="password" value={password} onChange={(event) => setPassword(event.target.value)} required />
        </label>
        {error ? <p className="alert">{error}</p> : null}
        <button className="primary-button" disabled={loading}>{loading ? "Đang đăng nhập" : "Đăng nhập"}</button>
      </form>
    </main>
  );
}

function DashboardView({ dashboard, payments, qrcodes }: { dashboard: AdminDashboard | null; payments: Payment[]; qrcodes: QRCodeItem[] }) {
  const pendingPayments = payments.filter((payment) => payment.status === "PENDING").length;
  const activeQRCodes = qrcodes.filter((qr) => qr.status === "ACTIVE").length;
  return (
    <div className="stack">
      <div className="metric-grid">
        <Metric label="Người dùng" value={dashboard?.users ?? 0} icon={Users} />
        <Metric label="QR đang quản lý" value={dashboard?.qrcodes ?? 0} icon={QrCode} />
        <Metric label="Lượt quét" value={dashboard?.scans ?? 0} icon={BarChart3} />
        <Metric label="Thanh toán thành công" value={dashboard?.successful_payments ?? 0} icon={CreditCard} />
      </div>
      <section className="panel split-panel">
        <div>
          <p className="panel-label">Tổng doanh thu</p>
          <strong>{formatMoney(dashboard?.revenue ?? 0)}</strong>
        </div>
        <div>
          <p className="panel-label">Thanh toán chờ xử lý</p>
          <strong>{pendingPayments}</strong>
        </div>
        <div>
          <p className="panel-label">QR đang hoạt động</p>
          <strong>{activeQRCodes}</strong>
        </div>
      </section>
    </div>
  );
}

function UsersView({ users, onStatusChange }: { users: User[]; onStatusChange: (user: User, status: UserStatus) => Promise<void> }) {
  return (
    <Table
      headers={["ID", "Tên", "Email", "Vai trò", "Trạng thái", "Thao tác"]}
      rows={users.map((user) => [
        user.id,
        user.full_name,
        user.email,
        user.roles?.map((role) => role.name).join(", ") || "USER",
        <StatusBadge status={user.status} key="status" />,
        <select key="action" value={user.status} onChange={(event) => onStatusChange(user, event.target.value as UserStatus)}>
          <option value="ACTIVE">ACTIVE</option>
          <option value="LOCKED">LOCKED</option>
          <option value="DELETED">DELETED</option>
        </select>
      ])}
    />
  );
}

function QRCodesView({ qrcodes, onStatusChange }: { qrcodes: QRCodeItem[]; onStatusChange: (qr: QRCodeItem, status: QRStatus) => Promise<void> }) {
  return (
    <Table
      headers={["ID", "User", "Tiêu đề", "Loại", "Dynamic", "Scan", "Trạng thái", "Thao tác"]}
      rows={qrcodes.map((qr) => [
        qr.id,
        qr.user_id,
        qr.title,
        qr.qr_type,
        qr.is_dynamic ? "Có" : "Không",
        qr.scan_count,
        <StatusBadge status={qr.status} key="status" />,
        <select key="action" value={qr.status} onChange={(event) => onStatusChange(qr, event.target.value as QRStatus)}>
          <option value="ACTIVE">ACTIVE</option>
          <option value="DISABLED">DISABLED</option>
          <option value="DELETED">DELETED</option>
        </select>
      ])}
    />
  );
}

function PaymentsView({ payments }: { payments: Payment[] }) {
  const revenue = payments.filter((payment) => payment.status === "SUCCESS").reduce((sum, payment) => sum + payment.amount, 0);
  return (
    <div className="stack">
      <section className="panel split-panel">
        <div>
          <p className="panel-label">Doanh thu đã ghi nhận</p>
          <strong>{formatMoney(revenue)}</strong>
        </div>
        <div>
          <p className="panel-label">Giao dịch</p>
          <strong>{payments.length}</strong>
        </div>
      </section>
      <Table
        headers={["ID", "User", "Mã giao dịch", "Số tiền", "Phương thức", "Trạng thái", "Ngày tạo"]}
        rows={payments.map((payment) => [
          payment.id,
          payment.user_id,
          payment.transaction_code,
          formatMoney(payment.amount),
          payment.payment_method,
          <StatusBadge status={payment.status} key="status" />,
          formatDate(payment.created_at)
        ])}
      />
    </div>
  );
}

function PlansView({ plans }: { plans: Plan[] }) {
  return (
    <Table
      headers={["ID", "Gói", "Giá", "Ngày", "QR tối đa", "Dynamic", "Logo", "Analytics", "Trạng thái"]}
      rows={plans.map((plan) => [
        plan.id,
        plan.name,
        formatMoney(plan.price),
        plan.duration_days,
        plan.max_qr_codes,
        plan.allow_dynamic_qr ? "Có" : "Không",
        plan.allow_logo ? "Có" : "Không",
        plan.allow_analytics ? "Có" : "Không",
        <StatusBadge status={plan.status} key="status" />
      ])}
    />
  );
}

function LogsView({ logs }: { logs: SystemLog[] }) {
  return (
    <Table
      headers={["ID", "Mức", "Hành động", "Đối tượng", "Thông điệp", "IP", "Thời gian"]}
      rows={logs.map((log) => [
        log.id,
        <StatusBadge status={log.level} key="level" />,
        log.action,
        log.entity_type ?? "",
        log.message,
        log.ip_address ?? "",
        formatDate(log.created_at)
      ])}
    />
  );
}

function Metric({ label, value, icon: Icon }: { label: string; value: number; icon: typeof Users }) {
  return (
    <section className="metric-card">
      <Icon size={22} />
      <span>{label}</span>
      <strong>{value.toLocaleString("vi-VN")}</strong>
    </section>
  );
}

function Table({ headers, rows }: { headers: string[]; rows: Array<Array<ReactNode>> }) {
  return (
    <section className="table-panel">
      <div className="table-scroll">
        <table>
          <thead>
            <tr>{headers.map((header) => <th key={header}>{header}</th>)}</tr>
          </thead>
          <tbody>
            {rows.length > 0 ? rows.map((row, index) => (
              <tr key={index}>{row.map((cell, cellIndex) => <td key={cellIndex}>{cell}</td>)}</tr>
            )) : (
              <tr><td colSpan={headers.length} className="empty-cell">Không có dữ liệu</td></tr>
            )}
          </tbody>
        </table>
      </div>
    </section>
  );
}

function StatusBadge({ status }: { status: string }) {
  return <span className={`status status-${status.toLowerCase()}`}>{status}</span>;
}

function createApi(token: string, onUnauthorized: () => void) {
  return {
    get: <T,>(path: string) => request<T>(path, { token, onUnauthorized }),
    put: <T,>(path: string, body: unknown) => request<T>(path, { method: "PUT", body: JSON.stringify(body), token, onUnauthorized })
  };
}

async function request<T>(path: string, options: RequestInit & { token?: string; onUnauthorized?: () => void } = {}): Promise<T> {
  const headers = new Headers(options.headers);
  headers.set("Content-Type", "application/json");
  if (options.token) {
    headers.set("Authorization", `Bearer ${options.token}`);
  }
  const response = await fetch(`${API_BASE_URL}${path}`, { ...options, headers });
  if (response.status === 401 && options.onUnauthorized) {
    options.onUnauthorized();
  }
  const payload = await response.json() as ApiEnvelope<T>;
  if (!response.ok) {
    throw new Error(payload.message || "Request failed");
  }
  return payload.data as T;
}

function filterItems<T>(items: T[], query: string, fields: (item: T) => string[]) {
  const normalized = query.trim().toLowerCase();
  if (!normalized) return items;
  return items.filter((item) => fields(item).some((field) => field.toLowerCase().includes(normalized)));
}

function formatMoney(value: number) {
  return `${value.toLocaleString("vi-VN")} VND`;
}

function formatDate(value?: string) {
  if (!value) return "";
  return new Intl.DateTimeFormat("vi-VN", { dateStyle: "short", timeStyle: "short" }).format(new Date(value));
}

function messageFromError(error: unknown) {
  return error instanceof Error ? error.message : "Đã có lỗi xảy ra";
}

export { App };
