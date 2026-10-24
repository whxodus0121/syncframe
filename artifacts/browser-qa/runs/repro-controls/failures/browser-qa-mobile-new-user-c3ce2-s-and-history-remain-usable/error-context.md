# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: browser-qa.spec.ts >> mobile: new-user editing, zoomed drag, properties and history remain usable
- Location: e2e\browser-qa.spec.ts:33:2

# Error details

```
TimeoutError: locator.click: Timeout 15000ms exceeded.
Call log:
  - waiting for getByRole('button', { name: 'Sticky note', exact: true })
    - locator resolved to <button class="" aria-pressed="false" title="Sticky note (N)" aria-label="Sticky note">…</button>
  - attempting click action
    2 × waiting for element to be visible, enabled and stable
      - element is visible, enabled and stable
      - scrolling into view if needed
      - done scrolling
      - <span>…</span> from <aside class="history-panel">…</aside> subtree intercepts pointer events
    - retrying click action
    - waiting 20ms
    2 × waiting for element to be visible, enabled and stable
      - element is visible, enabled and stable
      - scrolling into view if needed
      - done scrolling
      - <span>…</span> from <aside class="history-panel">…</aside> subtree intercepts pointer events
    - retrying click action
      - waiting 100ms
    28 × waiting for element to be visible, enabled and stable
       - element is visible, enabled and stable
       - scrolling into view if needed
       - done scrolling
       - <span>…</span> from <aside class="history-panel">…</aside> subtree intercepts pointer events
     - retrying click action
       - waiting 500ms

```

# Page snapshot

```yaml
- generic [ref=f2e3]:
  - banner [ref=f2e4]:
    - generic [ref=f2e5]:
      - link "SYNCFRAME home" [ref=f2e6] [cursor=pointer]:
        - /url: /
      - link "Back to workspace" [ref=f2e9] [cursor=pointer]:
        - /url: /
      - generic [ref=f2e12]:
        - strong [ref=f2e13]: QA mobile workspace
        - generic [ref=f2e14]: Shared workspace Canvas
    - generic [ref=f2e17]:
      - status [ref=f2e18]: Live
      - button "History" [expanded] [ref=f2e20] [cursor=pointer]
      - button "Share board" [ref=f2e25] [cursor=pointer]
  - generic [ref=f2e32]:
    - main [ref=f2e33]:
      - generic [ref=f2e34]:
        - generic [ref=f2e35]:
          - text: Live canvas
          - generic [ref=f2e40]: /
        - generic [ref=f2e41]: REV 0
      - application "Collaborative whiteboard" [ref=f2e43]:
        - generic:
          - generic: Every idea starts somewhere.
          - generic: Choose a tool below, then click to begin.
          - generic: Your canvas. Your pace.
      - generic:
        - generic [ref=f2e45]:
          - button "Zoom out" [ref=f2e46] [cursor=pointer]
          - generic [ref=f2e50]: 100%
          - button "Zoom in" [ref=f2e51] [cursor=pointer]
          - button "Reset zoom" [ref=f2e56] [cursor=pointer]
        - toolbar "Drawing tools" [ref=f2e62]:
          - button "Select" [pressed] [ref=f2e63] [cursor=pointer]:
            - generic [ref=f2e66]: V
          - button "Rectangle" [ref=f2e67] [cursor=pointer]:
            - generic [ref=f2e70]: R
          - button "Ellipse" [ref=f2e71] [cursor=pointer]:
            - generic [ref=f2e74]: O
          - button "Sticky note" [ref=f2e75] [cursor=pointer]:
            - generic [ref=f2e79]: "N"
          - button "Text" [ref=f2e80] [cursor=pointer]:
            - generic [ref=f2e83]: T
    - complementary [ref=f2e84]:
      - generic [ref=f2e85]:
        - heading "Time machine" [level=2] [ref=f2e90]
        - button "Close history" [ref=f2e91] [cursor=pointer]
      - paragraph [ref=f2e95]: Every idea has a journey.Revisit any moment of yours.
      - generic [ref=f2e96]:
        - generic [ref=f2e97]:
          - generic [ref=f2e98]: LIVE FRAME
          - generic [ref=f2e100]:
            - text: "0"
            - generic [ref=f2e101]: / 0
        - slider "History revision" [ref=f2e102] [cursor=pointer]: "0"
        - generic [ref=f2e103]:
          - generic [ref=f2e104]: Beginning
          - generic [ref=f2e105]: Now
        - generic [ref=f2e106]:
          - button "First event" [ref=f2e107] [cursor=pointer]
          - button "Previous event" [ref=f2e110] [cursor=pointer]
          - button "Play history" [disabled] [ref=f2e113]
          - button "Next event" [ref=f2e116] [cursor=pointer]
          - button "Latest event" [ref=f2e119] [cursor=pointer]
          - combobox "Playback speed" [ref=f2e122]:
            - option "0.5×"
            - option "1×" [selected]
            - option "2×"
            - option "4×"
      - generic [ref=f2e123]:
        - generic [ref=f2e124]: ACTIVITY
        - generic [ref=f2e125]: 0 events
      - generic [ref=f2e127]:
        - strong [ref=f2e132]: The story starts here.
        - paragraph [ref=f2e133]: Your edits will appear here,one frame at a time.
      - generic [ref=f2e134]: History is saved to MySQL.Exploring the past keeps the live canvas intact.
  - contentinfo [ref=f2e138]:
    - generic [ref=f2e139]:
      - text: All committed edits saved
      - generic [ref=f2e141]: /
      - text: Revision 0
```

# Test source

