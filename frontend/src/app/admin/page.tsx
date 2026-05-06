"use client";

import { AppShell } from "@/components/app-shell";
import { EmptyState, Panel, StatCard } from "@/components/ui";
import { apiFetch, asArray } from "@/lib/api";
import { cents } from "@/lib/format";
import type { AdminStats, UsageStat } from "@/types/api";
import { useEffect, useState } from "react";

function StatTable({ title, rows }: { title: string; rows: UsageStat[] }) {
  const keys = rows[0] ? Object.keys(rows[0]) : [];
  return (
    <Panel title={title}>
      {rows.length === 0 ? (
        <EmptyState text="暂无数据" />
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full min-w-[640px] text-left text-sm">
            <thead className="border-b text-slate-500">
              <tr>{keys.map((key) => <th className="py-3" key={key}>{key}</th>)}</tr>
            </thead>
            <tbody className="divide-y">
              {rows.map((row, index) => (
                <tr key={index}>
                  {keys.map((key) => (
                    <td className="py-3" key={key}>{String(row[key] ?? "-")}</td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </Panel>
  );
}

export default function AdminPage() {
  const [stats, setStats] = useState<AdminStats | null>(null);
  const [modelStats, setModelStats] = useState<UsageStat[]>([]);
  const [userStats, setUserStats] = useState<UsageStat[]>([]);
  const [dailyStats, setDailyStats] = useState<UsageStat[]>([]);
  const [error, setError] = useState("");

  useEffect(() => {
    Promise.all([
      apiFetch<AdminStats>("/admin/stats"),
      apiFetch<UsageStat[] | null>("/admin/stats/models"),
      apiFetch<UsageStat[] | null>("/admin/stats/users"),
      apiFetch<UsageStat[] | null>("/admin/stats/daily"),
    ])
      .then(([overview, models, users, daily]) => {
        setStats(overview);
        setModelStats(asArray(models));
        setUserStats(asArray(users));
        setDailyStats(asArray(daily));
      })
      .catch((err) => setError(err.message));
  }, []);

  return (
    <AppShell admin>
      <div className="grid gap-6">
        <div>
          <h1 className="text-2xl font-semibold text-slate-950">统计面板</h1>
          <p className="mt-1 text-sm text-slate-500">查看用户、订单、收入和调用成本。</p>
        </div>
        {error ? <div className="text-sm text-red-600">{error}</div> : null}
        <div className="grid gap-4 md:grid-cols-5">
          <StatCard label="用户数" value={stats?.user_count ?? "-"} />
          <StatCard label="订单数" value={stats?.order_count ?? "-"} />
          <StatCard label="已支付金额" value={cents(stats?.paid_amount_cents)} />
          <StatCard label="调用次数" value={stats?.usage_count ?? "-"} />
          <StatCard label="调用成本" value={cents(stats?.usage_cost_cents)} />
        </div>
        <StatTable title="模型用量" rows={modelStats} />
        <StatTable title="用户用量" rows={userStats} />
        <StatTable title="每日用量" rows={dailyStats} />
      </div>
    </AppShell>
  );
}
