# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: browser-qa.spec.ts >> desktop: new-user editing, zoomed drag, properties and history remain usable
- Location: e2e\browser-qa.spec.ts:33:2

# Error details

```
Error: expect(locator).toHaveAttribute(expected) failed

Locator:  locator('[data-object-id="05c783ac-26a5-4454-ade7-5345377d72dd"]')
Expected: "translate(120,170)"
Received: "translate(90,140)"
Timeout:  12000ms

Call log:
  - Expect "toHaveAttribute" locator('[data-object-id="05c783ac-26a5-4454-ade7-5345377d72dd"]') with timeout 12000ms
  - waiting for locator('[data-object-id="05c783ac-26a5-4454-ade7-5345377d72dd"]')
    27 × locator resolved to <g tabindex="0" role="button" class="canvas-object" data-testid="board-object" transform="translate(90,140)" aria-label="note: A readable idea on desktop" data-object-id="05c783ac-26a5-4454-ade7-5345377d72dd">…</g>
       - unexpected value "translate(90,140)"

```

```yaml
- 'button "note: A readable idea on desktop"': A readable idea on desktop
```

# Test source

```ts
  1  | import { test,expect,named,create,live,state,events,settled,invite,add,snap,assertPersisted } from "./qa/helpers";
  2  | import {execFileSync} from "node:child_process";
  3  | import path from "node:path";
  4  | import AxeBuilder from "@axe-core/playwright";
  5  | 
  6  | test("existing-board share survives reopening in a new browser tab",async({browser,audit},info)=>{
  7  |  const context=await browser.newContext();await named(context,"QA Alice");const p=await context.newPage();audit.watch(p);const id=await create(p,"QA existing-board sharing");await add(p,"Text",340,180);await p.getByRole("link",{name:"Back to workspace",exact:true}).click();await p.getByRole("link").filter({hasText:"QA existing-board sharing"}).click();await live(p);await expect(p.getByTestId("board-object")).toHaveCount(1);
  8  |  const reopened=await context.newPage();audit.watch(reopened);await reopened.goto(`http://localhost:8080/board/${id}`);await live(reopened);await snap(reopened,info,"existing-before-share");const link=await invite(reopened,context);
  9  |  const guest=await browser.newContext();await named(guest,"QA Guest");const joined=await guest.newPage();audit.watch(joined);await joined.goto(link);await live(joined);await expect(joined.getByTestId("board-object")).toHaveCount(1);await assertPersisted(joined);await guest.close();await context.close();
  10 | });
  11 | 
  12 | test("three named sessions converge across two servers, concurrent edits, reconnect and restart",async({browser,audit},info)=>{
  13 |  test.setTimeout(150000);const contexts=await Promise.all([browser.newContext(),browser.newContext(),browser.newContext()]);for(let i=0;i<3;i++)await named(contexts[i],["QA Alice","QA Bob","QA Cara"][i]);const [a,b,c]=await Promise.all(contexts.map(ctx=>ctx.newPage()));[a,b,c].forEach(audit.watch);const id=await create(a,"QA distributed collaboration");const link=await invite(a,contexts[0]);await b.goto(link.replace(":8080",":8081"));await c.goto(link);await Promise.all([live(b),live(c)]);await expect(a.locator(".board-footer")).toContainText("3 online",{timeout:20000});
  14 |  const objectId=await add(a,"Sticky note",350,170);const target=`[data-object-id="${objectId}"]`;await Promise.all([expect(b.locator(target)).toBeVisible(),expect(c.locator(target)).toBeVisible()]);await a.getByLabel("Close properties").click();await b.locator(target).click({position:{x:100,y:60}});await b.getByLabel("Object x",{exact:true}).fill("440");await b.getByLabel("Object x",{exact:true}).press("Tab");await expect(a.locator(target)).toHaveAttribute("transform","translate(440,170)");await expect(c.locator(target)).toHaveAttribute("transform","translate(440,170)");
  15 |  const bounds=await a.getByTestId("canvas").boundingBox();await a.mouse.move(bounds!.x+600,bounds!.y+300);await expect(b.getByTestId("canvas").getByText("QA Alice",{exact:true})).toBeVisible();
  16 |  await a.locator(target).click({position:{x:90,y:50}});await c.locator(target).click({position:{x:90,y:50}});const before=(await state(a)).revision;
  17 |  await Promise.all([a.getByLabel("Object text").fill("Alice's simultaneous idea"),c.getByLabel("Object text").fill("Cara's simultaneous idea")]);await Promise.all([a.getByLabel("Object text").press("Tab"),c.getByLabel("Object text").press("Tab")]);await expect.poll(async()=>(await state(a)).revision).toBe(before+2);const last=(await events(a)).at(-1)!;expect(["QA Alice","QA Cara"]).toContain(last.operation.actor);for(const page of[a,b,c])await expect(page.locator(target)).toContainText(last.object.text);await expect.poll(async()=>JSON.stringify((await state(b)).objects)).toBe(JSON.stringify((await state(a)).objects));
  18 |  audit.fault(true);await contexts[1].setOffline(true);await expect(b.locator(".connection")).toContainText("Reconnecting");await b.getByRole("button",{name:"Rectangle",exact:true}).click();await b.getByTestId("canvas").click({position:{x:760,y:380}});await expect(b.locator(".connection")).toContainText("1 pending");const count=Object.values((await state(a)).objects).filter(o=>!o.deleted).length;await a.getByLabel("Close properties").click();await add(a,"Ellipse",700,130);await snap(b,info,"offline-pending");await contexts[1].setOffline(false);await live(b);await expect(b.locator(".connection")).not.toContainText("pending");for(const page of[a,b,c])await expect(page.getByTestId("board-object")).toHaveCount(count+2);audit.fault(false);const persisted=await assertPersisted(a);expect(new Set((await events(a)).map(e=>e.operation.id)).size).toBe(persisted.revision);
  19 |  await b.close();await expect(a.locator(".board-footer")).toContainText("2 online",{timeout:20000});const again=await contexts[1].newPage();audit.watch(again);await again.goto(`http://localhost:8081/board/${id}`);await live(again);await expect(again.getByTestId("board-object")).toHaveCount(3);await expect(a.locator(".board-footer")).toContainText("3 online",{timeout:20000});
  20 |  audit.fault(true);execFileSync("docker",["compose","-p","syncframe","restart","api"],{cwd:path.resolve(".."),stdio:"pipe"});await live(a);await expect.poll(async()=>state(a)).toEqual(persisted);await expect(a.getByTestId("board-object")).toHaveCount(3);audit.fault(false);await snap(a,info,"three-session-recovery");await Promise.all(contexts.map(ctx=>ctx.close()));
  21 | });
  22 | 
  23 | test("time machine pause, all speeds, overlay and full replay preserve the changing live document",async({browser,audit},info)=>{
  24 |  const one=await browser.newContext(),two=await browser.newContext();await named(one,"QA Historian");await named(two,"QA Live editor");const a=await one.newPage(),b=await two.newPage();[a,b].forEach(audit.watch);await create(a,"QA history isolation");const link=await invite(a,one);await b.goto(link);await live(b);const id=await add(a,"Rectangle",350,160);await a.getByLabel("Object text").fill("Original frame");await a.getByLabel("Object text").press("Tab");await settled(a);for(let i=0;i<6;i++){await a.getByLabel("Object x",{exact:true}).fill(String(360+i*20));await a.getByLabel("Object x",{exact:true}).press("Tab");await expect(a.locator(`[data-object-id="${id}"]`)).toHaveAttribute("transform",`translate(${360+i*20},160)`)}const original=await assertPersisted(a);
  25 |  await a.getByRole("button",{name:"First event",exact:true}).click();await expect(a.getByTestId("revision")).toHaveText("REV 0");await expect(a.getByTestId("board-object")).toHaveCount(0);await a.getByRole("button",{name:"Next event",exact:true}).click();await expect(a.getByTestId("revision")).toHaveText("REV 1");await expect(a.getByTestId("board-object")).toHaveCount(1);await a.getByRole("button",{name:"Previous event",exact:true}).click();await expect(a.getByTestId("revision")).toHaveText("REV 0");
  26 |  for(const speed of[".5","1","2","4"]){await a.getByLabel("Playback speed").selectOption(speed);await expect(a.getByLabel("Playback speed")).toHaveValue(speed)}await a.getByLabel("Playback speed").selectOption(".5");await a.getByRole("button",{name:"Play history",exact:true}).click();await expect.poll(async()=>Number((await a.getByTestId("revision").innerText()).replace("REV ",""))).toBeGreaterThan(0);await a.getByRole("button",{name:"Pause history",exact:true}).click();const paused=await a.getByTestId("revision").innerText();
  27 |  // Advance beyond more than two playback intervals without a sleep-based oracle.
  28 |  await a.clock.install();await a.clock.fastForward(5000);await expect(a.getByTestId("revision")).toHaveText(paused);await a.clock.resume();await a.getByRole("checkbox",{name:"Overlay live positions"}).check();await expect(a.locator('svg.canvas > rect[stroke="#16a57a"]')).toHaveCount(1);
  29 |  await b.locator(`[data-object-id="${id}"]`).click({position:{x:90,y:50}});await b.getByLabel("Object text").fill("Live changed during playback");await b.getByLabel("Object text").press("Tab");await expect.poll(async()=>(await state(b)).revision).toBe(original.revision+1);await expect(a.getByTestId("revision")).toHaveText(paused);expect((await state(a)).objects[id].text).toBe("Live changed during playback");await snap(a,info,"history-overlay");await a.getByRole("button",{name:"Canvas",exact:true}).click();await expect(a.locator(".travel-banner")).toHaveCount(0);await expect(a.locator(`[data-object-id="${id}"]`)).toContainText("Live changed during playback");await assertPersisted(a);await one.close();await two.close();
  30 | });
  31 | 
  32 | for(const size of [{name:"desktop",width:1440,height:960},{name:"tablet",width:768,height:1024},{name:"mobile",width:390,height:844}]){
  33 |  test(`${size.name}: new-user editing, zoomed drag, properties and history remain usable`,async({browser,audit},info)=>{
  34 |   const context=await browser.newContext({viewport:size});await named(context,"QA "+size.name);const p=await context.newPage();audit.watch(p);await p.goto("/");await snap(p,info,size.name+"-home");expect(await p.evaluate(()=>document.documentElement.scrollWidth<=innerWidth)).toBe(true);await create(p,"QA "+size.name+" workspace");await snap(p,info,size.name+"-empty");await expect(p.getByRole("button",{name:"Sticky note",exact:true})).toBeInViewport();
  35 |   // A new user's first click must reach the canvas, including on a phone.
  36 |   const objectId=await add(p,"Sticky note",90,140);const object=p.locator(`[data-object-id="${objectId}"]`);await expect(object).toBeVisible();await p.getByLabel("Object text").fill("A readable idea on "+size.name);await p.getByLabel("Object text").press("Tab");await expect(object).toContainText("A readable idea");await p.getByRole("button",{name:"Color #e4dcfc",exact:true}).click();await expect(object.locator("rect").first()).toHaveAttribute("fill","#e4dcfc");await snap(p,info,size.name+"-properties");
> 37 |   const box=await object.boundingBox();await p.mouse.move(box!.x+75,box!.y+45);await p.mouse.down();await p.mouse.move(box!.x+105,box!.y+75,{steps:6});await p.mouse.up();await expect(object).toHaveAttribute("transform","translate(120,170)");await p.getByLabel("Close properties").click();await p.getByRole("button",{name:"Zoom out",exact:true}).click();await p.getByRole("button",{name:"Zoom out",exact:true}).click();await expect(p.locator(".zoom")).toContainText("80%");const zoomBox=await object.boundingBox();await p.mouse.move(zoomBox!.x+50,zoomBox!.y+35);await p.mouse.down();await p.mouse.move(zoomBox!.x+90,zoomBox!.y+59,{steps:5});await p.mouse.up();await expect(object).toHaveAttribute("transform","translate(170,200)");await object.press("ArrowRight");await expect(object).toHaveAttribute("transform","translate(180,200)");await snap(p,info,size.name+"-edited");
     |                                                                                                                                                                                                ^ Error: expect(locator).toHaveAttribute(expected) failed
  38 |   await p.getByLabel("Close properties").click();await p.getByRole("button",{name:"History",exact:true}).first().click();await expect(p.getByRole("heading",{name:"Time machine"})).toBeVisible();await snap(p,info,size.name+"-history");await p.getByRole("button",{name:"Close history",exact:true}).click();const scan=await new AxeBuilder({page:p}).withTags(["wcag2a","wcag2aa","wcag21aa"]).analyze();expect(scan.violations).toEqual([]);const durable=await assertPersisted(p);await p.reload();await live(p);expect(await state(p)).toEqual(durable);await p.locator(`[data-object-id="${objectId}"]`).press("Delete");await expect(p.getByTestId("board-object")).toHaveCount(0);await assertPersisted(p);await context.close();
  39 |  });
  40 | }
  41 | 
  42 | test("invalid invite shows an actionable error and offline history stays read-only",async({page,audit},info)=>{
  43 |  audit.watch(page);audit.fault(true);await page.goto("http://localhost:8080/board/qa-not-authorized");await expect(page.getByRole("heading",{name:"We couldn’t open this canvas."})).toBeVisible();await expect(page.getByRole("link",{name:"Back to workspace"})).toBeVisible();await snap(page,info,"unauthorized-error");await page.getByRole("link",{name:"Back to workspace"}).click();audit.fault(false);await expect(page.getByRole("button",{name:"New board",exact:true})).toBeVisible();
  44 | });
  45 | 
```