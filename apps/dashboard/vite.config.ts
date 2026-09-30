import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import { defineConfig, type Plugin } from "vite";

const API_URL = process.env.API_URL ?? "http://localhost:8080";

// Mirrors the /dashboard -> /dashboard/ redirect from the root vercel.json.
const redirectToBase: Plugin = {
  name: "redirect-to-base",
  configureServer(server) {
    server.middlewares.use((req, res, next) => {
      if (req.url !== "/dashboard") return next();
      res.writeHead(307, { Location: "/dashboard/" }).end();
    });
  },
};

// The dashboard is served under /dashboard/: on Vercel via the root
// vercel.json, locally via the web app's dev proxy. The build is written to
// dist/dashboard/ so static file paths match the public URL paths.
export default defineConfig({
  base: "/dashboard/",
  plugins: [react(), tailwindcss(), redirectToBase],
  build: {
    outDir: "dist/dashboard",
    emptyOutDir: true,
  },
  server: {
    // Lets the dashboard also work when opened directly in development.
    proxy: { "/api": API_URL },
  },
});
