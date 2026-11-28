(function () {
  if (window.__markataSTLViewerInstalled) return;
  window.__markataSTLViewerInstalled = true;

  const currentScript = document.currentScript;
  const assetBase = currentScript && currentScript.src
    ? new URL('../assets/vendor/three-r128/', currentScript.src).href
    : new URL('/assets/vendor/three-r128/', window.location.href).href;
  const viewerSelector = '[data-markata-stl-model]';
  const activeViewers = new WeakMap();
  let dependenciesPromise;

  function loadScript(name) {
    return new Promise((resolve, reject) => {
      const script = document.createElement('script');
      script.src = new URL(name, assetBase).href;
      script.async = false;
      script.dataset.markataThreeAsset = name;
      script.onload = resolve;
      script.onerror = () => reject(new Error(`Could not load ${name}`));
      document.head.appendChild(script);
    });
  }

  function ensureThree() {
    if (window.THREE && THREE.STLLoader && THREE.OrbitControls) {
      return Promise.resolve(THREE);
    }
    if (!dependenciesPromise) {
      dependenciesPromise = (async () => {
        await loadScript('three.min.js');
        await loadScript('STLLoader.js');
        await loadScript('OrbitControls.js');
        await loadScript('BufferGeometryUtils.js');
        return window.THREE;
      })().catch((error) => {
        dependenciesPromise = null;
        throw error;
      });
    }
    return dependenciesPromise;
  }

  function disposeViewer(state) {
    state.closed = true;
    if (state.animationFrame) cancelAnimationFrame(state.animationFrame);
    if (state.resizeObserver) state.resizeObserver.disconnect();
    if (state.fallbackResize) window.removeEventListener('resize', state.fallbackResize);
    if (state.abortController) state.abortController.abort();
    if (state.keyboardHandler) state.mount.removeEventListener('keydown', state.keyboardHandler);
    if (state.controls) state.controls.dispose();
    if (state.geometry) state.geometry.dispose();
    if (state.material) state.material.dispose();
    if (state.renderer) {
      state.renderer.dispose();
      state.renderer.domElement.remove();
    }
  }

  function setStatus(widget, message) {
    const status = widget.querySelector('.md-3d-model-status');
    if (status) status.textContent = message;
  }

  function handleViewerKeydown(event, camera, controls, THREE) {
    const step = Math.PI / 36;
    let handled = true;
    if (event.altKey || event.ctrlKey || event.metaKey) return;

    if (event.key.startsWith('Arrow')) {
      const up = camera.up.clone().normalize();
      const toYUp = new THREE.Quaternion().setFromUnitVectors(up, new THREE.Vector3(0, 1, 0));
      const fromYUp = toYUp.clone().invert();
      const offset = camera.position.clone().sub(controls.target).applyQuaternion(toYUp);
      const spherical = new THREE.Spherical().setFromVector3(offset);
      if (event.key === 'ArrowLeft') spherical.theta += step;
      else if (event.key === 'ArrowRight') spherical.theta -= step;
      else if (event.key === 'ArrowUp') spherical.phi -= step;
      else if (event.key === 'ArrowDown') spherical.phi += step;
      else handled = false;
      if (handled) {
        spherical.makeSafe();
        const nextOffset = new THREE.Vector3().setFromSpherical(spherical).applyQuaternion(fromYUp);
        camera.position.copy(controls.target).add(nextOffset);
        camera.lookAt(controls.target);
        controls.update();
      }
    } else if (event.key === '+' || event.key === '=' || event.key === 'Add') {
      if (camera.isOrthographicCamera) {
        camera.zoom = Math.min(controls.maxZoom, camera.zoom * 1.2);
        camera.updateProjectionMatrix();
      } else {
        const offset = camera.position.clone().sub(controls.target);
        offset.setLength(Math.max(controls.minDistance, offset.length() / 1.2));
        camera.position.copy(controls.target).add(offset);
      }
      controls.update();
    } else if (event.key === '-' || event.key === '_' || event.key === 'Subtract') {
      if (camera.isOrthographicCamera) {
        camera.zoom = Math.max(controls.minZoom, camera.zoom / 1.2);
        camera.updateProjectionMatrix();
      } else {
        const offset = camera.position.clone().sub(controls.target);
        offset.setLength(Math.min(controls.maxDistance, offset.length() * 1.2));
        camera.position.copy(controls.target).add(offset);
      }
      controls.update();
    } else if (event.key === '0') {
      controls.reset();
    } else {
      handled = false;
    }

    if (handled) {
      event.preventDefault();
      event.stopPropagation();
    }
  }

  async function openViewer(widget, trigger) {
    if (activeViewers.has(widget)) return;

    const stage = widget.querySelector('.md-3d-model-viewer');
    const mount = widget.querySelector('.md-3d-model-canvas');
    if (!stage || !mount) return;

    const state = { closed: false, mount };
    activeViewers.set(widget, state);
    trigger.hidden = true;
    trigger.setAttribute('aria-expanded', 'true');
    stage.hidden = false;
    setStatus(widget, 'Loading the 3D model…');

    try {
      const THREE = await ensureThree();
      if (state.closed) return;

      const source = new URL(widget.dataset.modelUrl, window.location.href);
      // Dropper may have an older edge-cached response for the bare file URL
      // without the current CORS headers. A viewer-specific query key gets a
      // fresh response while keeping the source URL and its access options.
      source.searchParams.set('markata_stl', '1');
      state.abortController = new AbortController();
      const response = await fetch(source.href, {
        credentials: 'include',
        signal: state.abortController.signal,
      });
      if (!response.ok) throw new Error(`Model request failed (${response.status})`);

      const buffer = await response.arrayBuffer();
      if (state.closed) return;
      let geometry = new THREE.STLLoader().parse(buffer);
      if (THREE.BufferGeometryUtils && !geometry.index) {
        const smoothed = THREE.BufferGeometryUtils.mergeVertices(geometry, 0.0001);
        geometry.dispose();
        geometry = smoothed;
        geometry.computeVertexNormals();
      }
      geometry.computeBoundingBox();
      geometry.computeBoundingSphere();
      const center = geometry.boundingBox.getCenter(new THREE.Vector3());
      geometry.translate(-center.x, -center.y, -center.z);

      const radius = Math.max(geometry.boundingSphere.radius, 0.01);
      const material = new THREE.MeshStandardMaterial({
        color: 0x8eb9df,
        metalness: 0.03,
        roughness: 0.68,
      });
      const model = new THREE.Mesh(geometry, material);
      model.rotation.order = 'XYZ';
      model.rotation.set(
        (Number(source.searchParams.get('rx')) || 0) * Math.PI / 180,
        (Number(source.searchParams.get('ry')) || 0) * Math.PI / 180,
        (Number(source.searchParams.get('rz')) || 0) * Math.PI / 180,
      );

      const scene = new THREE.Scene();
      scene.background = new THREE.Color(0x0b1118);
      scene.add(model);
      scene.add(new THREE.HemisphereLight(0xeaf4ff, 0x25303a, 1.6));
      const keyLight = new THREE.DirectionalLight(0xffffff, 2.2);
      keyLight.position.set(radius * -1.5, radius * -2, radius * 3);
      scene.add(keyLight);
      const fillLight = new THREE.DirectionalLight(0xb8d9ff, 0.8);
      fillLight.position.set(radius * 2, radius, radius * 0.5);
      scene.add(fillLight);

      const width = Math.max(mount.clientWidth, 1);
      const height = Math.max(mount.clientHeight, 1);
      const aspect = width / height;
      const projection = source.searchParams.get('projection');
      let camera;
      if (projection === 'perspective') {
        camera = new THREE.PerspectiveCamera(35, aspect, radius / 1000, radius * 40);
      } else {
        const halfHeight = radius * 1.5;
        camera = new THREE.OrthographicCamera(
          -halfHeight * aspect,
          halfHeight * aspect,
          halfHeight,
          -halfHeight,
          radius / 1000,
          radius * 40,
        );
      }
      camera.up.set(0, 0, 1);
      camera.position.set(radius * 2.1, -radius * 2.4, radius * 1.8);
      camera.lookAt(0, 0, 0);

      const renderer = new THREE.WebGLRenderer({ antialias: true, alpha: false });
      renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 2));
      renderer.setSize(width, height);
      renderer.outputEncoding = THREE.sRGBEncoding;
      renderer.domElement.setAttribute('aria-hidden', 'true');
      mount.replaceChildren(renderer.domElement);

      const controls = new THREE.OrbitControls(camera, renderer.domElement);
      controls.target.set(0, 0, 0);
      controls.enableDamping = true;
      controls.dampingFactor = 0.08;
      controls.minDistance = radius * 0.15;
      controls.maxDistance = radius * 20;
      controls.minZoom = 0.2;
      controls.maxZoom = 20;
      controls.update();

      state.keyboardHandler = (event) => handleViewerKeydown(event, camera, controls, THREE);
      mount.addEventListener('keydown', state.keyboardHandler);
      mount.focus({ preventScroll: true });

      Object.assign(state, { geometry, material, renderer, controls });
      const resize = () => {
        if (state.closed) return;
        const nextWidth = Math.max(mount.clientWidth, 1);
        const nextHeight = Math.max(mount.clientHeight, 1);
        renderer.setSize(nextWidth, nextHeight);
        if (camera.isPerspectiveCamera) {
          camera.aspect = nextWidth / nextHeight;
        } else {
          const halfHeight = radius * 1.5;
          camera.left = -halfHeight * nextWidth / nextHeight;
          camera.right = halfHeight * nextWidth / nextHeight;
          camera.top = halfHeight;
          camera.bottom = -halfHeight;
        }
        camera.updateProjectionMatrix();
      };
      if (window.ResizeObserver) {
        state.resizeObserver = new ResizeObserver(resize);
        state.resizeObserver.observe(mount);
      } else {
        state.fallbackResize = resize;
        window.addEventListener('resize', resize);
      }
      setStatus(widget, 'Drag or use arrow keys to rotate. Scroll or press +/− to zoom; press 0 to reset.');

      const render = () => {
        if (state.closed) return;
        controls.update();
        renderer.render(scene, camera);
        state.animationFrame = requestAnimationFrame(render);
      };
      render();
    } catch (error) {
      if (state.closed || error.name === 'AbortError') return;
      disposeViewer(state);
      activeViewers.delete(widget);
      stage.hidden = false;
      trigger.hidden = true;
      trigger.setAttribute('aria-expanded', 'false');
      setStatus(widget, 'Could not load this model. Check the file link and try again.');
      console.error('Markata STL viewer failed to load a model:', error);
    }
  }

  document.addEventListener('click', (event) => {
    const trigger = event.target.closest('.md-3d-model-trigger');
    if (trigger) {
      const widget = trigger.closest(viewerSelector);
      if (widget) openViewer(widget, trigger);
      return;
    }

    const close = event.target.closest('.md-3d-model-close');
    if (!close) return;
    const widget = close.closest(viewerSelector);
    const state = widget && activeViewers.get(widget);
    if (!widget) return;
    if (state) {
      disposeViewer(state);
      activeViewers.delete(widget);
    }
    widget.querySelector('.md-3d-model-viewer').hidden = true;
    const triggerButton = widget.querySelector('.md-3d-model-trigger');
    triggerButton.hidden = false;
    triggerButton.setAttribute('aria-expanded', 'false');
    triggerButton.focus({ preventScroll: true });
    setStatus(widget, '');
  });
})();
