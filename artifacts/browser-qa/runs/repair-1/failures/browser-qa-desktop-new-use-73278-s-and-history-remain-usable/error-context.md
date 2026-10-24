# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: browser-qa.spec.ts >> desktop: new-user editing, zoomed drag, properties and history remain usable
- Location: e2e\browser-qa.spec.ts:249:3

# Error details

```
Error: expect(locator).toBeVisible() failed

Locator: getByRole('heading', { name: 'Time machine' })
Expected: visible
Timeout: 12000ms
Error: element(s) not found

Call log:
  - Expect "toBeVisible" getByRole('heading', { name: 'Time machine' }) with timeout 12000ms
  - waiting for getByRole('heading', { name: 'Time machine' })

```

```yaml
- banner:
  - link "SYNCFRAME home":
    - /url: /
    - img
    - text: syncframe.
  - link "Back to workspace":
    - /url: /
    - img
  - strong: QA desktop workspace
  - text: Shared workspace
  - img
  - text: Canvas
  - status: Live
  - text: Q
  - button "History":
    - img
    - text: History
  - button "Share board":
    - img
    - text: Share board
- complementary:
  - button "Canvas":
    - img
    - text: Canvas
  - button "History":
    - img
    - text: History
  - text: Q You
- main:
  - img
  - text: Live canvas /
  - strong: 1 objects
  - text: REV 6
  - application "Collaborative whiteboard":
    - 'button "note: A readable idea on desktop"': A readable idea on desktop
  - button "Zoom out":
    - img
  - text: 80%
  - button "Zoom in":
    - img
  - button "Reset zoom":
    - img
  - toolbar "Drawing tools":
    - button "Select" [pressed]:
      - img
      - text: V
    - button "Rectangle":
      - img
      - text: R
    - button "Ellipse":
      - img
      - text: O
    - button "Sticky note":
      - img
      - text: "N"
    - button "Text":
      - img
      - text: T
  - text: Select, drag, make it yours
- contentinfo:
  - text: All committed edits saved / Revision 6
  - img
  - text: 1 online / You are QA desktop
```

# Test source

