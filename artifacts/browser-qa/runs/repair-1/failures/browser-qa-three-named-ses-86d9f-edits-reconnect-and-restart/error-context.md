# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: browser-qa.spec.ts >> three named sessions converge across two servers, concurrent edits, reconnect and restart
- Location: e2e\browser-qa.spec.ts:54:1

# Error details

```
Error: expect(received).toBe(expected) // Object.is equality

Expected: 4
Received: 3

Call Log:
- Timeout 12000ms exceeded while waiting on the predicate
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
        - strong [ref=f1e15]: QA distributed collaboration
        - generic [ref=f1e16]: Shared workspace Canvas
    - generic [ref=f1e19]:
      - status [ref=f1e20]: Live
      - generic "QA Cara, QA Bob, QA Alice" [ref=f1e22]:
        - generic [ref=f1e23]: Q
        - generic [ref=f1e24]: Q
        - generic [ref=f1e25]: Q
      - button "History" [expanded] [ref=f1e26] [cursor=pointer]
      - button "Share board" [ref=f1e31] [cursor=pointer]
  - generic [ref=f1e38]:
    - complementary [ref=f1e39]:
      - button "Canvas" [ref=f1e40] [cursor=pointer]
      - button "History" [ref=f1e43] [cursor=pointer]
      - generic [ref=f1e48]:
        - generic [ref=f1e49]: Q
        - generic [ref=f1e50]: You
    - main [ref=f1e51]:
      - generic [ref=f1e52]:
        - generic [ref=f1e53]:
          - text: Live canvas
          - generic [ref=f1e58]: /
          - strong [ref=f1e59]: 1 objects
        - generic [ref=f1e60]: REV 3
      - application "Collaborative whiteboard" [ref=f1e62]:
        - 'button "note: Cara''s simultaneous idea" [ref=f1e64]':
          - generic: Cara's simultaneous idea
        - generic: QA Bob
        - generic: QA Cara
      - generic [ref=f1e67]:
        - generic [ref=f1e68]:
          - strong [ref=f1e69]: Object properties
          - button "Close properties" [ref=f1e70] [cursor=pointer]
        - generic [ref=f1e74]:
          - text: Text
          - textbox "Object text" [ref=f1e75]: Cara's simultaneous idea
        - group "Object color" [ref=f1e76]:
          - 'button "Color #fff2b3" [active] [ref=f1e77] [cursor=pointer]'
          - 'button "Color #dce7ff" [ref=f1e78] [cursor=pointer]'
          - 'button "Color #e4dcfc" [ref=f1e79] [cursor=pointer]'
          - 'button "Color #d7f0e5" [ref=f1e80] [cursor=pointer]'
          - 'button "Color #ffdfe0" [ref=f1e81] [cursor=pointer]'
          - 'button "Color #ffffff" [ref=f1e82] [cursor=pointer]'
        - generic [ref=f1e83]:
          - generic [ref=f1e84]:
            - text: x
            - spinbutton "Object x" [ref=f1e85]: "440"
          - generic [ref=f1e86]:
            - text: "y"
            - spinbutton "Object y" [ref=f1e87]: "170"
          - generic [ref=f1e88]:
            - text: width
            - spinbutton "Object width" [ref=f1e89]: "200"
          - generic [ref=f1e90]:
            - text: height
            - spinbutton "Object height" [ref=f1e91]: "150"
        - button "Delete object" [ref=f1e92] [cursor=pointer]
      - generic:
        - generic [ref=f1e96]:
          - button "Zoom out" [ref=f1e97] [cursor=pointer]
          - generic [ref=f1e101]: 100%
          - button "Zoom in" [ref=f1e102] [cursor=pointer]
          - button "Reset zoom" [ref=f1e107] [cursor=pointer]
        - toolbar "Drawing tools" [ref=f1e113]:
          - button "Select" [pressed] [ref=f1e114] [cursor=pointer]:
            - generic [ref=f1e117]: V
          - button "Rectangle" [ref=f1e118] [cursor=pointer]:
            - generic [ref=f1e121]: R
          - button "Ellipse" [ref=f1e122] [cursor=pointer]:
            - generic [ref=f1e125]: O
          - button "Sticky note" [ref=f1e126] [cursor=pointer]:
            - generic [ref=f1e130]: "N"
          - button "Text" [ref=f1e131] [cursor=pointer]:
            - generic [ref=f1e134]: T
        - generic [ref=f1e135]: Select, drag, make it yours
    - complementary [ref=f1e136]:
      - generic [ref=f1e137]:
        - heading "Time machine" [level=2] [ref=f1e142]
        - button "Close history" [ref=f1e143] [cursor=pointer]
      - paragraph [ref=f1e147]: Every idea has a journey.Revisit any moment of yours.
      - generic [ref=f1e148]:
        - generic [ref=f1e149]:
          - generic [ref=f1e150]: LIVE FRAME
          - generic [ref=f1e152]:
            - text: "3"
            - generic [ref=f1e153]: / 3
        - slider "History revision" [ref=f1e154] [cursor=pointer]: "3"
        - generic [ref=f1e155]:
          - generic [ref=f1e156]: Beginning
          - generic [ref=f1e157]: Now
        - generic [ref=f1e158]:
          - button "First event" [ref=f1e159] [cursor=pointer]
          - button "Previous event" [ref=f1e162] [cursor=pointer]
          - button "Play history" [ref=f1e165] [cursor=pointer]
          - button "Next event" [ref=f1e168] [cursor=pointer]
          - button "Latest event" [ref=f1e171] [cursor=pointer]
          - combobox "Playback speed" [ref=f1e174]:
            - option "0.5×"
            - option "1×" [selected]
            - option "2×"
            - option "4×"
      - generic [ref=f1e175]:
        - generic [ref=f1e176]: ACTIVITY
        - generic [ref=f1e177]: 3 events
      - generic [ref=f1e178]:
        - button "QA Cara updated a note 07:05:46 PM · Revision 3" [ref=f1e179] [cursor=pointer]:
          - generic [ref=f1e183]:
            - strong [ref=f1e184]:
              - text: QA Cara
              - generic [ref=f1e185]: updated a note
            - generic [ref=f1e186]:
              - text: 07:05:46 PM
              - generic [ref=f1e187]: ·
              - text: Revision 3
        - button "QA Bob updated a note 07:05:46 PM · Revision 2" [ref=f1e190] [cursor=pointer]:
          - generic [ref=f1e194]:
            - strong [ref=f1e195]:
              - text: QA Bob
              - generic [ref=f1e196]: updated a note
            - generic [ref=f1e197]:
              - text: 07:05:46 PM
              - generic [ref=f1e198]: ·
              - text: Revision 2
        - button "QA Alice added a note 07:05:45 PM · Revision 1" [ref=f1e201] [cursor=pointer]:
          - generic [ref=f1e204]:
            - strong [ref=f1e205]:
              - text: QA Alice
              - generic [ref=f1e206]: added a note
            - generic [ref=f1e207]:
              - text: 07:05:45 PM
              - generic [ref=f1e208]: ·
              - text: Revision 1
      - generic [ref=f1e211]: History is saved to MySQL.Exploring the past keeps the live canvas intact.
  - contentinfo [ref=f1e215]:
    - generic [ref=f1e216]:
      - text: All committed edits saved
      - generic [ref=f1e218]: /
      - text: Revision 3
    - generic [ref=f1e219]:
      - text: 3 online
      - generic [ref=f1e225]: /
      - text: You are QA Alice
```

