/** @type {import('next').NextConfig} */
const nextConfig = {
  async rewrites() {
    return [
      {
        source: '/api/search',
        destination: 'http://localhost:8080/api/search', // Go backend
      },
    ];
  },
  reactStrictMode: true,
};

export default nextConfig;
