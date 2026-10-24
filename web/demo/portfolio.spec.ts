import {
  test,
  expect,
  named,
  create,
  live,
  state,
  settled,
  invite,
  add,
  assertPersisted,
} from "../e2e/qa/helpers";
import fs from "node:fs";
import path from "node:path";

test("record collaboration, offline recovery and read-only replay", async ({
  browser,
  audit,
}) => {
  const folder = path.resolve("../artifacts/portfolio/demo");
  fs.mkdirSync(folder, { recursive: true });
  const a = await browser.newContext({
    viewport: { width: 1440, height: 960 },
    recordVideo: { dir: folder, size: { width: 1440, height: 960 } },
  });
  const b = await browser.newContext({
    viewport: { width: 1440, height: 960 },
  });
  await named(a, "Alice");
  await named(b, "Bob");
  const p = await a.newPage(),
    q = await b.newPage();
  audit.watch(p);
  audit.watch(q);
  // Delays here are for legible video pacing; assertions establish correctness.
  const shot = async (name: string) => {
    await expect(
      p.getByText("Loading saved activity…", { exact: true }),
    ).not.toBeVisible();
    await p.screenshot({ path: path.join(folder, name + ".png") });
    await p.waitForTimeout(1000);
  };
  await create(p, "Designing reliable collaboration");
  await q.goto((await invite(p, a)).replace(":8080", ":8081"));
  await live(q);
  await expect(p.locator(".board-footer")).toContainText("2 online");
  const id = await add(p, "Sticky note", 350, 180);
  await p.getByLabel("Object text").fill("MySQL commits define the order");
  await p.getByLabel("Object text").press("Tab");
  await expect(q.getByTestId("board-object")).toContainText(
    "MySQL commits define the order",
  );
  await p.getByLabel("Close properties").click();
  await shot("01-collaboration");
  await q.locator(`[data-object-id="${id}"]`).click();
  await q.getByLabel("Object x", { exact: true }).fill("560");
  await q.getByLabel("Object x", { exact: true }).press("Tab");
  await expect(p.locator(`[data-object-id="${id}"]`)).toHaveAttribute(
    "transform",
    "translate(560,180)",
  );
  const box = await q.getByTestId("canvas").boundingBox();
  await q.mouse.move(box!.x + 480, box!.y + 330);
  await expect(
    p.getByTestId("canvas").getByText("Bob", { exact: true }),
  ).toBeVisible();
  await shot("02-remote-move");
  audit.fault(true);
  await a.setOffline(true);
  await expect(p.locator(".connection")).toContainText("Reconnecting");
  await p.getByRole("button", { name: "Rectangle", exact: true }).click();
  await p.getByTestId("canvas").click({ position: { x: 350, y: 390 } });
  await expect(p.locator(".connection")).toContainText("1 pending");
  await shot("03-offline-pending");
  await a.setOffline(false);
  await live(p);
  await settled(p);
  await expect(q.getByTestId("board-object")).toHaveCount(2);
  audit.fault(false);
  await p.getByLabel("Close properties").click();
  const final = await assertPersisted(p);
  await shot("04-recovered");
  await p.getByRole("button", { name: "First event", exact: true }).click();
  await expect(p.getByTestId("revision")).toHaveText("REV 0");
  await shot("05-history-start");
  await p.getByRole("button", { name: "Next event", exact: true }).click();
  await expect(p.getByTestId("revision")).toHaveText("REV 1");
  await p.getByRole("checkbox", { name: "Overlay live positions" }).check();
  await expect(p.locator('svg.canvas > rect[stroke="#16a57a"]')).toHaveCount(2);
  await shot("06-history-overlay");
  await p.getByLabel("Playback speed").selectOption(".5");
  await p.getByRole("button", { name: "Play history", exact: true }).click();
  await expect(p.getByTestId("revision")).not.toHaveText("REV 1");
  await p.getByRole("button", { name: "Pause history", exact: true }).click();
  expect(await state(p)).toEqual(final);
  await shot("07-history-paused");
  await p.getByRole("button", { name: "Canvas", exact: true }).click();
  await expect(p.getByTestId("revision")).toHaveText(`REV ${final.revision}`);
  await shot("08-live");
  const video = p.video()!;
  await a.close();
  await b.close();
  await video.saveAs(path.join(folder, "syncframe-demo.webm"));
  await video.delete();
  fs.writeFileSync(
    path.join(folder, "capture.json"),
    JSON.stringify(
      {
        at: new Date().toISOString(),
        browser: browser.version(),
        viewport: "1440x960",
        scenarios: [
          "two independent sessions through 8080/8081",
          "remote movement and cursor",
          "offline pending and same-ID recovery",
          "read-only replay and overlay",
        ],
        liveRevision: final.revision,
        rawVideo: "syncframe-demo.webm",
      },
      null,
      2,
    ),
  );
});
