export type Role = "user" | "admin";
export type Status = "active" | "disabled" | "pending" | "paid";

export type ApiEnvelope<T> = {
  code: number;
  message: string;
  data: T;
};

export type User = {
  id: number;
  email: string;
  name: string;
  role: Role;
  status: Status;
  balance_cents: number;
};

export type LoginResponse = {
  user: User;
  token: string;
};

export type ApiKey = {
  id: number;
  user_id?: number;
  name: string;
  key_prefix: string;
  status: Status;
  last_used_at: string | null;
  created_at: string;
};

export type CreatedApiKey = {
  api_key: {
    key: string;
    key_prefix: string;
  };
  record: ApiKey;
};

export type Order = {
  id: number;
  user_id: number;
  order_no: string;
  amount_cents: number;
  status: Status;
  paid_at: string | null;
  created_at: string;
};

export type UsageLog = {
  id: number;
  user_id: number;
  api_key_id: number;
  model: string;
  prompt_tokens: number;
  completion_tokens: number;
  total_tokens: number;
  cost_cents: number;
  created_at: string;
};

export type BalanceTransaction = {
  id: number;
  user_id: number;
  type: string;
  amount_cents: number;
  balance_before_cents: number;
  balance_after_cents: number;
  remark: string;
  created_at: string;
};

export type Model = {
  id: number;
  name: string;
  provider: string;
  input_price_per_1k_cents: number;
  output_price_per_1k_cents: number;
  status: Status;
  created_at: string;
};

export type AdminStats = {
  user_count: number;
  order_count: number;
  paid_amount_cents: number;
  usage_count: number;
  usage_cost_cents: number;
};

export type UsageStat = Record<string, string | number | null>;
