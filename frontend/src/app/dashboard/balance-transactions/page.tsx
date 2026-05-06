"use client";

import { AppShell } from "@/components/app-shell";
import { EmptyState, Panel } from "@/components/ui";
import { apiFetch, asArray } from "@/lib/api";
import { cents, dateTime, statusLabel } from "@/lib/format";
import type { BalanceTransaction } from "@/types/api";
import { useEffect, useState } from "react";

export default function BalanceTransactionsPage() {
  const [transactions, setTransactions] = useState<BalanceTransaction[]>([]);
  const [error, setError] = useState("");

  useEffect(() => {
    apiFetch<BalanceTransaction[] | null>("/balance-transactions?page=1&page_size=50")
      .then((data) => setTransactions(asArray(data)))
      .catch((err) => setError(err.message));
  }, []);

  return (
    <AppShell>
      <Panel title="余额流水">
        {error ? <div className="mb-3 text-sm text-red-600">{error}</div> : null}
        {transactions.length === 0 ? (
          <EmptyState text="暂无余额流水" />
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full min-w-190 text-left text-sm">
              <thead className="border-b text-slate-500">
                <tr>
                  <th className="py-3">类型</th>
                  <th>金额</th>
                  <th>变动前</th>
                  <th>变动后</th>
                  <th>备注</th>
                  <th>时间</th>
                </tr>
              </thead>
              <tbody className="divide-y">
                {transactions.map((item) => (
                  <tr key={item.id}>
                    <td className="py-3 font-medium text-slate-950">{statusLabel(item.type)}</td>
                    <td>{cents(item.amount_cents)}</td>
                    <td>{cents(item.balance_before_cents)}</td>
                    <td>{cents(item.balance_after_cents)}</td>
                    <td>{item.remark || "-"}</td>
                    <td>{dateTime(item.created_at)}</td>
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