```ts
  208 |     .poll(async () =>
  209 |       Number((await a.getByTestId("revision").innerText()).replace("REV ", "")),
  210 |     )
  211 |     .toBeGreaterThan(0);
  212 |   await a.getByRole("button", { name: "Pause history", exact: true }).click();
  213 |   const paused = await a.getByTestId("revision").innerText();
  214 |   // Advance beyond more than two playback intervals without a sleep-based oracle.
  215 |   await a.clock.install();
  216 |   await a.clock.fastForward(5000);
  217 |   await expect(a.getByTestId("revision")).toHaveText(paused);
  218 |   await a.clock.resume();
  219 |   await a.getByRole("checkbox", { name: "Overlay live positions" }).check();
  220 |   await expect(a.locator('svg.canvas > rect[stroke="#16a57a"]')).toHaveCount(1);
  221 |   await b
  222 |     .locator(`[data-object-id="${id}"]`)
  223 |     .click({ position: { x: 90, y: 50 } });
  224 |   await b.getByLabel("Object text").fill("Live changed during playback");
  225 |   await b.getByLabel("Object text").press("Tab");
  226 |   await expect
  227 |     .poll(async () => (await state(b)).revision)
  228 |     .toBe(original.revision + 1);
  229 |   await expect(a.getByTestId("revision")).toHaveText(paused);
  230 |   expect((await state(a)).objects[id].text).toBe(
  231 |     "Live changed during playback",
  232 |   );
  233 |   await snap(a, info, "history-overlay");
  234 |   await a.getByRole("button", { name: "Canvas", exact: true }).click();
  235 |   await expect(a.locator(".travel-banner")).toHaveCount(0);
  236 |   await expect(a.locator(`[data-object-id="${id}"]`)).toContainText(
  237 |     "Live changed during playback",
  238 |   );
  239 |   await assertPersisted(a);
  240 |   await one.close();
  241 |   await two.close();
  242 | });
  243 | 
  244 | for (const size of [
  245 |   { name: "desktop", width: 1440, height: 960 },
  246 |   { name: "tablet", width: 768, height: 1024 },
  247 |   { name: "mobile", width: 390, height: 844 },
  248 | ]) {
  249 |   test(`${size.name}: new-user editing, zoomed drag, properties and history remain usable`, async ({
  250 |     browser,
  251 |     audit,
  252 |   }, info) => {
  253 |     const context = await browser.newContext({ viewport: size });
  254 |     await named(context, "QA " + size.name);
  255 |     const p = await context.newPage();
  256 |     audit.watch(p);
  257 |     await p.goto("/");
  258 |     await snap(p, info, size.name + "-home");
  259 |     expect(
  260 |       await p.evaluate(
  261 |         () => document.documentElement.scrollWidth <= innerWidth,
  262 |       ),
  263 |     ).toBe(true);
  264 |     await create(p, "QA " + size.name + " workspace");
  265 |     await snap(p, info, size.name + "-empty");
  266 |     await expect(
  267 |       p.getByRole("button", { name: "Sticky note", exact: true }),
  268 |     ).toBeInViewport();
  269 |     // A new user's first click must reach the canvas, including on a phone.
  270 |     const objectId = await add(p, "Sticky note", 90, 140);
  271 |     const object = p.locator(`[data-object-id="${objectId}"]`);
  272 |     await expect(object).toBeVisible();
  273 |     await p.getByLabel("Object text").fill("A readable idea on " + size.name);
  274 |     await p.getByLabel("Object text").press("Tab");
  275 |     await expect(object).toContainText("A readable idea");
  276 |     await p.getByRole("button", { name: "Color #e4dcfc", exact: true }).click();
  277 |     await expect(object.locator("rect").first()).toHaveAttribute(
  278 |       "fill",
  279 |       "#e4dcfc",
  280 |     );
  281 |     await snap(p, info, size.name + "-properties");
  282 |     const box = await object.boundingBox();
  283 |     await p.mouse.move(box!.x + 75, box!.y + 45);
  284 |     await p.mouse.down();
  285 |     await p.mouse.move(box!.x + 105, box!.y + 75, { steps: 6 });
  286 |     await p.mouse.up();
  287 |     await expect(object).toHaveAttribute("transform", "translate(120,170)");
  288 |     await p.getByLabel("Close properties").click();
  289 |     await p.getByRole("button", { name: "Zoom out", exact: true }).click();
  290 |     await p.getByRole("button", { name: "Zoom out", exact: true }).click();
  291 |     await expect(p.locator(".zoom")).toContainText("80%");
  292 |     const zoomBox = await object.boundingBox();
  293 |     await p.mouse.move(zoomBox!.x + 50, zoomBox!.y + 35);
  294 |     await p.mouse.down();
  295 |     await p.mouse.move(zoomBox!.x + 90, zoomBox!.y + 59, { steps: 5 });
  296 |     await p.mouse.up();
  297 |     await expect(object).toHaveAttribute("transform", "translate(170,200)");
  298 |     await object.press("ArrowRight");
  299 |     await expect(object).toHaveAttribute("transform", "translate(180,200)");
  300 |     await snap(p, info, size.name + "-edited");
  301 |     await p.getByLabel("Close properties").click();
  302 |     await p
  303 |       .getByRole("button", { name: "History", exact: true })
  304 |       .first()
  305 |       .click();
  306 |     await expect(
  307 |       p.getByRole("heading", { name: "Time machine" }),
> 308 |     ).toBeVisible();
      |       ^ Error: expect(locator).toBeVisible() failed
  309 |     await snap(p, info, size.name + "-history");
  310 |     await p.getByRole("button", { name: "Close history", exact: true }).click();
  311 |     const scan = await new AxeBuilder({ page: p })
  312 |       .withTags(["wcag2a", "wcag2aa", "wcag21aa"])
  313 |       .analyze();
  314 |     expect(scan.violations).toEqual([]);
  315 |     const durable = await assertPersisted(p);
  316 |     await p.reload();
  317 |     await live(p);
  318 |     expect(await state(p)).toEqual(durable);
  319 |     await p.locator(`[data-object-id="${objectId}"]`).press("Delete");
  320 |     await expect(p.getByTestId("board-object")).toHaveCount(0);
  321 |     await assertPersisted(p);
  322 |     await context.close();
  323 |   });
  324 | }
  325 | 
  326 | test("invalid invite shows an actionable error and offline history stays read-only", async ({
  327 |   page,
  328 |   audit,
  329 | }, info) => {
  330 |   audit.watch(page);
  331 |   audit.fault(true);
  332 |   await page.goto("http://localhost:8080/board/qa-not-authorized");
  333 |   await expect(
  334 |     page.getByRole("heading", { name: "We couldn’t open this canvas." }),
  335 |   ).toBeVisible();
  336 |   await expect(
  337 |     page.getByRole("link", { name: "Back to workspace" }),
  338 |   ).toBeVisible();
  339 |   await snap(page, info, "unauthorized-error");
  340 |   await page.getByRole("link", { name: "Back to workspace" }).click();
  341 |   audit.fault(false);
  342 |   await expect(
  343 |     page.getByRole("button", { name: "New board", exact: true }),
  344 |   ).toBeVisible();
  345 | });
  346 | 
```