# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: browser-qa.spec.ts >> unsubmitted text and dimension drafts survive another user's committed patch
- Location: e2e\browser-qa.spec.ts:349:1

# Error details

```
Error: expect(locator).toHaveValue(expected) failed

Locator:  getByLabel('Object text')
Expected: "Local draft must survive"
Received: "Remote commit while Alice is typing"
Timeout:  12000ms

Call log:
  - Expect "toHaveValue" getByLabel('Object text') with timeout 12000ms
  - waiting for getByLabel('Object text')
    27 × locator resolved to <textarea maxlength="2000" aria-label="Object text">Remote commit while Alice is typing</textarea>
       - unexpected value "Remote commit while Alice is typing"

```

```yaml
- textbox "Object text": Remote commit while Alice is typing
```

# Test source

```ts
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
  302 |     if (await p.getByRole("heading", { name: "Time machine" }).isVisible())
  303 |       await p.getByRole("button", { name:"Close history", exact:true }).click();
  304 |     await p
  305 |       .getByRole("button", { name: "History", exact: true })
  306 |       .first()
  307 |       .click();
  308 |     await expect(
  309 |       p.getByRole("heading", { name: "Time machine" }),
  310 |     ).toBeVisible();
  311 |     await snap(p, info, size.name + "-history");
  312 |     await p.getByRole("button", { name: "Close history", exact: true }).click();
  313 |     const scan = await new AxeBuilder({ page: p })
  314 |       .withTags(["wcag2a", "wcag2aa", "wcag21aa"])
  315 |       .analyze();
  316 |     expect(scan.violations).toEqual([]);
  317 |     const durable = await assertPersisted(p);
  318 |     await p.reload();
  319 |     await live(p);
  320 |     expect(await state(p)).toEqual(durable);
  321 |     await p.locator(`[data-object-id="${objectId}"]`).press("Delete");
  322 |     await expect(p.getByTestId("board-object")).toHaveCount(0);
  323 |     await assertPersisted(p);
  324 |     await context.close();
  325 |   });
  326 | }
  327 | 
  328 | test("invalid invite shows an actionable error and offline history stays read-only", async ({
  329 |   page,
  330 |   audit,
  331 | }, info) => {
  332 |   audit.watch(page);
  333 |   audit.fault(true);
  334 |   await page.goto("http://localhost:8080/board/qa-not-authorized");
  335 |   await expect(
  336 |     page.getByRole("heading", { name: "We couldn’t open this canvas." }),
  337 |   ).toBeVisible();
  338 |   await expect(
  339 |     page.getByRole("link", { name: "Back to workspace" }),
  340 |   ).toBeVisible();
  341 |   await snap(page, info, "unauthorized-error");
  342 |   await page.getByRole("link", { name: "Back to workspace" }).click();
  343 |   audit.fault(false);
  344 |   await expect(
  345 |     page.getByRole("button", { name: "New board", exact: true }),
  346 |   ).toBeVisible();
  347 | });
  348 | 
  349 | test("unsubmitted text and dimension drafts survive another user's committed patch",async({browser,audit},info)=>{
  350 |  const ac=await browser.newContext(),bc=await browser.newContext();await named(ac,"QA Draft owner");await named(bc,"QA Remote writer");const a=await ac.newPage(),b=await bc.newPage();[a,b].forEach(audit.watch);await create(a,"QA draft preservation");await b.goto(await invite(a,ac));await live(b);const id=await add(a,"Sticky note",350,150);const object=`[data-object-id="${id}"]`;await expect(b.locator(object)).toBeVisible();await b.locator(object).click({position:{x:90,y:50}});
> 351 |  await a.getByLabel("Object text").fill("Local draft must survive");await b.getByLabel("Object text").fill("Remote commit while Alice is typing");await b.getByLabel("Object text").press("Tab");await expect(a.locator(object)).toContainText("Remote commit while Alice is typing");await snap(a,info,"local-draft-during-remote-change");await expect(a.getByLabel("Object text")).toHaveValue("Local draft must survive");await a.getByLabel("Object text").press("Tab");await expect(b.locator(object)).toContainText("Local draft must survive");
      |                                                                                                                                                                                                                                                                                                                                                                                       ^ Error: expect(locator).toHaveValue(expected) failed
  352 |  await a.getByLabel("Object width",{exact:true}).fill("280");await b.getByLabel("Object width",{exact:true}).fill("240");await b.getByLabel("Object width",{exact:true}).press("Tab");await expect(a.locator(object).locator("rect").first()).toHaveAttribute("width","240");await expect(a.getByLabel("Object width",{exact:true})).toHaveValue("280");await a.getByLabel("Object width",{exact:true}).press("Tab");await expect(b.locator(object).locator("rect").first()).toHaveAttribute("width","280");await assertPersisted(a);await ac.close();await bc.close();
  353 | });
  354 | 
  355 | test("a delayed real history response shows loading rather than the live document",async({page,audit},info)=>{
  356 |  audit.watch(page);await create(page,"QA delayed history");await add(page,"Rectangle",350,170);let release:()=>void=()=>{};let ready:()=>void=()=>{};const held=new Promise<void>(resolve=>{ready=resolve});await page.route("**/api/boards/*?revision=0",async route=>{await new Promise<void>(resolve=>{release=resolve;ready()});await route.continue()});await page.getByRole("button",{name:"First event",exact:true}).click();await held;
  357 |  try{await snap(page,info,"history-response-in-flight");await expect(page.getByRole("status").filter({hasText:"Loading revision 0"})).toBeVisible();await expect(page.getByTestId("board-object")).toHaveCount(0)}finally{release()}
  358 |  await expect(page.getByTestId("revision")).toHaveText("REV 0");await expect(page.getByTestId("board-object")).toHaveCount(0);expect((await state(page)).revision).toBe(1);await assertPersisted(page);
  359 | });
  360 | 
```