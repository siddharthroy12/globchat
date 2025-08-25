import tailwindcss from "@tailwindcss/vite";
import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vite";

const devOrigin = "localhost:4000";

export default defineConfig({
  plugins: [tailwindcss(), sveltekit()],
  server: {
    allowedHosts: true,
    proxy: {
      "/api/v1/ws": {
        target: `ws://${devOrigin}`,
        ws: true,
        rewriteWsOrigin: true,
        changeOrigin: true,
      },
      "/api/": {
        ws: true,
        target: `http://${devOrigin}`,
        changeOrigin: true,
      },
      "/media/": {
        target: `http://${devOrigin}`,
        changeOrigin: true,
      },
    },
  },
});
