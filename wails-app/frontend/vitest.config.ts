import { defineConfig } from "vitest/config";
import { svelte } from "@sveltejs/vite-plugin-svelte";

// Separate from vite.config.ts so test-only settings never affect the app
// build. The svelte plugin is required to compile `.svelte.ts` rune modules
// (e.g. the timer store) that the tests import.
export default defineConfig({
  plugins: [svelte({ hot: false })],
  resolve: { conditions: ["browser"] },
  test: {
    environment: "jsdom",
    include: ["src/**/*.{test,spec}.ts"],
  },
});
