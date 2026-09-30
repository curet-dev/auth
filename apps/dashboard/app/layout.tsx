import type { Metadata } from "next";
import { Logo } from "@repo/ui/logo";
import { LogoutButton } from "./logout-button";
import "./globals.css";

export const metadata: Metadata = {
  title: "Dashboard – Auth",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="de">
      <body className="flex min-h-dvh flex-col">
        <header className="border-b border-zinc-200 dark:border-zinc-800">
          <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-4">
            {/* Plain <a>: the landing page lives in another zone. */}
            <a href="/">
              <Logo />
            </a>
            <LogoutButton />
          </div>
        </header>
        <main className="mx-auto w-full max-w-6xl flex-1 px-4 py-12">{children}</main>
      </body>
    </html>
  );
}