```ts
  1  | import { test as base, expect, type Page, type BrowserContext, type TestInfo } from "@playwright/test";
  2  | import fs from "node:fs";
  3  | import path from "node:path";
  4  | import { execFileSync } from "node:child_process";
  5  | import type { State, EditEvent } from "../../src/model";
  6  | 
  7  | export type Diagnostic = { kind:string; message:string; expected:boolean };
  8  | const redact=(value:string)=>value.replace(/[a-f0-9]{64}/gi,"[REDACTED_TEST_CAPABILITY]");
  9  | export const test=base.extend<{audit:{watch:(page:Page)=>void;fault:(active:boolean)=>void}}>({
  10 |  audit: [async({},use,info)=>{
  11 |   const messages:Diagnostic[]=[];let fault=false;
  12 |   await use({fault:(active)=>{fault=active},watch:(page)=>{
  13 |    page.on("console",m=>{if(m.type()==="error")messages.push({kind:"console",message:redact(m.text()),expected:fault})});
  14 |    page.on("pageerror",e=>messages.push({kind:"pageerror",message:e.message,expected:false}));
  15 |    page.on("requestfailed",r=>{if(/\/api\//.test(r.url()))messages.push({kind:"network",message:redact(r.url()+" "+r.failure()?.errorText),expected:fault||r.failure()?.errorText==="net::ERR_ABORTED"})});
  16 |    page.on("response",r=>{if(r.status()>=400)messages.push({kind:"http",message:redact(`${r.status()} ${r.url()}`),expected:fault})});
  17 |   }});
  18 |   fs.mkdirSync(info.outputDir,{recursive:true});
  19 |   fs.writeFileSync(info.outputPath("diagnostics.json"),JSON.stringify(messages,null,2));
  20 |   await info.attach("Console and network diagnostics",{body:JSON.stringify(messages,null,2),contentType:"application/json"});
  21 |   expect(messages.filter(m=>!m.expected),"No unexpected browser errors").toEqual([]);
  22 |  },{auto:true}]
  23 | });
  24 | export { expect };
  25 | export const live=async(p:Page)=>expect(p.locator(".connection")).toContainText("Live");
  26 | export async function named(context:BrowserContext,name:string){await context.addInitScript(name=>{if(!sessionStorage.getItem("sf:name"))sessionStorage.setItem("sf:name",name)},name)}
  27 | export async function create(p:Page,title:string){await p.goto("/");await p.getByLabel("Board title").fill(title);await p.getByRole("button",{name:"New board",exact:true}).click();await live(p);return p.url().split("/").at(-1)!}
  28 | export async function state(p:Page):Promise<State>{return p.evaluate(async()=>{const id=location.pathname.split('/').at(-1);const response=await fetch(`/api/boards/${id}`);if(!response.ok)throw Error("State API failed: "+response.status);return response.json()})}
  29 | export async function events(p:Page):Promise<EditEvent[]>{return p.evaluate(async()=>{const id=location.pathname.split('/').at(-1);let all:any[]=[];for(;;){const r=await fetch(`/api/boards/${id}/events?after=${all.at(-1)?.revision||0}`);if(!r.ok)throw Error("Events API failed");const page=await r.json();all=all.concat(page);if(page.length<500)return all}})}
  30 | export async function settled(p:Page){await expect(p.locator(".connection")).not.toContainText("pending");await expect.poll(async()=>Number((await p.getByTestId("revision").innerText()).replace("REV ",""))).toBe((await state(p)).revision)}
  31 | export async function invite(p:Page,context:BrowserContext){await context.grantPermissions(["clipboard-read","clipboard-write"],{origin:new URL(p.url()).origin});await p.getByRole("button",{name:"Share board",exact:true}).click();await expect(p.getByRole("status").filter({hasText:"Invite link copied"})).toBeVisible();const link=await p.evaluate(()=>navigator.clipboard.readText());expect(link).toMatch(/\/board\/[^#]+#key=[a-f0-9]{64}$/);return link}
> 32 | export async function add(p:Page,kind:string,x:number,y:number){await p.getByRole("button",{name:kind,exact:true}).click();const before=(await state(p)).revision;await p.getByTestId("canvas").click({position:{x,y}});await expect.poll(async()=>(await state(p)).revision).toBe(before+1);await settled(p);const st=await state(p);return (await events(p)).at(-1)!.object.id}
     |                                                                                                                    ^ TimeoutError: locator.click: Timeout 15000ms exceeded.
  33 | export async function snap(p:Page,info:TestInfo,label:string){const folder=path.resolve("../artifacts/browser-qa/screenshots",process.env.QA_RUN_ID||"latest");fs.mkdirSync(folder,{recursive:true});const file=path.join(folder,label+".png");await p.screenshot({path:file,fullPage:true,mask:[p.getByLabel("Invite link")]});await info.attach(label,{path:file,contentType:"image/png"})}
  34 | export async function assertPersisted(p:Page){const st=await state(p);const es=await events(p);expect(es.map(e=>e.revision)).toEqual(Array.from({length:st.revision},(_,i)=>i+1));expect(new Set(es.map(e=>e.operation.id)).size).toBe(es.length);const replay:State={...st,revision:0,objects:{}};for(const e of es){replay.objects[e.object.id]=e.object;replay.revision=e.revision}expect(replay).toEqual(st);
  35 |  // Compare the actual SQL state, not just the HTTP representation.
  36 |  const result=execFileSync("docker",["compose","-p","syncframe","exec","-T","mysql","mysql","-N","-B","--raw","-usyncframe","-psyncframe","-D","syncframe","-e",`SELECT state FROM boards WHERE id='${st.id}'`],{cwd:path.resolve(".."),encoding:"utf8",stdio:["ignore","pipe","pipe"]});expect(JSON.parse(result.trim())).toEqual(st);return st}
  37 | 
```