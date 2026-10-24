import { defineConfig } from "@playwright/test";
export default defineConfig({
  testDir: "./demo",
  workers: 1,
  timeout: 90000,
  expect: { timeout: 15000 },
  outputDir: "../artifacts/portfolio/demo-run",
  reporter: [
    ["list"],
    ["json", { outputFile: "../artifacts/portfolio/demo-results.json" }],
  ],
  use: {
    baseURL: "http://localhost:8080",
    viewport: { width: 1440, height: 960 },
  },
});
