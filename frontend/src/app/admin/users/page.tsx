"use client";

import { AppShell } from "@/components/app-shell";
import { Button, EmptyState, Field, Panel, Select, StatusBadge, TextInput } from "@/components/ui";
import { apiFetch, asArray } from "@/lib/api";
import { cents, statusLabel } from "@/lib/format";
import type { User } from "@/types/api";
import { FormEvent, useEffect, useState } from "react";

export default function AdminUsersPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [selectedUserId, setSelectedUserId] = useState("");
  const [amount, setAmount] = useState("100");
  const [error, setError] = useState("");

  async function load() {
    setUsers(asArray(await apiFetch<User[] | null>("/admin/users?page=1&page_size=100")));
  }

  useEffect(() => {
    load().catch((err) => setError(err.message));
  }, []);

  async function addBalance(event: FormEvent) {
    event.preventDefault();
    if (!selectedUserId) return;
    await apiFetch<User>(`/admin/users/${selectedUserId}/balance`, {
      method: "POST",
      body: { amount_cents: Math.round(Number(amount) * 100) },
    });
    await load();
  }

  async function patchUser(id: number, path: "status" | "role", value: string) {
    await apiFetch<User>(`/admin/users/${id}/${path}`, {
      method: "PATCH",
      body: { [path]: value },
    });
    await load();
  }

  return (
    <AppShell admin>
      <div className="grid gap-6">
        <Panel title="加余额">
          <form className="flex flex-wrap items-end gap-3" onSubmit={addBalance}>
            <Field label="用户">
              <Select onChange={(event) => setSelectedUserId(event.target.value)} value={selectedUserId}>
                <option value="">选择用户</option>
                {users.map((user) => (
                  <option key={user.id} value={user.id}>
                    {user.email}
                  </option>
                ))}
              </Select>
            </Field>
            <Field label="金额">
              <TextInput
                min="1"
                onChange={(event) => setAmount(event.target.value)}
                step="1"
                type="number"
                value={amount}
              />
            </Field>
            <Button type="submit">确认加余额</Button>
          </form>
          {error ? <div className="mt-3 text-sm text-red-600">{error}</div> : null}
        </Panel>
        <Panel title="用户管理">
          {users.length === 0 ? (
            <EmptyState text="暂无用户" />
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full min-w-[900px] text-left text-sm">
                <thead className="border-b text-slate-500">
                  <tr>
                    <th className="py-3">ID</th>
                    <th>邮箱</th>
                    <th>名称</th>
                    <th>角色</th>
                    <th>状态</th>
                    <th>余额</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody className="divide-y">
                  {users.map((user) => (
                    <tr key={user.id}>
                      <td className="py-3">{user.id}</td>
                      <td className="font-medium text-slate-950">{user.email}</td>
                      <td>{user.name}</td>
                      <td>
                        <StatusBadge value={statusLabel(user.role)} />
                      </td>
                      <td>
                        <StatusBadge value={statusLabel(user.status)} />
                      </td>
                      <td>{cents(user.balance_cents)}</td>
                      <td>
                        <div className="flex flex-wrap gap-2">
                          <Button
                            variant="secondary"
                            onClick={() =>
                              patchUser(user.id, "status", user.status === "active" ? "disabled" : "active")
                            }
                          >
                            {user.status === "active" ? "禁用" : "启用"}
                          </Button>
                          <Button
                            variant="secondary"
                            onClick={() => patchUser(user.id, "role", user.role === "admin" ? "user" : "admin")}
                          >
                            {user.role === "admin" ? "设为用户" : "设为管理员"}
                          </Button>
                        </div>
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
