# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: browser-qa.spec.ts >> time machine pause, all speeds, overlay and full replay preserve the changing live document
- Location: e2e\browser-qa.spec.ts:23:1

# Error details

```
Error: expect(locator).toHaveCount(expected) failed

Locator:  locator('.travel-banner')
Expected: 0
Received: 1
Timeout:  12000ms

Call log:
  - Expect "toHaveCount" locator('.travel-banner') with timeout 12000ms
  - waiting for locator('.travel-banner')
    27 × locator resolved to 1 element
       - unexpected value "1"

```

# Page snapshot

```yaml
- generic [ref=f1e3]:
  - banner [ref=f1e4]:
    - generic [ref=f1e5]:
      - link "SYNCFRAME home" [ref=f1e6] [cursor=pointer]:
        - /url: /
        - generic [ref=f1e9]: syncframe.
      - link "Back to workspace" [ref=f1e11] [cursor=pointer]:
        - /url: /
      - generic [ref=f1e14]:
        - strong [ref=f1e15]: QA history isolation
        - generic [ref=f1e16]: Shared workspace Canvas
    - generic [ref=f1e19]:
      - status [ref=f1e20]: Live
      - generic "QA Live editor, QA Historian" [ref=f1e22]:
        - generic [ref=f1e23]: Q
        - generic [ref=f1e24]: Q
      - button "History" [expanded] [ref=f1e25] [cursor=pointer]
      - button "Share board" [ref=f1e30] [cursor=pointer]
  - generic [ref=f1e37]:
    - complementary [ref=f1e38]:
      - button "Canvas" [active] [ref=f1e39] [cursor=pointer]
      - button "History" [ref=f1e42] [cursor=pointer]
      - generic [ref=f1e47]:
        - generic [ref=f1e48]: Q
        - generic [ref=f1e49]: You
    - main [ref=f1e50]:
      - generic [ref=f1e51]:
        - generic [ref=f1e52]:
          - text: History preview
          - generic [ref=f1e57]: /
          - strong [ref=f1e58]: 1 objects
        - generic [ref=f1e59]: REV 1
      - generic [ref=f1e60]:
        - generic [ref=f1e65]: You're looking at revision 1. Live edits continue in the background.
        - button "Return to live" [ref=f1e66] [cursor=pointer]
      - application "Collaborative whiteboard" [ref=f1e70]:
        - 'button "rect: 59c4da69" [ref=f1e72] [cursor=pointer]'
      - generic:
        - generic [ref=f1e74]:
          - button "Zoom out" [ref=f1e75] [cursor=pointer]
          - generic [ref=f1e79]: 100%
          - button "Zoom in" [ref=f1e80] [cursor=pointer]
          - button "Reset zoom" [ref=f1e85] [cursor=pointer]
        - toolbar "Drawing tools" [ref=f1e91]:
          - button "Select" [disabled] [pressed] [ref=f1e92]:
            - generic [ref=f1e95]: V
          - button "Rectangle" [disabled] [ref=f1e96]:
            - generic [ref=f1e99]: R
          - button "Ellipse" [disabled] [ref=f1e100]:
            - generic [ref=f1e103]: O
          - button "Sticky note" [disabled] [ref=f1e104]:
            - generic [ref=f1e108]: "N"
          - button "Text" [disabled] [ref=f1e109]:
            - generic [ref=f1e112]: T
        - generic [ref=f1e113]: Read-only history preview
    - complementary [ref=f1e114]:
      - generic [ref=f1e115]:
        - heading "Time machine" [level=2] [ref=f1e120]
        - button "Close history" [ref=f1e121] [cursor=pointer]
      - paragraph [ref=f1e125]: Every idea has a journey.Revisit any moment of yours.
      - generic [ref=f1e126]:
        - generic [ref=f1e127]:
          - generic [ref=f1e128]: PAST FRAME
          - generic [ref=f1e130]:
            - text: "1"
            - generic [ref=f1e131]: / 9
        - slider "History revision" [ref=f1e132] [cursor=pointer]: "1"
        - generic [ref=f1e133]:
          - generic [ref=f1e134]: Beginning
          - generic [ref=f1e135]: Now
        - generic [ref=f1e136]:
          - button "First event" [ref=f1e137] [cursor=pointer]
          - button "Previous event" [ref=f1e140] [cursor=pointer]
          - button "Play history" [ref=f1e143] [cursor=pointer]
          - button "Next event" [ref=f1e146] [cursor=pointer]
          - button "Latest event" [ref=f1e149] [cursor=pointer]
          - combobox "Playback speed" [ref=f1e152]:
            - option "0.5×" [selected]
            - option "1×"
            - option "2×"
            - option "4×"
      - generic [ref=f1e153]:
        - generic [ref=f1e154]:
          - checkbox "Overlay live positions" [checked] [ref=f1e155]
          - text: Overlay live positions
        - generic [ref=f1e156]: "Past: 1 · Live: 1 objects"
        - button "Back to live canvas" [ref=f1e157] [cursor=pointer]
      - generic [ref=f1e161]:
        - generic [ref=f1e162]: ACTIVITY
        - generic [ref=f1e163]: 9 events
      - generic [ref=f1e164]:
        - button "QA Live editor updated a shape 07:02:18 PM · Revision 9" [ref=f1e165] [cursor=pointer]:
          - generic [ref=f1e169]:
            - strong [ref=f1e170]:
              - text: QA Live editor
              - generic [ref=f1e171]: updated a shape
            - generic [ref=f1e172]:
              - text: 07:02:18 PM
              - generic [ref=f1e173]: ·
              - text: Revision 9
        - button "QA Historian updated a shape 07:02:15 PM · Revision 8" [ref=f1e176] [cursor=pointer]:
          - generic [ref=f1e180]:
            - strong [ref=f1e181]:
              - text: QA Historian
              - generic [ref=f1e182]: updated a shape
            - generic [ref=f1e183]:
              - text: 07:02:15 PM
              - generic [ref=f1e184]: ·
              - text: Revision 8
        - button "QA Historian updated a shape 07:02:15 PM · Revision 7" [ref=f1e187] [cursor=pointer]:
          - generic [ref=f1e191]:
            - strong [ref=f1e192]:
              - text: QA Historian
              - generic [ref=f1e193]: updated a shape
            - generic [ref=f1e194]:
              - text: 07:02:15 PM
              - generic [ref=f1e195]: ·
              - text: Revision 7
        - button "QA Historian updated a shape 07:02:15 PM · Revision 6" [ref=f1e198] [cursor=pointer]:
          - generic [ref=f1e202]:
            - strong [ref=f1e203]:
              - text: QA Historian
              - generic [ref=f1e204]: updated a shape
            - generic [ref=f1e205]:
              - text: 07:02:15 PM
              - generic [ref=f1e206]: ·
              - text: Revision 6
        - button "QA Historian updated a shape 07:02:15 PM · Revision 5" [ref=f1e209] [cursor=pointer]:
          - generic [ref=f1e213]:
            - strong [ref=f1e214]:
              - text: QA Historian
              - generic [ref=f1e215]: updated a shape
            - generic [ref=f1e216]:
              - text: 07:02:15 PM
              - generic [ref=f1e217]: ·
              - text: Revision 5
        - button "QA Historian updated a shape 07:02:15 PM · Revision 4" [ref=f1e220] [cursor=pointer]:
          - generic [ref=f1e224]:
            - strong [ref=f1e225]:
              - text: QA Historian
              - generic [ref=f1e226]: updated a shape
            - generic [ref=f1e227]:
              - text: 07:02:15 PM
              - generic [ref=f1e228]: ·
              - text: Revision 4
        - button "QA Historian updated a shape 07:02:15 PM · Revision 3" [ref=f1e231] [cursor=pointer]:
          - generic [ref=f1e235]:
            - strong [ref=f1e236]:
              - text: QA Historian
              - generic [ref=f1e237]: updated a shape
            - generic [ref=f1e238]:
              - text: 07:02:15 PM
              - generic [ref=f1e239]: ·
              - text: Revision 3
        - button "QA Historian updated a shape 07:02:15 PM · Revision 2" [ref=f1e242] [cursor=pointer]:
          - generic [ref=f1e246]:
            - strong [ref=f1e247]:
              - text: QA Historian
              - generic [ref=f1e248]: updated a shape
            - generic [ref=f1e249]:
              - text: 07:02:15 PM
              - generic [ref=f1e250]: ·
              - text: Revision 2
        - button "QA Historian added a shape 07:02:15 PM · Revision 1" [ref=f1e253] [cursor=pointer]:
          - generic [ref=f1e256]:
            - strong [ref=f1e257]:
              - text: QA Historian
              - generic [ref=f1e258]: added a shape
            - generic [ref=f1e259]:
              - text: 07:02:15 PM
              - generic [ref=f1e260]: ·
              - text: Revision 1
      - generic [ref=f1e263]: History is saved to MySQL.Exploring the past keeps the live canvas intact.
  - contentinfo [ref=f1e267]:
    - generic [ref=f1e268]:
      - text: Exploring history
      - generic [ref=f1e270]: /
      - text: Revision 9
    - generic [ref=f1e271]:
      - text: 2 online
      - generic [ref=f1e277]: /
      - text: You are QA Historian
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
> 29 |  await b.locator(`[data-object-id="${id}"]`).click({position:{x:90,y:50}});await b.getByLabel("Object text").fill("Live changed during playback");await b.getByLabel("Object text").press("Tab");await expect.poll(async()=>(await state(b)).revision).toBe(original.revision+1);await expect(a.getByTestId("revision")).toHaveText(paused);expect((await state(a)).objects[id].text).toBe("Live changed during playback");await snap(a,info,"history-overlay");await a.getByRole("button",{name:"Canvas",exact:true}).click();await expect(a.locator(".travel-banner")).toHaveCount(0);await expect(a.locator(`[data-object-id="${id}"]`)).toContainText("Live changed during playback");await assertPersisted(a);await one.close();await two.close();
     |                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          ^ Error: expect(locator).toHaveCount(expected) failed
  30 | });
  31 | 
  32 | for(const size of [{name:"desktop",width:1440,height:960},{name:"tablet",width:768,height:1024},{name:"mobile",width:390,height:844}]){
  33 |  test(`${size.name}: new-user editing, zoomed drag, properties and history remain usable`,async({browser,audit},info)=>{
  34 |   const context=await browser.newContext({viewport:size});await named(context,"QA "+size.name);const p=await context.newPage();audit.watch(p);await p.goto("/");await snap(p,info,size.name+"-home");expect(await p.evaluate(()=>document.documentElement.scrollWidth<=innerWidth)).toBe(true);await create(p,"QA "+size.name+" workspace");await snap(p,info,size.name+"-empty");await expect(p.getByRole("button",{name:"Sticky note",exact:true})).toBeInViewport();
  35 |   // A new user's first click must reach the canvas, including on a phone.
  36 |   const objectId=await add(p,"Sticky note",90,140);const object=p.locator(`[data-object-id="${objectId}"]`);await expect(object).toBeVisible();await p.getByLabel("Object text").fill("A readable idea on "+size.name);await p.getByLabel("Object text").press("Tab");await expect(object).toContainText("A readable idea");await p.getByRole("button",{name:"Color #e4dcfc",exact:true}).click();await expect(object.locator("rect").first()).toHaveAttribute("fill","#e4dcfc");await snap(p,info,size.name+"-properties");
  37 |   const box=await object.boundingBox();await p.mouse.move(box!.x+75,box!.y+45);await p.mouse.down();await p.mouse.move(box!.x+105,box!.y+75,{steps:6});await p.mouse.up();await expect(object).toHaveAttribute("transform","translate(120,170)");await p.getByLabel("Close properties").click();await p.getByRole("button",{name:"Zoom out",exact:true}).click();await p.getByRole("button",{name:"Zoom out",exact:true}).click();await expect(p.locator(".zoom")).toContainText("80%");const zoomBox=await object.boundingBox();await p.mouse.move(zoomBox!.x+50,zoomBox!.y+35);await p.mouse.down();await p.mouse.move(zoomBox!.x+90,zoomBox!.y+59,{steps:5});await p.mouse.up();await expect(object).toHaveAttribute("transform","translate(170,200)");await object.press("ArrowRight");await expect(object).toHaveAttribute("transform","translate(180,200)");await snap(p,info,size.name+"-edited");
  38 |   await p.getByLabel("Close properties").click();await p.getByRole("button",{name:"History",exact:true}).first().click();await expect(p.getByRole("heading",{name:"Time machine"})).toBeVisible();await snap(p,info,size.name+"-history");await p.getByRole("button",{name:"Close history",exact:true}).click();const scan=await new AxeBuilder({page:p}).withTags(["wcag2a","wcag2aa","wcag21aa"]).analyze();expect(scan.violations).toEqual([]);const durable=await assertPersisted(p);await p.reload();await live(p);expect(await state(p)).toEqual(durable);await p.locator(`[data-object-id="${objectId}"]`).press("Delete");await expect(p.getByTestId("board-object")).toHaveCount(0);await assertPersisted(p);await context.close();
  39 |  });
  40 | }
  41 | 
  42 | test("invalid invite shows an actionable error and offline history stays read-only",async({page,audit},info)=>{
  43 |  audit.watch(page);audit.fault(true);await page.goto("http://localhost:8080/board/qa-not-authorized");await expect(page.getByRole("heading",{name:"We couldn’t open this canvas."})).toBeVisible();await expect(page.getByRole("link",{name:"Back to workspace"})).toBeVisible();await snap(page,info,"unauthorized-error");await page.getByRole("link",{name:"Back to workspace"}).click();audit.fault(false);await expect(page.getByRole("button",{name:"New board",exact:true})).toBeVisible();
  44 | });
  45 | 
```