import crypto from "node:crypto";
import fs from "node:fs";
import path from "node:path";
import { createRequire } from "node:module";
import { fileURLToPath } from "node:url";
import { performance } from "node:perf_hooks";
import { JSDOM } from "jsdom";
import { exportScene } from "@excalidraw/utils";

const here = path.dirname(fileURLToPath(import.meta.url));
const require = createRequire(import.meta.url);
const fixtureDir = path.resolve(here, "../../../pkg/plugins/testdata/excalidraw/renderer-corpus");
const outputDir = process.env.EXCALIDRAW_SPIKE_OUTPUT || "/tmp/markata-excalidraw-spike/node-official";
const runsPerFixture = Number(process.env.EXCALIDRAW_SPIKE_RUNS || 6);
const onlyFixture = process.env.EXCALIDRAW_SPIKE_FIXTURE;

const utilsDir = path.dirname(require.resolve("@excalidraw/utils"));
// The published ESM bundle uses __dirname to find its packaged fonts. Supply
// its own bundle directory in Node, then keep font reads local to that package.
globalThis.__dirname = utilsDir;
globalThis.FontFace = class {
  constructor(family, source, descriptors = {}) {
    this.family = family;
    this.source = source;
    this.unicodeRange = descriptors.unicodeRange || "U+0-10FFFF";
    this.weight = descriptors.weight || "normal";
    this.style = descriptors.style || "normal";
    this.display = descriptors.display || "auto";
    this.status = "loaded";
  }
  async load() { return this; }
};

const dom = new JSDOM("<!doctype html><html><body></body></html>");
Object.defineProperty(dom.window.document, "fonts", {
  value: { add() {}, check() { return true; } },
});
const environment = {
  document: dom.window.document,
  createPath: (pathData) => ({ pathData }),
  textMetrics: {
    getLineWidth(text, font) {
      const size = Number.parseFloat(font.match(/([0-9.]+)px/)?.[1] || "16");
      return [...text].reduce((width, char) => width + size * (char === " " ? 0.3 : 0.55), 0);
    },
  },
  loadFonts: async () => {},
};

async function main() {
  fs.mkdirSync(outputDir, { recursive: true });
  const results = [];
  for (const name of fs.readdirSync(fixtureDir).sort()) {
    if (onlyFixture && name !== onlyFixture) continue;
    if (!name.endsWith(".excalidraw") && name !== "malformed.json") continue;
    const source = fs.readFileSync(path.join(fixtureDir, name), "utf8");
    const timings = [];
    const digests = [];
    let svg = "";
    let error = "";
    try {
      const scene = JSON.parse(source);
      const data = { elements: scene.elements, appState: scene.appState, files: scene.files };
      for (let index = 0; index < runsPerFixture; index++) {
        const start = performance.now();
        const result = await exportScene({
          type: "svg",
          data,
          config: {
            theme: "dark",
            padding: 20,
            environment,
            skipInliningFonts: process.env.EXCALIDRAW_SPIKE_SKIP_INLINE === "1",
          },
        });
        svg = result.outerHTML;
        timings.push(performance.now() - start);
        digests.push(crypto.createHash("sha256").update(svg).digest("hex"));
      }
      fs.writeFileSync(path.join(outputDir, `${name}.svg`), svg);
    } catch (caught) {
      error = caught instanceof Error ? `${caught.name}: ${caught.message}` : String(caught);
    }
    results.push({
      fixture: name,
      first_ms: timings[0] === undefined ? null : +timings[0].toFixed(3),
      warm_median_ms: timings.length < 2 ? null : +median(timings.slice(1)).toFixed(3),
      svg_bytes: svg ? Buffer.byteLength(svg) : null,
      deterministic: digests.length > 0 && new Set(digests).size === 1,
      external_references: /(?:href|src)=["'](?:https?:)?\/\//i.test(svg) || /url\(["']?(?:https?:)?\/\//i.test(svg),
      error: error || undefined,
    });
  }
  const report = `${JSON.stringify({
    renderer: "@excalidraw/utils",
    version: JSON.parse(fs.readFileSync(path.resolve(utilsDir, "../..", "package.json"), "utf8")).version,
    node: process.version,
    process_max_rss_kb: process.resourceUsage().maxRSS,
    results,
  }, null, 2)}\n`;
  await new Promise((resolve) => process.stdout.write(report, resolve));
  process.exit(0);
}

function median(values) {
  const sorted = [...values].sort((a, b) => a - b);
  const middle = Math.floor(sorted.length / 2);
  return sorted.length % 2 ? sorted[middle] : (sorted[middle - 1] + sorted[middle]) / 2;
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
