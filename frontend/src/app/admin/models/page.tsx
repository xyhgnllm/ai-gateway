"use client";

import { AppShell } from "@/components/app-shell";
import { Button, EmptyState, Field, Panel, StatusBadge, TextInput } from "@/components/ui";
import { apiFetch, asArray } from "@/lib/api";
import { cents, dateTime, statusLabel } from "@/lib/format";
import type { Model } from "@/types/api";
import { FormEvent, useEffect, useState } from "react";

export default function AdminModelsPage() {
  const [models, setModels] = useState<Model[]>([]);
  const [name, setName] = useState("");
  const [provider, setProvider] = useState("openai");
  const [inputPrice, setInputPrice] = useState("0");
  const [outputPrice, setOutputPrice] = useState("0");
  const [error, setError] = useState("");

  async function load() {
    setModels(asArray(await apiFetch<Model[] | null>("/admin/models")));
  }

  useEffect(() => {
    load().catch((err) => setError(err.message));
  }, []);

  async function create(event: FormEvent) {
    event.preventDefault();
    await apiFetch<Model>("/admin/models", {
      method: "POST",
      body: {
        name,
        provider,
        input_price_per_1k_cents: Math.round(Number(inputPrice) * 100),
        output_price_per_1k_cents: Math.round(Number(outputPrice) * 100),
      },
    });
    setName("");
    await load();
  }

  async function toggle(model: Model) {
    await apiFetch<Model>(`/admin/models/${model.id}/status`, {
      method: "PATCH",
      body: { status: model.status === "active" ? "disabled" : "active" },
    });
    await load();
  }

  return (
    <AppShell admin>
      <div className="grid gap-6">
        <Panel title="创建模型">
          <form className="grid gap-3 md:grid-cols-5 md:items-end" onSubmit={create}>
            <Field label="模型名称">
              <TextInput onChange={(event) => setName(event.target.value)} required value={name} />
            </Field>
            <Field label="供应商">
              <TextInput onChange={(event) => setProvider(event.target.value)} value={provider} />
            </Field>
            <Field label="输入单价/1k">
              <TextInput onChange={(event) => setInputPrice(event.target.value)} type="number" value={inputPrice} />
            </Field>
            <Field label="输出单价/1k">
              <TextInput onChange={(event) => setOutputPrice(event.target.value)} type="number" value={outputPrice} />
            </Field>
            <Button type="submit">创建</Button>
          </form>
          {error ? <div className="mt-3 text-sm text-red-600">{error}</div> : null}
        </Panel>
        <Panel title="模型管理">
          {models.length === 0 ? (
            <EmptyState text="暂无模型" />
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full min-w-[860px] text-left text-sm">
                <thead className="border-b text-slate-500">
                  <tr>
                    <th className="py-3">名称</th>
                    <th>供应商</th>
                    <th>输入单价</th>
                    <th>输出单价</th>
                    <th>状态</th>
                    <th>创建时间</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody className="divide-y">
                  {models.map((model) => (
                    <tr key={model.id}>
                      <td className="py-3 font-medium text-slate-950">{model.name}</td>
                      <td>{model.provider}</td>
                      <td>{cents(model.input_price_per_1k_cents)}</td>
                      <td>{cents(model.output_price_per_1k_cents)}</td>
                      <td>
                        <StatusBadge value={statusLabel(model.status)} />
                      </td>
                      <td>{dateTime(model.created_at)}</td>
                      <td>
                        <Button variant="secondary" onClick={() => toggle(model)}>
                          {model.status === "active" ? "禁用" : "启用"}
                        </Button>
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
