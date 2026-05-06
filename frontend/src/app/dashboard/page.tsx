"use client";

import { AppShell } from "@/components/app-shell";
import { Panel, StatCard } from "@/components/ui";
import { cents, statusLabel } from "@/lib/format";
import { useAuth } from "@/components/auth-provider";

export default function DashboardPage() {
  const { user } = useAuth();

  return (
    <AppShell>
      <div className="grid gap-6">
        <div>
          <h1 className="text-2xl font-semibold text-slate-950">账户概览</h1>
          <p className="mt-1 text-sm text-slate-500">查看账户状态、余额和基础资料。</p>
        </div>
        <div className="grid gap-4 md:grid-cols-3">
          <StatCard label="当前余额" value={cents(user?.balance_cents)} />
          <StatCard label="账户状态" value={statusLabel(user?.status)} />
          <StatCard label="账户角色" value={statusLabel(user?.role)} />
        </div>
        <Panel title="个人信息">
          <dl className="grid gap-4 text-sm md:grid-cols-2">
            <div>
              <dt className="text-slate-500">用户 ID</dt>
              <dd className="mt-1 font-medium text-slate-950">{user?.id}</dd>
            </div>
            <div>
              <dt className="text-slate-500">名称</dt>
              <dd className="mt-1 font-medium text-slate-950">{user?.name}</dd>
            </div>
            <div>
              <dt className="text-slate-500">邮箱</dt>
              <dd className="mt-1 font-medium text-slate-950">{user?.email}</dd>
            </div>
            <div>
              <dt className="text-slate-500">余额</dt>
              <dd className="mt-1 font-medium text-slate-950">{cents(user?.balance_cents)}</dd>
            </div>
          </dl>
        </Panel>
      </div>
    </AppShell>
  );
}
