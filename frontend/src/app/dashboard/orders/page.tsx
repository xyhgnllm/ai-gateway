"use client";

import { AppShell } from "@/components/app-shell";
import { Button, EmptyState, Field, Panel, StatusBadge, TextInput } from "@/components/ui";
import { apiFetch, asArray } from "@/lib/api";
import { cents, dateTime, statusLabel } from "@/lib/format";
import type { Order } from "@/types/api";
import { FormEvent, useEffect, useState } from "react";

export default function OrdersPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [amount, setAmount] = useState("100");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function load() {
    const data = await apiFetch<Order[] | null>("/orders?page=1&page_size=50");
    setOrders(asArray(data));
  }

  useEffect(() => {
    load().catch((err) => setError(err.message));
  }, []);

  async function submit(event: FormEvent) {
    event.preventDefault();
    setError("");
    setLoading(true);
    try {
      await apiFetch<Order>("/orders", {
        method: "POST",
        body: { amount_cents: Math.round(Number(amount) * 100) },
      });
      setAmount("100");
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "创建订单失败");
    } finally {
      setLoading(false);
    }
  }

  return (
    <AppShell>
      <div className="grid gap-6">
        <Panel title="创建充值订单">
          <form className="flex flex-wrap items-end gap-3" onSubmit={submit}>
            <Field label="充值金额">
              <TextInput
                min="1"
                onChange={(event) => setAmount(event.target.value)}
                step="1"
                type="number"
                value={amount}
              />
            </Field>
            <Button disabled={loading} type="submit">
              {loading ? "创建中..." : "创建订单"}
            </Button>
          </form>
          {error ? <div className="mt-3 text-sm text-red-600">{error}</div> : null}
        </Panel>
        <Panel title="订单列表">
          {orders.length === 0 ? (
            <EmptyState text="暂无订单" />
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full min-w-[720px] text-left text-sm">
                <thead className="border-b text-slate-500">
                  <tr>
                    <th className="py-3">订单号</th>
                    <th>金额</th>
                    <th>状态</th>
                    <th>创建时间</th>
                    <th>支付时间</th>
                  </tr>
                </thead>
                <tbody className="divide-y">
                  {orders.map((order) => (
                    <tr key={order.id}>
                      <td className="py-3 font-medium text-slate-950">{order.order_no}</td>
                      <td>{cents(order.amount_cents)}</td>
                      <td>
                        <StatusBadge value={statusLabel(order.status)} />
                      </td>
                      <td>{dateTime(order.created_at)}</td>
                      <td>{dateTime(order.paid_at)}</td>
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
