import Link from "next/link";
import { buttonClass } from "@repo/ui/button";
import { Logo } from "@repo/ui/logo";

export function Header() {
  return (
    <header className="border-b border-zinc-200 dark:border-zinc-800">
      <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-4">
        <Link href="/">
          <Logo />
        </Link>
        <nav className="flex items-center gap-2">
          <Link href="/login" className={buttonClass("secondary")}>
            Login
          </Link>
          <Link href="/register" className={buttonClass("primary")}>
            Registrieren
          </Link>
        </nav>
      </div>
    </header>
  );
}
