"use client";

import { AppShell } from "@/components/app-shell";
import { EmptyState, Panel } from "@/components/ui";
import { apiFetch, asArray } from "@/lib/api";
import { cents, dateTime } from "@/lib/format";
import type { UsageLog } from "@/types/api";
import { useEffect, useState } from "react";

export default function UsageLogsPage() {
  const [logs, setLogs] = useState<UsageLog[]>([]);
  const [error, setError] = useState("");

  useEffect(() => {
    apiFetch<UsageLog[] | null>("/usage-logs?page=1&page_size=50")
      .then((data) => setLogs(asArray(data)))
      .catch((err) => setError(err.message));
  }, []);

  return (
    <AppShell>
      <Panel title="调用日志">
        {error ? <div className="mb-3 text-sm text-red-600">{error}</div> : null}
        {logs.length === 0 ? (
          <EmptyState text="暂无调用日志" />
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full min-w-[760px] text-left text-sm">
              <thead className="border-b text-slate-500">
                <tr>
                  <th className="py-3">模型</th>
                  <th>输入</th>
                  <th>输出</th>
                  <th>总量</th>
                  <th>费用</th>
                  <th>时间</th>
                </tr>
              </thead>
              <tbody className="divide-y">
                {logs.map((log) => (
                  <tr key={log.id}>
                    <td className="py-3 font-medium text-slate-950">{log.model}</td>
                    <td>{log.prompt_tokens}</td>
                    <td>{log.completion_tokens}</td>
                    <td>{log.total_tokens}</td>
                    <td>{cents(log.cost_cents)}</td>
                    <td>{dateTime(log.created_at)}</td>
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
