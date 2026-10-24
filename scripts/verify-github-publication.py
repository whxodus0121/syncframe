"""Verify anonymous public access and exact Git tree/asset integrity after a push.
Usage: python scripts/verify-github-publication.py --repo owner/name
No GitHub credentials are read or sent by this verifier.
"""
import argparse
from concurrent.futures import ThreadPoolExecutor
from datetime import datetime, timezone
import hashlib
import json
from pathlib import Path
import re
import subprocess
import time
from urllib.parse import quote
from urllib.request import Request, urlopen

root = Path(__file__).resolve().parent.parent
parser = argparse.ArgumentParser()
parser.add_argument("--repo", required=True)
parser.add_argument("--output", default="artifacts/evidence/remote-verification.json")
args = parser.parse_args()
assert re.fullmatch(r"[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+", args.repo)
def git(*args):
    return subprocess.check_output(["git", *args], cwd=root)
def get(url):
    for attempt in range(3):
        try:
            with urlopen(Request(url, headers={"User-Agent": "syncframe-publication-verification"}), timeout=25) as response:
                return response.status, response.headers.get("Content-Type"), response.read()
        except Exception:
            if attempt == 2:
                raise
            time.sleep(1 + attempt)

base = "https://api.github.com/repos/" + args.repo
metadata = json.loads(get(base)[2])
assert metadata["private"] is False
assert metadata["owner"]["type"] == "User"
assert metadata["default_branch"] == "main"
assert metadata["description"] == "Real-time collaborative whiteboard with Go, WebSocket, MySQL, Redis, Kafka, and revision-based history replay."
expected_topics = {"golang", "websocket", "distributed-systems", "real-time", "mysql", "redis", "kafka", "collaboration", "event-sourcing"}
assert set(metadata["topics"]) == expected_topics
head = git("rev-parse", "HEAD").decode().strip()
remote = json.loads(get(base + "/commits/main")[2])
assert head == remote["sha"]
tree = json.loads(get(base + "/git/trees/" + remote["commit"]["tree"]["sha"] + "?recursive=1")[2])
assert not tree["truncated"]
assert tree["sha"] == git("rev-parse", "HEAD^{tree}").decode().strip()
local_files = {}
for row in git("ls-tree", "-r", "HEAD").decode().splitlines():
    meta, name = row.split("\t", 1)
    mode, kind, sha = meta.split()
    local_files[name] = {"mode": mode, "sha": sha}
remote_files = {item["path"]: {"mode": item["mode"], "sha": item["sha"]} for item in tree["tree"] if item["type"] == "blob"}
assert local_files == remote_files
assert not any(path.startswith("artifacts/private/") for path in remote_files)
assert not {"AGENTS.md", "PLAN.md", "PROGRESS.md"}.intersection(remote_files)

paths = ["README.md", *[p.relative_to(root).as_posix() for p in sorted((root / "docs").glob("*.md"))],
         "artifacts/portfolio/demo/syncframe-demo.gif", "artifacts/portfolio/demo/syncframe-demo.webm",
         "artifacts/evidence/screenshots/mobile.png"]
def verify(path):
    url = "https://raw.githubusercontent.com/" + args.repo + "/" + head + "/" + quote(path, safe="/")
    status, content_type, content = get(url)
    expected = git("show", "HEAD:" + path)
    assert status == 200 and content == expected, path
    return {"path": path, "url": url, "status": status, "contentType": content_type, "bytes": len(content),
            "sha256": hashlib.sha256(content).hexdigest(), "matchesGitBlob": True}
with ThreadPoolExecutor(max_workers=4) as workers:
    assets = list(workers.map(verify, paths))
page_status = get(metadata["html_url"])[0]
assert page_status == 200
report = {"at": datetime.now(timezone.utc).isoformat(), "status": "PASS", "repository": metadata["html_url"],
          "public": True, "defaultBranch": "main", "verifiedCommit": head, "verifiedTree": tree["sha"],
          "description": metadata["description"], "topics": sorted(metadata["topics"]), "anonymousHTTP": True,
          "repositoryHTTPStatus": page_status, "allGitBlobsAndModesMatch": True, "fileCount": len(local_files),
          "privateBackupNotUploaded": True, "anonymousAssets": assets,
          "scope": "exact immutable commit at check time; later evidence-only commits are checked again separately"}
output = root / args.output
output.parent.mkdir(parents=True, exist_ok=True)
output.write_text(json.dumps(report, ensure_ascii=False, indent=2), encoding="utf-8")
print(json.dumps({key: report[key] for key in ["status", "repository", "verifiedCommit", "fileCount", "anonymousHTTP", "allGitBlobsAndModesMatch"]}, indent=2))
