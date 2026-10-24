# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: browser-qa.spec.ts >> initial board loading is visible while a real state response is delayed
- Location: e2e\browser-qa.spec.ts:450:1

# Error details

```
Test timeout of 90000ms exceeded.
```

```
Error: page.evaluate: Test timeout of 90000ms exceeded.
```

# Page snapshot

```yaml
- generic [ref=f3e3]:
  - banner [ref=f3e4]:
    - generic [ref=f3e5]:
      - link "SYNCFRAME home" [ref=f3e6] [cursor=pointer]:
        - /url: /
        - generic [ref=f3e9]: syncframe.
      - link "Back to workspace" [ref=f3e11] [cursor=pointer]:
        - /url: /
      - generic [ref=f3e14]:
        - strong [ref=f3e15]: QA initial loading
        - generic [ref=f3e16]: Shared workspace Canvas
    - generic [ref=f3e19]:
      - status [ref=f3e20]: Live
      - generic "Maker 270" [ref=f3e22]: M
      - button "History" [expanded] [ref=f3e24] [cursor=pointer]
      - button "Share board" [ref=f3e29] [cursor=pointer]
  - generic [ref=f3e36]:
    - complementary [ref=f3e37]:
      - button "Canvas" [ref=f3e38] [cursor=pointer]
      - button "History" [ref=f3e41] [cursor=pointer]
      - generic [ref=f3e46]:
        - generic [ref=f3e47]: M
        - generic [ref=f3e48]: You
    - main [ref=f3e49]:
      - generic [ref=f3e50]:
        - generic [ref=f3e51]:
          - text: Live canvas
          - generic [ref=f3e56]: /
          - strong [ref=f3e57]: 0 objects
        - generic [ref=f3e58]: REV 0
      - application "Collaborative whiteboard" [ref=f3e60]
      - generic:
        - strong: Every idea starts somewhere.
        - generic: Choose a tool below, then click to begin.
      - generic:
        - generic [ref=f3e62]:
          - button "Zoom out" [ref=f3e63] [cursor=pointer]
          - generic [ref=f3e67]: 100%
          - button "Zoom in" [ref=f3e68] [cursor=pointer]
          - button "Reset zoom" [ref=f3e73] [cursor=pointer]
        - toolbar "Drawing tools" [ref=f3e79]:
          - button "Select" [pressed] [ref=f3e80] [cursor=pointer]:
            - generic [ref=f3e83]: V
          - button "Rectangle" [ref=f3e84] [cursor=pointer]:
            - generic [ref=f3e87]: R
          - button "Ellipse" [ref=f3e88] [cursor=pointer]:
            - generic [ref=f3e91]: O
          - button "Sticky note" [ref=f3e92] [cursor=pointer]:
            - generic [ref=f3e96]: "N"
          - button "Text" [ref=f3e97] [cursor=pointer]:
            - generic [ref=f3e100]: T
        - generic [ref=f3e101]: Select, drag, make it yours
    - complementary [ref=f3e102]:
      - generic [ref=f3e103]:
        - heading "Time machine" [level=2] [ref=f3e108]
        - button "Close history" [ref=f3e109] [cursor=pointer]
      - paragraph [ref=f3e113]: Every idea has a journey.Revisit any moment of yours.
      - generic [ref=f3e114]:
        - generic [ref=f3e115]:
          - generic [ref=f3e116]: LIVE FRAME
          - generic [ref=f3e118]:
            - text: "0"
            - generic [ref=f3e119]: / 0
        - slider "History revision" [ref=f3e120] [cursor=pointer]: "0"
        - generic [ref=f3e121]:
          - generic [ref=f3e122]: Beginning
          - generic [ref=f3e123]: Now
        - generic [ref=f3e124]:
          - button "First event" [ref=f3e125] [cursor=pointer]
          - button "Previous event" [ref=f3e128] [cursor=pointer]
          - button "Play history" [disabled] [ref=f3e131]
          - button "Next event" [ref=f3e134] [cursor=pointer]
          - button "Latest event" [ref=f3e137] [cursor=pointer]
          - combobox "Playback speed" [ref=f3e140]:
            - option "0.5×"
            - option "1×" [selected]
            - option "2×"
            - option "4×"
      - generic [ref=f3e141]:
        - generic [ref=f3e142]: ACTIVITY
        - generic [ref=f3e143]: 0 events
      - status [ref=f3e145]:
        - strong [ref=f3e150]: The story starts here.
        - paragraph [ref=f3e151]: Your edits will appear here,one frame at a time.
      - generic [ref=f3e152]: History is saved to MySQL.Exploring the past keeps the live canvas intact.
  - contentinfo [ref=f3e156]:
    - generic [ref=f3e157]:
      - text: All committed edits saved
      - generic [ref=f3e159]: /
      - text: Revision 0
    - generic [ref=f3e160]:
      - text: 1 online
      - generic [ref=f3e166]: /
      - text: You are Maker 270
```

