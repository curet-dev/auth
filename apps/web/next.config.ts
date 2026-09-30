import type { NextConfig } from "next";

// On Vercel the root vercel.json routes /api and /dashboard to their services
// before a request reaches this app. For local development the same routing
// is done here, so the dashboard (Next.js multi-zone) and the Go API share
// one origin and the session cookie just works.
const API_URL = process.env.API_URL ?? "http://localhost:8080";
const DASHBOARD_URL = process.env.DASHBOARD_URL ?? "http://localhost:3001";

const nextConfig: NextConfig = {
  transpilePackages: ["@repo/ui"],
  async rewrites() {
    return [
      { source: "/api/:path*", destination: `${API_URL}/api/:path*` },
      { source: "/dashboard", destination: `${DASHBOARD_URL}/dashboard` },
      { source: "/dashboard/:path*", destination: `${DASHBOARD_URL}/dashboard/:path*` },
    ];
  },
};

export default nextConfig;
