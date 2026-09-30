import { useEffect, useState } from "react";
import { Button } from "@repo/ui/button";
import { Logo } from "@repo/ui/logo";
import { fetchMe, goTo, logout, type User } from "./api";

const cards = ["Projekte", "Aktivität", "Einstellungen"];

export function App() {
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    fetchMe()
      .then((u) => (u ? setUser(u) : goTo("/login")))
      .catch(() => goTo("/login"));
  }, []);

  async function onLogout() {
    await logout();
    goTo("/");
  }

  // Render nothing until the session is confirmed to avoid a flash of the
  // dashboard for signed-out users.
  if (!user) return null;

  return (
    <>
      <header className="border-b border-zinc-200 dark:border-zinc-800">
        <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-4">
          <a href="/">
            <Logo />
          </a>
          <Button variant="secondary" onClick={onLogout}>
            Logout
          </Button>
        </div>
      </header>
      <main className="mx-auto w-full max-w-6xl flex-1 px-4 py-12">
        <h1 className="text-3xl font-bold tracking-tight">Hallo, {user.name}!</h1>
        <p className="mt-2 text-zinc-600 dark:text-zinc-400">Eingeloggt als {user.email}</p>
        <div className="mt-10 grid gap-4 sm:grid-cols-3">
          {cards.map((title) => (
            <div key={title} className="rounded-lg border border-zinc-200 p-6 dark:border-zinc-800">
              <h2 className="font-semibold">{title}</h2>
              <p className="mt-2 text-sm text-zinc-500">Hier kommt bald Inhalt hin.</p>
            </div>
          ))}
        </div>
      </main>
    </>
  );
}
