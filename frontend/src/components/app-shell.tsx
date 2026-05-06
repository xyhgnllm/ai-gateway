"use client";

import { useAuth } from "@/components/auth-provider";
import { NavLink } from "@/components/ui";
import { apiFetch } from "@/lib/api";
import type { User } from "@/types/api";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useEffect, useState } from "react";

const userLinks = [
  ["/dashboard", "概览"],
  ["/dashboard/orders", "订单"],
  ["/dashboard/api-keys", "API Key"],
  ["/dashboard/usage-logs", "调用日志"],
  ["/dashboard/balance-transactions", "余额流水"],
];

const adminLinks = [
  ["/admin", "统计"],
  ["/admin/users", "用户"],
  ["/admin/orders", "订单"],
  ["/admin/models", "模型"],
];

export function AppShell({
  children,
  admin = false,
}: {
  children: React.ReactNode;
  admin?: boolean;
}) {
  const auth = useAuth();
  const router = useRouter();
  const pathname = usePathname();
  const { ready, signIn, signOut, token } = auth;
  const [authError, setAuthError] = useState("");

  useEffect(() => {
    if (!ready) return;
    if (!token) {
      router.replace("/login");
      return;
    }
    apiFetch<User>("/me")
      .then((user) => {
        setAuthError("");
        signIn(token, user);
        if (admin && user.role !== "admin") router.replace("/dashboard");
      })
      .catch((error) => {
        setAuthError(error instanceof Error ? error.message : "认证失败");
      });
  }, [admin, ready, router, signIn, token]);

  if (!ready || !token) {
    return (
      <div className="p-8 text-sm text-slate-500" suppressHydrationWarning>
        加载中...
      </div>
    );
  }

  const links = admin ? adminLinks : userLinks;

  return (
    <div className="min-h-screen bg-slate-100">
      <header className="border-b border-slate-200 bg-white">
        <div className="mx-auto flex max-w-7xl flex-wrap items-center justify-between gap-3 px-4 py-4">
          <Link
            className="text-lg font-semibold text-slate-950"
            href={admin ? "/admin" : "/dashboard"}
          >
            AI Gateway
          </Link>
          <div className="flex items-center gap-3 text-sm text-slate-600">
            <span>{auth.user?.name ?? auth.user?.email}</span>
            <button
              className="rounded-md border border-slate-300 px-3 py-2 text-slate-700 hover:bg-slate-50"
              onClick={() => {
                auth.signOut();
                router.replace("/login");
              }}
            >
              退出
            </button>
          </div>
        </div>
      </header>
      <div className="mx-auto grid max-w-7xl gap-6 px-4 py-6 lg:grid-cols-[220px_1fr]">
        <aside className="rounded-lg border border-slate-200 bg-white p-3">
          <nav className="grid gap-1">
            {links.map(([href, label]) => (
              <Link
                key={href}
                className={`rounded-md px-3 py-2 text-sm font-medium ${
                  pathname === href
                    ? "bg-slate-950 text-white"
                    : "text-slate-600 hover:bg-slate-100 hover:text-slate-950"
                }`}
                href={href}
              >
                {label}
              </Link>
            ))}
            {auth.user?.role === "admin" && !admin ? (
              <NavLink href="/admin">进入管理端</NavLink>
            ) : null}
            {admin ? <NavLink href="/dashboard">返回用户端</NavLink> : null}
          </nav>
        </aside>
        <main>
          {authError ? (
            <div className="mb-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">
              登录态校验失败：{authError}
              <button
                className="ml-3 font-medium underline"
                onClick={() => {
                  signOut();
                  router.replace("/login");
                }}
              >
                重新登录
              </button>
            </div>
          ) : null}
          {children}
        </main>
      </div>
    </div>
  );
}
