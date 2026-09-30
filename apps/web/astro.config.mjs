import react from "@astrojs/react";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "astro/config";

// On Vercel the root vercel.json routes /api and /dashboard to their
// services. For local development the dev server proxies them instead, so
// everything shares one origin and the session cookie just works.
const API_URL = process.env.API_URL ?? "http://localhost:8080";
const DASHBOARD_URL = process.env.DASHBOARD_URL ?? "http://localhost:3001";

export default defineConfig({
  integrations: [react()],
  vite: {
    plugins: [tailwindcss()],
    server: {
      proxy: {
        "/api": API_URL,
        "/dashboard": { target: DASHBOARD_URL, ws: true },
      },
    },
  },
});
