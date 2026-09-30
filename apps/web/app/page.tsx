import Link from "next/link";
import { buttonClass } from "@repo/ui/button";

const features = [
  { title: "Next.js", text: "Landing Page und Dashboard als eigenständige Apps." },
  { title: "Go Backend", text: "Schlanke API mit JWT-Session-Cookies." },
  { title: "Vercel ready", text: "Jede App als eigenes Vercel-Projekt deploybar." },
];

export default function Home() {
  return (
    <div className="mx-auto w-full max-w-6xl px-4">
      <section className="py-24 text-center">
        <h1 className="text-4xl font-bold tracking-tight sm:text-6xl">
          Starte dein Projekt in Minuten
        </h1>
        <p className="mx-auto mt-6 max-w-xl text-lg text-zinc-600 dark:text-zinc-400">
          Ein Turborepo mit Next.js Frontend, Dashboard und Go Backend – bereit für Vercel.
        </p>
        <div className="mt-10 flex justify-center gap-3">
          <Link href="/register" className={buttonClass("primary")}>
            Jetzt registrieren
          </Link>
          <Link href="/login" className={buttonClass("secondary")}>
            Ich habe schon einen Account
          </Link>
        </div>
      </section>
      <section className="grid gap-4 pb-24 sm:grid-cols-3">
        {features.map((f) => (
          <div key={f.title} className="rounded-lg border border-zinc-200 p-6 dark:border-zinc-800">
            <h2 className="font-semibold">{f.title}</h2>
            <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">{f.text}</p>
          </div>
        ))}
      </section>
    </div>
  );
}
