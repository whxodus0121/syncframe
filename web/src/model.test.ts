import { describe, it, expect } from "vitest";
import { applyEvent, receive, type State, type EditEvent } from "./model";
const initial = (): State => ({
  id: "board",
  title: "Test",
  revision: 0,
  objects: {},
});
const event = (revision: number, x = revision): EditEvent => ({
  boardId: "board",
  revision,
  operation: {
    id: "op-" + revision,
    actor: "A",
    kind: "update",
    objectId: "o",
    baseRevision: 0,
  },
  object: {
    id: "o",
    kind: "rect",
    x,
    y: 0,
    width: 100,
    height: 80,
    color: "#ffffff",
    text: "",
    deleted: false,
  },
  at: "2026-10-24T00:00:00Z",
});
describe("ordered recovery", () => {
  it("buffers reversed notifications and converges when missing revisions arrive", () => {
    let state = initial();
    const buffer = new Map<number, EditEvent>();
    state = receive(state, buffer, event(3));
    state = receive(state, buffer, event(2));
    expect(state.revision).toBe(0);
    state = receive(state, buffer, event(1));
    expect(state.revision).toBe(3);
    expect(state.objects.o.x).toBe(3);
    expect(buffer.size).toBe(0);
  });
  it("ignores duplicates and stale notifications after recovery", () => {
    let state = initial();
    state = applyEvent(state, event(1));
    state = applyEvent(state, event(1, 99));
    expect(state.objects.o.x).toBe(1);
    const buffer = new Map<number, EditEvent>();
    expect(receive(state, buffer, event(1))).toEqual(state);
    expect(buffer.size).toBe(0);
  });
  it("rejects a direct gap", () =>
    expect(() => applyEvent(initial(), event(2))).toThrow("Revision gap"));
  it("replay retains deletion tombstones", () => {
    let state = applyEvent(initial(), event(1));
    const deletion = event(2);
    deletion.object.deleted = true;
    state = applyEvent(state, deletion);
    state = applyEvent(state, event(1));
    expect(state.objects.o.deleted).toBe(true);
  });
  it("converges for every permutation of four events", () => {
    function permute<T>(xs: T[]): T[][] {
      return xs.length === 0
        ? [[]]
        : xs.flatMap((x, i) =>
            permute(xs.filter((_, j) => i !== j)).map((rest) => [x, ...rest]),
          );
    }
    for (const permutation of permute([1, 2, 3, 4])) {
      let state = initial();
      const buffer = new Map<number, EditEvent>();
      for (const revision of permutation)
        state = receive(state, buffer, event(revision));
      expect(state.revision).toBe(4);
      expect(state.objects.o.x).toBe(4);
    }
  });
});
