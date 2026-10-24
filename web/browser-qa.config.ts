import { defineConfig } from "@playwright/test";
import base from "./playwright.config";
const run =
  process.env.QA_RUN_ID || new Date().toISOString().replace(/[:.]/g, "-");
export default defineConfig(base, {
  timeout: 90000,
  outputDir: `../artifacts/browser-qa/runs/${run}/failures`,
  reporter: [
    ["list"],
    [
      "json",
      { outputFile: `../artifacts/browser-qa/runs/${run}/results.json` },
    ],
  ],
  use: {
    baseURL: process.env.BASE_URL || "http://localhost:8080",
    actionTimeout: 15000,
    navigationTimeout: 20000,
    trace: "retain-on-failure",
    video: "retain-on-failure",
    screenshot: "only-on-failure",
  },
});
