import {
  test as base,
  expect,
  type Page,
  type BrowserContext,
  type TestInfo,
} from "@playwright/test";
import fs from "node:fs";
import path from "node:path";
import { execFileSync } from "node:child_process";
import type { State, EditEvent } from "../../src/model";

export type Diagnostic = { kind: string; message: string; expected: boolean };
const redact = (value: string) =>
  value.replace(/[a-f0-9]{64}/gi, "[REDACTED_TEST_CAPABILITY]");
export const test = base.extend<{
  audit: { watch: (page: Page) => void; fault: (active: boolean) => void };
}>({
  audit: [
    async ({}, use, info) => {
      const messages: Diagnostic[] = [];
      let fault = false;
      await use({
        fault: (active) => {
          fault = active;
        },
        watch: (page) => {
          page.on("console", (m) => {
            if (m.type() === "error")
              messages.push({
                kind: "console",
                message: redact(m.text()),
                expected:
                  fault &&
                  /ERR_INTERNET_DISCONNECTED|ERR_CONNECTION_REFUSED|ERR_NETWORK_CHANGED|ERR_CONNECTION_RESET|WebSocket.*(closed|failed)|status of (401|502|503)/i.test(
                    m.text(),
                  ),
              });
          });
          page.on("pageerror", (e) =>
            messages.push({
              kind: "pageerror",
              message: e.message,
              expected: false,
            }),
          );
          page.on("requestfailed", (r) => {
            if (/\/api\//.test(r.url()))
              messages.push({
                kind: "network",
                message: redact(r.url() + " " + r.failure()?.errorText),
                expected:
                  r.failure()?.errorText === "net::ERR_ABORTED" ||
                  (fault &&
                    /ERR_INTERNET_DISCONNECTED|ERR_CONNECTION_REFUSED|ERR_NETWORK_CHANGED|ERR_CONNECTION_RESET|ERR_EMPTY_RESPONSE/.test(
                      r.failure()?.errorText ?? "",
                    )),
              });
          });
          page.on("response", (r) => {
            if (r.status() >= 400)
              messages.push({
                kind: "http",
                message: redact(`${r.status()} ${r.url()}`),
                expected:
                  fault &&
                  ([502, 503].includes(r.status()) ||
                    (r.status() === 401 &&
                      r.url().includes("qa-not-authorized"))),
              });
          });
        },
      });
      fs.mkdirSync(info.outputDir, { recursive: true });
      fs.writeFileSync(
        info.outputPath("diagnostics.json"),
        JSON.stringify(messages, null, 2),
      );
      await info.attach("Console and network diagnostics", {
        body: JSON.stringify(messages, null, 2),
        contentType: "application/json",
      });
      expect(
        messages.filter((m) => !m.expected),
        "No unexpected browser errors",
      ).toEqual([]);
    },
    { auto: true },
  ],
});
export { expect };
export const live = async (p: Page) =>
  expect(p.locator(".connection")).toContainText("Live");
export async function named(context: BrowserContext, name: string) {
  await context.addInitScript((name) => {
    if (!sessionStorage.getItem("sf:name"))
      sessionStorage.setItem("sf:name", name);
  }, name);
}
export async function create(p: Page, title: string) {
  await p.goto("/");
  await p.getByLabel("Board title").fill(title);
  await p.getByRole("button", { name: "New board", exact: true }).click();
  await live(p);
  return p.url().split("/").at(-1)!;
}
export async function state(p: Page): Promise<State> {
  return p.evaluate(async () => {
    const id = location.pathname.split("/").at(-1);
    const response = await fetch(`/api/boards/${id}`);
    if (!response.ok) throw Error("State API failed: " + response.status);
    return response.json();
  });
}
export async function events(p: Page): Promise<EditEvent[]> {
  return p.evaluate(async () => {
    const id = location.pathname.split("/").at(-1);
    let all: any[] = [];
    for (;;) {
      const r = await fetch(
        `/api/boards/${id}/events?after=${all.at(-1)?.revision || 0}`,
      );
      if (!r.ok) throw Error("Events API failed");
      const page = await r.json();
      all = all.concat(page);
      if (page.length < 500) return all;
    }
  });
}
export async function settled(p: Page) {
  await expect(p.locator(".connection")).not.toContainText("pending");
  await expect
    .poll(async () =>
      Number((await p.getByTestId("revision").innerText()).replace("REV ", "")),
    )
    .toBe((await state(p)).revision);
}
export async function invite(p: Page, context: BrowserContext) {
  await context.grantPermissions(["clipboard-read", "clipboard-write"], {
    origin: new URL(p.url()).origin,
  });
  await p.getByRole("button", { name: "Share board", exact: true }).click();
  await expect(
    p.getByRole("status").filter({ hasText: "Invite link copied" }),
  ).toBeVisible();
  const link = await p.evaluate(() => navigator.clipboard.readText());
  expect(link).toMatch(/\/board\/[^#]+#key=[a-f0-9]{64}$/);
  return link;
}
export async function add(p: Page, kind: string, x: number, y: number) {
  await p.getByRole("button", { name: kind, exact: true }).click();
  const before = (await state(p)).revision;
  await p.getByTestId("canvas").click({ position: { x, y } });
  await expect.poll(async () => (await state(p)).revision).toBe(before + 1);
  await settled(p);
  return (await events(p)).at(-1)!.object.id;
}
export async function snap(p: Page, info: TestInfo, label: string) {
  const folder = path.resolve(
    "../artifacts/browser-qa/screenshots",
    process.env.QA_RUN_ID || "latest",
  );
  fs.mkdirSync(folder, { recursive: true });
  const file = path.join(folder, label + ".png");
  await p.screenshot({
    path: file,
    fullPage: true,
    mask: [p.getByLabel("Invite link")],
  });
  await info.attach(label, { path: file, contentType: "image/png" });
}
export async function assertPersisted(p: Page) {
  const st = await state(p);
  const es = await events(p);
  expect(es.map((e) => e.revision)).toEqual(
    Array.from({ length: st.revision }, (_, i) => i + 1),
  );
  expect(new Set(es.map((e) => e.operation.id)).size).toBe(es.length);
  const replay: State = { ...st, revision: 0, objects: {} };
  for (const e of es) {
    replay.objects[e.object.id] = e.object;
    replay.revision = e.revision;
  }
  expect(replay).toEqual(st);
  // Compare the actual SQL state, not just the HTTP representation.
  const result = execFileSync(
    "docker",
    [
      "compose",
      "-p",
      "syncframe",
      "exec",
      "-T",
      "mysql",
      "mysql",
      "-N",
      "-B",
      "--raw",
      "-usyncframe",
      "-psyncframe",
      "-D",
      "syncframe",
      "-e",
      `SELECT state FROM boards WHERE id='${st.id}'`,
    ],
    {
      cwd: path.resolve(".."),
      encoding: "utf8",
      stdio: ["ignore", "pipe", "pipe"],
    },
  );
  expect(JSON.parse(result.trim())).toEqual(st);
  return st;
}
