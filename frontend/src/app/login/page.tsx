"use client";

import { useAuth } from "@/components/auth-provider";
import { Button, Field, TextInput } from "@/components/ui";
import { apiFetch } from "@/lib/api";
import type { LoginResponse, User } from "@/types/api";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";

export default function LoginPage() {
  const auth = useAuth();
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function submit(event: FormEvent) {
    event.preventDefault();
    setError("");
    setLoading(true);
    try {
      const data = await apiFetch<LoginResponse>("/login", {
        method: "POST",
        body: { email, password },
        token: null,
      });
      const freshUser = await apiFetch<User>("/me", { token: data.token });
      auth.signIn(data.token, freshUser);
      router.replace(freshUser.role === "admin" ? "/admin" : "/dashboard");
    } catch (err) {
      setError(err instanceof Error ? err.message : "登录失败");
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="grid min-h-screen place-items-center bg-slate-100 px-4">
      <section className="w-full max-w-md rounded-lg border border-slate-200 bg-white p-6">
        <h1 className="text-2xl font-semibold text-slate-950">AI Gateway</h1>
        <p className="mt-2 text-sm text-slate-500">登录控制台继续管理账户和模型调用。</p>
        <form className="mt-6 grid gap-4" onSubmit={submit}>
          <Field label="邮箱">
            <TextInput
              autoComplete="email"
              onChange={(event) => setEmail(event.target.value)}
              required
              type="email"
              value={email}
            />
          </Field>
          <Field label="密码">
            <TextInput
              autoComplete="current-password"
              minLength={8}
              onChange={(event) => setPassword(event.target.value)}
              required
              type="password"
              value={password}
            />
          </Field>
          {error ? <div className="text-sm text-red-600">{error}</div> : null}
          <Button disabled={loading} type="submit">
            {loading ? "登录中..." : "登录"}
          </Button>
        </form>
        <div className="mt-4 text-sm text-slate-600">
          没有账号？{" "}
          <Link className="font-medium text-slate-950 underline" href="/register">
            注册
          </Link>
        </div>
      </section>
    </main>
  );
}
