"""Create a compact slideshow from assertion-checked real Chromium screenshots.
This GIF is a sequence of captured states, not a continuous motion recording.
The unedited browser video is preserved separately as syncframe-demo.webm.
Requires: python -m pip install Pillow
"""
from pathlib import Path
from PIL import Image

folder = Path(__file__).resolve().parent.parent / "artifacts/portfolio/demo"
frames = []
for file in sorted(folder.glob("[0-9][0-9]-*.png")):
    image = Image.open(file).convert("RGB")
    image.thumbnail((960, 640))
    frames.append(image.quantize(colors=128))
if len(frames) != 8:
    raise SystemExit("Expected all eight assertion-checked capture stages")
frames[0].save(folder / "syncframe-demo.gif", save_all=True, append_images=frames[1:], duration=1600, loop=0, optimize=True)
print(f"Saved {len(frames)} actual browser states; 12.8s slideshow")
