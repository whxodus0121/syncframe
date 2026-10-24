"""Check README/docs local link targets and GitHub source line anchors."""
import json
from pathlib import Path
import re
from urllib.parse import unquote, urlsplit

root = Path(__file__).resolve().parent.parent
checked = []
errors = []
external = set()
for source in [root / "README.md", *sorted((root / "docs").glob("*.md"))]:
    for target in re.findall(r"\]\(([^)]+)\)", source.read_text(encoding="utf-8")):
        url = urlsplit(target)
        if url.scheme:
            external.add(target)
            continue
        path = (source.parent / unquote(url.path)).resolve()
        ok = path.exists()
        if ok and re.fullmatch(r"L[0-9]+", url.fragment):
            ok = path.is_file() and 1 <= int(url.fragment[1:]) <= len(path.read_text(encoding="utf-8").splitlines())
        row = {"source": source.relative_to(root).as_posix(), "target": target, "exists": ok}
        checked.append(row)
        if not ok: errors.append(row)
result = {"status": "PASS" if not errors else "FAIL", "localLinks": len(checked), "errors": errors,
          "externalLinksNotCheckedByThisScript": sorted(external), "checked": checked}
output = root / "artifacts/portfolio/document-links.json"
output.write_text(json.dumps(result, ensure_ascii=False, indent=2), encoding="utf-8")
print(json.dumps({k: v for k, v in result.items() if k != "checked"}, ensure_ascii=False, indent=2))
raise SystemExit(bool(errors))
