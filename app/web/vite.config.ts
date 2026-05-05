import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vite";
import { VitePWA } from "vite-plugin-pwa";

export default defineConfig({
  server: {
    fs: {
      allow: ["../.."]
    },
    proxy: {
      "/api": "http://localhost:3000"
    }
  },
  plugins: [
    sveltekit(),
    VitePWA({
      registerType: "autoUpdate",
      manifest: {
        name: "Apple Fest POS",
        short_name: "Apple POS",
        description: "Offline-capable point of sale for Apple Fest booths",
        theme_color: "#7c2d12",
        background_color: "#fff7ed",
        display: "standalone",
        start_url: "/pos",
        icons: [
          {
            src: "/icon.svg",
            sizes: "any",
            type: "image/svg+xml",
            purpose: "any maskable"
          }
        ]
      },
      workbox: {
        navigateFallback: "/index.html",
        runtimeCaching: [
          {
            urlPattern: ({ url }) => url.pathname.startsWith("/api/"),
            handler: "NetworkOnly"
          }
        ]
      }
    })
  ]
});
