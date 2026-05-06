"use client";

import { Button, Field, TextInput } from "@/components/ui";
import { apiFetch } from "@/lib/api";
import type { User } from "@/types/api";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";

export default function RegisterPage() {
  const router = useRouter();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function submit(event: FormEvent) {
    event.preventDefault();
    setError("");
    setLoading(true);
    try {
      await apiFetch<User>("/users", {
        method: "POST",
        body: { name, email, password },
        token: null,
      });
      router.replace("/login");
    } catch (err) {
      setError(err instanceof Error ? err.message : "注册失败");
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="grid min-h-screen place-items-center bg-slate-100 px-4">
      <section className="w-full max-w-md rounded-lg border border-slate-200 bg-white p-6">
        <h1 className="text-2xl font-semibold text-slate-950">创建账号</h1>
        <form className="mt-6 grid gap-4" onSubmit={submit}>
          <Field label="名称">
            <TextInput onChange={(event) => setName(event.target.value)} required value={name} />
          </Field>
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
              autoComplete="new-password"
              minLength={8}
              onChange={(event) => setPassword(event.target.value)}
              required
              type="password"
              value={password}
            />
          </Field>
          {error ? <div className="text-sm text-red-600">{error}</div> : null}
          <Button disabled={loading} type="submit">
            {loading ? "提交中..." : "注册"}
          </Button>
        </form>
        <div className="mt-4 text-sm text-slate-600">
          已有账号？{" "}
          <Link className="font-medium text-slate-950 underline" href="/login">
            登录
          </Link>
        </div>
      </section>
    </main>
  );
}
