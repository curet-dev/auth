export type User = { id: string; email: string; name: string };

// /api and /login live on the same origin (see the root vercel.json).
export async function fetchMe(): Promise<User | null> {
  const res = await fetch("/api/auth/me");
  if (res.status === 401) return null;
  if (!res.ok) throw new Error(`GET /api/auth/me: ${res.status}`);
  const data = (await res.json()) as { user: User };
  return data.user;
}

export async function logout() {
  await fetch("/api/auth/logout", { method: "POST" });
}

export function goTo(path: "/" | "/login") {
  window.location.assign(path);
}
