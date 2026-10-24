import React, { useEffect, useRef, useState } from "react";
import { createRoot } from "react-dom/client";
import {
  ArrowDownLeft,
  ArrowRight,
  Check,
  ChevronLeft,
  ChevronRight,
  Circle,
  Clock3,
  Copy,
  ExternalLink,
  Frame,
  GitBranch,
  History,
  Layers,
  Loader2,
  Maximize2,
  MousePointer2,
  Pause,
  Play,
  Plus,
  RotateCcw,
  Share2,
  SkipBack,
  SkipForward,
  Square,
  StickyNote,
  Trash2,
  Type,
  Users,
  WifiOff,
  X,
  ZoomIn,
  ZoomOut,
} from "lucide-react";
import { api, post, uid, type State, type CanvasObject } from "./model";
import { useBoard } from "./useBoard";
import "./style.css";

const palette = [
  "#fff2b3",
  "#dce7ff",
  "#e4dcfc",
  "#d7f0e5",
  "#ffdfe0",
  "#ffffff",
];
const author = () => {
  let name = sessionStorage.getItem("sf:name");
  if (!name) {
    name = "Maker " + Math.floor(100 + Math.random() * 900);
    sessionStorage.setItem("sf:name", name);
  }
  return name;
};
function Logo() {
  return (
    <a href="/" className="brand" aria-label="SYNCFRAME home">
      <span className="brand-mark">
        <Frame size={20} />
      </span>
      <span>
        syncframe<span className="brand-dot">.</span>
      </span>
    </a>
  );
}

async function createBoard(title: string, demo = false) {
  const { board, token } = await api<{ board: State; token: string }>(
    "/api/boards",
    post({ title }),
  );
  sessionStorage.setItem("sf:key:" + board.id, token);
  if (demo) {
    const specs: [
      CanvasObject["kind"],
      number,
      number,
      number,
      number,
      string,
      string,
    ][] = [
      [
        "text",
        110,
        70,
        600,
        60,
        "#ffffff",
        "A little space for the next big idea.",
      ],
      [
        "note",
        120,
        200,
        210,
        170,
        "#fff2b3",
        "01 / THE CHALLENGE\n\nMake working together feel effortless.",
      ],
      [
        "note",
        400,
        200,
        210,
        170,
        "#e4dcfc",
        "02 / THE IDEA\n\nOne shared canvas. Every change remembered.",
      ],
      [
        "note",
        680,
        200,
        210,
        170,
        "#d7f0e5",
        "03 / WHAT’S NEXT\n\nInvite someone. Move an idea forward.",
      ],
      [
        "rect",
        270,
        440,
        450,
        90,
        "#dce7ff",
        "Start a conversation. Build something together.",
      ],
    ];
    let rev = 0;
    for (const [kind, x, y, width, height, color, text] of specs) {
      const id = uid();
      await api(
        `/api/boards/${board.id}/operations`,
        post({
          id: uid(),
          actor: author(),
          kind: "create",
          objectId: id,
          baseRevision: rev++,
          object: {
            id,
            kind,
            x,
            y,
            width,
            height,
            color,
            text,
            deleted: false,
          },
        }),
      );
    }
  }
  location.href = `/board/${board.id}`;
}
function Home() {
  const [title, setTitle] = useState("Untitled canvas"),
    [busy, setBusy] = useState(false),
    [error, setError] = useState("");
  const recents: { id: string; title: string }[] = JSON.parse(
    localStorage.getItem("sf:recent") || "[]",
  );
  const create = async (demo = false) => {
    setBusy(true);
    try {
      await createBoard(demo ? "The next big idea" : title, demo);
    } catch (e) {
      setError(String(e));
      setBusy(false);
    }
  };
  return (
    <div className="home">
      <nav className="home-nav">
        <Logo />
        <span className="nav-note">
          <span className="live-dot" /> A shared canvas with a memory
        </span>
        <a className="subtle-link" href="#workspace">
          Your workspace <ArrowDownLeft size={15} />
        </a>
      </nav>
      <main className="home-main">
        <div className="eyebrow">
          <span className="tiny-grid" /> ROOM FOR YOUR NEXT IDEA
        </div>
        <h1>
          Think together.
          <br />
          <span>Keep every frame.</span>
        </h1>
        <p className="hero-copy">
          A quiet place for your team's loudest ideas.
          <br />
          Create in real time. Revisit every moment.
        </p>
        <div className="create-box" id="workspace">
          <label className="sr-only" htmlFor="board-title">
            Board title
          </label>
          <input
            id="board-title"
            value={title}
            maxLength={100}
            onChange={(e) => setTitle(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter" && !busy) void create();
            }}
          />
          <button
            className="primary"
            disabled={busy || !title.trim()}
            onClick={() => void create()}
          >
            {busy ? <Loader2 className="spin" size={17} /> : <Plus size={17} />}{" "}
            New board <ArrowRight size={16} />
          </button>
        </div>
        <button
          className="demo-link"
          disabled={busy}
          onClick={() => void create(true)}
        >
          Or open a real, editable starter board <ArrowRight size={14} />
        </button>
        {error && (
          <div className="error" role="alert">
            {error}
          </div>
        )}
        <div className="home-features">
          <span>
            <Users size={17} /> Together, in real time
          </span>
          <span>
            <History size={17} /> Every edit, remembered
          </span>
          <span>
            <GitBranch size={17} /> One consistent canvas
          </span>
        </div>
        <div
          className="canvas-preview"
          aria-label="Canvas concept illustration"
        >
          <div className="preview-top">
            <span className="preview-dots">● ● ●</span>
            <span>Your next big idea</span>
            <span className="preview-live">● Live canvas</span>
          </div>
          <div className="preview-grid">
            <div className="preview-heading">
              Good things start with a little space.
            </div>
            <div className="preview-note yellow">
              <small>THE WHAT IF</small>
              <p>
                What could we
                <br />
                make together?
              </p>
              <span>↗ Start here</span>
            </div>
            <div className="preview-arrow">⤳</div>
            <div className="preview-note lavender">
              <small>THE POSSIBILITY</small>
              <p>
                Something worth
                <br />
                coming back to.
              </p>
              <span>One idea at a time.</span>
            </div>
            <div className="preview-label">
              <MousePointer2 size={16} fill="currentColor" /> Your imagination
            </div>
            <div className="preview-footer">
              <Frame size={14} /> An open space. An unbroken history.
            </div>
          </div>
        </div>
        {recents.length > 0 && (
          <section className="recent">
            <div className="section-title">
              <h2>Pick up where you left off</h2>
              <span>On this device</span>
            </div>
            <div className="recent-grid">
              {recents.map((b) => (
                <a href={`/board/${b.id}`} className="recent-card" key={b.id}>
                  <span className="board-thumb">
                    <Layers size={25} />
                  </span>
                  <strong>{b.title}</strong>
                  <span>
                    Open canvas <ArrowUpRightIcon />
                  </span>
                </a>
              ))}
            </div>
          </section>
        )}
        <footer className="home-footer">
          <span>SYNCFRAME / A SPACE TO THINK TOGETHER</span>
          <span>Built on real events. Designed for shared ideas.</span>
        </footer>
      </main>
    </div>
  );
}
function ArrowUpRightIcon() {
  return <ExternalLink size={13} />;
}

