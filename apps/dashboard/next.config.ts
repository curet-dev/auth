import type { NextConfig } from "next";

// The dashboard is served under /dashboard: on Vercel via the root
// vercel.json, locally via the web app's rewrites (Next.js multi-zone).
// /api is proxied so the app also works when opened directly in development.
const API_URL = process.env.API_URL ?? "http://localhost:8080";

const nextConfig: NextConfig = {
  basePath: "/dashboard",
  transpilePackages: ["@repo/ui"],
  async rewrites() {
    return [{ source: "/api/:path*", destination: `${API_URL}/api/:path*`, basePath: false }];
  },
};

export default nextConfig;
