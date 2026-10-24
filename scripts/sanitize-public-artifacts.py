"""Redact local identity metadata in public test evidence, preserving private originals.
Does not change screenshots, event/state values, timings or test assertions.
Run only after test processes have finished writing their output.
"""
from datetime import datetime, timezone
import io
import json
from pathlib import Path
import re
import shutil
import subprocess
import zipfile

root = Path(__file__).resolve().parent.parent
backup = root / "artifacts/private/publication/originals"
workspace = str(root)
variants = {workspace, workspace.replace("\\", "/"), str(Path.home()), str(Path.home()).replace("\\", "/")}
for _ in range(3):
    variants |= {json.dumps(p, ensure_ascii=ascii_only)[1:-1] for p in list(variants) for ascii_only in (True, False)}
private_email = subprocess.check_output(["git", "config", "user.email"], cwd=root).decode().strip()
patterns = [
    (re.compile(r"[A-Za-z]:[\\/]+Users[\\/]+[^\\/\s\"'<>:]+"), "[USER_HOME]"),
    (re.compile(r"/(?:home|Users)/[A-Za-z0-9_.-]+"), "[USER_HOME]"),
]
def safe(data):
    try:
        text = data.decode("utf-8")
    except UnicodeDecodeError:
        return data
    for value in sorted(variants, key=len, reverse=True):
        text = text.replace(value, "[WORKSPACE]")
    for pattern, replacement in patterns:
        text = pattern.sub(replacement, text)
    if private_email and not private_email.endswith("@users.noreply.github.com"):
        text = text.replace(private_email, "[PRIVATE_EMAIL]")
    return text.encode("utf-8")

changed = []
names = subprocess.check_output(["git", "ls-files", "--cached", "--others", "--exclude-standard", "-z"], cwd=root).decode().split("\0")
for name in names:
    if not name.startswith("artifacts/") or name.startswith("artifacts/private/"):
        continue
    file = root / name
    if not file.is_file() or file.suffix.lower() in {".png", ".gif", ".jpg", ".webm"}:
        continue
    data = file.read_bytes()
    entries = []
    if zipfile.is_zipfile(file):
        output = io.BytesIO()
        with zipfile.ZipFile(io.BytesIO(data)) as source, zipfile.ZipFile(output, "w") as target:
            target.comment = source.comment
            for entry in source.infolist():
                original = source.read(entry)
                redacted = safe(original)
                if redacted != original:
                    entries.append(entry.filename)
                target.writestr(entry, redacted)
        cleaned = output.getvalue() if entries else data
    else:
        cleaned = safe(data)
    if cleaned == data:
        continue
    original = backup / name
    original.parent.mkdir(parents=True, exist_ok=True)
    if not original.exists():
        shutil.copy2(file, original)
    file.write_bytes(cleaned)
    changed.append({"path": name, "zipEntries": entries})

report = {"at": datetime.now(timezone.utc).isoformat(), "changedFiles": len(changed),
          "redactions": ["workspace paths", "user home paths", "configured private email"],
          "preserved": ["private originals outside publication", "screenshots/video", "event/state values", "assertions", "timings"],
          "files": changed}
path = root / "artifacts/publication/privacy-redactions.json"
path.parent.mkdir(parents=True, exist_ok=True)
path.write_text(json.dumps(report, ensure_ascii=False, indent=2), encoding="utf-8")
print(json.dumps({"changedFiles": len(changed), "report": path.relative_to(root).as_posix()}))
