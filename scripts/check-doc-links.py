"""Check README/docs local link targets and GitHub source line anchors."""
import json
import argparse
from pathlib import Path
import re
import subprocess
from urllib.parse import unquote, urlsplit

root = Path(__file__).resolve().parent.parent
parser = argparse.ArgumentParser()
parser.add_argument("--output", default="artifacts/evidence/document-links.json")
args = parser.parse_args()
tracked = set(subprocess.check_output(["git", "ls-files", "-z"], cwd=root).decode().split("\0"))
checked = []
errors = []
external = set()
sources = [root / "README.md", *sorted((root / "docs").glob("*.md"))]
sources += [root / name for name in sorted(tracked) if name.startswith("artifacts/") and name.endswith(".md")]
for source in sources:
    for target in re.findall(r"\]\(([^)]+)\)", source.read_text(encoding="utf-8")):
        url = urlsplit(target)
        if url.scheme:
            external.add(target)
            continue
        path = (source.parent / unquote(url.path)).resolve()
        try:
            relative = path.relative_to(root).as_posix()
        except ValueError:
            relative = ""
        ok = path.exists() and bool(relative) and (relative in tracked if path.is_file() else any(p.startswith(relative + "/") for p in tracked))
        if ok and re.fullmatch(r"L[0-9]+", url.fragment):
            ok = path.is_file() and 1 <= int(url.fragment[1:]) <= len(path.read_text(encoding="utf-8").splitlines())
        elif ok and url.fragment and path.suffix == ".md":
            headings = re.findall(r"^#{1,6}\s+(.+)$", path.read_text(encoding="utf-8"), re.MULTILINE)
            anchors = {re.sub(r"[^\w -]", "", title.lower()).replace(" ", "-") for title in headings}
            ok = unquote(url.fragment) in anchors
        row = {"source": source.relative_to(root).as_posix(), "target": target, "exists": ok}
        checked.append(row)
        if not ok: errors.append(row)
result = {"status": "PASS" if not errors else "FAIL", "localLinks": len(checked), "errors": errors,
          "externalLinksNotCheckedByThisScript": sorted(external), "checked": checked}
output = root / args.output
output.parent.mkdir(parents=True, exist_ok=True)
output.write_text(json.dumps(result, ensure_ascii=False, indent=2), encoding="utf-8")
print(json.dumps({k: v for k, v in result.items() if k != "checked"}, ensure_ascii=False, indent=2))
raise SystemExit(bool(errors))
