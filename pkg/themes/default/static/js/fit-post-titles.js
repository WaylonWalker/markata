(() => {
  const selector = '.post-header__title-row > h1[data-title-size]';
  let observed;
  let frame;
  let observedWidth = 0;
  let followingFonts = false;
  const resizeObserver = window.ResizeObserver ? new ResizeObserver((entries) => {
    followFonts();
    const width = entries[0]?.contentRect.width || 0;
    if (width !== observedWidth) { observedWidth = width; schedule(); }
  }) : null;

  function fit() {
    frame = null;
    const title = observed;
    if (!title?.isConnected) return;
    title.style.fontSize = '';
    title.style.whiteSpace = 'nowrap';
    const available = title.clientWidth;
    const computed = getComputedStyle(title);
    const starting = parseFloat(computed.fontSize);
    const minimum = parseFloat(getComputedStyle(document.documentElement).fontSize) * 1.25;
    if (!available || !Number.isFinite(starting)) { title.style.whiteSpace = ''; return; }
    const range = document.createRange();
    range.selectNodeContents(title);
    const fits = (size) => {
      title.style.fontSize = `${size}px`;
      return range.getBoundingClientRect().width <= available - 2;
    };
    if (!fits(starting)) {
      if (fits(minimum)) {
        let low = minimum;
        let high = starting;
        for (let i = 0; i < 9; i++) {
          const middle = (low + high) / 2;
          if (fits(middle)) low = middle;
          else high = middle;
        }
        title.style.fontSize = `${low}px`;
      } else {
        // When a title cannot fit on one line at a readable size, keep the
        // server-rendered size and let normal wrapping make room.
        title.style.fontSize = '';
      }
    }
    title.style.whiteSpace = '';
  }

  function schedule() {
    if (frame) cancelAnimationFrame(frame);
    frame = requestAnimationFrame(fit);
  }

  function followFonts() {
    if (followingFonts) return;
    followingFonts = true;
    // ResizeObserver delivery follows layout, so reading font readiness here
    // cannot pull the first full-page layout into script evaluation.
    document.fonts?.ready.then(schedule);
  }

  function connect() {
    const title = document.querySelector(selector);
    if (title === observed) return;
    resizeObserver?.disconnect();
    observed = title;
    observedWidth = 0;
    // The initial observer delivery follows normal layout. An eager frame
    // callback would measure before that layout and force it synchronously.
    if (title) { if (resizeObserver) resizeObserver.observe(title); else schedule(); }
  }

  connect();
  if (!resizeObserver) requestAnimationFrame(() => requestAnimationFrame(followFonts));
  window.addEventListener('resize', schedule);
  window.addEventListener('view-transition-complete', connect);
})();
