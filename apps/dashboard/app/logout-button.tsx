"use client";

import { Button } from "@repo/ui/button";

export function LogoutButton() {
  async function logout() {
    await fetch("/api/auth/logout", { method: "POST" });
    window.location.assign("/");
  }

  return (
    <Button variant="secondary" onClick={logout}>
      Logout
    </Button>
  );
}
