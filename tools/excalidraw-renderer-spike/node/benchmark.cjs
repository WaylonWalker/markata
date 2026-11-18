const crypto = require("node:crypto");
const fs = require("node:fs");
const path = require("node:path");
const { performance } = require("node:perf_hooks");
const render = require("@aldinokemal2104/excalidraw-to-svg");

const fixtureDir = path.resolve(__dirname, "../../../pkg/plugins/testdata/excalidraw/renderer-corpus");
const outputDir = process.env.EXCALIDRAW_SPIKE_OUTPUT || "/tmp/markata-excalidraw-spike/node";
const runsPerFixture = Number(process.env.EXCALIDRAW_SPIKE_RUNS || 6);
const onlyFixture = process.env.EXCALIDRAW_SPIKE_FIXTURE;

async function main() {
  fs.mkdirSync(outputDir, { recursive: true });
  const results = [];
  for (const name of fs.readdirSync(fixtureDir).sort()) {
    if (onlyFixture && name !== onlyFixture) continue;
    const file = path.join(fixtureDir, name);
    if (!name.endsWith(".excalidraw") && name !== "malformed.json") continue;
    const source = fs.readFileSync(file, "utf8");
    const timings = [];
    const digests = [];
    let svg = "";
    let error = "";
    try {
      const scene = JSON.parse(source);
      scene.appState = { ...scene.appState, exportWithDarkMode: true };
      for (let index = 0; index < runsPerFixture; index++) {
        const start = performance.now();
        const result = await render(scene);
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
    renderer: "@aldinokemal2104/excalidraw-to-svg",
    version: require("@aldinokemal2104/excalidraw-to-svg/package.json").version,
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
