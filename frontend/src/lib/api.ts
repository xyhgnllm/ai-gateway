import type { ApiEnvelope } from "@/types/api";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL?.replace(/\/$/, "") ||
  "http://localhost:8080";

export class ApiError extends Error {
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.status = status;
  }
}

export function getToken() {
  if (typeof window === "undefined") return null;
  return localStorage.getItem("ai_gateway_token");
}

export function setToken(token: string) {
  localStorage.setItem("ai_gateway_token", token);
}

export function clearToken() {
  localStorage.removeItem("ai_gateway_token");
}

type RequestOptions = {
  method?: "GET" | "POST" | "PATCH" | "DELETE";
  body?: unknown;
  token?: string | null;
};

export async function apiFetch<T>(path: string, options: RequestOptions = {}) {
  const token = options.token === undefined ? getToken() : options.token;
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (token) headers.Authorization = `Bearer ${token}`;

  const response = await fetch(`${API_BASE_URL}${path}`, {
    method: options.method ?? "GET",
    headers,
    body: options.body ? JSON.stringify(options.body) : undefined,
    cache: "no-store",
  });

  const contentType = response.headers.get("content-type") ?? "";
  const payload = contentType.includes("application/json")
    ? ((await response.json()) as ApiEnvelope<T>)
    : null;

  if (!response.ok || !payload || payload.code !== 0) {
    throw new ApiError(payload?.message || response.statusText, response.status);
  }

  return normalizeApiData(payload.data) as T;
}

export function asArray<T>(value: T[] | null | undefined): T[] {
  return Array.isArray(value) ? value : [];
}

function normalizeApiData(value: unknown): unknown {
  if (Array.isArray(value)) {
    return value.map(normalizeApiData);
  }

  if (!value || typeof value !== "object") {
    return value;
  }

  const record = value as Record<string, unknown>;
  if ("Valid" in record && "Time" in record) {
    return record.Valid ? record.Time : null;
  }
  if ("Valid" in record && "Int64" in record) {
    return record.Valid ? record.Int64 : null;
  }

  return Object.fromEntries(
    Object.entries(record).map(([key, item]) => [
      toSnakeCase(key),
      normalizeApiData(item),
    ]),
  );
}

function toSnakeCase(key: string) {
  return key
    .replace(/API/g, "Api")
    .replace(/ID/g, "Id")
    .replace(/([A-Z]+)([A-Z][a-z])/g, "$1_$2")
    .replace(/([a-z0-9])([A-Z])/g, "$1_$2")
    .toLowerCase();
}
