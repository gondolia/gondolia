/** @type {import('next').NextConfig} */
const nextConfig = {
  output: "standalone",

  // Proxy API calls during local development (when running outside Docker).
  // In Docker, Traefik sits in front and routes /api/* to the backend services
  // directly — no rewrite needed there.
  // Locally (npm run dev), set API_PROXY_TARGET in .env.local to point at the
  // running Traefik/Docker gateway (e.g. http://localhost:3001).
  async rewrites() {
    const target = process.env.API_PROXY_TARGET;
    if (!target) return [];
    return [
      {
        source: "/api/:path*",
        destination: `${target}/api/:path*`,
      },
    ];
  },
};

export default nextConfig;
