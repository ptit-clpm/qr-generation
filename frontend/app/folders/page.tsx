"use client";

import { useEffect, useState } from "react";
import { FolderPlus, Trash2 } from "lucide-react";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { Button } from "@/components/common/Button";
import { EmptyState } from "@/components/common/State";
import { api, messageFromError } from "@/lib/api";
import type { ApiEnvelope } from "@/types";

interface FolderItem {
  id: number;
  name: string;
  description?: string;
}

export default function FoldersPage() {
  const [folders, setFolders] = useState<FolderItem[]>([]);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [error, setError] = useState("");

  async function load() {
    try {
      const res = await api.get<ApiEnvelope<FolderItem[]>>("/folders");
      setFolders(res.data.data ?? []);
    } catch (err) {
      setError(messageFromError(err));
    }
  }

  useEffect(() => { load(); }, []);

  async function create() {
    setError("");
    try {
      await api.post("/folders", { name, description });
      setName("");
      setDescription("");
      load();
    } catch (err) {
      setError(messageFromError(err));
    }
  }

  async function remove(id: number) {
    await api.delete(`/folders/${id}`);
    load();
  }

  return (
    <DashboardShell>
      <h1 className="text-3xl font-bold text-ink">Folders</h1>
      <section className="mt-6 rounded-md border border-slate-200 bg-white p-5 shadow-soft">
        <div className="grid gap-3 md:grid-cols-[1fr_1fr_auto]">
          <input value={name} onChange={(e) => setName(e.target.value)} placeholder="Name" className="focus-ring rounded-md border border-slate-200 px-3 py-2" />
          <input value={description} onChange={(e) => setDescription(e.target.value)} placeholder="Description" className="focus-ring rounded-md border border-slate-200 px-3 py-2" />
          <Button onClick={create}><FolderPlus className="h-4 w-4" />Create</Button>
        </div>
        {error ? <p className="mt-3 rounded-md bg-coral/10 px-3 py-2 text-sm text-coral">{error}</p> : null}
      </section>
      <div className="mt-5 grid gap-4 md:grid-cols-3">
        {folders.map((folder) => (
          <div key={folder.id} className="rounded-md border border-slate-200 bg-white p-4 shadow-soft">
            <p className="font-semibold text-ink">{folder.name}</p>
            <p className="mt-1 text-sm text-muted">{folder.description || "No description"}</p>
            <Button className="mt-4" tone="danger" onClick={() => remove(folder.id)}><Trash2 className="h-4 w-4" /></Button>
          </div>
        ))}
      </div>
      {folders.length === 0 ? <div className="mt-5"><EmptyState title="No folders" description="Folders help organize QR codes for Pro workflows." /></div> : null}
    </DashboardShell>
  );
}
