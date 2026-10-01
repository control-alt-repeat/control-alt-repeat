// Loads the official schedule page in a headless browser and captures the gig
// records it pulls from Knack. There is no static JSON URL, so we watch the
// page's own network traffic, then replay the records request ourselves to
// page through every result.
//
// Usage: node scrape.mjs <out gigs.json>
// Env:   SCHEDULE_URL (default https://lancastermusicfestival.com/schedule.html#json)
//        MIN_RECORDS  (default 100) – fewer than this is treated as a failed scrape

import { readFile, writeFile } from 'node:fs/promises';
import { chromium } from 'playwright';

const outPath = process.argv[2];
if (!outPath) {
  console.error('usage: node scrape.mjs <out gigs.json>');
  process.exit(2);
}
const scheduleUrl = process.env.SCHEDULE_URL || 'https://lancastermusicfestival.com/schedule.html#json';
const minRecords = Number(process.env.MIN_RECORDS || 100);

const looksLikeGigs = (body) =>
  Array.isArray(body?.records) && body.records.length > 0 && body.records.some((r) => r && 'field_143_raw' in r);

const found = new Map(); // record id -> record
const add = (records) => records.forEach((r) => r?.id && found.set(r.id, r));
let apiRequest = null; // { url, headers } of a Knack records call we can replay

const browser = await chromium.launch();
const context = await browser.newContext({ acceptDownloads: true });
const page = await context.newPage();

page.on('response', async (res) => {
  if (!/json|javascript/.test(res.headers()['content-type'] || '')) return;
  let body;
  try {
    body = await res.json();
  } catch {
    return;
  }
  if (!looksLikeGigs(body)) return;
  add(body.records);
  console.log(`captured ${body.records.length} records from ${res.url()}`);
  if (!apiRequest && /\/records/.test(res.url())) {
    apiRequest = { url: res.url(), headers: await res.request().allHeaders() };
  }
});

// Some pages build the JSON in the browser and offer it as a download.
page.on('download', async (download) => {
  try {
    const path = await download.path();
    const body = JSON.parse(await readFile(path, 'utf8'));
    const records = Array.isArray(body) ? body : body.records;
    if (looksLikeGigs({ records })) {
      add(records);
      console.log(`captured ${records.length} records from download ${download.suggestedFilename()}`);
    }
  } catch (err) {
    console.warn(`ignored download: ${err.message}`);
  }
});

console.log(`loading ${scheduleUrl}`);
await page.goto(scheduleUrl, { waitUntil: 'networkidle', timeout: 90_000 });

// Give lazy views a moment, and nudge infinite-scroll lists.
for (let i = 0; i < 10 && !apiRequest; i++) {
  await page.mouse.wheel(0, 5000);
  await page.waitForTimeout(1500);
}

if (apiRequest) {
  // Replay with a big page size and walk every page.
  const headers = Object.fromEntries(
    Object.entries(apiRequest.headers).filter(([k]) => /^x-knack|^authorization$|^accept$|^content-type$/i.test(k)),
  );
  for (let pageNo = 1, totalPages = 1; pageNo <= totalPages; pageNo++) {
    const url = new URL(apiRequest.url);
    url.searchParams.set('rows_per_page', '1000');
    url.searchParams.set('page', String(pageNo));
    const res = await context.request.get(url.toString(), { headers });
    if (!res.ok()) {
      console.warn(`replay of ${url} failed: ${res.status()}`);
      break;
    }
    const body = await res.json();
    if (!looksLikeGigs(body)) break;
    add(body.records);
    totalPages = Number(body.total_pages || 1);
    console.log(`replayed page ${pageNo}/${totalPages}: ${body.records.length} records`);
  }
}

await browser.close();

if (found.size < minRecords) {
  console.error(`Only found ${found.size} records (expected at least ${minRecords}); not overwriting data.`);
  process.exit(1);
}
await writeFile(outPath, JSON.stringify({ records: [...found.values()] }));
console.log(`Wrote ${found.size} records to ${outPath}`);
