/*
 * Roads - turn a `::: road STYLE` container into a scroll-drawn path.
 *
 * List styles: winding (default), timeline, steps, snake, cycle, constellation, git.
 *   Each list item is a link plus optional label text: `- [Title](url) Label`.
 *   In git, the label names the branch: `feat/x`, `fix/y from feat/x` (fork),
 *   `main <- feat/x` (merge).
 * Editorial: a heading per entry (`### [Title](url) Label`) followed by any
 *   markdown. A road whose first child is a heading always uses this style.
 *
 * Safe to load more than once and re-initialises after view transitions.
 */
(() => {
  if (window.MarkataRoads) {
    window.MarkataRoads.init();
    return;
  }

  const STYLES = ["winding", "timeline", "steps", "snake", "cycle", "constellation", "git", "editorial"];
  const SPACING = { winding: 210, timeline: 110, snake: 170, constellation: 150, git: 64 };
  const LANE = 28;
  const RAIL = 24;
  const BRANCH_COLORS = ["var(--road-accent)", "#5ec2b7", "#c678dd", "#e06c75", "#61afef", "#98c379"];
  const reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
  const svgNS = "http://www.w3.org/2000/svg";
  const roads = new Set();
  let clipCounter = 0;

  const clamp = (v, lo, hi) => Math.min(Math.max(v, lo), hi);
  const rand = (i) => { const s = Math.sin(i * 12.9898 + 78.233) * 43758.5453; return s - Math.floor(s); };

  const collect = (list, depth, out) => {
    for (const li of list.children) {
      if (li.tagName !== "LI") continue;
      const link = li.querySelector(":scope > a, :scope > p > a");
      const nested = li.querySelector(":scope > ul, :scope > ol");
      if (link) {
        const clone = li.cloneNode(true);
        clone.querySelectorAll("ul, ol, a").forEach((n) => n.remove());
        out.push({ link, depth, label: clone.textContent.replace(/^[\s\-–—:·|]+/, "").trim() });
      }
      if (nested) collect(nested, depth + 1, out);
    }
    return out;
  };

  const smooth = (pts) => {
    let d = `M ${pts[0].x} ${pts[0].y}`;
    for (let i = 0; i < pts.length - 1; i++) {
      const p0 = pts[Math.max(i - 1, 0)], p1 = pts[i], p2 = pts[i + 1], p3 = pts[Math.min(i + 2, pts.length - 1)];
      d += ` C ${p1.x + (p2.x - p0.x) / 6} ${p1.y + (p2.y - p0.y) / 6}, ${p2.x - (p3.x - p1.x) / 6} ${p2.y - (p3.y - p1.y) / 6}, ${p2.x} ${p2.y}`;
    }
    return d;
  };
  const polyline = (pts) => pts.map((p, i) => `${i ? "L" : "M"} ${p.x} ${p.y}`).join(" ");
  const sideOf = (x, w) => (x < w / 2 ? "right" : "left");
  const sideWidth = (x, w, side) => clamp((side === "right" ? w - x : x) - 34, 96, 220);

  // Each layout returns the height, a stop per point, and the path to draw.
  const LAYOUTS = {
    winding(n, w, sp) {
      const h = n * sp;
      const stops = Array.from({ length: n }, (_, i) => {
        const x = w * (0.5 + (w < 640 ? 0.13 : 0.22) * Math.sin(i * 0.95));
        const card = sideOf(x, w);
        return { x, y: h * (i + 0.5) / n, card, maxW: sideWidth(x, w, card) };
      });
      return { h, stops, d: smooth([{ x: w / 2, y: -h * 0.02 }, ...stops, { x: w / 2, y: h * 1.02 }]) };
    },
    timeline(n, w, sp) {
      const h = n * sp;
      const stops = Array.from({ length: n }, (_, i) => ({ x: RAIL, y: sp * (i + 0.5), card: "right", maxW: w - RAIL - 40 }));
      return { h, stops, d: `M ${RAIL} 0 L ${RAIL} ${h}` };
    },
    steps(n, w, sp) {
      const col = w / n;
      if (col < 110) return LAYOUTS.timeline(n, w, SPACING.timeline);
      const y = 24;
      const stops = Array.from({ length: n }, (_, i) => ({ x: col * (i + 0.5), y, card: "below", maxW: col - 12 }));
      return { h: 140, stops, d: `M ${col * 0.15} ${y} L ${w - col * 0.15} ${y}` };
    },
    snake(n, w, sp) {
      const perRow = clamp(Math.floor((w - 60) / 170), 2, 5);
      const pad = 34, col = (w - 2 * pad) / perRow;
      const rows = Math.ceil(n / perRow);
      const stops = Array.from({ length: n }, (_, i) => {
        const r = Math.floor(i / perRow), c = i % perRow;
        const x = pad + col * ((r % 2 ? perRow - 1 - c : c) + 0.5);
        return { x, y: sp * (r + 0.3), card: "below", maxW: col - 16 };
      });
      let d = `M ${stops[0].x - col / 2} ${stops[0].y} L ${stops[0].x} ${stops[0].y}`;
      for (let i = 1; i < n; i++) {
        const a = stops[i - 1], b = stops[i];
        if (a.y === b.y) { d += ` L ${b.x} ${b.y}`; continue; }
        const dir = Math.floor((i - 1) / perRow) % 2 ? -1 : 1;
        const k = Math.min(col * 0.5 + pad - 6, 90) * dir;
        d += ` C ${a.x + k} ${a.y}, ${b.x + k} ${b.y}, ${b.x} ${b.y}`;
      }
      return { h: rows * sp, stops, d };
    },
    cycle(n, w) {
      const R = clamp(Math.min(w / 2 - 140, 200), 70, 200);
      const cx = w / 2, cy = R + 70;
      const stops = Array.from({ length: n }, (_, i) => {
        const a = -Math.PI / 2 + (2 * Math.PI * i) / n;
        const x = cx + R * Math.cos(a), y = cy + R * Math.sin(a);
        let card = Math.abs(x - cx) < 4 ? (y < cy ? "above" : "below") : x < cx ? "left" : "right";
        return { x, y, card, maxW: card === "above" || card === "below" ? 200 : sideWidth(x, w, card === "left" ? "left" : "right") };
      });
      const d = `M ${cx} ${cy - R} A ${R} ${R} 0 1 1 ${cx} ${cy + R} A ${R} ${R} 0 1 1 ${cx} ${cy - R}`;
      return { h: 2 * R + 140, stops, d };
    },
    constellation(n, w, sp) {
      const stops = Array.from({ length: n }, (_, i) => {
        const x = w * (0.15 + 0.7 * rand(i + 1));
        const card = sideOf(x, w);
        return { x, y: sp * (i + 0.5) + (rand(i + 99) - 0.5) * sp * 0.4, card, maxW: sideWidth(x, w, card) };
      });
      return { h: n * sp, stops, d: polyline(stops) };
    },
    git(n, w, sp, r) {
      const { lanes, edges } = r.git;
      const maxLane = Math.max(...lanes.map((l) => l.lane));
      const stops = lanes.map((l, i) => {
        const x = RAIL + l.lane * LANE;
        const offset = (maxLane - l.lane) * LANE + 26;
        return { x, y: sp * (i + 0.5), card: "right", offset, maxW: w - x - offset - 8 };
      });
      const byBranch = new Map();
      const all = edges.map((e) => {
        const a = stops[e.from], b = stops[e.to];
        let seg;
        if (a.x === b.x) seg = `M ${a.x} ${a.y} L ${b.x} ${b.y}`;
        else if (e.kind === "merge") seg = `M ${a.x} ${a.y} L ${a.x} ${b.y - sp} C ${a.x} ${b.y - sp / 2}, ${b.x} ${b.y - sp / 2}, ${b.x} ${b.y}`;
        else seg = `M ${a.x} ${a.y} C ${a.x} ${a.y + sp / 2}, ${b.x} ${a.y + sp / 2}, ${b.x} ${a.y + sp} L ${b.x} ${b.y}`;
        byBranch.set(e.branch, (byBranch.get(e.branch) || "") + " " + seg);
        return seg;
      });
      return { h: n * sp, stops, d: all.join(" ") || `M ${RAIL} 0`, branches: byBranch };
    },
  };

  // Git: lanes per branch, forks from the parent's latest commit, merges join two lines.
  const parseGit = (items) => {
    const re = /^(\S+)?(?:\s+from\s+(\S+))?(?:\s*(?:<-|←)\s*(\S+))?/;
    const laneOf = new Map(), lastOn = new Map(), colorOf = new Map(), used = new Set(), edges = [];
    let trunk = null;
    const lanes = items.map((it, i) => {
      const m = (it.label || "").match(re) || [];
      const b = m[1] && m[1] !== "<-" && m[1] !== "←" ? m[1] : trunk || "main";
      trunk = trunk || b;
      if (!colorOf.has(b)) colorOf.set(b, colorOf.size);
      if (!laneOf.has(b)) {
        let lane = 0;
        while (used.has(lane)) lane++;
        laneOf.set(b, lane);
        used.add(lane);
        const parent = m[2] || trunk;
        if (parent !== b && lastOn.has(parent)) edges.push({ from: lastOn.get(parent), to: i, branch: b, kind: "fork" });
      } else {
        edges.push({ from: lastOn.get(b), to: i, branch: b, kind: "line" });
      }
      const merged = m[3];
      if (merged && merged !== b && lastOn.has(merged)) {
        edges.push({ from: lastOn.get(merged), to: i, branch: merged, kind: "merge" });
        if (laneOf.has(merged)) { used.delete(laneOf.get(merged)); laneOf.delete(merged); }
      }
      lastOn.set(b, i);
      return { lane: laneOf.get(b), branch: b, color: BRANCH_COLORS[colorOf.get(b) % BRANCH_COLORS.length] };
    });
    const colors = new Map([...colorOf].map(([b, c]) => [b, BRANCH_COLORS[c % BRANCH_COLORS.length]]));
    return { lanes, edges, colors };
  };

  const buildEditorial = (root) => {
    root.dataset.road = "editorial";
    const first = [...root.children].find((c) => /^H[1-6]$/.test(c.tagName));
    if (!first) return null;
    const level = first.tagName;
    const rail = document.createElement("div");
    rail.className = "road-rail";
    rail.setAttribute("aria-hidden", "true");
    const fill = document.createElement("div");
    fill.className = "road-fill";
    const marker = document.createElement("div");
    marker.className = "road-marker";
    rail.append(fill, marker);

    const entries = [];
    let entry = null;
    for (const node of [...root.childNodes]) {
      if (node.nodeType === 1 && node.tagName === level) {
        entry = document.createElement("section");
        entry.className = "road-entry";
        const dot = document.createElement("span");
        dot.className = "road-entry-dot";
        dot.setAttribute("aria-hidden", "true");
        const body = document.createElement("div");
        body.className = "road-entry-body";
        node.classList.add("road-entry-heading");
        // Text beside the title link becomes an eyebrow label above it.
        const link = [...node.children].find((c) => c.tagName === "A" && !c.classList.contains("heading-anchor"));
        if (link) {
          const rest = [...node.childNodes].filter((c) => c !== link && !(c.classList && c.classList.contains("heading-anchor")));
          const text = rest.map((c) => c.textContent).join("").trim();
          rest.forEach((c) => c.remove());
          if (text) {
            const label = document.createElement("span");
            label.className = "road-entry-label";
            label.textContent = text;
            node.prepend(label, " ");
          }
          link.classList.add("road-entry-title");
        }
        root.replaceChild(entry, node);
        entry.append(dot, node, body);
        entries.push({ el: entry, dot, body });
      } else if (entry) {
        entry.querySelector(".road-entry-body").appendChild(node);
      }
    }
    entries.forEach((e) => { if (!e.body.textContent.trim() && !e.body.querySelector("img, video, iframe")) e.body.remove(); });
    root.prepend(rail);
    return { root, style: "editorial", editorial: true, fill, marker, entries, length: 1 };
  };

  const build = (root) => {
    const lead = root.firstElementChild;
    if (root.dataset.road === "editorial" || root.classList.contains("editorial") || (lead && /^H[1-6]$/.test(lead.tagName))) {
      return buildEditorial(root);
    }
    const fromClass = [...root.classList].find((c) => STYLES.includes(c));
    const style = STYLES.includes(root.dataset.road) ? root.dataset.road : fromClass || "winding";
    root.dataset.road = style;
    const list = root.querySelector(":scope > ul, :scope > ol");
    if (!list) return null;
    const items = collect(list, 0, []);
    if (!items.length) return null;

    const inner = document.createElement("div");
    inner.className = "road-inner";
    const svg = document.createElementNS(svgNS, "svg");
    svg.setAttribute("class", "road-svg");
    svg.setAttribute("aria-hidden", "true");
    const mk = (cls, parent = svg) => { const p = document.createElementNS(svgNS, "path"); p.setAttribute("class", cls); parent.appendChild(p); return p; };
    const base = mk("road-path"), glow = mk("road-glow"), progress = mk("road-progress");
    const marker = document.createElement("div");
    marker.className = "road-marker";
    marker.setAttribute("aria-hidden", "true");
    inner.append(svg, marker);

    const r = { root, inner, svg, style, base, glow, progress, marker, stops: [], length: 0, h: 0 };
    if (style === "git") {
      r.git = parseGit(items);
      const id = `road-clip-${++clipCounter}`;
      const defs = document.createElementNS(svgNS, "defs");
      const clip = document.createElementNS(svgNS, "clipPath");
      clip.id = id;
      r.clipRect = document.createElementNS(svgNS, "rect");
      clip.appendChild(r.clipRect);
      defs.appendChild(clip);
      svg.prepend(defs);
      const group = document.createElementNS(svgNS, "g");
      group.setAttribute("clip-path", `url(#${id})`);
      svg.appendChild(group);
      r.branchPaths = new Map([...r.git.colors].map(([b, color]) => {
        const p = mk("road-branch", group);
        p.style.stroke = color;
        return [b, p];
      }));
    }

    r.points = items.map(({ link, label }, i) => {
      const title = link.textContent.trim();
      if (!link.matches("[data-title], [data-local-preview], [data-hover-title]")) {
        link.dataset.hoverTitle = title;
        if (label) link.dataset.hoverDescription = label;
      }
      link.classList.add("road-point");
      link.setAttribute("role", "listitem");
      link.textContent = "";
      if (r.git) {
        link.style.setProperty("--branch-color", r.git.lanes[i].color);
      }
      const dot = document.createElement("span");
      dot.className = "road-dot";
      const card = document.createElement("span");
      card.className = "road-card";
      if (label) {
        const num = document.createElement("span");
        num.className = "road-label";
        num.textContent = label;
        card.appendChild(num);
      }
      const t = document.createElement("span");
      t.className = "road-title";
      t.textContent = title;
      card.appendChild(t);
      link.append(dot, card);
      inner.appendChild(link);
      return { el: link };
    });

    list.remove();
    root.appendChild(inner);
    root.setAttribute("role", "list");
    return r;
  };

  const layout = (r) => {
    if (r.editorial) return update(r);
    const n = r.points.length;
    const sp = Number(r.root.dataset.roadSpacing) || SPACING[r.style] || 150;
    const w = r.inner.clientWidth;
    const out = LAYOUTS[r.style](n, w, sp, r);
    r.h = out.h;
    r.stops = out.stops;
    r.inner.style.height = `${out.h}px`;
    r.svg.setAttribute("viewBox", `0 0 ${w} ${out.h}`);

    r.points.forEach((p, i) => {
      const s = out.stops[i];
      p.el.style.left = `${s.x}px`;
      p.el.style.top = `${s.y}px`;
      for (const side of ["left", "right", "below", "above"]) p.el.classList.toggle(`road-card-${side}`, s.card === side);
      p.el.style.setProperty("--card-offset", `${s.offset || 26}px`);
      p.el.querySelector(".road-card").style.maxWidth = `${Math.max(80, s.maxW)}px`;
    });

    r.base.setAttribute("d", out.d);
    if (r.git) {
      r.clipRect.setAttribute("x", "-50");
      r.clipRect.setAttribute("y", "-50");
      r.clipRect.setAttribute("width", `${w + 100}`);
      for (const [b, path] of r.branchPaths) path.setAttribute("d", out.branches.get(b) || "");
      r.length = out.h;
    } else {
      r.glow.setAttribute("d", out.d);
      r.progress.setAttribute("d", out.d);
      r.length = r.progress.getTotalLength();
      for (const p of [r.glow, r.progress]) p.style.strokeDasharray = `${r.length}`;
      // Where along the path each stop sits, so "active" follows the path, not the page.
      const samples = Array.from({ length: 400 }, (_, k) => {
        const len = (r.length * k) / 399;
        const pt = r.progress.getPointAtLength(len);
        return { len, x: pt.x, y: pt.y };
      });
      r.stops.forEach((s) => {
        let best = samples[0];
        for (const q of samples) if (Math.hypot(q.x - s.x, q.y - s.y) < Math.hypot(best.x - s.x, best.y - s.y)) best = q;
        s.at = best.len;
      });
    }
    update(r);
  };

  const update = (r) => {
    if (!r.length) return;
    const rect = (r.inner || r.root).getBoundingClientRect(), vh = window.innerHeight;
    if (rect.bottom < -vh || rect.top > vh * 2) return;
    if (r.editorial) {
      // The fill tracks the reading line, so entries light up as you reach them.
      const rail = r.fill.parentElement.getBoundingClientRect();
      const reach = reduceMotion ? rail.height : clamp(vh * 0.55 - rail.top, 0, rail.height);
      r.fill.style.height = `${reach}px`;
      r.marker.style.transform = `translate(-50%, ${reach}px) translateY(-50%)`;
      r.marker.style.opacity = reach > 0 && reach < rail.height ? 1 : 0;
      r.entries.forEach((e) => {
        const d = e.dot.getBoundingClientRect();
        e.el.classList.toggle("active", d.top + d.height / 2 - rail.top <= reach + 1);
      });
      return;
    }
    const t = reduceMotion ? 1 : clamp((vh * 0.8 - rect.top) / (rect.height + vh * 0.3), 0, 1);
    if (r.git) {
      const reveal = r.h * t;
      r.clipRect.setAttribute("height", `${reveal + 50}`);
      r.points.forEach((p, i) => p.el.classList.toggle("active", r.stops[i].y <= reveal + 2));
      return;
    }
    const len = r.length * t;
    r.progress.style.strokeDashoffset = `${r.length - len}`;
    r.glow.style.strokeDashoffset = `${r.length - len}`;
    const pt = r.progress.getPointAtLength(len);
    r.marker.style.transform = `translate(${pt.x}px, ${pt.y}px) translate(-50%, -50%)`;
    r.points.forEach((p, i) => p.el.classList.toggle("active", r.stops[i].at <= len + 2));
  };

  const resizer = new ResizeObserver((entries) => {
    for (const r of roads) if (entries.some((e) => e.target === r.root)) layout(r);
  });

  const init = () => {
    for (const r of roads) {
      if (!r.root.isConnected) { resizer.unobserve(r.root); roads.delete(r); }
    }
    document.querySelectorAll(".road:not([data-road-ready]), [data-road]:not([data-road-ready])").forEach((root) => {
      root.setAttribute("data-road-ready", "");
      const r = build(root);
      if (!r) return;
      roads.add(r);
      resizer.observe(root);
      layout(r);
    });
  };

  let ticking = false;
  window.addEventListener("scroll", () => {
    if (ticking) return;
    ticking = true;
    requestAnimationFrame(() => { roads.forEach(update); ticking = false; });
  }, { passive: true });
  window.addEventListener("view-transition-complete", init);
  window.addEventListener("pageshow", init);

  window.MarkataRoads = { init };
  if (document.readyState === "loading") document.addEventListener("DOMContentLoaded", init);
  else init();
})();
