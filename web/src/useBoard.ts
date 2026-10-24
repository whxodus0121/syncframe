import { useCallback, useEffect, useRef, useState } from "react";
import {
  api,
  post,
  receive,
  uid,
  type State,
  type Operation,
  type EditEvent,
  type Peer,
} from "./model";

export function useBoard(id: string, name: string) {
  const [state, setState] = useState<State | null>(null),
    [events, setEvents] = useState<EditEvent[]>([]),
    [status, setStatus] = useState("Connecting"),
    [error, setError] = useState(""),
    [peers, setPeers] = useState<Peer[]>([]),
    [redisOK, setRedisOK] = useState(true),
    [pendingCount, setPendingCount] = useState(0),
    [cursors, setCursors] = useState<
      Record<string, { peer: Peer; x: number; y: number; at: number }>
    >({});
  const current = useRef<State | null>(null),
    socket = useRef<WebSocket | null>(null),
    self = useRef(""),
    buffer = useRef(new Map<number, EditEvent>()),
    pending = useRef<Operation[]>([]),
    lastCursor = useRef(0);
  const storageKey = "sf:pending:" + id + ":" + name;
  const persist = useCallback(() => {
    try {
      localStorage.setItem(storageKey, JSON.stringify(pending.current));
    } catch {
      setError(
        "Device storage is full. Keep this tab open until pending edits are confirmed.",
      );
    }
    setPendingCount(pending.current.length);
  }, [storageKey]);
  useEffect(() => {
    let disposed = false,
      timer: ReturnType<typeof setTimeout>,
      attempt = 0,
      recovering = false,
      lastMessage = Date.now();
    try {
      pending.current = JSON.parse(localStorage.getItem(storageKey) || "[]");
    } catch {
      pending.current = [];
    }
    setPendingCount(pending.current.length);
    const update = (value: State) => {
      current.current = value;
      if (!disposed) setState(value);
    };
    const accept = (event: EditEvent) => {
      if (!current.current || disposed) return;
      update(receive(current.current, buffer.current, event));
      if (buffer.current.size > 0) void recover();
    };
    const recover = async () => {
      if (recovering || !current.current) return;
      recovering = true;
      try {
        let more = true;
        while (more && !disposed) {
          const page = await api<EditEvent[]>(
            `/api/boards/${id}/events?after=${current.current!.revision}`,
          );
          for (const event of page) accept(event);
          more = page.length === 500;
        }
      } catch (e) {
        if (!disposed) setError(String(e));
      } finally {
        recovering = false;
      }
    };
    const connect = () => {
      if (disposed) return;
      if (!navigator.onLine) {
        setStatus("Reconnecting");
        timer = setTimeout(connect, 2000);
        return;
      }
      setStatus(attempt ? "Reconnecting" : "Connecting");
      const ws = new WebSocket(
        `${location.protocol === "https:" ? "wss" : "ws"}://${location.host}/api/boards/${id}/ws?after=${current.current?.revision ?? 0}&name=${encodeURIComponent(name)}`,
      );
      socket.current = ws;
      lastMessage = Date.now();
      ws.onopen = () => {
        attempt = 0;
        setStatus("Live");
        setError("");
        for (const op of pending.current)
          ws.send(JSON.stringify({ type: "operation", operation: op }));
      };
      ws.onmessage = (msg) => {
        if (disposed) return;
        lastMessage = Date.now();
        const data = JSON.parse(msg.data);
        switch (data.type) {
          case "hello":
            self.current = data.selfId;
            break;
          case "snapshot":
            self.current = data.selfId;
            buffer.current.clear();
            update(data.state);
            break;
          case "event":
            accept(data.event);
            break;
          case "ack":
            pending.current = pending.current.filter(
              (op) => op.id !== data.operationId,
            );
            persist();
            break;
          case "rejected":
            if (data.retryable) {
              setError(
                "Storage unavailable. Your pending edit will retry on reconnect.",
              );
              ws.close();
            } else {
              pending.current = pending.current.filter(
                (op) => op.id !== data.operationId,
              );
              persist();
              setError(data.error);
            }
            break;
          case "presence":
            setPeers(data.users);
            setRedisOK(data.redisAvailable);
            setCursors((old) =>
              Object.fromEntries(
                Object.entries(old).filter(
                  ([key, value]) =>
                    data.users.some((p: Peer) => p.id === key) &&
                    Date.now() - value.at < 15000,
                ),
              ),
            );
            break;
          case "cursor":
            if (data.peer.id !== self.current)
              setCursors((old) => ({
                ...old,
                [data.peer.id]: {
                  peer: data.peer,
                  x: data.x,
                  y: data.y,
                  at: Date.now(),
                },
              }));
            break;
        }
      };
      ws.onerror = () => ws.close();
      ws.onclose = () => {
        if (disposed || socket.current !== ws) return;
        setStatus("Reconnecting");
        setPeers([]);
        attempt++;
        clearTimeout(timer);
        timer = setTimeout(
          connect,
          Math.min(500 * 2 ** Math.min(attempt, 4), 8000),
        );
      };
    };
    const offline = () => {
      setStatus("Reconnecting");
      socket.current?.close();
    };
    const online = () => {
      if (
        socket.current?.readyState === WebSocket.OPEN ||
        socket.current?.readyState === WebSocket.CONNECTING
      )
        return;
      clearTimeout(timer);
      connect();
    };
    window.addEventListener("offline", offline);
    window.addEventListener("online", online);
    const watchdog = setInterval(() => {
      if (
        socket.current?.readyState === WebSocket.OPEN &&
        Date.now() - lastMessage > 18000
      )
        socket.current.close();
    }, 3000);
    (async () => {
      try {
        const token = new URLSearchParams(location.hash.slice(1)).get("key");
        if (token) {
          await api(`/api/boards/${id}/session`, post({ token }));
          sessionStorage.setItem("sf:key:" + id, token);
          history.replaceState(null, "", location.pathname);
        }
        const initial = await api<State>(`/api/boards/${id}`);
        if (disposed) return;
        update(initial);
        const recents = JSON.parse(
          localStorage.getItem("sf:recent") || "[]",
        ).filter((b: { id: string }) => b.id !== id);
        localStorage.setItem(
          "sf:recent",
          JSON.stringify(
            [{ id, title: initial.title }, ...recents].slice(0, 12),
          ),
        );
        connect();
      } catch (e) {
        if (!disposed) {
          setError(String(e));
          setStatus("Unavailable");
        }
      }
    })();
    return () => {
      disposed = true;
      clearTimeout(timer);
      clearInterval(watchdog);
      window.removeEventListener("offline", offline);
      window.removeEventListener("online", online);
      if (socket.current) {
        socket.current.onclose = null;
        socket.current.close();
      }
    };
  }, [id, name, persist, storageKey]);
  useEffect(() => {
    if (!state) return;
    let disposed = false;
    const timer = setTimeout(() => {
      (async () => {
        try {
          let all: EditEvent[] = [];
          let more = true;
          while (more) {
            const page = await api<EditEvent[]>(
              `/api/boards/${id}/events?after=${all.at(-1)?.revision ?? 0}`,
            );
            all = all.concat(page);
            more = page.length === 500;
          }
          if (!disposed) setEvents(all);
        } catch (e) {
          if (!disposed) setError(String(e));
        }
      })();
    }, 250);
    return () => {
      disposed = true;
      clearTimeout(timer);
    };
  }, [id, state?.revision]);
  const send = useCallback(
    (part: Omit<Operation, "id" | "actor" | "baseRevision">) => {
      if (!current.current) return;
      if (pending.current.length >= 200) {
        setError("Pending edit limit reached. Reconnect before editing more.");
        return;
      }
      const op: Operation = {
        ...part,
        id: uid(),
        actor: name,
        baseRevision: current.current.revision,
      };
      pending.current.push(op);
      persist();
      if (socket.current?.readyState === WebSocket.OPEN)
        socket.current.send(
          JSON.stringify({ type: "operation", operation: op }),
        );
    },
    [name, persist],
  );
  const cursor = useCallback((x: number, y: number) => {
    if (
      Date.now() - lastCursor.current < 60 ||
      socket.current?.readyState !== WebSocket.OPEN
    )
      return;
    lastCursor.current = Date.now();
    socket.current.send(JSON.stringify({ type: "cursor", x, y }));
  }, []);
  return {
    state,
    events,
    status,
    error,
    setError,
    peers,
    cursors,
    redisOK,
    pendingCount,
    send,
    cursor,
  };
}