function Board({ id }: { id: string }) {
  const [name] = useState(author);
  const b = useBoard(id, name);
  const [selected, setSelected] = useState<string | null>(null),
    [tool, setTool] = useState<CanvasObject["kind"] | "select">("select"),
    [historyOpen, setHistoryOpen] = useState(
      () => window.matchMedia("(min-width: 801px)").matches,
    ),
    [revision, setRevision] = useState<number | null>(null),
    [past, setPast] = useState<State | null>(null),
    [playing, setPlaying] = useState(false),
    [speed, setSpeed] = useState(1),
    [zoom, setZoom] = useState(1),
    [toast, setToast] = useState(""),
    [share, setShare] = useState(false),
    [draft, setDraft] = useState<CanvasObject | null>(null),
    [textDraft, setTextDraft] = useState(""),
    [dimensionDrafts, setDimensionDrafts] = useState<Record<string, string>>(
      {},
    ),
    [comparing, setComparing] = useState(false);
  const textDirty = useRef(false),
    textSelection = useRef<string | null>(null);
  const svg = useRef<SVGSVGElement>(null),
    drag = useRef<{
      id: string;
      x: number;
      y: number;
      object: CanvasObject;
      resize: boolean;
    } | null>(null),
    toastTimer = useRef<ReturnType<typeof setTimeout>>(undefined);
  const travel = revision !== null;
  const historyLoading = travel && past?.revision !== revision;
  const shown = travel ? (historyLoading ? null : past) : b.state;
  const objects = Object.values(shown?.objects ?? {}).filter((o) => !o.deleted);
  const object = selected ? shown?.objects[selected] : undefined;
  const inspectorOnLeft =
    object &&
    (object.x + object.width / 2) * zoom -
      (svg.current?.parentElement?.scrollLeft ?? 0) >
      (svg.current?.parentElement?.clientWidth ?? innerWidth) / 2;
  const liveCount = Object.values(b.state?.objects ?? {}).filter(
    (o) => !o.deleted,
  ).length;
  const notify = (message: string) => {
    setToast(message);
    clearTimeout(toastTimer.current);
    toastTimer.current = setTimeout(() => setToast(""), 3500);
  };
  useEffect(() => {
    if (revision === null) {
      setPast(null);
      return;
    }
    let active = true;
    api<State>(`/api/boards/${id}?revision=${revision}`)
      .then((s) => {
        if (active) setPast(s);
      })
      .catch((e) => {
        if (active) b.setError(String(e));
      });
    return () => {
      active = false;
    };
  }, [id, revision]);
  useEffect(() => {
    if (!playing) return;
    const timer = setInterval(
      () =>
        setRevision((r) => {
          const next = (r ?? 0) + 1;
          if (next >= (b.state?.revision ?? 0)) {
            setPlaying(false);
            return b.state?.revision ?? 0;
          }
          return next;
        }),
      900 / speed,
    );
    return () => clearInterval(timer);
  }, [playing, speed, b.state?.revision]);
  useEffect(() => {
    if (textSelection.current !== selected) {
      textSelection.current = selected;
      textDirty.current = false;
      setDimensionDrafts({});
      setTextDraft(object?.text ?? "");
    } else if (!textDirty.current) setTextDraft(object?.text ?? "");
  }, [object?.text, selected]);
  useEffect(() => {
    const handle = (e: KeyboardEvent) => {
      if ((e.target as HTMLElement).closest("input,textarea,select")) return;
      if (e.key === "Escape") {
        setSelected(null);
        setTool("select");
        setShare(false);
      }
      if (e.key === "Delete" || e.key === "Backspace") {
        if (selected && !travel) {
          e.preventDefault();
          b.send({ kind: "delete", objectId: selected });
          setSelected(null);
        }
      }
      if (e.key === "v") setTool("select");
      if (e.key === "r") setTool("rect");
      if (e.key === "o") setTool("ellipse");
      if (e.key === "n") setTool("note");
      if (e.key === "t") setTool("text");
    };
    window.addEventListener("keydown", handle);
    return () => window.removeEventListener("keydown", handle);
  }, [selected, travel, b.send]);
  const point = (e: React.PointerEvent) => {
    const rect = svg.current!.getBoundingClientRect();
    return {
      x: (e.clientX - rect.left) / zoom,
      y: (e.clientY - rect.top) / zoom,
    };
  };
  const add = (x: number, y: number) => {
    if (tool === "select" || travel) return;
    const objectId = uid();
    b.send({
      kind: "create",
      objectId,
      object: {
        id: objectId,
        kind: tool,
        x: Math.max(0, x),
        y: Math.max(0, y),
        width: tool === "text" ? 300 : 200,
        height: tool === "text" ? 64 : 150,
        color:
          tool === "note"
            ? palette[0]
            : tool === "text"
              ? "#ffffff"
              : palette[1],
        text:
          tool === "text"
            ? "Your text"
            : tool === "note"
              ? "A fresh idea…"
              : "",
        deleted: false,
      },
    });
    setSelected(objectId);
    setTool("select");
  };
  const shareLink = () => {
    const key = sessionStorage.getItem("sf:key:" + id);
    return `${location.origin}/board/${id}${key ? "#key=" + key : ""}`;
  };
  const copy = async () => {
    try {
      const { token } = await api<{ token: string }>(
        `/api/boards/${id}/share`,
        post({}),
      );
      sessionStorage.setItem("sf:key:" + id, token);
      try {
        await navigator.clipboard.writeText(shareLink());
        notify("Invite link copied. Anyone with this link can edit.");
      } catch {
        setShare(true);
      }
    } catch (error) {
      b.setError(String(error));
    }
  };
  if (!b.state)
    return (
      <div className="loading-page">
        <Logo />
        {b.status === "Unavailable" ? (
          <>
            <WifiOff size={32} />
            <h1>We couldn’t open this canvas.</h1>
            <p role="alert">{b.error}</p>
            <p>Use the complete invite link to join a shared board.</p>
            <a href="/" className="primary">
              Back to workspace
            </a>
          </>
        ) : (
          <>
            <Loader2 size={28} className="spin" />
            <h1>Making room for your ideas…</h1>
            <span>Connecting to your canvas</span>
          </>
        )}
      </div>
    );
  return (
    <div className="board-app">
      <header className="board-header">
        <div className="board-header-left">
          <Logo />
          <span className="divider" />
          <a href="/" className="icon-button" aria-label="Back to workspace">
            <ChevronLeft size={17} />
          </a>
          <div className="board-title">
            <strong>{b.state.title}</strong>
            <span>
              Shared workspace <ChevronRight size={10} /> Canvas
            </span>
          </div>
        </div>
        <div className="board-header-right">
          <div
            className={`connection ${b.status === "Live" ? "" : "disconnected"}`}
            role="status"
          >
            <span className="live-dot" />
            {b.status}
            {b.pendingCount > 0 && ` · ${b.pendingCount} pending`}
          </div>
          <div
            className="avatars"
            title={b.peers.map((p) => p.name).join(", ")}
          >
            {b.peers.slice(0, 4).map((p, i) => (
              <span
                style={{
                  background: ["#e4dcfc", "#d7f0e5", "#fff2b3", "#ffdfe0"][
                    i % 4
                  ],
                }}
                key={p.id}
              >
                {p.name.slice(0, 1).toUpperCase()}
              </span>
            ))}
          </div>
          <button
            className="secondary history-button"
            onClick={() => setHistoryOpen(!historyOpen)}
            aria-expanded={historyOpen}
          >
            <History size={16} /> History
          </button>
          <button className="primary" onClick={() => void copy()}>
            <Share2 size={15} /> Share board
          </button>
        </div>
      </header>
      {(b.error || !b.redisOK || b.status !== "Live") && (
        <div className="notice" role="status">
          <WifiOff size={15} />
          {b.error ||
            (!b.redisOK
              ? "Presence is limited. Edits still save to MySQL and sync through recovery polling."
              : "Connection interrupted. Pending edits are kept on this device and resent with the same IDs.")}
          <button aria-label="Dismiss message" onClick={() => b.setError("")}>
            <X size={14} />
          </button>
        </div>
      )}
      <div className="board-body">
        <aside className="left-rail">
          <button
            className="rail-item active"
            title="Canvas"
            onClick={() => {
              setRevision(null);
              setPlaying(false);
              setSelected(null);
            }}
          >
            <Frame size={20} />
            <span>Canvas</span>
          </button>
          <button
            className="rail-item"
            onClick={() => setHistoryOpen(!historyOpen)}
            title="History"
          >
            <Clock3 size={20} />
            <span>History</span>
          </button>
          <div className="rail-bottom">
            <span className="own-avatar">{name[0]}</span>
            <span>You</span>
          </div>
        </aside>
        <main className="canvas-area">
          <div className="canvas-topline">
            <span>
              <Layers size={14} /> {travel ? "History preview" : "Live canvas"}{" "}
              <span className="slash">/</span>{" "}
              <strong>{objects.length} objects</strong>
            </span>
            <span className="revision-badge" data-testid="revision">
              {historyLoading ? "LOADING" : `REV ${shown?.revision ?? 0}`}
            </span>
          </div>
          {travel && (
            <div className="travel-banner">
              <History size={15} />
              <span>
                You're looking at revision {revision}. Live edits continue in
                the background.
              </span>
              <button
                onClick={() => {
                  setRevision(null);
                  setPlaying(false);
                }}
              >
                Return to live <ArrowRight size={13} />
              </button>
            </div>
          )}
          <div className={`canvas-scroll tool-${tool}`}>
            <svg
              ref={svg}
              className="canvas"
              width={1400 * zoom}
              height={950 * zoom}
              viewBox="0 0 1400 950"
              role="application"
              aria-label="Collaborative whiteboard"
              data-testid="canvas"
              onPointerDown={(e) => {
                if (
                  e.target === svg.current ||
                  ((e.target as Element).tagName === "rect" &&
                    (e.target as Element).getAttribute("data-bg") === "true")
                ) {
                  const p = point(e);
                  if (tool !== "select") add(p.x, p.y);
                  else setSelected(null);
                }
              }}
              onPointerMove={(e) => {
                const p = point(e);
                b.cursor(p.x, p.y);
                const d = drag.current;
                if (!d) return;
                const dx = p.x - d.x,
                  dy = p.y - d.y;
                setDraft(
                  d.resize
                    ? {
                        ...d.object,
                        width: Math.min(
                          2000,
                          Math.max(20, d.object.width + dx),
                        ),
                        height: Math.min(
                          2000,
                          Math.max(20, d.object.height + dy),
                        ),
                      }
                    : {
                        ...d.object,
                        x: Math.min(10000, Math.max(0, d.object.x + dx)),
                        y: Math.min(10000, Math.max(0, d.object.y + dy)),
                      },
                );
              }}
              onPointerUp={() => {
                const d = drag.current;
                if (d && draft) {
                  b.send({
                    kind: "update",
                    objectId: d.id,
                    patch: d.resize
                      ? { width: draft.width, height: draft.height }
                      : { x: draft.x, y: draft.y },
                  });
                }
                drag.current = null;
                setDraft(null);
              }}
              onPointerCancel={() => {
                drag.current = null;
                setDraft(null);
              }}
            >
              <defs>
                <pattern
                  id="dots"
                  width="24"
                  height="24"
                  patternUnits="userSpaceOnUse"
                >
                  <circle cx="1" cy="1" r="0.85" fill="#d8dce7" />
                </pattern>
                <filter
                  id="note-shadow"
                  x="-20%"
                  y="-20%"
                  width="140%"
                  height="150%"
                >
                  <feDropShadow
                    dx="0"
                    dy="3"
                    stdDeviation="4"
                    floodColor="#26304d"
                    floodOpacity=".07"
                  />
                </filter>
              </defs>
              <rect
                data-bg="true"
                width="1400"
                height="950"
                fill="url(#dots)"
              />
              {objects.map((original) => {
                const o = draft?.id === original.id ? draft : original;
                const active = selected === o.id;
                return (
                  <g
                    key={o.id}
                    data-testid="board-object"
                    data-object-id={o.id}
                    role="button"
                    tabIndex={0}
                    aria-label={`${o.kind}: ${o.text || o.id.slice(0, 8)}`}
                    onFocus={() => setSelected(o.id)}
                    onKeyDown={(e) => {
                      if (e.key === "Enter" || e.key === " ") {
                        e.preventDefault();
                        setSelected(o.id);
                        return;
                      }
                      if (
                        travel ||
                        ![
                          "ArrowLeft",
                          "ArrowRight",
                          "ArrowUp",
                          "ArrowDown",
                        ].includes(e.key)
                      )
                        return;
                      e.preventDefault();
                      e.stopPropagation();
                      const step = e.shiftKey ? 1 : 10;
                      b.send({
                        kind: "update",
                        objectId: o.id,
                        patch: {
                          x: Math.max(
                            0,
                            o.x +
                              (e.key === "ArrowRight"
                                ? step
                                : e.key === "ArrowLeft"
                                  ? -step
                                  : 0),
                          ),
                          y: Math.max(
                            0,
                            o.y +
                              (e.key === "ArrowDown"
                                ? step
                                : e.key === "ArrowUp"
                                  ? -step
                                  : 0),
                          ),
                        },
                      });
                    }}
                    transform={`translate(${o.x},${o.y})`}
                    className={travel ? "object-readonly" : "canvas-object"}
                    onPointerDown={(e) => {
                      e.stopPropagation();
                      setSelected(o.id);
                      if (travel) return;
                      if (tool !== "select") {
                        const p = point(e);
                        add(p.x, p.y);
                        return;
                      }
                      const p = point(e);
                      drag.current = {
                        id: o.id,
                        x: p.x,
                        y: p.y,
                        object: o,
                        resize: false,
                      };
                      svg.current!.setPointerCapture(e.pointerId);
                    }}
                  >
                    {o.kind === "ellipse" ? (
                      <ellipse
                        cx={o.width / 2}
                        cy={o.height / 2}
                        rx={o.width / 2}
                        ry={o.height / 2}
                        fill={o.color}
                        stroke={active ? "#7964e8" : "#c5cee1"}
                        strokeWidth={active ? 2 : 1}
                      />
                    ) : (
                      <rect
                        width={o.width}
                        height={o.height}
                        rx={o.kind === "note" ? 3 : 8}
                        fill={o.kind === "text" ? "transparent" : o.color}
                        stroke={
                          active
                            ? "#7964e8"
                            : o.kind === "note"
                              ? "#00000009"
                              : o.kind === "text"
                                ? "transparent"
                                : "#c5cee1"
                        }
                        strokeWidth={active ? 2 : 1}
                        filter={
                          o.kind === "note" ? "url(#note-shadow)" : undefined
                        }
                      />
                    )}
                    <foreignObject
                      x="14"
                      y={o.kind === "note" ? 18 : 10}
                      width={Math.max(1, o.width - 28)}
                      height={Math.max(1, o.height - 20)}
                      pointerEvents="none"
                    >
                      <div
                        className={`object-text ${o.kind}`}
                        style={{ fontSize: o.kind === "text" ? 25 : 15 }}
                      >
                        {o.text}
                      </div>
                    </foreignObject>
                    {active && !travel && (
                      <>
                        <rect
                          x="-5"
                          y="-5"
                          width={o.width + 10}
                          height={o.height + 10}
                          fill="none"
                          stroke="#7964e8"
                          strokeDasharray="4 4"
                          pointerEvents="none"
                        />
                        <rect
                          data-testid="resize-handle"
                          x={o.width - 6}
                          y={o.height - 6}
                          width="12"
                          height="12"
                          fill="white"
                          stroke="#7964e8"
                          className="resize-handle"
                          onPointerDown={(e) => {
                            e.stopPropagation();
                            const p = point(e);
                            drag.current = {
                              id: o.id,
                              x: p.x,
                              y: p.y,
                              object: o,
                              resize: true,
                            };
                            svg.current!.setPointerCapture(e.pointerId);
                          }}
                        />
                      </>
                    )}
                  </g>
                );
              })}
              {comparing &&
                travel &&
                Object.values(b.state.objects)
                  .filter((o) => !o.deleted)
                  .map((o) => (
                    <rect
                      key={"compare" + o.id}
                      x={o.x}
                      y={o.y}
                      width={o.width}
                      height={o.height}
                      fill="none"
                      stroke="#16a57a"
                      strokeDasharray="7 5"
                      strokeWidth="2"
                      opacity=".65"
                      pointerEvents="none"
                    />
                  ))}
              {!travel &&
                Object.entries(b.cursors).map(([key, c]) => (
                  <g
                    key={key}
                    transform={`translate(${c.x},${c.y})`}
                    pointerEvents="none"
                  >
                    <path
                      d="M0 0L4 18L9 12L16 10Z"
                      fill={c.peer.color}
                      stroke="white"
                      strokeWidth="1.5"
                    />
                    <rect
                      x="14"
                      y="12"
                      width={c.peer.name.length * 7 + 18}
                      height="24"
                      rx="5"
                      fill={c.peer.color}
                    />
                    <text x="22" y="28" fill="white" fontSize="11">
                      {c.peer.name}
                    </text>
                  </g>
                ))}
            </svg>
          </div>
          {historyLoading && (
            <div className="empty-canvas-hint" role="status">
              <Loader2 className="spin" size={24} />
              <strong>Loading revision {revision}…</strong>
              <span>The live canvas is unchanged.</span>
            </div>
          )}
          {!historyLoading && objects.length === 0 && (
            <div className="empty-canvas-hint">
              <strong>Every idea starts somewhere.</strong>
              <span>Choose a tool below, then click to begin.</span>
            </div>
          )}
          {object && !object.deleted && !travel && (
            <div
              className="inspector"
              style={
                {
                  "--inspector-left": inspectorOnLeft ? "22px" : "auto",
                  "--inspector-right": inspectorOnLeft ? "auto" : "22px",
                } as React.CSSProperties
              }
            >
              <div className="inspector-title">
                <strong>Object properties</strong>
                <button
                  className="icon-button"
                  aria-label="Close properties"
                  onClick={() => setSelected(null)}
                >
                  <X size={14} />
                </button>
              </div>
              <label>
                Text
                <textarea
                  aria-label="Object text"
                  value={textDraft}
                  maxLength={2000}
                  onChange={(e) => {
                    textDirty.current = true;
                    setTextDraft(e.target.value);
                  }}
                  onBlur={() => {
                    if (textDirty.current && textDraft !== object.text)
                      b.send({
                        kind: "update",
                        objectId: object.id,
                        patch: { text: textDraft },
                      });
                    textDirty.current = false;
                  }}
                />
              </label>
              <div className="color-row" role="group" aria-label="Object color">
                {palette.map((color) => (
                  <button
                    aria-label={`Color ${color}`}
                    title={color}
                    key={color}
                    className={object.color === color ? "chosen" : ""}
                    style={{ background: color }}
                    onClick={() =>
                      b.send({
                        kind: "update",
                        objectId: object.id,
                        patch: { color },
                      })
                    }
                  />
                ))}
              </div>
              <div className="dimensions">
                {(["x", "y", "width", "height"] as const).map((key) => (
                  <label key={key}>
                    {key}
                    <input
                      aria-label={`Object ${key}`}
                      key={`${object.id}:${key}`}
                      type="number"
                      step="any"
                      value={dimensionDrafts[key] ?? Math.round(object[key])}
                      onChange={(e) =>
                        setDimensionDrafts((old) => ({
                          ...old,
                          [key]: e.target.value,
                        }))
                      }
                      onBlur={(e) => {
                        const value = Number(e.target.value);
                        if (Number.isFinite(value) && value !== object[key])
                          b.send({
                            kind: "update",
                            objectId: object.id,
                            patch: { [key]: value },
                          });
                        setDimensionDrafts((old) => {
                          const next = { ...old };
                          delete next[key];
                          return next;
                        });
                      }}
                    />
                  </label>
                ))}
              </div>
              <button
                className="delete-object"
                onClick={() => {
                  b.send({ kind: "delete", objectId: object.id });
                  setSelected(null);
                }}
              >
                <Trash2 size={14} /> Delete object
              </button>
            </div>
          )}
          <div className="canvas-bottom">
            <div className="zoom">
              <button
                aria-label="Zoom out"
                onClick={() => setZoom((z) => Math.max(0.5, z - 0.1))}
              >
                <ZoomOut size={15} />
              </button>
              <span>{Math.round(zoom * 100)}%</span>
              <button
                aria-label="Zoom in"
                onClick={() => setZoom((z) => Math.min(1.5, z + 0.1))}
              >
                <ZoomIn size={15} />
              </button>
              <span className="mini-divider" />
              <button aria-label="Reset zoom" onClick={() => setZoom(1)}>
                <Maximize2 size={14} />
              </button>
            </div>
            <div className="toolbar" role="toolbar" aria-label="Drawing tools">
              {(
                [
                  {
                    id: "select",
                    icon: MousePointer2,
                    label: "Select",
                    key: "V",
                  },
                  { id: "rect", icon: Square, label: "Rectangle", key: "R" },
                  { id: "ellipse", icon: Circle, label: "Ellipse", key: "O" },
                  {
                    id: "note",
                    icon: StickyNote,
                    label: "Sticky note",
                    key: "N",
                  },
                  { id: "text", icon: Type, label: "Text", key: "T" },
                ] as const
              ).map((t) => (
                <button
                  aria-label={t.label}
                  title={`${t.label} (${t.key})`}
                  aria-pressed={tool === t.id}
                  disabled={travel}
                  className={tool === t.id ? "selected" : ""}
                  key={t.id}
                  onClick={() => setTool(t.id)}
                >
                  <t.icon size={21} />
                  <small>{t.key}</small>
                </button>
              ))}
            </div>
            <span className="canvas-help">
              {travel
                ? "Read-only history preview"
                : tool === "select"
                  ? "Select, drag, make it yours"
                  : "Click anywhere to create"}
            </span>
          </div>
        </main>
        {historyOpen && (
          <aside className="history-panel">
            <div className="history-title">
              <div>
                <Clock3 size={18} />
                <h2>Time machine</h2>
              </div>
              <button
                className="icon-button"
                aria-label="Close history"
                onClick={() => setHistoryOpen(false)}
              >
                <X size={15} />
              </button>
            </div>
            <p className="history-description">
              Every idea has a journey.
              <br />
              Revisit any moment of yours.
            </p>
            <div className="timeline-card">
              <div className="timeline-label">
                <span className={travel ? "past-indicator" : "live-indicator"}>
                  <span className="live-dot" />
                  {travel ? "PAST FRAME" : "LIVE FRAME"}
                </span>
                <span>
                  {revision ?? b.state.revision}{" "}
                  <span className="muted">/ {b.state.revision}</span>
                </span>
              </div>
              <input
                className="timeline-slider"
                aria-label="History revision"
                type="range"
                min="0"
                max={b.state.revision}
                value={revision ?? b.state.revision}
                onChange={(e) => {
                  setRevision(Number(e.target.value));
                  setPlaying(false);
                  setSelected(null);
                }}
              />
              <div className="timeline-endpoints">
                <span>Beginning</span>
                <span>Now</span>
              </div>
              <div className="playback-controls">
                <button aria-label="First event" onClick={() => setRevision(0)}>
                  <SkipBack size={15} />
                </button>
                <button
                  aria-label="Previous event"
                  onClick={() =>
                    setRevision((r) =>
                      Math.max(0, (r ?? b.state!.revision) - 1),
                    )
                  }
                >
                  <ChevronLeft size={18} />
                </button>
                <button
                  className="play-button"
                  aria-label={playing ? "Pause history" : "Play history"}
                  disabled={!b.state.revision}
                  onClick={() => {
                    if (
                      !playing &&
                      (revision === null || revision >= b.state!.revision)
                    )
                      setRevision(0);
                    setPlaying(!playing);
                  }}
                >
                  {playing ? (
                    <Pause size={17} />
                  ) : (
                    <Play size={17} fill="currentColor" />
                  )}
                </button>
                <button
                  aria-label="Next event"
                  onClick={() =>
                    setRevision((r) =>
                      Math.min(b.state!.revision, (r ?? 0) + 1),
                    )
                  }
                >
                  <ChevronRight size={18} />
                </button>
                <button
                  aria-label="Latest event"
                  onClick={() => setRevision(b.state!.revision)}
                >
                  <SkipForward size={15} />
                </button>
                <select
                  aria-label="Playback speed"
                  value={speed}
                  onChange={(e) => setSpeed(Number(e.target.value))}
                >
                  <option value=".5">0.5×</option>
                  <option value="1">1×</option>
                  <option value="2">2×</option>
                  <option value="4">4×</option>
                </select>
              </div>
            </div>
            {travel && (
              <div className="compare-controls">
                <label>
                  <input
                    type="checkbox"
                    checked={comparing}
                    onChange={(e) => setComparing(e.target.checked)}
                  />{" "}
                  Overlay live positions
                </label>
                <span>
                  Past: {objects.length} · Live: {liveCount} objects
                </span>
                <button
                  className="return-live"
                  onClick={() => {
                    setRevision(null);
                    setPlaying(false);
                  }}
                >
                  <RotateCcw size={13} /> Back to live canvas
                </button>
              </div>
            )}
            <div className="event-section-title">
              <span>ACTIVITY</span>
              <span>{b.events.length} events</span>
            </div>
            <div className="event-list">
              {[...b.events].reverse().map((e) => (
                <button
                  className={`event-row ${revision === e.revision ? "event-selected" : ""}`}
                  key={e.revision}
                  onClick={() => {
                    setRevision(e.revision);
                    setPlaying(false);
                    setSelected(null);
                  }}
                >
                  <span className={`event-symbol ${e.operation.kind}`}>
                    {e.operation.kind === "create" ? (
                      <Plus size={14} />
                    ) : e.operation.kind === "delete" ? (
                      <Trash2 size={13} />
                    ) : (
                      <MousePointer2 size={13} />
                    )}
                  </span>
                  <span className="event-content">
                    <strong>
                      {e.operation.actor}
                      <span>
                        {" "}
                        {e.operation.kind === "create"
                          ? "added"
                          : e.operation.kind === "delete"
                            ? "removed"
                            : "updated"}{" "}
                        {e.object.kind === "note"
                          ? "a note"
                          : e.object.kind === "text"
                            ? "text"
                            : "a shape"}
                      </span>
                    </strong>
                    <small>
                      {new Date(e.at).toLocaleTimeString([], {
                        hour: "2-digit",
                        minute: "2-digit",
                        second: "2-digit",
                      })}{" "}
                      <span>·</span> Revision {e.revision}
                    </small>
                  </span>
                  <ChevronRight size={12} />
                </button>
              ))}
              {b.events.length === 0 && (
                <div className="empty-history" role="status">
                  {b.state.revision > 0 ? (
                    <>
                      <Loader2 size={25} className="spin" />
                      <strong>Loading saved activity…</strong>
                    </>
                  ) : (
                    <>
                      <History size={25} />
                      <strong>The story starts here.</strong>
                      <p>
                        Your edits will appear here,
                        <br />
                        one frame at a time.
                      </p>
                    </>
                  )}
                </div>
              )}
            </div>
            <div className="history-footnote">
              <Check size={13} />
              <span>
                History is saved to MySQL.
                <br />
                Exploring the past keeps the live canvas intact.
              </span>
            </div>
          </aside>
        )}
      </div>
      <footer className="board-footer">
        <span>
          <span className="live-dot" />
          {travel ? "Exploring history" : "All committed edits saved"}{" "}
          <span className="slash">/</span>{" "}
          {b.pendingCount
            ? `${b.pendingCount} awaiting confirmation`
            : `Revision ${b.state.revision}`}
        </span>
        <span>
          <Users size={12} />
          {b.peers.length} online <span className="slash">/</span> You are{" "}
          {name}
        </span>
      </footer>
      {toast && (
        <div className="toast" role="status">
          <Check size={15} />
          {toast}
        </div>
      )}
      {share && (
        <div className="modal-backdrop" onClick={() => setShare(false)}>
          <div
            className="share-modal"
            role="dialog"
            aria-modal="true"
            aria-label="Share board"
            onClick={(e) => e.stopPropagation()}
          >
            <button
              className="modal-close icon-button"
              aria-label="Close share"
              onClick={() => setShare(false)}
            >
              <X size={17} />
            </button>
            <Share2 size={25} />
            <h2>Make room for someone else.</h2>
            <p>
              Anyone with this invite link can view and edit the board. Keep it
              within your team.
            </p>
            <input
              aria-label="Invite link"
              readOnly
              value={shareLink()}
              onFocus={(e) => e.target.select()}
            />
            <button className="primary" onClick={() => void copy()}>
              <Copy size={15} /> Copy invite link
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
const match = location.pathname.match(/^\/board\/([a-zA-Z0-9_-]+)$/);
createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    {match ? <Board id={match[1]} /> : <Home />}
  </React.StrictMode>,
);
