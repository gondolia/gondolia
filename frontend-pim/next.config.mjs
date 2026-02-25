/** @type {import('next').NextConfig} */
const nextConfig = {
  output: "standalone",

  // Allow access from custom host (aioc) for remote development
  experimental: {
    allowedOrigins: ['aioc', 'localhost', '127.0.0.1'],
  },

  // Proxy API calls during local development
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
