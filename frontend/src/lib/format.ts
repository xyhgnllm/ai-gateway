export function cents(value?: number | null) {
  return `¥${((value ?? 0) / 100).toFixed(2)}`;
}

export function dateTime(value?: string | null) {
  if (!value) return "-";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "-";
  return date.toLocaleString();
}

export function statusLabel(value?: string | null) {
  if (!value) return "-";
  const labels: Record<string, string> = {
    active: "启用",
    disabled: "禁用",
    pending: "待支付",
    paid: "已支付",
    user: "用户",
    admin: "管理员",
    recharge: "充值",
    usage: "调用扣费",
  };
  return labels[value] ?? value;
}
