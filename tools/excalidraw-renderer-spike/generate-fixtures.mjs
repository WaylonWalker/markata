#!/usr/bin/env node

import { mkdir, writeFile } from "node:fs/promises";
import { fileURLToPath } from "node:url";
import path from "node:path";

const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "../..");
const outputDir = path.join(repoRoot, "pkg/plugins/testdata/excalidraw/renderer-corpus");
const pngDataURL = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAIAAAACCAIAAAD91JpzAAAAG0lEQVR4nAAOAPH/AuZGUDyW8AE8tGS+GeIDACt/Bk/MAciQAAAAAElFTkSuQmCC";
const elementBase = {
  angle: 0,
  backgroundColor: "transparent",
  fillStyle: "solid",
  groupIds: [],
  isDeleted: false,
  locked: false,
  opacity: 100,
  roughness: 1,
  seed: 11,
  strokeColor: "#1e1e1e",
  strokeStyle: "solid",
  strokeWidth: 2,
  updated: 1,
  version: 1,
  versionNonce: 11,
};

function makeElement(id, type, values = {}) {
  return { ...elementBase, id, type, ...values };
}

function scene(elements, files = {}, version = 2) {
  return {
    type: "excalidraw",
    version,
    source: "https://excalidraw.com",
    elements,
    appState: { viewBackgroundColor: "#ffffff", gridSize: null },
    files,
  };
}

const fixtures = new Map();

fixtures.set("shapes-rough.excalidraw", scene([
  makeElement("rect", "rectangle", { x: -20, y: 10, width: 180, height: 88, angle: 0.2, backgroundColor: "#a5d8ff", fillStyle: "hachure", roundness: { type: 3 }, boundElements: [{ id: "rect-label", type: "text" }] }),
  makeElement("rect-label", "text", { x: 15, y: 35, width: 115, height: 32, text: "Rough box", originalText: "Rough box", fontSize: 20, fontFamily: 5, textAlign: "center", verticalAlign: "middle", baseline: 18, lineHeight: 1.25, containerId: "rect", autoResize: false }),
  makeElement("ellipse", "ellipse", { x: 220, y: 0, width: 110, height: 75, strokeColor: "#1971c2", backgroundColor: "#d0bfff", fillStyle: "cross-hatch", strokeStyle: "dashed", roughness: 2, seed: 23 }),
  makeElement("diamond", "diamond", { x: 360, y: 4, width: 100, height: 80, strokeColor: "#2f9e44", backgroundColor: "#b2f2bb", fillStyle: "solid", seed: 37 }),
  makeElement("arrow", "arrow", { x: 10, y: 150, width: 320, height: 30, points: [[0, 0], [100, 45], [320, 30]], startArrowhead: "triangle", endArrowhead: "arrow", seed: 41 }),
  makeElement("line", "line", { x: 340, y: 150, width: 120, height: 50, points: [[0, 0], [40, 50], [120, 20]], strokeStyle: "dotted", seed: 53 }),
]));

fixtures.set("text-fonts.excalidraw", scene([
  ...[5, 6, 7, 8, 9].map((fontFamily, index) => makeElement(`font-${fontFamily}`, "text", {
    x: -10 + (index % 2) * 260,
    y: -20 + Math.floor(index / 2) * 85,
    width: 250,
    height: 70,
    angle: index === 2 ? -0.3 : 0,
    text: index === 0 ? "Excalifont\nmultiline" : `Font family ${fontFamily}`,
    originalText: index === 0 ? "Excalifont\nmultiline" : `Font family ${fontFamily}`,
    fontSize: 25,
    fontFamily,
    textAlign: "left",
    verticalAlign: "top",
    baseline: 24,
    lineHeight: 1.25,
  })),
]));

fixtures.set("freedraw-frames-groups.excalidraw", scene([
  makeElement("frame", "frame", { x: 0, y: 0, width: 320, height: 220, name: "Group frame", strokeColor: "#4263eb", roundness: { type: 3 }, seed: 67 }),
  makeElement("frame-rect", "rectangle", { x: 30, y: 48, width: 240, height: 115, frameId: "frame", groupIds: ["group-a"], backgroundColor: "#fff3bf", fillStyle: "solid", roundness: { type: 3 }, seed: 71 }),
  makeElement("frame-text", "text", { x: 60, y: 85, width: 180, height: 50, frameId: "frame", groupIds: ["group-a"], text: "Bound in a frame", originalText: "Bound in a frame", fontSize: 20, fontFamily: 6, textAlign: "center", verticalAlign: "middle", baseline: 18, lineHeight: 1.25, containerId: "frame-rect", autoResize: false }),
  makeElement("stroke", "freedraw", { x: 25, y: 180, width: 260, height: 58, frameId: "frame", points: [[0, 40], [35, 6], [70, 44], [110, 12], [160, 50], [220, 4], [260, 32]], pressures: [], simulatePressure: true, lastCommittedPoint: null, seed: 83, strokeColor: "#e03131", strokeWidth: 3 }),
]));

fixtures.set("embedded-raster-image.excalidraw", scene([
  makeElement("image", "image", { x: -10, y: 20, width: 320, height: 240, status: "saved", fileId: "fixture-png", scale: [1, 1], crop: null }),
  makeElement("image-frame", "rectangle", { x: -10, y: 20, width: 320, height: 240, strokeColor: "#1e1e1e", backgroundColor: "transparent", fillStyle: "solid", strokeWidth: 2, roughness: 0, seed: 89 }),
], {
  "fixture-png": { id: "fixture-png", mimeType: "image/png", dataURL: pngDataURL, created: 1, lastRetrieved: 1 },
}));
fixtures.set("external-image-source.excalidraw", scene([
  makeElement("remote-image", "image", { x: 0, y: 0, width: 320, height: 240, status: "saved", fileId: "remote-image", scale: [1, 1], crop: null }),
], {
  "remote-image": { id: "remote-image", mimeType: "image/png", dataURL: "https://example.invalid/image.png", created: 1, lastRetrieved: 1 },
}));

const largeElements = [];
for (let i = 0; i < 360; i++) {
  const col = i % 24;
  const row = Math.floor(i / 24);
  const x = -50000 + col * 4200;
  const y = 100000 + row * 3100;
  largeElements.push(makeElement(`large-${String(i).padStart(3, "0")}`, i % 4 === 0 ? "ellipse" : "rectangle", {
    x,
    y,
    width: 700 + (i % 7) * 23,
    height: 420 + (i % 5) * 19,
    angle: (i % 11) / 37,
    strokeColor: ["#1e1e1e", "#1971c2", "#e03131", "#2f9e44"][i % 4],
    backgroundColor: i % 3 === 0 ? "#d0bfff" : "transparent",
    fillStyle: i % 3 === 0 ? "hachure" : "solid",
    seed: 1000 + i,
  }));
}
fixtures.set("large-unusual-coordinates.excalidraw", scene(largeElements));
fixtures.set("newer-scene-version.excalidraw", scene([
  makeElement("future", "rectangle", { x: 0, y: 0, width: 120, height: 80 }),
], {}, 3));
fixtures.set("unsupported-element.excalidraw", scene([
  makeElement("future-type", "magic-portal", { x: 0, y: 0, width: 120, height: 80 }),
]));

await mkdir(outputDir, { recursive: true });
for (const [name, value] of fixtures) {
  await writeFile(path.join(outputDir, name), `${JSON.stringify(value, null, 2)}\n`);
}
await writeFile(path.join(outputDir, "malformed.json"), '{"type":"excalidraw","version":2,"elements":[}\n');
console.log(`Wrote ${fixtures.size + 1} fixtures to ${outputDir}`);
