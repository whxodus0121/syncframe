import { test, expect } from "./qa/helpers";
import AxeBuilder from "@axe-core/playwright";
test("landing and canvas meet automated WCAG A/AA checks", async ({
  page,
  audit,
}) => {
  audit.watch(page);
  await page.goto("/");
  let scan = await new AxeBuilder({ page })
    .withTags(["wcag2a", "wcag2aa", "wcag21aa"])
    .analyze();
  expect(
    scan.violations.map((v) => ({
      id: v.id,
      nodes: v.nodes.map((n) => ({
        target: n.target,
        summary: n.failureSummary,
      })),
    })),
  ).toEqual([]);
  await page
    .getByRole("button", { name: "Or open a real, editable starter board" })
    .click();
  await expect(page.locator(".connection")).toContainText("Live");
  scan = await new AxeBuilder({ page })
    .withTags(["wcag2a", "wcag2aa", "wcag21aa"])
    .analyze();
  expect(
    scan.violations.map((v) => ({
      id: v.id,
      nodes: v.nodes.map((n) => ({
        target: n.target,
        summary: n.failureSummary,
      })),
    })),
  ).toEqual([]);
});
