# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: browser-qa.spec.ts >> mobile: new-user editing, zoomed drag, properties and history remain usable
- Location: e2e\browser-qa.spec.ts:33:2

# Error details

```
Test timeout of 90000ms exceeded.
```

# Page snapshot

```yaml
- generic [ref=f2e3]:
  - banner [ref=f2e4]:
    - generic [ref=f2e5]:
      - link "SYNCFRAME home" [ref=f2e6] [cursor=pointer]:
        - /url: /
      - link "Back to workspace" [ref=f2e9] [cursor=pointer]:
        - /url: /
      - generic [ref=f2e12]:
        - strong [ref=f2e13]: QA mobile workspace
        - generic [ref=f2e14]: Shared workspace Canvas
    - generic [ref=f2e17]:
      - status [ref=f2e18]: Live
      - button "History" [expanded] [ref=f2e20] [cursor=pointer]
      - button "Share board" [ref=f2e25] [cursor=pointer]
  - generic [ref=f2e32]:
    - main [ref=f2e33]:
      - generic [ref=f2e34]:
        - generic [ref=f2e35]:
          - text: Live canvas
          - generic [ref=f2e40]: /
        - generic [ref=f2e41]: REV 0
      - application "Collaborative whiteboard" [ref=f2e43]:
        - generic:
          - generic: Every idea starts somewhere.
          - generic: Choose a tool below, then click to begin.
          - generic: Your canvas. Your pace.
      - generic:
        - generic [ref=f2e45]:
          - button "Zoom out" [ref=f2e46] [cursor=pointer]
          - generic [ref=f2e50]: 100%
          - button "Zoom in" [ref=f2e51] [cursor=pointer]
          - button "Reset zoom" [ref=f2e56] [cursor=pointer]
        - toolbar "Drawing tools" [ref=f2e62]:
          - button "Select" [pressed] [ref=f2e63] [cursor=pointer]:
            - generic [ref=f2e66]: V
          - button "Rectangle" [ref=f2e67] [cursor=pointer]:
            - generic [ref=f2e70]: R
          - button "Ellipse" [ref=f2e71] [cursor=pointer]:
            - generic [ref=f2e74]: O
          - button "Sticky note" [ref=f2e75] [cursor=pointer]:
            - generic [ref=f2e79]: "N"
          - button "Text" [ref=f2e80] [cursor=pointer]:
            - generic [ref=f2e83]: T
    - complementary [ref=f2e84]:
      - generic [ref=f2e85]:
        - heading "Time machine" [level=2] [ref=f2e90]
        - button "Close history" [ref=f2e91] [cursor=pointer]
      - paragraph [ref=f2e95]: Every idea has a journey.Revisit any moment of yours.
      - generic [ref=f2e96]:
        - generic [ref=f2e97]:
          - generic [ref=f2e98]: LIVE FRAME
          - generic [ref=f2e100]:
            - text: "0"
            - generic [ref=f2e101]: / 0
        - slider "History revision" [ref=f2e102] [cursor=pointer]: "0"
        - generic [ref=f2e103]:
          - generic [ref=f2e104]: Beginning
          - generic [ref=f2e105]: Now
        - generic [ref=f2e106]:
          - button "First event" [ref=f2e107] [cursor=pointer]
          - button "Previous event" [ref=f2e110] [cursor=pointer]
          - button "Play history" [disabled] [ref=f2e113]
          - button "Next event" [ref=f2e116] [cursor=pointer]
          - button "Latest event" [ref=f2e119] [cursor=pointer]
          - combobox "Playback speed" [ref=f2e122]:
            - option "0.5×"
            - option "1×" [selected]
            - option "2×"
            - option "4×"
      - generic [ref=f2e123]:
        - generic [ref=f2e124]: ACTIVITY
        - generic [ref=f2e125]: 0 events
      - generic [ref=f2e127]:
        - strong [ref=f2e132]: The story starts here.
        - paragraph [ref=f2e133]: Your edits will appear here,one frame at a time.
      - generic [ref=f2e134]: History is saved to MySQL.Exploring the past keeps the live canvas intact.
  - contentinfo [ref=f2e138]:
    - generic [ref=f2e139]:
      - text: All committed edits saved
      - generic [ref=f2e141]: /
      - text: Revision 0
```