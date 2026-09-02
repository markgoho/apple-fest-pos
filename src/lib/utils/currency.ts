const usdFormatter = new Intl.NumberFormat("en-US", {
  style: "currency",
  currency: "USD",
  minimumFractionDigits: 0,
  maximumFractionDigits: 0
});

export function formatCents(cents: number): string {
  return usdFormatter.format(cents / 100);
}

export function parseCashInput(value: string): number | null {
  const normalized = value.trim().replace(/[$,]/g, "");
  if (!normalized) {
    return null;
  }

  const dollars = Number(normalized);
  if (!Number.isFinite(dollars) || dollars < 0) {
    return null;
  }

  return Math.round(dollars * 100);
}
