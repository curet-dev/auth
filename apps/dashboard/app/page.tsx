"use client";

import { useEffect, useState } from "react";

type User = { id: string; email: string; name: string };

export default function DashboardPage() {
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    fetch("/api/auth/me")
      .then((res) => {
        if (res.status === 401) {
          window.location.assign("/login");
          return null;
        }
        return res.json() as Promise<{ user: User }>;
      })
      .then((data) => data && setUser(data.user))
      .catch(() => window.location.assign("/login"));
  }, []);

  return (
    <div>
      <h1 className="text-3xl font-bold tracking-tight">
        {user ? `Hallo, ${user.name}!` : "Dashboard"}
      </h1>
      <p className="mt-2 text-zinc-600 dark:text-zinc-400">
        {user ? `Eingeloggt als ${user.email}` : "Lade…"}
      </p>
      <div className="mt-10 grid gap-4 sm:grid-cols-3">
        {["Projekte", "Aktivität", "Einstellungen"].map((title) => (
          <div key={title} className="rounded-lg border border-zinc-200 p-6 dark:border-zinc-800">
            <h2 className="font-semibold">{title}</h2>
            <p className="mt-2 text-sm text-zinc-500">Hier kommt bald Inhalt hin.</p>
          </div>
        ))}
      </div>
    </div>
  );
}
