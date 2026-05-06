"use client";

import { AppShell } from "@/components/app-shell";
import { Button, EmptyState, Field, Panel, StatusBadge, TextInput } from "@/components/ui";
import { apiFetch, asArray } from "@/lib/api";
import { dateTime, statusLabel } from "@/lib/format";
import type { ApiKey, CreatedApiKey } from "@/types/api";
import { FormEvent, useEffect, useState } from "react";

export default function ApiKeysPage() {
  const [keys, setKeys] = useState<ApiKey[]>([]);
  const [name, setName] = useState("default");
  const [createdKey, setCreatedKey] = useState("");
  const [error, setError] = useState("");

  async function load() {
    setKeys(asArray(await apiFetch<ApiKey[] | null>("/api-keys")));
  }

  useEffect(() => {
    load().catch((err) => setError(err.message));
  }, []);

  async function submit(event: FormEvent) {
    event.preventDefault();
    setError("");
    try {
      const data = await apiFetch<CreatedApiKey>("/api-keys", {
        method: "POST",
        body: { name },
      });
      setCreatedKey(data.api_key.key);
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "创建 API Key 失败");
    }
  }

  async function disable(id: number) {
    await apiFetch<ApiKey>(`/api-keys/${id}/disable`, { method: "PATCH" });
    await load();
  }

  return (
    <AppShell>
      <div className="grid gap-6">
        <Panel title="创建 API Key">
          <form className="flex flex-wrap items-end gap-3" onSubmit={submit}>
            <Field label="名称">
              <TextInput onChange={(event) => setName(event.target.value)} value={name} />
            </Field>
            <Button type="submit">创建</Button>
          </form>
          {createdKey ? (
            <div className="mt-4 rounded-md border border-emerald-200 bg-emerald-50 p-3 text-sm text-emerald-800">
              <div className="font-medium">新 Key 仅显示一次</div>
              <code className="mt-2 block break-all rounded bg-white p-2 text-slate-950">{createdKey}</code>
            </div>
          ) : null}
          {error ? <div className="mt-3 text-sm text-red-600">{error}</div> : null}
        </Panel>
        <Panel title="API Key 列表">
          {keys.length === 0 ? (
            <EmptyState text="暂无 API Key" />
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full min-w-[720px] text-left text-sm">
                <thead className="border-b text-slate-500">
                  <tr>
                    <th className="py-3">名称</th>
                    <th>前缀</th>
                    <th>状态</th>
                    <th>最后使用</th>
                    <th>创建时间</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody className="divide-y">
                  {keys.map((key) => (
                    <tr key={key.id}>
                      <td className="py-3 font-medium text-slate-950">{key.name}</td>
                      <td>{key.key_prefix}</td>
                      <td>
                        <StatusBadge value={statusLabel(key.status)} />
                      </td>
                      <td>{dateTime(key.last_used_at)}</td>
                      <td>{dateTime(key.created_at)}</td>
                      <td>
                        {key.status === "active" ? (
                          <Button variant="secondary" onClick={() => disable(key.id)}>
                            禁用
                          </Button>
                        ) : (
                          "-"
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </Panel>
      </div>
    </AppShell>
  );
}
