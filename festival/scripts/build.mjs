// Normalises the raw Knack export (gigs.json) into the compact data.json the
// site loads, and optionally reports what changed against a previous data.json.
//
// Usage: node build.mjs <raw gigs.json> <out data.json> [previous data.json or URL]

import { readFile, writeFile, appendFile } from 'node:fs/promises';

const [rawPath, outPath, previous] = process.argv.slice(2);
if (!rawPath || !outPath) {
  console.error('usage: node build.mjs <raw gigs.json> <out data.json> [previous]');
  process.exit(2);
}

const stripHtml = (s) =>
  String(s ?? '')
    .replace(/<[^>]*>/g, '')
    .replace(/&amp;/g, '&')
    .replace(/&#39;|&apos;/g, "'")
    .replace(/&quot;/g, '"')
    .replace(/&lt;/g, '<')
    .replace(/&gt;/g, '>')
    .replace(/\s+/g, ' ')
    .trim();

const connection = (raw) => (Array.isArray(raw) ? raw.map((c) => stripHtml(c.identifier)).join(', ') : stripHtml(raw));

// "dd/mm/yyyy" + "HH:MM" -> "yyyy-mm-ddTHH:MM" (festival local time, no zone).
const localIso = (part) => {
  if (!part?.date_formatted || !part?.time_formatted) return null;
  const [dd, mm, yyyy] = part.date_formatted.split('/');
  return `${yyyy}-${mm}-${dd}T${part.time_formatted}`;
};

// Share links use a short key; the last 8 hex chars of a Knack id are its
// per-process counter + random bytes, so they are unique in practice.
const shortKey = (id) => parseInt(id.slice(-8), 16).toString(36);

function normalise(raw) {
  const records = Array.isArray(raw) ? raw : raw.records;
  if (!Array.isArray(records)) throw new Error('No records array found in input');

  const gigs = records.map((r) => {
    const when = r.field_143_raw || {};
    const start = localIso(when);
    return {
      id: r.id,
      k: shortKey(r.id),
      a: connection(r.field_47_raw ?? r.field_47),
      v: stripHtml(r.field_406),
      d: stripHtml(r.field_174 || r.field_144), // festival day: after-midnight sets stay on the night before
      st: start,
      en: localIso(when.to) || start,
      t: connection(r.field_581_raw ?? r.field_581),
      x: stripHtml(r.field_48),
      u: (r.field_56_raw && r.field_56_raw.url) || (String(r.field_325).match(/href="([^"]+)"/) || [])[1] || '',
    };
  });

  const keys = new Set(gigs.map((g) => g.k));
  if (keys.size !== gigs.length) throw new Error('Short keys collided; share links would be ambiguous');
  const missing = gigs.filter((g) => !g.st || !g.a);
  if (missing.length) console.warn(`warning: ${missing.length} records missing a start time or artist`);

  gigs.sort((x, y) => (x.st || '').localeCompare(y.st || '') || x.v.localeCompare(y.v) || x.a.localeCompare(y.a));
  return gigs;
}

async function loadPrevious(src) {
  try {
    if (/^https?:/.test(src)) {
      const res = await fetch(src);
      if (!res.ok) return null;
      return await res.json();
    }
    return JSON.parse(await readFile(src, 'utf8'));
  } catch {
    return null;
  }
}

function diff(before, after) {
  const label = (g) => `**${g.a}** — ${g.d} ${g.st?.slice(11)} @ ${g.v}`;
  const old = new Map(before.map((g) => [g.id, g]));
  const now = new Map(after.map((g) => [g.id, g]));
  const added = after.filter((g) => !old.has(g.id));
  const removed = before.filter((g) => !now.has(g.id));
  const changed = after
    .filter((g) => old.has(g.id))
    .map((g) => {
      const o = old.get(g.id);
      const fields = ['a', 'v', 'st', 'en', 'x', 'u', 't'].filter((f) => o[f] !== g[f]);
      return fields.length ? { o, g, fields } : null;
    })
    .filter(Boolean);

  const names = { a: 'artist', v: 'venue', st: 'start', en: 'end', x: 'description', u: 'link', t: 'type' };
  const lines = [`### Schedule changes`, '', `${added.length} added · ${removed.length} removed · ${changed.length} changed`, ''];
  if (added.length) lines.push('#### Added', ...added.map((g) => `- ${label(g)}`), '');
  if (removed.length) lines.push('#### Removed', ...removed.map((g) => `- ${label(g)}`), '');
  if (changed.length)
    lines.push(
      '#### Changed',
      ...changed.map(({ o, g, fields }) => `- ${label(g)}: ${fields.map((f) => `${names[f]} \`${o[f]}\` → \`${g[f]}\``).join(', ')}`),
      '',
    );
  return { count: added.length + removed.length + changed.length, markdown: lines.join('\n') };
}

const gigs = normalise(JSON.parse(await readFile(rawPath, 'utf8')));
const data = { fetchedAt: new Date().toISOString(), count: gigs.length, gigs };
await writeFile(outPath, JSON.stringify(data));
console.log(`Wrote ${gigs.length} gigs to ${outPath}`);

if (previous) {
  const prev = await loadPrevious(previous);
  const summary = prev?.gigs ? diff(prev.gigs, gigs).markdown : '### Schedule changes\n\nNo previous data to compare against.';
  console.log(summary);
  if (process.env.GITHUB_STEP_SUMMARY) await appendFile(process.env.GITHUB_STEP_SUMMARY, summary + '\n');
}