# Test source

```ts
  10  |   invite,
  11  |   add,
  12  |   snap,
  13  |   assertPersisted,
  14  | } from "./qa/helpers";
  15  | import { execFileSync } from "node:child_process";
  16  | import path from "node:path";
  17  | import AxeBuilder from "@axe-core/playwright";
  18  | 
  19  | test("existing-board share survives reopening in a new browser tab", async ({
  20  |   browser,
  21  |   audit,
  22  | }, info) => {
  23  |   const context = await browser.newContext();
  24  |   await named(context, "QA Alice");
  25  |   const p = await context.newPage();
  26  |   audit.watch(p);
  27  |   const id = await create(p, "QA existing-board sharing");
  28  |   await add(p, "Text", 340, 180);
  29  |   await p.getByRole("link", { name: "Back to workspace", exact: true }).click();
  30  |   await p
  31  |     .getByRole("link")
  32  |     .filter({ hasText: "QA existing-board sharing" })
  33  |     .click();
  34  |   await live(p);
  35  |   await expect(p.getByTestId("board-object")).toHaveCount(1);
  36  |   const reopened = await context.newPage();
  37  |   audit.watch(reopened);
  38  |   await reopened.goto(`http://localhost:8080/board/${id}`);
  39  |   await live(reopened);
  40  |   await snap(reopened, info, "existing-before-share");
  41  |   const link = await invite(reopened, context);
  42  |   const guest = await browser.newContext();
  43  |   await named(guest, "QA Guest");
  44  |   const joined = await guest.newPage();
  45  |   audit.watch(joined);
  46  |   await joined.goto(link);
  47  |   await live(joined);
  48  |   await expect(joined.getByTestId("board-object")).toHaveCount(1);
  49  |   await assertPersisted(joined);
  50  |   await guest.close();
  51  |   await context.close();
  52  | });
  53  | 
  54  | test("three named sessions converge across two servers, concurrent edits, reconnect and restart", async ({
  55  |   browser,
  56  |   audit,
  57  | }, info) => {
  58  |   test.setTimeout(150000);
  59  |   const contexts = await Promise.all([
  60  |     browser.newContext(),
  61  |     browser.newContext(),
  62  |     browser.newContext(),
  63  |   ]);
  64  |   for (let i = 0; i < 3; i++)
  65  |     await named(contexts[i], ["QA Alice", "QA Bob", "QA Cara"][i]);
  66  |   const [a, b, c] = await Promise.all(contexts.map((ctx) => ctx.newPage()));
  67  |   [a, b, c].forEach(audit.watch);
  68  |   const id = await create(a, "QA distributed collaboration");
  69  |   const link = await invite(a, contexts[0]);
  70  |   await b.goto(link.replace(":8080", ":8081"));
  71  |   await c.goto(link);
  72  |   await Promise.all([live(b), live(c)]);
  73  |   await expect(a.locator(".board-footer")).toContainText("3 online", {
  74  |     timeout: 20000,
  75  |   });
  76  |   const objectId = await add(a, "Sticky note", 350, 170);
  77  |   const target = `[data-object-id="${objectId}"]`;
  78  |   await Promise.all([
  79  |     expect(b.locator(target)).toBeVisible(),
  80  |     expect(c.locator(target)).toBeVisible(),
  81  |   ]);
  82  |   await a.getByLabel("Close properties").click();
  83  |   await b.locator(target).click({ position: { x: 100, y: 60 } });
  84  |   await b.getByLabel("Object x", { exact: true }).fill("440");
  85  |   await b.getByLabel("Object x", { exact: true }).press("Tab");
  86  |   await expect(a.locator(target)).toHaveAttribute(
  87  |     "transform",
  88  |     "translate(440,170)",
  89  |   );
  90  |   await expect(c.locator(target)).toHaveAttribute(
  91  |     "transform",
  92  |     "translate(440,170)",
  93  |   );
  94  |   const bounds = await a.getByTestId("canvas").boundingBox();
  95  |   await a.mouse.move(bounds!.x + 600, bounds!.y + 300);
  96  |   await expect(
  97  |     b.getByTestId("canvas").getByText("QA Alice", { exact: true }),
  98  |   ).toBeVisible();
  99  |   await a.locator(target).click({ position: { x: 90, y: 50 } });
  100 |   await c.locator(target).click({ position: { x: 90, y: 50 } });
  101 |   const before = (await state(a)).revision;
  102 |   await Promise.all([
  103 |     a.getByLabel("Object text").fill("Alice's simultaneous idea"),
  104 |     c.getByLabel("Object text").fill("Cara's simultaneous idea"),
  105 |   ]);
  106 |   await Promise.all([
  107 |     a.getByLabel("Object text").press("Tab"),
  108 |     c.getByLabel("Object text").press("Tab"),
  109 |   ]);
> 110 |   await expect.poll(async () => (await state(a)).revision).toBe(before + 2);
      |                                                            ^ Error: expect(received).toBe(expected) // Object.is equality
  111 |   const last = (await events(a)).at(-1)!;
  112 |   expect(["QA Alice", "QA Cara"]).toContain(last.operation.actor);
  113 |   for (const page of [a, b, c])
  114 |     await expect(page.locator(target)).toContainText(last.object.text);
  115 |   await expect
  116 |     .poll(async () => JSON.stringify((await state(b)).objects))
  117 |     .toBe(JSON.stringify((await state(a)).objects));
  118 |   audit.fault(true);
  119 |   await contexts[1].setOffline(true);
  120 |   await expect(b.locator(".connection")).toContainText("Reconnecting");
  121 |   await b.getByRole("button", { name: "Rectangle", exact: true }).click();
  122 |   await b.getByTestId("canvas").click({ position: { x: 760, y: 380 } });
  123 |   await expect(b.locator(".connection")).toContainText("1 pending");
  124 |   const count = Object.values((await state(a)).objects).filter(
  125 |     (o) => !o.deleted,
  126 |   ).length;
  127 |   await a.getByLabel("Close properties").click();
  128 |   await add(a, "Ellipse", 700, 130);
  129 |   await snap(b, info, "offline-pending");
  130 |   await contexts[1].setOffline(false);
  131 |   await live(b);
  132 |   await expect(b.locator(".connection")).not.toContainText("pending");
  133 |   for (const page of [a, b, c])
  134 |     await expect(page.getByTestId("board-object")).toHaveCount(count + 2);
  135 |   audit.fault(false);
  136 |   const persisted = await assertPersisted(a);
  137 |   expect(new Set((await events(a)).map((e) => e.operation.id)).size).toBe(
  138 |     persisted.revision,
  139 |   );
  140 |   await b.close();
  141 |   await expect(a.locator(".board-footer")).toContainText("2 online", {
  142 |     timeout: 20000,
  143 |   });
  144 |   const again = await contexts[1].newPage();
  145 |   audit.watch(again);
  146 |   await again.goto(`http://localhost:8081/board/${id}`);
  147 |   await live(again);
  148 |   await expect(again.getByTestId("board-object")).toHaveCount(3);
  149 |   await expect(a.locator(".board-footer")).toContainText("3 online", {
  150 |     timeout: 20000,
  151 |   });
  152 |   audit.fault(true);
  153 |   execFileSync("docker", ["compose", "-p", "syncframe", "restart", "api"], {
  154 |     cwd: path.resolve(".."),
  155 |     stdio: "pipe",
  156 |   });
  157 |   await live(a);
  158 |   await expect.poll(async () => state(a)).toEqual(persisted);
  159 |   await expect(a.getByTestId("board-object")).toHaveCount(3);
  160 |   audit.fault(false);
  161 |   await snap(a, info, "three-session-recovery");
  162 |   await Promise.all(contexts.map((ctx) => ctx.close()));
  163 | });
  164 | 
  165 | test("time machine pause, all speeds, overlay and full replay preserve the changing live document", async ({
  166 |   browser,
  167 |   audit,
  168 | }, info) => {
  169 |   const one = await browser.newContext(),
  170 |     two = await browser.newContext();
  171 |   await named(one, "QA Historian");
  172 |   await named(two, "QA Live editor");
  173 |   const a = await one.newPage(),
  174 |     b = await two.newPage();
  175 |   [a, b].forEach(audit.watch);
  176 |   await create(a, "QA history isolation");
  177 |   const link = await invite(a, one);
  178 |   await b.goto(link);
  179 |   await live(b);
  180 |   const id = await add(a, "Rectangle", 350, 160);
  181 |   await a.getByLabel("Object text").fill("Original frame");
  182 |   await a.getByLabel("Object text").press("Tab");
  183 |   await settled(a);
  184 |   for (let i = 0; i < 6; i++) {
  185 |     await a.getByLabel("Object x", { exact: true }).fill(String(360 + i * 20));
  186 |     await a.getByLabel("Object x", { exact: true }).press("Tab");
  187 |     await expect(a.locator(`[data-object-id="${id}"]`)).toHaveAttribute(
  188 |       "transform",
  189 |       `translate(${360 + i * 20},160)`,
  190 |     );
  191 |   }
  192 |   const original = await assertPersisted(a);
  193 |   await a.getByRole("button", { name: "First event", exact: true }).click();
  194 |   await expect(a.getByTestId("revision")).toHaveText("REV 0");
  195 |   await expect(a.getByTestId("board-object")).toHaveCount(0);
  196 |   await a.getByRole("button", { name: "Next event", exact: true }).click();
  197 |   await expect(a.getByTestId("revision")).toHaveText("REV 1");
  198 |   await expect(a.getByTestId("board-object")).toHaveCount(1);
  199 |   await a.getByRole("button", { name: "Previous event", exact: true }).click();
  200 |   await expect(a.getByTestId("revision")).toHaveText("REV 0");
  201 |   for (const speed of [".5", "1", "2", "4"]) {
  202 |     await a.getByLabel("Playback speed").selectOption(speed);
  203 |     await expect(a.getByLabel("Playback speed")).toHaveValue(speed);
  204 |   }
  205 |   await a.getByLabel("Playback speed").selectOption(".5");
  206 |   await a.getByRole("button", { name: "Play history", exact: true }).click();
  207 |   await expect
  208 |     .poll(async () =>
  209 |       Number((await a.getByTestId("revision").innerText()).replace("REV ", "")),
  210 |     )
```