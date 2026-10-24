"""Independently compare measured boards' SQL rows and full event replay.
Run after load and service-fault tests have finished. Never deletes data.
"""
import json
import argparse
from pathlib import Path
import re
import subprocess
from datetime import datetime, timezone

root = Path(__file__).resolve().parent.parent
folder = root / "artifacts/portfolio"
def query(sql):
    return subprocess.check_output(["docker", "compose", "-p", "syncframe", "exec", "-T", "mysql",
        "mysql", "-N", "-B", "--raw", "-usyncframe", "-psyncframe", "-D", "syncframe", "-e", sql],
        cwd=root, stderr=subprocess.DEVNULL).decode("utf-8").strip()

parser = argparse.ArgumentParser()
parser.add_argument("reports", nargs="*", type=Path, help="optional report paths relative to the repository")
args = parser.parse_args()
files = [root / p for p in args.reports] if args.reports else sorted(folder.glob("burst-final-*.json")) + [folder / "soak-20-5m.json"]
results = []
for file in files:
    report = json.loads(file.read_text(encoding="utf-8-sig"))
    board = report["boardId"]
    assert re.fullmatch("[a-f0-9]{32}", board)
    state = json.loads(query(f"SELECT state FROM boards WHERE id='{board}'"))
    events = [json.loads(line) for line in query(f"SELECT payload FROM events WHERE board_id='{board}' ORDER BY revision").splitlines()]
    replay = dict(state, revision=0, objects={})
    identities = set()
    for event in events:
        assert event["revision"] == replay["revision"] + 1
        assert event["operation"]["id"] not in identities
        identities.add(event["operation"]["id"])
        replay["objects"][event["object"]["id"]] = event["object"]
        replay["revision"] = event["revision"]
    counts = query(f"SELECT (SELECT COUNT(*) FROM outbox WHERE board_id='{board}' AND published_at IS NULL), (SELECT COUNT(*) FROM audit_events WHERE board_id='{board}'), (SELECT COUNT(*) FROM outbox WHERE board_id='{board}')").split("\t")
    pending, audit, outbox = map(int, counts)
    assert state == replay
    assert state["revision"] == report["succeeded"] == len(events)
    assert len(state["objects"]) == report["users"]
    assert pending == 0 and audit == outbox == len(events)
    results.append({"report": file.name, "boardId": board, "sqlRevision": state["revision"], "events": len(events), "uniqueOperations": len(identities), "objects": len(state["objects"]), "sqlEqualsReplay": state == replay, "outbox": outbox, "pending": pending, "audit": audit, "status": "PASS"})
output = {"at": datetime.now(timezone.utc).isoformat(), "method": "direct MySQL state and events; full ordered replay; unique operation IDs; outbox/audit counts; no writes", "results": results}
(folder / "sql-verification.json").write_text(json.dumps(output, indent=2), encoding="utf-8")
print(json.dumps(output, indent=2))
