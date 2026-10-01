# Lancaster Music Festival planner

A static site for browsing the [official schedule](https://lancastermusicfestival.com/schedule.html). It lets you search by typing, star the gigs you want, spot clashes, share your picks with a link or QR code, and export them to a calendar (`.ics`).

## How the data works

- The festival publishes no static JSON URL. Instead, `scripts/scrape.mjs` loads the schedule page in headless Chromium and captures the Knack `records` response the page requests. It then replays that request to fetch every page of results.
- `scripts/build.mjs` normalises the raw export (~940 KB, 562 gigs) into `site/data.json` (~120 KB, ~30 KB gzipped). The browser downloads the whole file once, and all searching and filtering happen client-side.
- `data/gigs.json` is a committed fallback. The site uses it if a scrape fails, or if you tick **use committed data** when running the workflow manually.

## Updating

The **Lancaster Music Festival planner** workflow (`.github/workflows/festival.yaml`) scrapes the schedule, builds the site and deploys it to GitHub Pages. It runs:

- daily at 05:17 UTC
- on demand from **Actions → Lancaster Music Festival planner → Run workflow**
- on every push to `main` that touches `festival/`

Each run's summary lists the gigs added, removed or changed since the last deploy.

To update the data by hand, download the JSON from the official site, replace `data/gigs.json`, and push to `main`.

## Local development

```sh
npm run build   # data/gigs.json -> site/data.json
npm run serve   # http://localhost:8080
```
