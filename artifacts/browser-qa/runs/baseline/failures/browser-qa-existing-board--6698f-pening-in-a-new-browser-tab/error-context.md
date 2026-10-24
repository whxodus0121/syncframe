# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: browser-qa.spec.ts >> existing-board share survives reopening in a new browser tab
- Location: e2e\browser-qa.spec.ts:6:1

# Error details

```
Error: expect(locator).toBeVisible() failed

Locator: getByRole('status').filter({ hasText: 'Invite link copied' })
Expected: visible
Timeout: 12000ms
Error: element(s) not found

Call log:
  - Expect "toBeVisible" getByRole('status').filter({ hasText: 'Invite link copied' }) with timeout 12000ms
  - waiting for getByRole('status').filter({ hasText: 'Invite link copied' })

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
  - strong: QA existing-board sharing
  - text: Shared workspace
  - img
  - text: Canvas
  - status: Live
  - text: Q Q
  - button "History" [expanded]:
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
  - text: REV 1
  - application "Collaborative whiteboard":
    - 'button "text: Your text"': Your text
  - button "Zoom out":
    - img
  - text: 100%
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
- complementary:
  - img
  - heading "Time machine" [level=2]
  - button "Close history":
    - img
  - paragraph: Every idea has a journey. Revisit any moment of yours.
  - text: LIVE FRAME 1 / 1
  - slider "History revision": "1"
  - text: Beginning Now
  - button "First event":
    - img
  - button "Previous event":
    - img
  - button "Play history":
    - img
  - button "Next event":
    - img
  - button "Latest event":
    - img
  - combobox "Playback speed":
    - option "0.5×"
    - option "1×" [selected]
    - option "2×"
    - option "4×"
  - text: ACTIVITY 1 events
  - button "QA Alice added text 06:59:09 PM · Revision 1":
    - img
    - strong: QA Alice added text
    - text: 06:59:09 PM · Revision 1
    - img
  - img
  - text: History is saved to MySQL. Exploring the past keeps the live canvas intact.
- contentinfo:
  - text: All committed edits saved / Revision 1
  - img
  - text: 2 online / You are QA Alice
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
> 31 | export async function invite(p:Page,context:BrowserContext){await context.grantPermissions(["clipboard-read","clipboard-write"],{origin:new URL(p.url()).origin});await p.getByRole("button",{name:"Share board",exact:true}).click();await expect(p.getByRole("status").filter({hasText:"Invite link copied"})).toBeVisible();const link=await p.evaluate(()=>navigator.clipboard.readText());expect(link).toMatch(/\/board\/[^#]+#key=[a-f0-9]{64}$/);return link}
     |                                                                                                                                                                                                                                                                                                                  ^ Error: expect(locator).toBeVisible() failed
  32 | export async function add(p:Page,kind:string,x:number,y:number){await p.getByRole("button",{name:kind,exact:true}).click();const before=(await state(p)).revision;await p.getByTestId("canvas").click({position:{x,y}});await expect.poll(async()=>(await state(p)).revision).toBe(before+1);await settled(p);const st=await state(p);return (await events(p)).at(-1)!.object.id}
  33 | export async function snap(p:Page,info:TestInfo,label:string){const folder=path.resolve("../artifacts/browser-qa/screenshots",process.env.QA_RUN_ID||"latest");fs.mkdirSync(folder,{recursive:true});const file=path.join(folder,label+".png");await p.screenshot({path:file,fullPage:true,mask:[p.getByLabel("Invite link")]});await info.attach(label,{path:file,contentType:"image/png"})}
  34 | export async function assertPersisted(p:Page){const st=await state(p);const es=await events(p);expect(es.map(e=>e.revision)).toEqual(Array.from({length:st.revision},(_,i)=>i+1));expect(new Set(es.map(e=>e.operation.id)).size).toBe(es.length);const replay:State={...st,revision:0,objects:{}};for(const e of es){replay.objects[e.object.id]=e.object;replay.revision=e.revision}expect(replay).toEqual(st);
  35 |  // Compare the actual SQL state, not just the HTTP representation.
  36 |  const result=execFileSync("docker",["compose","-p","syncframe","exec","-T","mysql","mysql","-N","-B","--raw","-usyncframe","-psyncframe","-D","syncframe","-e",`SELECT state FROM boards WHERE id='${st.id}'`],{cwd:path.resolve(".."),encoding:"utf8",stdio:["ignore","pipe","pipe"]});expect(JSON.parse(result.trim())).toEqual(st);return st}
  37 | 
```