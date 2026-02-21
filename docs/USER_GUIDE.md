# User Guide

## First run
1. Start engine and desktop.
2. Use **Add Root** on welcome/navigator pane.
3. Watch indexing status in viewer pane.

## Search & navigate
- Press **Ctrl/Cmd+K** to focus search.
- Select a file to load metadata and links.

## Linking
- Choose source and target node IDs in Links panel.
- Click **Link**.
- Backlinks appear in IN list.

## Graph view
- Click **Open in Graph**.
- Pan/zoom, click node to select, double-click to refocus.
- Use depth/direction/relation/ext filters.
- If partial graph banner appears, use **Load more** or **Narrow filters**.

## Themes
- Switch theme from Graph toolbar selector: Light/Dark/Neon/Sakura.
- Theme applies instantly and persists locally.

## Export / import / backup
- UI Backup button exports graph backup JSON in config folder.
- CLI export/import:
  - `nodeveil-engine --graph-export graph.json`
  - `nodeveil-engine --graph-import graph.json`
