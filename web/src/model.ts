export type CanvasObject = {
  id: string;
  kind: "rect" | "ellipse" | "text" | "note";
  x: number;
  y: number;
  width: number;
  height: number;
  color: string;
  text: string;
  deleted: boolean;
};
export type State = {
  id: string;
  title: string;
  revision: number;
  objects: Record<string, CanvasObject>;
};
export type Operation = {
  id: string;
  actor: string;
  kind: "create" | "update" | "delete";
  objectId: string;
  baseRevision: number;
  object?: CanvasObject;
  patch?: Partial<
    Pick<CanvasObject, "x" | "y" | "width" | "height" | "color" | "text">
  >;
};
export type EditEvent = {
  boardId: string;
  revision: number;
  operation: Operation;
  object: CanvasObject;
  at: string;
};
export type Peer = { id: string; name: string; color: string };
export function applyEvent(state: State, event: EditEvent): State {
  if (event.revision <= state.revision) return state;
  if (event.revision !== state.revision + 1) throw new Error("Revision gap");
  return {
    ...state,
    revision: event.revision,
    objects: { ...state.objects, [event.object.id]: event.object },
  };
}
// Buffer out-of-order events; a caller must fetch missing revisions from the server.
export function receive(
  state: State,
  pending: Map<number, EditEvent>,
  event: EditEvent,
): State {
  if (event.revision > state.revision) pending.set(event.revision, event);
  while (pending.has(state.revision + 1)) {
    const next = pending.get(state.revision + 1)!;
    pending.delete(next.revision);
    state = applyEvent(state, next);
  }
  for (const revision of pending.keys())
    if (revision <= state.revision) pending.delete(revision);
  return state;
}
export async function api<T>(path: string, options?: RequestInit): Promise<T> {
  const response = await fetch(path, {
    ...options,
    headers: { "Content-Type": "application/json", ...options?.headers },
  });
  const data = await response.json();
  if (!response.ok)
    throw new Error(data.error ?? `Request failed (${response.status})`);
  return data as T;
}
export const post = (data: unknown): RequestInit => ({
  method: "POST",
  body: JSON.stringify(data),
});
export const uid = () => crypto.randomUUID();
