import type { NextConfig } from "next";

// The web app is the public entrypoint. It proxies the dashboard (Next.js
// multi-zone) and the Go API so everything shares one origin and the
// session cookie just works.
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
