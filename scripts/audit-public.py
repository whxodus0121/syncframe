"""Scan publishable Git blobs, including historical trace ZIP contents.
Only locations/categories are reported; matched credentials are never printed.
Known fixed loopback demo passwords are documented separately, not hidden.
"""
import argparse
import io
import json
from pathlib import Path
import re
import subprocess
import zipfile
from datetime import datetime, timezone

root = Path(__file__).resolve().parent.parent
parser = argparse.ArgumentParser()
parser.add_argument("--ref", help="inspect only history reachable from this publication ref; default is every local ref")
parser.add_argument("--working-only", action="store_true", help="preflight before creating a clean publication root")
parser.add_argument("--output", default="artifacts/portfolio/public-audit.json")
args = parser.parse_args()
def git(*args):
    return subprocess.check_output(["git", *args], cwd=root)

patterns = {
    "private-key": re.compile(rb"-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----"),
    "github-token": re.compile(rb"(?:gh[pousr]_[A-Za-z0-9]{30,}|github_pat_[A-Za-z0-9_]{40,})"),
    "cloud-access-key": re.compile(rb"(?:AKIA|ASIA)[A-Z0-9]{16}"),
    "capability-in-link": re.compile(rb"(?:#key=|%23key%3[Dd]|key%3[Dd])[a-fA-F0-9]{64}"),
    "capability-in-json": re.compile(rb'"token"\s*:\s*"[a-fA-F0-9]{64}"'),
    "capability-in-cookie": re.compile(rb"sf_[a-zA-Z0-9_-]+=[a-fA-F0-9]{64}"),
    "bearer-token": re.compile(rb"Bearer\s+[A-Za-z0-9_-]{20,}(?:\.[A-Za-z0-9_-]+){2}"),
    "windows-user-path": re.compile(rb"[A-Za-z]:[\\/]+Users[\\/]+[^\\/\s\"'<>:]+"),
    "unix-user-path": re.compile(rb"/(?:home|Users)/[A-Za-z0-9_.-]+"),
    "local-profile-name": re.compile(re.escape(Path.home().name.encode("utf-8"))),
}
findings = []
zip_entries = 0
def scan(data, location):
    global zip_entries
    for kind, regex in patterns.items():
        if regex.search(data): findings.append({"kind": kind, "location": location})
    for match in re.finditer(rb"([A-Za-z0-9_.-]+):([^\s\"'<>@]+)@tcp\(", data):
        if match.group(1, 2) not in {(b"syncframe", b"syncframe"), (b"root", b"local-root-only")}:
            findings.append({"kind": "non-demo-database-credential", "location": location})
    if data.startswith(b"PK\x03\x04"):
        with zipfile.ZipFile(io.BytesIO(data)) as z:
            for entry in z.infolist():
                zip_entries += 1
                scan(z.read(entry), location + "!" + entry.filename)

# Rev-list walks all local refs, so an earlier credential cannot hide behind a clean HEAD.
objects = [] if args.working_only else git("rev-list", "--objects", *([args.ref] if args.ref else ["--all"])).decode().splitlines()
blobs = 0
for line in objects:
    oid, _, name = line.partition(" ")
    if git("cat-file", "-t", oid).strip() != b"blob": continue
    blobs += 1
    scan(git("cat-file", "blob", oid), f"git:{oid[:12]}:{name}")
files = git("ls-files", "--cached", "--others", "--exclude-standard", "-z").decode().split("\0")
working = 0
for name in files:
    if not name: continue
    file = root / name
    if not file.is_file() or name == args.output: continue
    working += 1
    scan(file.read_bytes(), "worktree:" + name)
report = {"at": datetime.now(timezone.utc).isoformat(), "historicalBlobs": blobs,
          "workingFiles": working, "zipEntries": zip_entries, "findings": findings,
          "status": "PASS" if not findings else "FAIL",
          "scope": ("working preflight only" if args.working_only else "history reachable from " + args.ref if args.ref else "all local Git refs") + "; tracked/unignored working files and ZIP entries; pattern-based, not a proof of absence",
          "knownPublicDemoCredentials": ["syncframe:syncframe", "root:local-root-only"],
          "notVerified": ["remote refs not present locally", "third-party hosted copies", "video OCR", "arbitrary proprietary token formats"]}
output = root / args.output
output.parent.mkdir(parents=True, exist_ok=True)
output.write_text(json.dumps(report, ensure_ascii=False, indent=2), encoding="utf-8")
print(json.dumps(report, ensure_ascii=False, indent=2))
raise SystemExit(bool(findings))
