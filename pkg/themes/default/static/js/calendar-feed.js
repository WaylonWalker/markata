(() => {
  'use strict';

  const monthNames = [
    'January', 'February', 'March', 'April', 'May', 'June',
    'July', 'August', 'September', 'October', 'November', 'December'
  ];
  const weekdayNames = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
  const dropperHosts = new Set([
    'dropper.wayl.one',
    'dropper.waylonwalker.com',
    'dropper-dev.wayl.one'
  ]);

  const pad2 = (value) => String(value).padStart(2, '0');
  const dateKey = (year, month, day) => `${year}-${pad2(month + 1)}-${pad2(day)}`;

  function parseDate(value) {
    const match = /^(\d{4})-(\d{2})-(\d{2})$/.exec(value || '');
    if (!match) return null;

    const year = Number(match[1]);
    const month = Number(match[2]) - 1;
    const day = Number(match[3]);
    const candidate = new Date(year, month, day);

    if (
      candidate.getFullYear() !== year ||
      candidate.getMonth() !== month ||
      candidate.getDate() !== day
    ) {
      return null;
    }

    return { year, month, day };
  }

  function element(tag, className, text) {
    const node = document.createElement(tag);
    if (className) node.className = className;
    if (text !== undefined) node.textContent = text;
    return node;
  }

  function dropperThumbURL(raw) {
    if (!raw) return '';
    try {
      const url = new URL(raw, window.location.href);
      if (!dropperHosts.has(url.hostname.toLowerCase())) return '';

      const slash = url.pathname.lastIndexOf('/');
      const dot = url.pathname.lastIndexOf('.');
      if (dot <= slash) return '';

      const stem = url.pathname.slice(0, dot);
      if (!stem.endsWith('_thumb')) {
        url.pathname = `${stem}_thumb${url.pathname.slice(dot)}`;
      }
      for (const key of ['w', 'h', 'width', 'height']) url.searchParams.delete(key);
      return url.toString();
    } catch (error) {
      return '';
    }
  }

  function previewImageCandidates(raw) {
    const thumb = dropperThumbURL(raw);
    if (thumb && thumb !== raw) return [thumb, raw];
    return raw ? [thumb || raw] : [];
  }

  function postPreview(post) {
    const item = element('li', 'calendar-preview-post');
    const link = element('a', 'calendar-preview-link');
    link.href = post.href;

    const imageCandidates = previewImageCandidates(post.image);
    if (imageCandidates.length > 0) {
      const image = document.createElement('img');
      image.className = 'calendar-preview-image';
      image.alt = '';
      image.decoding = 'async';
      image.width = 72;
      image.height = 54;
      if (imageCandidates.length > 1) {
        image.addEventListener('error', () => {
          if (image.dataset.fallbackUsed === 'true') {
            image.hidden = true;
            return;
          }
          image.dataset.fallbackUsed = 'true';
          image.src = imageCandidates[1];
        });
      } else {
        image.addEventListener('error', () => { image.hidden = true; });
      }
      image.src = imageCandidates[0];
      link.append(image);
    }

    const copy = element('span', 'calendar-preview-copy');
    copy.append(element('strong', 'calendar-preview-title', post.title));
    if (post.description) {
      copy.append(element('span', 'calendar-preview-description', post.description));
    }
    link.append(copy);
    item.append(link);
    return item;
  }

  function renderPreview(details, year, month, day, posts) {
    const preview = element('div', 'calendar-day-preview');
    preview.style.insetBlockStart = 'calc(100% + 0.12rem)';
    const dateLabel = `${monthNames[month]} ${day}, ${year}`;
    const postWord = posts.length === 1 ? 'post' : 'posts';
    preview.append(element('div', 'calendar-preview-summary', `${dateLabel} · ${posts.length} ${postWord}`));

    const postList = element('ul', 'calendar-day-posts');
    for (const post of posts) postList.append(postPreview(post));
    preview.append(postList);
    details.append(preview);
  }

  function renderMonth(year, month, postsByDate, previewData) {
    const section = element('section', 'calendar-month');
    section.setAttribute('aria-labelledby', `calendar-${year}-${month + 1}`);

    const heading = element('h3', '', monthNames[month]);
    heading.id = `calendar-${year}-${month + 1}`;
    section.append(heading);

    const weekdayRow = element('div', 'calendar-weekdays');
    weekdayRow.setAttribute('aria-hidden', 'true');
    for (const weekday of weekdayNames) weekdayRow.append(element('span', '', weekday));
    section.append(weekdayRow);

    const grid = element('div', 'calendar-days');
    const firstWeekday = new Date(year, month, 1).getDay();
    const daysInMonth = new Date(year, month + 1, 0).getDate();

    for (let i = 0; i < firstWeekday; i += 1) {
      const blank = element('span', 'calendar-day-empty');
      blank.setAttribute('aria-hidden', 'true');
      grid.append(blank);
    }

    for (let day = 1; day <= daysInMonth; day += 1) {
      const key = dateKey(year, month, day);
      const posts = postsByDate.get(key) || [];
      const cell = element('div', 'calendar-day');

      if (posts.length === 0) {
        const quiet = element('span', 'calendar-day-quiet', String(day));
        quiet.setAttribute('aria-label', `${monthNames[month]} ${day}, ${year}`);
        cell.append(quiet);
        grid.append(cell);
        continue;
      }

      const details = document.createElement('details');
      const summary = document.createElement('summary');
      const postWord = posts.length === 1 ? 'post' : 'posts';
      details.dataset.postCount = String(posts.length);
      summary.setAttribute('aria-label', `${monthNames[month]} ${day}, ${year}: ${posts.length} ${postWord}`);
      summary.append(element('span', 'calendar-day-number', String(day)));

      const marker = element('span', 'calendar-day-marker');
      marker.setAttribute('aria-hidden', 'true');
      marker.style.setProperty('--calendar-density', String(Math.min(posts.length, 6)));
      summary.append(marker);

      if (posts.length > 1) {
        const count = element('span', 'calendar-day-count', String(posts.length));
        count.setAttribute('aria-hidden', 'true');
        summary.append(count);
      }

      details.append(summary);
      previewData.set(details, { year, month, day, posts });
      cell.append(details);
      grid.append(cell);
    }

    const usedCells = firstWeekday + daysInMonth;
    const trailingCells = (7 - (usedCells % 7)) % 7;
    for (let i = 0; i < trailingCells; i += 1) {
      const blank = element('span', 'calendar-day-empty');
      blank.setAttribute('aria-hidden', 'true');
      grid.append(blank);
    }

    section.append(grid);
    return section;
  }

  function renderCalendar(root, calendar, datedPosts) {
    const postsByDate = new Map();
    const years = new Set();
    const previewData = new WeakMap();

    for (const post of datedPosts) {
      const parsed = parseDate(post.date);
      if (!parsed) continue;

      years.add(parsed.year);
      const existing = postsByDate.get(post.date) || [];
      existing.push(post);
      postsByDate.set(post.date, existing);
    }

    const sortedYears = Array.from(years).sort((a, b) => b - a);
    if (sortedYears.length === 0) return false;

    const yearNav = element('nav', 'calendar-year-nav calendar-view-switch');
    yearNav.setAttribute('aria-label', 'Calendar year');
    const olderButton = element('button', 'calendar-year-step calendar-year-older');
    olderButton.type = 'button';
    const yearSelect = document.createElement('select');
    yearSelect.className = 'calendar-year-select';
    yearSelect.setAttribute('aria-label', 'Calendar year');
    const newerButton = element('button', 'calendar-year-step calendar-year-newer');
    newerButton.type = 'button';

    for (const year of sortedYears) {
      const option = document.createElement('option');
      option.value = String(year);
      option.textContent = String(year);
      yearSelect.append(option);
    }
    yearNav.append(olderButton, yearSelect, newerButton);

    const syncYearControls = (year) => {
      const index = sortedYears.indexOf(year);
      const olderYear = index >= 0 ? sortedYears[index + 1] : undefined;
      const newerYear = index > 0 ? sortedYears[index - 1] : undefined;

      yearSelect.value = String(year);
      olderButton.disabled = olderYear === undefined;
      newerButton.disabled = newerYear === undefined;
      olderButton.textContent = olderYear === undefined ? '← Older' : `← ${olderYear}`;
      newerButton.textContent = newerYear === undefined ? 'Newer →' : `${newerYear} →`;
    };

    const renderYear = (year, persistURL = true) => {
      if (!years.has(year)) year = sortedYears[0];

      const yearSection = element('section', 'calendar-year');
      const yearHeading = element('h2', '', String(year));
      yearHeading.id = `calendar-year-${year}`;
      yearSection.setAttribute('aria-labelledby', yearHeading.id);
      const yearHeader = element('header', 'calendar-year-header');
      yearHeader.append(yearHeading);
      const yearPosts = datedPosts.filter(post => post.date.startsWith(`${year}-`));
      const activeDays = new Set(yearPosts.map(post => post.date)).size;
      yearHeader.append(element('p', 'calendar-year-stats', `${yearPosts.length} posts · ${activeDays} publishing days`));
      yearSection.append(yearHeader);

      const months = element('div', 'calendar-months');
      for (let month = 0; month < 12; month += 1) {
        months.append(renderMonth(year, month, postsByDate, previewData));
      }
      yearSection.append(months);

      calendar.replaceChildren();
      if (sortedYears.length > 1) calendar.append(yearNav);
      calendar.append(yearSection);
      syncYearControls(year);
      if (persistURL) updateURLYear(year);
    };

    olderButton.addEventListener('click', () => {
      const currentIndex = sortedYears.indexOf(Number(yearSelect.value));
      const olderYear = sortedYears[currentIndex + 1];
      if (olderYear !== undefined) renderYear(olderYear);
    });
    newerButton.addEventListener('click', () => {
      const currentIndex = sortedYears.indexOf(Number(yearSelect.value));
      const newerYear = sortedYears[currentIndex - 1];
      if (newerYear !== undefined) renderYear(newerYear);
    });
    yearSelect.addEventListener('change', () => renderYear(Number(yearSelect.value)));

    const initialYear = calendarYearFromURL(sortedYears) || sortedYears[0];
    renderYear(initialYear, false);
    window.addEventListener('popstate', () => {
      renderYear(calendarYearFromURL(sortedYears) || sortedYears[0], false);
    });

    const ensurePreview = (details) => {
      if (details.querySelector('.calendar-day-preview')) return;
      const data = previewData.get(details);
      if (!data) return;
      renderPreview(details, data.year, data.month, data.day, data.posts);
    };

    const openExclusive = (details, interaction = 'interactive') => {
      ensurePreview(details);
      for (const candidate of root.querySelectorAll('.calendar-day details[open]')) {
        if (candidate !== details) {
          candidate.open = false;
          delete candidate.dataset.calendarHoverPreview;
        }
      }
      const preview = details.querySelector('.calendar-day-preview');
      if (interaction === 'hover') {
        details.dataset.calendarHoverPreview = 'true';
        if (preview) preview.style.pointerEvents = 'none';
      } else {
        delete details.dataset.calendarHoverPreview;
        if (preview) preview.style.pointerEvents = 'auto';
      }
      details.open = true;
    };

    calendar.addEventListener('pointerover', (event) => {
      if (!window.matchMedia('(hover: hover) and (pointer: fine)').matches) return;
      const details = event.target.closest('.calendar-day details');
      if (!details || !calendar.contains(details)) return;
      openExclusive(details, 'hover');
    });

    calendar.addEventListener('focusin', (event) => {
      const details = event.target.closest('.calendar-day details');
      if (details && calendar.contains(details)) openExclusive(details, 'interactive');
    });

    calendar.addEventListener('click', (event) => {
      const summary = event.target.closest('.calendar-day details > summary');
      if (!summary || !calendar.contains(summary)) return;
      const details = summary.parentElement;
      if (details && details.dataset.calendarHoverPreview === 'true') {
        event.preventDefault();
        openExclusive(details, 'interactive');
      }
    });

    calendar.addEventListener('toggle', (event) => {
      const opened = event.target;
      if (!(opened instanceof HTMLDetailsElement) || !opened.open) return;
      ensurePreview(opened);
      for (const details of root.querySelectorAll('.calendar-day details[open]')) {
        if (details !== opened) details.open = false;
      }
    }, true);

    return true;
  }

  function calendarModeFromURL(defaultMode = 'list') {
    try {
      const mode = new URLSearchParams(window.location.search).get('view');
      return mode === 'calendar' || mode === 'list' ? mode : defaultMode;
    } catch (error) {
      return defaultMode;
    }
  }

  function calendarYearFromURL(years) {
    try {
      const raw = new URLSearchParams(window.location.search).get('year');
      if (!/^\d{4}$/.test(raw || '')) return null;
      const year = Number(raw);
      return years.includes(year) ? year : null;
    } catch (error) {
      return null;
    }
  }

  function updateURLMode(mode, defaultMode = 'list') {
    try {
      const url = new URL(window.location.href);
      if (mode === 'calendar') url.searchParams.set('view', 'calendar');
      else if (defaultMode === 'calendar') url.searchParams.set('view', 'list');
      else if (['calendar', 'list'].includes(url.searchParams.get('view'))) url.searchParams.delete('view');
      if (url.toString() !== window.location.href) window.history.pushState(window.history.state, '', url.toString());
    } catch (error) {
      // URL state is progressive enhancement only.
    }
  }

  function updateURLYear(year) {
    try {
      const url = new URL(window.location.href);
      url.searchParams.set('year', String(year));
      if (url.toString() !== window.location.href) window.history.pushState(window.history.state, '', url.toString());
    } catch (error) {
      // URL state is progressive enhancement only.
    }
  }

  function primaryNodesFor(root, sourceList) {
    const nodes = Array.from(root.querySelectorAll('[data-calendar-primary]'));
    for (const child of root.children) {
      if (child.classList.contains('pagination') || child.classList.contains('pagination-infinite') || child.classList.contains('feed-empty')) {
        if (!nodes.includes(child)) nodes.push(child);
      }
    }
    if (nodes.length === 0 && sourceList) nodes.push(sourceList);
    return nodes;
  }

  function initialize(root) {
    const sourceList = root.querySelector('[data-calendar-list]');
    const calendar = root.querySelector('[data-calendar-years]');
    const switcher = root.querySelector('[data-calendar-switch]');
    const calendarSource = root.getAttribute('data-calendar-src') || '';
    const errorMessage = root.querySelector('[data-calendar-error]');
    if ((!sourceList && !calendarSource) || !calendar || !switcher) return;

    let datedPosts = null;
    let postsPromise = null;
    const loadDatedPosts = async () => {
      if (datedPosts) return datedPosts;

      if (sourceList) {
        datedPosts = [];
        for (const item of sourceList.querySelectorAll('[data-calendar-post]')) {
          const date = item.getAttribute('data-date') || '';
          const link = item.querySelector('a[href]');
          if (!parseDate(date) || !link) continue;
          datedPosts.push({
            date,
            href: link.href,
            title: link.textContent.trim() || link.href,
            description: item.getAttribute('data-description') || '',
            image: item.getAttribute('data-image') || ''
          });
        }
        return datedPosts;
      }

      if (!postsPromise) {
        postsPromise = fetch(calendarSource, { headers: { Accept: 'application/json' } })
          .then((response) => {
            if (!response.ok) throw new Error(`Calendar request failed (${response.status})`);
            return response.json();
          })
          .then((payload) => {
            const items = Array.isArray(payload) ? payload : payload.items;
            if (!Array.isArray(items)) throw new Error('Calendar response has no items');
            return items.filter((item) => parseDate(item.date) && item.href).map((item) => ({
              date: item.date,
              href: new URL(item.href, window.location.href).href,
              title: item.title || item.href,
              description: item.description || '',
              image: item.image || ''
            }));
          });
      }

      try {
        datedPosts = await postsPromise;
        return datedPosts;
      } catch (error) {
        postsPromise = null;
        throw error;
      }
    };

    let calendarRendered = false;
    let calendarPromise = null;
    const ensureCalendar = async () => {
      if (calendarRendered) return true;
      if (!calendarPromise) {
        calendarPromise = loadDatedPosts().then((posts) => renderCalendar(root, calendar, posts));
      }
      try {
        calendarRendered = await calendarPromise;
        if (errorMessage) errorMessage.hidden = true;
        return calendarRendered;
      } catch (error) {
        calendarPromise = null;
        if (errorMessage) errorMessage.hidden = false;
        return false;
      }
    };

    const defaultMode = root.getAttribute('data-calendar-default') === 'calendar' ? 'calendar' : 'list';
    const primaryNodes = primaryNodesFor(root, sourceList);
    const sourceIsPrimary = primaryNodes.includes(sourceList);
    const buttons = Array.from(switcher.querySelectorAll('[data-calendar-mode]'));
    const setMode = async (mode, persistURL = true) => {
      const calendarMode = mode === 'calendar';
      if (calendarMode && !await ensureCalendar()) return;
      calendar.hidden = !calendarMode;
      for (const node of primaryNodes) node.hidden = calendarMode;
      if (!sourceIsPrimary && sourceList) sourceList.hidden = true;
      if (!calendarMode && errorMessage) errorMessage.hidden = true;

      for (const button of buttons) {
        button.setAttribute('aria-pressed', button.getAttribute('data-calendar-mode') === mode ? 'true' : 'false');
      }

      if (persistURL && root.hasAttribute('data-calendar-url-state')) updateURLMode(mode, defaultMode);
    };

    for (const button of buttons) button.addEventListener('click', () => {
      void setMode(button.getAttribute('data-calendar-mode'));
    });

    const initialMode = root.hasAttribute('data-calendar-url-state') ? calendarModeFromURL(defaultMode) : defaultMode;

    switcher.hidden = false;
    void setMode(initialMode, false);
    window.addEventListener('popstate', () => {
      const mode = calendarModeFromURL(defaultMode);
      void setMode(mode, false);
    });
  }

  for (const root of document.querySelectorAll('[data-calendar-feed]')) initialize(root);
})();