# Test source

```ts
  1   | import {
  2   |   test as base,
  3   |   expect,
  4   |   type Page,
  5   |   type BrowserContext,
  6   |   type TestInfo,
  7   | } from "@playwright/test";
  8   | import fs from "node:fs";
  9   | import path from "node:path";
  10  | import { execFileSync } from "node:child_process";
  11  | import type { State, EditEvent } from "../../src/model";
  12  | 
  13  | export type Diagnostic = { kind: string; message: string; expected: boolean };
  14  | const redact = (value: string) =>
  15  |   value.replace(/[a-f0-9]{64}/gi, "[REDACTED_TEST_CAPABILITY]");
  16  | export const test = base.extend<{
  17  |   audit: { watch: (page: Page) => void; fault: (active: boolean) => void };
  18  | }>({
  19  |   audit: [
  20  |     async ({}, use, info) => {
  21  |       const messages: Diagnostic[] = [];
  22  |       let fault = false;
  23  |       await use({
  24  |         fault: (active) => {
  25  |           fault = active;
  26  |         },
  27  |         watch: (page) => {
  28  |           page.on("console", (m) => {
  29  |             if (m.type() === "error")
  30  |               messages.push({
  31  |                 kind: "console",
  32  |                 message: redact(m.text()),
  33  |                 expected: fault,
  34  |               });
  35  |           });
  36  |           page.on("pageerror", (e) =>
  37  |             messages.push({
  38  |               kind: "pageerror",
  39  |               message: e.message,
  40  |               expected: false,
  41  |             }),
  42  |           );
  43  |           page.on("requestfailed", (r) => {
  44  |             if (/\/api\//.test(r.url()))
  45  |               messages.push({
  46  |                 kind: "network",
  47  |                 message: redact(r.url() + " " + r.failure()?.errorText),
  48  |                 expected:
  49  |                   fault || r.failure()?.errorText === "net::ERR_ABORTED",
  50  |               });
  51  |           });
  52  |           page.on("response", (r) => {
  53  |             if (r.status() >= 400)
  54  |               messages.push({
  55  |                 kind: "http",
  56  |                 message: redact(`${r.status()} ${r.url()}`),
  57  |                 expected: fault,
  58  |               });
  59  |           });
  60  |         },
  61  |       });
  62  |       fs.mkdirSync(info.outputDir, { recursive: true });
  63  |       fs.writeFileSync(
  64  |         info.outputPath("diagnostics.json"),
  65  |         JSON.stringify(messages, null, 2),
  66  |       );
  67  |       await info.attach("Console and network diagnostics", {
  68  |         body: JSON.stringify(messages, null, 2),
  69  |         contentType: "application/json",
  70  |       });
  71  |       expect(
  72  |         messages.filter((m) => !m.expected),
  73  |         "No unexpected browser errors",
  74  |       ).toEqual([]);
  75  |     },
  76  |     { auto: true },
  77  |   ],
  78  | });
  79  | export { expect };
  80  | export const live = async (p: Page) =>
  81  |   expect(p.locator(".connection")).toContainText("Live");
  82  | export async function named(context: BrowserContext, name: string) {
  83  |   await context.addInitScript((name) => {
  84  |     if (!sessionStorage.getItem("sf:name"))
  85  |       sessionStorage.setItem("sf:name", name);
  86  |   }, name);
  87  | }
  88  | export async function create(p: Page, title: string) {
  89  |   await p.goto("/");
  90  |   await p.getByLabel("Board title").fill(title);
  91  |   await p.getByRole("button", { name: "New board", exact: true }).click();
  92  |   await live(p);
  93  |   return p.url().split("/").at(-1)!;
  94  | }
  95  | export async function state(p: Page): Promise<State> {
> 96  |   return p.evaluate(async () => {
      |            ^ Error: page.evaluate: Test timeout of 90000ms exceeded.
  97  |     const id = location.pathname.split("/").at(-1);
  98  |     const response = await fetch(`/api/boards/${id}`);
  99  |     if (!response.ok) throw Error("State API failed: " + response.status);
  100 |     return response.json();
  101 |   });
  102 | }
  103 | export async function events(p: Page): Promise<EditEvent[]> {
  104 |   return p.evaluate(async () => {
  105 |     const id = location.pathname.split("/").at(-1);
  106 |     let all: any[] = [];
  107 |     for (;;) {
  108 |       const r = await fetch(
  109 |         `/api/boards/${id}/events?after=${all.at(-1)?.revision || 0}`,
  110 |       );
  111 |       if (!r.ok) throw Error("Events API failed");
  112 |       const page = await r.json();
  113 |       all = all.concat(page);
  114 |       if (page.length < 500) return all;
  115 |     }
  116 |   });
  117 | }
  118 | export async function settled(p: Page) {
  119 |   await expect(p.locator(".connection")).not.toContainText("pending");
  120 |   await expect
  121 |     .poll(async () =>
  122 |       Number((await p.getByTestId("revision").innerText()).replace("REV ", "")),
  123 |     )
  124 |     .toBe((await state(p)).revision);
  125 | }
  126 | export async function invite(p: Page, context: BrowserContext) {
  127 |   await context.grantPermissions(["clipboard-read", "clipboard-write"], {
  128 |     origin: new URL(p.url()).origin,
  129 |   });
  130 |   await p.getByRole("button", { name: "Share board", exact: true }).click();
  131 |   await expect(
  132 |     p.getByRole("status").filter({ hasText: "Invite link copied" }),
  133 |   ).toBeVisible();
  134 |   const link = await p.evaluate(() => navigator.clipboard.readText());
  135 |   expect(link).toMatch(/\/board\/[^#]+#key=[a-f0-9]{64}$/);
  136 |   return link;
  137 | }
  138 | export async function add(p: Page, kind: string, x: number, y: number) {
  139 |   await p.getByRole("button", { name: kind, exact: true }).click();
  140 |   const before = (await state(p)).revision;
  141 |   await p.getByTestId("canvas").click({ position: { x, y } });
  142 |   await expect.poll(async () => (await state(p)).revision).toBe(before + 1);
  143 |   await settled(p);
  144 |   const st = await state(p);
  145 |   return (await events(p)).at(-1)!.object.id;
  146 | }
  147 | export async function snap(p: Page, info: TestInfo, label: string) {
  148 |   const folder = path.resolve(
  149 |     "../artifacts/browser-qa/screenshots",
  150 |     process.env.QA_RUN_ID || "latest",
  151 |   );
  152 |   fs.mkdirSync(folder, { recursive: true });
  153 |   const file = path.join(folder, label + ".png");
  154 |   await p.screenshot({
  155 |     path: file,
  156 |     fullPage: true,
  157 |     mask: [p.getByLabel("Invite link")],
  158 |   });
  159 |   await info.attach(label, { path: file, contentType: "image/png" });
  160 | }
  161 | export async function assertPersisted(p: Page) {
  162 |   const st = await state(p);
  163 |   const es = await events(p);
  164 |   expect(es.map((e) => e.revision)).toEqual(
  165 |     Array.from({ length: st.revision }, (_, i) => i + 1),
  166 |   );
  167 |   expect(new Set(es.map((e) => e.operation.id)).size).toBe(es.length);
  168 |   const replay: State = { ...st, revision: 0, objects: {} };
  169 |   for (const e of es) {
  170 |     replay.objects[e.object.id] = e.object;
  171 |     replay.revision = e.revision;
  172 |   }
  173 |   expect(replay).toEqual(st);
  174 |   // Compare the actual SQL state, not just the HTTP representation.
  175 |   const result = execFileSync(
  176 |     "docker",
  177 |     [
  178 |       "compose",
  179 |       "-p",
  180 |       "syncframe",
  181 |       "exec",
  182 |       "-T",
  183 |       "mysql",
  184 |       "mysql",
  185 |       "-N",
  186 |       "-B",
  187 |       "--raw",
  188 |       "-usyncframe",
  189 |       "-psyncframe",
  190 |       "-D",
  191 |       "syncframe",
  192 |       "-e",
  193 |       `SELECT state FROM boards WHERE id='${st.id}'`,
  194 |     ],
  195 |     {
  196 |       cwd: path.resolve(".."),
```