import { type Page } from "@playwright/test";
import { test, expect } from "./qa/helpers";
import fs from "node:fs";
const base = process.env.BASE_URL || "http://127.0.0.1:8080";
const waitLive = async (page: Page) =>
  expect(page.locator(".connection")).toContainText("Live");
test("two browser contexts edit, resize, recolor, recover offline edits and replay history", async ({
  browser,
  audit,
}) => {
  const a = await browser.newContext({
      viewport: { width: 1440, height: 960 },
    }),
    b = await browser.newContext({ viewport: { width: 1440, height: 960 } });
  const pa = await a.newPage(),
    pb = await b.newPage();
  [pa, pb].forEach(audit.watch);
  const errors: string[] = [];
  pa.on("pageerror", (err) => errors.push(err.message));
  pb.on("pageerror", (err) => errors.push(err.message));
  await pa.goto("/");
  await pa.getByLabel("Board title").fill("Product discovery");
  await pa.getByRole("button", { name: "New board", exact: true }).click();
  await waitLive(pa);
  const id = pa.url().split("/").at(-1)!;
  const token = await pa.evaluate(
    (id) => sessionStorage.getItem("sf:key:" + id),
    id,
  );
  expect(token).toBeTruthy();
  await pb.goto(`${base}/board/${id}#key=${token}`);
  await waitLive(pb);
  await expect(pa.locator(".board-footer")).toContainText("2 online", {
    timeout: 15000,
  });
  const displayName = await pa.evaluate(() =>
    sessionStorage.getItem("sf:name"),
  );
  const canvasBounds = await pa.getByTestId("canvas").boundingBox();
  await pa.mouse.move(canvasBounds!.x + 500, canvasBounds!.y + 150);
  await expect(
    pb.getByTestId("canvas").getByText(displayName!, { exact: true }),
  ).toBeVisible();
  await pa.getByRole("button", { name: "Sticky note", exact: true }).click();
  await pa.getByTestId("canvas").click({ position: { x: 330, y: 190 } });
  await expect(pa.getByTestId("board-object")).toHaveCount(1);
  await expect(pb.getByTestId("board-object")).toHaveCount(1);
  await pa.getByLabel("Object text").fill("Design a calmer collaboration flow");
  await pa.getByLabel("Object text").press("Tab");
  await expect(pb.getByTestId("board-object")).toContainText(
    "Design a calmer collaboration flow",
  );
  await pa.getByRole("button", { name: "Color #d7f0e5", exact: true }).click();
  await expect(
    pb.getByTestId("board-object").locator("rect").first(),
  ).toHaveAttribute("fill", "#d7f0e5");
  await pa.getByLabel("Object width", { exact: true }).fill("245");
  await pa.getByLabel("Object width", { exact: true }).press("Tab");
  await expect(
    pb.getByTestId("board-object").locator("rect").first(),
  ).toHaveAttribute("width", "245");
  const resize = await pa.getByTestId("resize-handle").boundingBox();
  await pa.mouse.move(resize!.x + 6, resize!.y + 6);
  await pa.mouse.down();
  await pa.mouse.move(resize!.x + 36, resize!.y + 26, { steps: 5 });
  await pa.mouse.up();
  await expect(
    pb.getByTestId("board-object").locator("rect").first(),
  ).toHaveAttribute("width", "275");
  const beforeKeyboard = await pa
    .getByTestId("board-object")
    .getAttribute("transform");
  await pa.getByTestId("board-object").press("ArrowRight");
  await expect(pa.getByTestId("board-object")).not.toHaveAttribute(
    "transform",
    beforeKeyboard!,
  );
  const initial = await pa
    .getByTestId("board-object")
    .getAttribute("transform");
  const box = await pa.getByTestId("board-object").boundingBox();
  expect(box).toBeTruthy();
  await pa.mouse.move(box!.x + 100, box!.y + 50);
  await pa.mouse.down();
  await pa.mouse.move(box!.x + 180, box!.y + 100, { steps: 8 });
  await pa.mouse.up();
  await expect(pa.getByTestId("board-object")).not.toHaveAttribute(
    "transform",
    initial!,
  );
  await expect(pb.getByTestId("board-object")).toHaveAttribute(
    "transform",
    (await pa.getByTestId("board-object").getAttribute("transform"))!,
  );
  const before = Number(
    (await pb.getByTestId("revision").innerText()).replace("REV ", ""),
  );
  audit.fault(true);
  await b.setOffline(true);
  await expect(pb.locator(".connection")).toContainText("Reconnecting");
  // Create through actual UI while disconnected; operation must be queued then committed once.
  await pb.getByRole("button", { name: "Rectangle", exact: true }).click();
  await pb.getByTestId("canvas").click({ position: { x: 650, y: 420 } });
  await expect(pb.locator(".connection")).toContainText("pending");
  await pa.getByRole("button", { name: "Ellipse", exact: true }).click();
  await pa.getByTestId("canvas").click({ position: { x: 650, y: 180 } });
  await expect(pa.getByTestId("board-object")).toHaveCount(2);
  await b.setOffline(false);
  await waitLive(pb);
  await expect(pa.getByTestId("board-object")).toHaveCount(3);
  await expect(pb.getByTestId("board-object")).toHaveCount(3);
  await expect(pb.locator(".connection")).not.toContainText("pending");
  await expect(pb.getByTestId("revision")).toHaveText(`REV ${before + 2}`);
  audit.fault(false);
  await pa.getByRole("button", { name: "First event", exact: true }).click();
  await expect(pa.getByTestId("board-object")).toHaveCount(0);
  await expect(pa.getByTestId("revision")).toHaveText("REV 0");
  await expect(pb.getByTestId("board-object")).toHaveCount(3);
  await pa.getByRole("button", { name: "Next event", exact: true }).click();
  await expect(pa.getByTestId("board-object")).toHaveCount(1);
  await pa.getByRole("combobox", { name: "Playback speed" }).selectOption("4");
  await pa.getByRole("button", { name: "Play history", exact: true }).click();
  await expect(
    pa.getByRole("button", { name: "Play history", exact: true }),
  ).toBeVisible();
  await expect(pa.getByTestId("board-object")).toHaveCount(3);
  await pa.getByRole("checkbox", { name: "Overlay live positions" }).check();
  await pa
    .getByRole("button", { name: "Back to live canvas", exact: true })
    .click();
  await pa
    .getByTestId("board-object")
    .first()
    .click({ position: { x: 100, y: 50 } });
  await pa.getByRole("button", { name: "Delete object", exact: true }).click();
  await expect(pb.getByTestId("board-object")).toHaveCount(2);
  await pa.reload();
  await waitLive(pa);
  await expect(pa.getByTestId("board-object")).toHaveCount(2);
  expect(errors).toEqual([]);
  fs.mkdirSync("../artifacts/screenshots", { recursive: true });
  await pa.screenshot({
    path: "../artifacts/screenshots/collaboration.png",
    fullPage: true,
  });
  await a.close();
  await b.close();
});
test("landing page and starter board are functional at desktop and phone sizes", async ({
  browser,
  audit,
}) => {
  const page = await browser.newPage({
    viewport: { width: 1440, height: 1080 },
  });
  audit.watch(page);
  await page.goto("/");
  await expect(
    page.getByRole("heading", { name: "Think together. Keep every frame." }),
  ).toBeVisible();
  await page.screenshot({
    path: "../artifacts/screenshots/home.png",
    fullPage: true,
  });
  await page
    .getByRole("button", { name: "Or open a real, editable starter board" })
    .click();
  await waitLive(page);
  await expect(page.getByTestId("board-object")).toHaveCount(5);
  await expect(page.locator(".event-row")).toHaveCount(5);
  await page.screenshot({
    path: "../artifacts/screenshots/starter-board.png",
    fullPage: true,
  });
  await page.setViewportSize({ width: 390, height: 844 });
  await page
    .getByRole("button", { name: "Close history", exact: true })
    .click();
  await expect(
    page.getByRole("toolbar", { name: "Drawing tools" }),
  ).toBeVisible();
  await expect(
    page.getByRole("button", { name: "Share board", exact: true }),
  ).toBeVisible();
  expect(
    await page.evaluate(
      () => document.documentElement.scrollWidth <= window.innerWidth,
    ),
  ).toBe(true);
  await page.screenshot({
    path: "../artifacts/screenshots/mobile.png",
    fullPage: true,
  });
  await page.close();
});
