import React from 'react';
import {createRoot} from 'react-dom/client';
import {Excalidraw, serializeAsJSON} from '@excalidraw/excalidraw';

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
