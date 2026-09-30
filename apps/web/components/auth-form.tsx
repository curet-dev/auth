"use client";

import Link from "next/link";
import { useState, type FormEvent } from "react";
import { Button } from "@repo/ui/button";

type Mode = "login" | "register";

const copy = {
  login: {
    title: "Willkommen zurück",
    submit: "Einloggen",
    switchText: "Noch keinen Account?",
    switchLink: "Registrieren",
    switchHref: "/register",
  },
  register: {
    title: "Account erstellen",
    submit: "Registrieren",
    switchText: "Schon registriert?",
    switchLink: "Einloggen",
    switchHref: "/login",
  },
} as const;

const inputClass =
  "mt-1 block w-full rounded-md border border-zinc-300 bg-white px-3 py-2 text-sm outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 dark:border-zinc-700 dark:bg-zinc-900";

export function AuthForm({ mode }: { mode: Mode }) {
  const [error, setError] = useState<string | null>(null);
  const [pending, setPending] = useState(false);
  const t = copy[mode];

  async function onSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setError(null);
    setPending(true);
    try {
      const res = await fetch(`/api/auth/${mode}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(Object.fromEntries(new FormData(e.currentTarget))),
      });
      if (!res.ok) {
        const data = (await res.json().catch(() => ({}))) as { error?: string };
        throw new Error(data.error ?? "Etwas ist schiefgelaufen");
      }
      // The dashboard is a separate Next.js zone, so use a full navigation.
      window.location.assign("/dashboard");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Etwas ist schiefgelaufen");
      setPending(false);
    }
  }

  return (
    <div className="mx-auto w-full max-w-sm px-4 py-24">
      <h1 className="text-2xl font-bold tracking-tight">{t.title}</h1>
      <form onSubmit={onSubmit} className="mt-8 space-y-4">
        {mode === "register" && (
          <label className="block text-sm font-medium">
            Name
            <input name="name" required autoComplete="name" className={inputClass} />
          </label>
        )}
        <label className="block text-sm font-medium">
          E-Mail
          <input name="email" type="email" required autoComplete="email" className={inputClass} />
        </label>
        <label className="block text-sm font-medium">
          Passwort
          <input
            name="password"
            type="password"
            required
            minLength={mode === "register" ? 8 : undefined}
            autoComplete={mode === "register" ? "new-password" : "current-password"}
            className={inputClass}
          />
        </label>
        {error && (
          <p role="alert" className="text-sm text-red-600">
            {error}
          </p>
        )}
        <Button type="submit" disabled={pending} className="w-full">
          {pending ? "Bitte warten…" : t.submit}
        </Button>
      </form>
      <p className="mt-6 text-center text-sm text-zinc-600 dark:text-zinc-400">
        {t.switchText}{" "}
        <Link href={t.switchHref} className="font-medium text-indigo-600 hover:underline">
          {t.switchLink}
        </Link>
      </p>
    </div>
  );
}
