import tailwindcss from "@tailwindcss/vite";
import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vite";

const devServer = "http://localhost:4000";

export default defineConfig({
  plugins: [tailwindcss(), sveltekit()],
  server: {
    allowedHosts: true,
    proxy: {
      "/api/": {
        target: devServer,
        changeOrigin: true,
      },
      "/media/": {
        target: devServer,
        changeOrigin: true,
      },
    },
  },
});
