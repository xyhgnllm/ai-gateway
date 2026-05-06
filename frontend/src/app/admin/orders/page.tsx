"use client";

import { AppShell } from "@/components/app-shell";
import { Button, EmptyState, Panel, StatusBadge } from "@/components/ui";
import { apiFetch, asArray } from "@/lib/api";
import { cents, dateTime, statusLabel } from "@/lib/format";
import type { Order } from "@/types/api";
import { useEffect, useState } from "react";

export default function AdminOrdersPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [error, setError] = useState("");

  async function load() {
    setOrders(asArray(await apiFetch<Order[] | null>("/admin/orders?page=1&page_size=100")));
  }

  useEffect(() => {
    load().catch((err) => setError(err.message));
  }, []);

  async function pay(id: number) {
    await apiFetch(`/admin/orders/${id}/pay`, { method: "POST" });
    await load();
  }

  return (
    <AppShell admin>
      <Panel title="订单管理">
        {error ? <div className="mb-3 text-sm text-red-600">{error}</div> : null}
        {orders.length === 0 ? (
          <EmptyState text="暂无订单" />
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full min-w-215 text-left text-sm">
              <thead className="border-b text-slate-500">
                <tr>
                  <th className="py-3">订单号</th>
                  <th>用户 ID</th>
                  <th>金额</th>
                  <th>状态</th>
                  <th>创建时间</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody className="divide-y">
                {orders.map((order) => (
                  <tr key={order.id}>
                    <td className="py-3 font-medium text-slate-950">{order.order_no}</td>
                    <td>{order.user_id}</td>
                    <td>{cents(order.amount_cents)}</td>
                    <td>
                      <StatusBadge value={statusLabel(order.status)} />
                    </td>
                    <td>{dateTime(order.created_at)}</td>
                    <td>
                      {order.status === "pending" ? (
                        <Button variant="secondary" onClick={() => pay(order.id)}>
                          标记支付
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
    </AppShell>
  );
}
