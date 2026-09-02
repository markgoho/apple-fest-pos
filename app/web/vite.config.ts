import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vite";

export default defineConfig({
  server: {
    fs: {
      allow: ["../.."]
    },
    proxy: {
      "/api": "http://localhost:3000"
    }
  },
  plugins: [sveltekit()]
});
