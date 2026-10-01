(() => {
  'use strict';

  const DAYS = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];
  const store = {
    get(key, fallback) {
      try {
        const v = localStorage.getItem(key);
        return v == null ? fallback : JSON.parse(v);
      } catch {
        return fallback;
      }
    },
    set(key, value) {
      try {
        localStorage.setItem(key, JSON.stringify(value));
      } catch {}
    },
  };

  const $ = (id) => document.getElementById(id);
  const el = {
    q: $('q'), list: $('list'), status: $('status'), days: $('days'), meta: $('meta'),
    fPicks: $('f-picks'), fShared: $('f-shared'), fPast: $('f-hide-past'),
    picksCount: $('picks-count'), sharedCount: $('shared-count'), banner: $('shared-banner'),
    shareBtn: $('share-btn'), icsBtn: $('ics-btn'), dialog: $('share-dialog'), qr: $('qr'),
    shareName: $('share-name'), shareUrl: $('share-url'), copyBtn: $('copy-btn'), nativeShare: $('native-share-btn'),
  };

  let gigs = [];
  const byKey = new Map();
  const picks = new Set(store.get('lmf.picks', []));
  let shared = null; // { name, ids: Set }
  const state = { q: '', day: store.get('lmf.day', 'All'), onlyPicks: false, onlyShared: false, hidePast: store.get('lmf.hidePast', false) };

  // ---------- helpers ----------
  const esc = (s) => String(s).replace(/[&<>"']/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' })[c]);
  const fold = (s) => String(s).normalize('NFD').replace(/[̀-ͯ]/g, '').toLowerCase();
  const time = (iso) => (iso ? iso.slice(11, 16) : '');
  // Current wall-clock time at the festival, as "yyyy-mm-ddTHH:MM".
  const nowLocal = () => new Intl.DateTimeFormat('sv-SE', {
    timeZone: 'Europe/London', year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', hourCycle: 'h23',
  }).format(new Date()).replace(' ', 'T');

  function highlight(text, terms) {
    if (!terms.length) return esc(text);
    const f = fold(text);
    const marks = new Array(text.length).fill(false);
    for (const t of terms) {
      let i = f.indexOf(t);
      while (i !== -1) {
        for (let j = i; j < i + t.length; j++) marks[j] = true;
        i = f.indexOf(t, i + t.length);
      }
    }
    // fold() can change length for some characters; only highlight when it didn't.
    if (f.length !== text.length) return esc(text);
    let out = '', open = false;
    for (let i = 0; i < text.length; i++) {
      if (marks[i] !== open) { out += open ? '</mark>' : '<mark>'; open = marks[i]; }
      out += esc(text[i]);
    }
    return out + (open ? '</mark>' : '');
  }

  const overlaps = (a, b) => a.st < b.en && b.st < a.en;

  function clashes() {
    const mine = gigs.filter((g) => picks.has(g.id));
    const map = new Map();
    for (let i = 0; i < mine.length; i++) {
      for (let j = i + 1; j < mine.length; j++) {
        if (overlaps(mine[i], mine[j])) {
          (map.get(mine[i].id) || map.set(mine[i].id, []).get(mine[i].id)).push(mine[j]);
          (map.get(mine[j].id) || map.set(mine[j].id, []).get(mine[j].id)).push(mine[i]);
        }
      }
    }
    return map;
  }

  // ---------- rendering ----------
  function renderDays() {
    const counts = {};
    gigs.forEach((g) => (counts[g.d] = (counts[g.d] || 0) + 1));
    const days = ['All', ...DAYS.filter((d) => counts[d])];
    if (!days.includes(state.day)) state.day = 'All';
    el.days.innerHTML = days
      .map((d) => `<button class="chip" type="button" data-day="${d}" aria-pressed="${state.day === d}">${d === 'All' ? 'All days' : d.slice(0, 3)} <span class="n">${d === 'All' ? gigs.length : counts[d]}</span></button>`)
      .join('');
  }

  function render() {
    const terms = fold(state.q).split(/\s+/).filter(Boolean);
    const now = nowLocal();
    const clashMap = clashes();

    const visible = gigs.filter((g) => {
      if (state.day !== 'All' && g.d !== state.day) return false;
      if (state.onlyPicks && !picks.has(g.id)) return false;
      if (state.onlyShared && !(shared && shared.ids.has(g.id))) return false;
      if (state.hidePast && g.en <= now) return false;
      return terms.every((t) => g._s.includes(t));
    });

    el.picksCount.textContent = picks.size || '';
    el.fPicks.setAttribute('aria-pressed', state.onlyPicks);
    el.fPast.setAttribute('aria-pressed', state.hidePast);
    el.fShared.setAttribute('aria-pressed', state.onlyShared);
    el.days.querySelectorAll('.chip').forEach((c) => c.setAttribute('aria-pressed', c.dataset.day === state.day));

    if (!visible.length) {
      el.status.hidden = false;
      el.status.textContent = state.onlyPicks && !picks.size
        ? 'No picks yet — tap ☆ next to a band to add it.'
        : 'Nothing matches. Try fewer words, or another day.';
      el.list.innerHTML = '';
      return;
    }
    el.status.hidden = false;
    el.status.textContent = `${visible.length} of ${gigs.length} gigs` + (terms.length ? ` matching “${state.q.trim()}”` : '');

    let html = '', day = null;
    for (const g of visible) {
      if (g.d !== day) {
        day = g.d;
        html += `<h2 class="day">${esc(day)} <span class="n">${esc(g.st.slice(8, 10))}/${esc(g.st.slice(5, 7))}</span></h2>`;
      }
      const picked = picks.has(g.id);
      const tags = [];
      if (g.st <= now && now < g.en) tags.push('<span class="tag now">On now</span>');
      if (g.t && g.t !== 'Performance') tags.push(`<span class="tag">${esc(g.t)}</span>`);
      if (shared && shared.ids.has(g.id)) tags.push(`<span class="tag shared">♥ ${esc(shared.name || 'Shared')}</span>`);
      if (picked && clashMap.has(g.id)) {
        tags.push(`<span class="tag clash" title="${esc(clashMap.get(g.id).map((c) => `${c.a} ${time(c.st)} @ ${c.v}`).join('\n'))}">Clashes with ${esc(clashMap.get(g.id).map((c) => c.a).join(', '))}</span>`);
      }
      if (g.u) tags.push(`<span class="tag"><a href="${esc(g.u)}" target="_blank" rel="noopener">Tickets</a></span>`);
      html += `<article class="gig${picked ? ' picked' : ''}${g.en <= now ? ' past' : ''}" data-id="${g.id}">
        <div class="time">${time(g.st)}<small>–${time(g.en)}</small></div>
        <div class="who">
          <div class="artist"><button type="button" data-search="${esc(g.a)}" title="Show all sets by ${esc(g.a)}">${highlight(g.a, terms)}</button></div>
          <div class="venue"><button type="button" data-search="${esc(g.v)}" title="Show everything at ${esc(g.v)}">${highlight(g.v, terms)}</button></div>
          ${g.x ? `<div class="desc">${highlight(g.x, terms)}</div>` : ''}
          ${tags.length ? `<div class="tags">${tags.join('')}</div>` : ''}
        </div>
        <button class="star" type="button" data-pick="${g.id}" aria-pressed="${picked}" aria-label="${picked ? 'Remove from' : 'Add to'} my picks: ${esc(g.a)}">${picked ? '★' : '☆'}</button>
      </article>`;
    }
    el.list.innerHTML = html;
  }

  // ---------- sharing ----------
  function shareLink() {
    const keys = gigs.filter((g) => picks.has(g.id)).map((g) => g.k);
    const params = new URLSearchParams();
    params.set('s', keys.join('.'));
    const name = el.shareName.value.trim();
    if (name) params.set('n', name);
    return `${location.origin}${location.pathname}#${params.toString()}`;
  }

  function updateShareDialog() {
    const url = shareLink();
    el.shareUrl.value = url;
    try {
      const qr = qrcode(0, 'L');
      qr.addData(url);
      qr.make();
      el.qr.innerHTML = qr.createSvgTag({ cellSize: 4, margin: 0, scalable: true });
    } catch {
      el.qr.textContent = 'Too many picks to fit in a QR code — use the link instead.';
    }
  }

  function readSharedFromHash() {
    const params = new URLSearchParams(location.hash.slice(1));
    const s = params.get('s');
    if (s == null) return;
    const ids = new Set(s.split('.').map((k) => byKey.get(k)).filter(Boolean).map((g) => g.id));
    shared = { name: params.get('n') || '', ids };
    renderBanner();
  }

  function renderBanner() {
    if (!shared) {
      el.banner.hidden = true;
      el.fShared.hidden = true;
      return;
    }
    const who = shared.name ? esc(shared.name) + '’s' : 'Shared';
    const inCommon = [...shared.ids].filter((id) => picks.has(id)).length;
    el.banner.hidden = false;
    el.fShared.hidden = false;
    el.sharedCount.textContent = shared.ids.size;
    el.banner.innerHTML = `<span><strong>${who} picks</strong> · ${shared.ids.size} gigs · ${inCommon} in common with yours</span>
      <button class="btn" type="button" data-banner="only">${state.onlyShared ? 'Show everything' : 'Show only these'}</button>
      <button class="btn" type="button" data-banner="merge">Add all to my picks</button>
      <button class="btn" type="button" data-banner="close">Dismiss</button>`;
  }

  // ---------- calendar export ----------
  function downloadIcs() {
    const mine = gigs.filter((g) => picks.has(g.id));
    if (!mine.length) {
      alert('Pick some gigs first (tap ☆).');
      return;
    }
    const icsTime = (iso) => iso.replace(/[-:]/g, '') + '00';
    const text = (s) => String(s).replace(/[\\;,]/g, (c) => '\\' + c).replace(/\n/g, '\\n');
    const stamp = new Date().toISOString().replace(/[-:]/g, '').replace(/\.\d+/, '');
    const lines = ['BEGIN:VCALENDAR', 'VERSION:2.0', 'PRODID:-//LMF Planner//EN', 'CALSCALE:GREGORIAN'];
    for (const g of mine) {
      lines.push(
        'BEGIN:VEVENT',
        `UID:${g.id}@lmf-planner`,
        `DTSTAMP:${stamp}`,
        `DTSTART;TZID=Europe/London:${icsTime(g.st)}`,
        `DTEND;TZID=Europe/London:${icsTime(g.en)}`,
        `SUMMARY:${text(g.a)}`,
        `LOCATION:${text(g.v + ', Lancaster')}`,
        `DESCRIPTION:${text([g.x, g.u].filter(Boolean).join('\n'))}`,
        'END:VEVENT',
      );
    }
    lines.push('END:VCALENDAR');
    const blob = new Blob([lines.join('\r\n')], { type: 'text/calendar' });
    const a = Object.assign(document.createElement('a'), { href: URL.createObjectURL(blob), download: 'lmf-picks.ics' });
    a.click();
    setTimeout(() => URL.revokeObjectURL(a.href), 1000);
  }

  // ---------- events ----------
  function savePicks() {
    store.set('lmf.picks', [...picks]);
  }

  el.q.addEventListener('input', () => {
    state.q = el.q.value;
    render();
  });

  // Type anywhere to search.
  document.addEventListener('keydown', (e) => {
    if (el.dialog.open || e.ctrlKey || e.metaKey || e.altKey) return;
    const t = e.target;
    if (t !== el.q && (t.tagName === 'INPUT' || t.tagName === 'TEXTAREA')) return;
    if (e.key === 'Escape') {
      el.q.value = '';
      state.q = '';
      render();
      el.q.blur();
      return;
    }
    if (t !== el.q && e.key.length === 1 && e.key !== ' ') {
      el.q.focus();
    }
  });

  el.days.addEventListener('click', (e) => {
    const b = e.target.closest('[data-day]');
    if (!b) return;
    state.day = b.dataset.day;
    store.set('lmf.day', state.day);
    render();
  });

  el.fPicks.addEventListener('click', () => {
    state.onlyPicks = !state.onlyPicks;
    render();
  });
  el.fShared.addEventListener('click', () => {
    state.onlyShared = !state.onlyShared;
    renderBanner();
    render();
  });
  el.fPast.addEventListener('click', () => {
    state.hidePast = !state.hidePast;
    store.set('lmf.hidePast', state.hidePast);
    render();
  });

  el.list.addEventListener('click', (e) => {
    const star = e.target.closest('[data-pick]');
    if (star) {
      const id = star.dataset.pick;
      picks.has(id) ? picks.delete(id) : picks.add(id);
      savePicks();
      renderBanner();
      render();
      return;
    }
    const s = e.target.closest('[data-search]');
    if (s) {
      el.q.value = s.dataset.search;
      state.q = el.q.value;
      state.day = 'All';
      render();
      window.scrollTo({ top: 0 });
    }
  });

  el.banner.addEventListener('click', (e) => {
    const action = e.target.closest('[data-banner]')?.dataset.banner;
    if (action === 'only') state.onlyShared = !state.onlyShared;
    if (action === 'merge') {
      shared.ids.forEach((id) => picks.add(id));
      savePicks();
    }
    if (action === 'close') {
      shared = null;
      state.onlyShared = false;
      history.replaceState(null, '', location.pathname + location.search);
    }
    renderBanner();
    render();
  });

  el.shareBtn.addEventListener('click', () => {
    if (!picks.size) {
      alert('Pick some gigs first (tap ☆), then share them.');
      return;
    }
    el.shareName.value = store.get('lmf.name', '');
    updateShareDialog();
    el.nativeShare.hidden = !navigator.share;
    el.dialog.showModal();
  });
  el.shareName.addEventListener('input', () => {
    store.set('lmf.name', el.shareName.value.trim());
    updateShareDialog();
  });
  el.copyBtn.addEventListener('click', async () => {
    try {
      await navigator.clipboard.writeText(el.shareUrl.value);
      el.copyBtn.textContent = 'Copied!';
    } catch {
      el.shareUrl.select();
      el.copyBtn.textContent = 'Press Ctrl+C';
    }
    setTimeout(() => (el.copyBtn.textContent = 'Copy link'), 1500);
  });
  el.nativeShare.addEventListener('click', () => {
    navigator.share({ title: 'My Lancaster Music Festival picks', url: el.shareUrl.value }).catch(() => {});
  });
  el.icsBtn.addEventListener('click', downloadIcs);

  window.addEventListener('hashchange', () => {
    readSharedFromHash();
    render();
  });

  // Keep the day headers stuck just below the (variable-height) header.
  const header = document.querySelector('.top');
  new ResizeObserver(() => document.documentElement.style.setProperty('--header-h', header.offsetHeight + 'px')).observe(header);

  // Refresh "on now" / finished state every minute.
  setInterval(render, 60_000);

  // ---------- load ----------
  fetch('data.json', { cache: 'no-cache' })
    .then((r) => {
      if (!r.ok) throw new Error(r.status);
      return r.json();
    })
    .then((data) => {
      gigs = data.gigs;
      for (const g of gigs) {
        g._s = fold([g.a, g.v, g.x, g.t, g.d].join(' '));
        byKey.set(g.k, g);
      }
      // Forget picks for gigs that no longer exist.
      const ids = new Set(gigs.map((g) => g.id));
      [...picks].forEach((id) => ids.has(id) || picks.delete(id));
      savePicks();

      const fetched = new Date(data.fetchedAt);
      el.meta.textContent = `${gigs.length} gigs · data updated ${fetched.toLocaleString('en-GB', { dateStyle: 'medium', timeStyle: 'short' })}`;
      renderDays();
      readSharedFromHash();
      render();
    })
    .catch((err) => {
      el.status.textContent = `Couldn't load the schedule (${err.message}).`;
    });
})();
