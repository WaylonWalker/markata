import React from 'react';
import {createRoot} from 'react-dom/client';
import {Excalidraw, exportToSvg, serializeAsJSON} from '@excalidraw/excalidraw';

export function mountExcalidraw(host, initialScene, onChange) {
  const root = createRoot(host);
  let latest = JSON.stringify(initialScene, null, 2);

  root.render(React.createElement(Excalidraw, {
    initialData: initialScene,
    UIOptions: {canvasActions:{export:{saveFileToDisk:true}}},
    handleKeyboardGlobally: false,
    onChange(elements, appState, files) {
      latest = serializeAsJSON(elements, appState, files, 'local');
      onChange?.(latest);
    },
  }));

  return {
    serialize: () => latest,
    destroy: () => root.unmount(),
  };
}

export async function renderExcalidrawSvg(scene) {
  if (!scene || scene.type !== 'excalidraw' || !Array.isArray(scene.elements)) {
    throw new Error('expected an Excalidraw scene with an elements array');
  }
  await document.fonts.ready;
  const svg = await exportToSvg({
    elements: scene.elements.filter(element => !element.isDeleted),
    appState: scene.appState || {},
    files: scene.files || {},
    exportPadding: 20,
  });
  return new XMLSerializer().serializeToString(svg);
}
