"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { Download, Folder as FolderIcon, Search, Trash2 } from "lucide-react";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { CreateQRForm } from "@/components/qrcode/CreateQRForm";
import { EmptyState } from "@/components/common/State";
import { Button } from "@/components/common/Button";
import { api, messageFromError } from "@/lib/api";
import type { ApiEnvelope, QRCode } from "@/types";

interface FolderItem {
  id: number;
  name: string;
}

export default function QRCodesPage() {
  const [items, setItems] = useState<QRCode[]>([]);
  const [folders, setFolders] = useState<FolderItem[]>([]);
  const [selectedFolder, setSelectedFolder] = useState<string>("");
  const [q, setQ] = useState("");
  const [error, setError] = useState("");

  const loadFolders = useCallback(async () => {
    try {
      const res = await api.get<ApiEnvelope<FolderItem[]>>("/folders");
      setFolders(res.data.data ?? []);
    } catch {
      // ignore
    }
  }, []);

  const load = useCallback(async () => {
    try {
      const params: Record<string, string> = {};
      if (q) params.q = q;
      if (selectedFolder) params.folder_id = selectedFolder;

      const res = await api.get<ApiEnvelope<{ items: QRCode[] }>>("/qrcodes", { params });
      setItems(res.data.data?.items ?? []);
    } catch (err) {
      setError(messageFromError(err));
    }
  }, [q, selectedFolder]);

  useEffect(() => {
    loadFolders();
  }, [loadFolders]);

  useEffect(() => {
    load();
  }, [load]);

  async function remove(id: number) {
    await api.delete(`/qrcodes/${id}`);
    load();
  }

  async function download(id: number) {
    const res = await api.get(`/qrcodes/${id}/download`, { responseType: "blob" });
    const url = URL.createObjectURL(res.data);
    const link = document.createElement("a");
    link.href = url;
    link.download = "qr-code.png";
    link.click();
    URL.revokeObjectURL(url);
  }

  return (
    <DashboardShell>
      <div className="flex flex-wrap items-center justify-between gap-3">
        <h1 className="text-3xl font-bold text-ink">QR Codes</h1>
        <div className="flex flex-wrap gap-2">
          {/* Folder filter dropdown */}
          <select
            value={selectedFolder}
            onChange={(e) => setSelectedFolder(e.target.value)}
            className="focus-ring rounded-md border border-slate-200 bg-white px-3 py-2 text-sm"
          >
            <option value="">Tất cả Thư mục</option>
            <option value="uncategorized">Chưa phân loại (Uncategorized)</option>
            {folders.map((f) => (
              <option key={f.id} value={f.id}>
                📁 {f.name}
              </option>
            ))}
          </select>

          <input
            value={q}
            onChange={(e) => setQ(e.target.value)}
            className="focus-ring rounded-md border border-slate-200 px-3 py-2 text-sm"
            placeholder="Search QR..."
          />
          <Button tone="secondary" onClick={load}>
            <Search className="h-4 w-4" />
          </Button>
        </div>
      </div>

      <div className="mt-6">
        <CreateQRForm onSaved={load} />
      </div>

      <section className="mt-8">
        <div className="flex items-center justify-between">
          <h2 className="text-xl font-bold text-ink">Saved QR</h2>
          {selectedFolder ? (
            <button
              onClick={() => setSelectedFolder("")}
              className="text-xs text-teal hover:underline"
            >
              Bỏ lọc thư mục
            </button>
          ) : null}
        </div>

        {error ? <p className="mt-3 rounded-md bg-coral/10 px-3 py-2 text-sm text-coral">{error}</p> : null}
        {items.length === 0 ? (
          <div className="mt-4">
            <EmptyState title="No QR codes found" description="Saved QR codes will appear here." />
          </div>
        ) : null}

        <div className="mt-4 grid gap-4 md:grid-cols-2">
          {items.map((qr) => {
            const folderName = folders.find((f) => f.id === qr.folder_id)?.name;
            return (
              <div key={qr.id} className="rounded-md border border-slate-200 bg-white p-4 shadow-soft">
                <div className="flex items-start justify-between gap-3">
                  <div>
                    <Link href={`/qrcodes/${qr.id}`} className="font-semibold text-ink hover:text-teal">
                      {qr.title}
                    </Link>
                    <p className="mt-1 text-sm text-muted">
                      {qr.qr_type} · {qr.status} · {qr.scan_count} scans
                    </p>
                    {folderName ? (
                      <p className="mt-1 flex items-center gap-1 text-xs text-teal font-medium">
                        <FolderIcon className="h-3 w-3" />
                        {folderName}
                      </p>
                    ) : null}
                  </div>
                  <span className="rounded-md bg-panel px-2 py-1 text-xs font-semibold text-muted">
                    {qr.is_dynamic ? "Dynamic" : "Static"}
                  </span>
                </div>
                <div className="mt-4 flex gap-2">
                  <Button tone="secondary" onClick={() => download(qr.id)}>
                    <Download className="h-4 w-4" />
                  </Button>
                  <Button tone="danger" onClick={() => remove(qr.id)}>
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            );
          })}
        </div>
      </section>
    </DashboardShell>
  );
}
