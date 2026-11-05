---
name: browser-check
description: Run Roommates and check it in a real browser with Playwright MCP. Use to run, screenshot or verify the app, to check a UI change, an HTMX action or the phone/desktop layout, and before calling any UI work verified.
---

# Browser check

1. [Start](#start) the app and seed it.
2. [Sign in](#sign-in) with one call.
3. [Check](#check-each-page) every page and action the change touches.
4. [Report](#report) what was and wasn't verified.
5. [Clean up](#clean-up).

## Start

Docker Desktop must be running. If `task dev` fails with `dockerDesktopLinuxEngine` in the error, stop and ask the user to start Docker Desktop. If `https://localhost:8090` already answers, the user's `task dev` is running: reuse it and never stop it. Otherwise run `task dev` in the background from `app/` (it starts db, rustfs and mailpit). Then wait with one `browser_run_code_unsafe` call (curl isn't allowed here):

```js
async (page) => {
  for (let i = 0; i < 60; i++) {
    try { const r = await page.goto('https://localhost:8090/login'); if (r?.status() === 200) return `ready after ${i * 2}s`; } catch {}
    await page.waitForTimeout(2000);
  }
  return 'not ready after 120s';
}
```

The `ERR_CONNECTION_REFUSED` entries this leaves in `browser_network_requests` are from the wait, not the app.

Then run `task seed` from `app/`. Re-running it is safe.

After editing Go, templ or `web/assets`, `restart-watch` restarts the server. Run the same wait loop and reload the page before checking.

Always use `https://localhost:8090`. The session cookie is `Secure` and `SameSite=Strict`, so don't mix hosts.

## Sign in

| User | Email | Role |
| --- | --- | --- |
| Mari | `mari@roommates.test` | owner |
| Jaan | `jaan@roommates.test` | member |

Both use the password `Roommates-dev-1`. They share the household "Tähe 3 korter", and Jaan has a payment waiting for Mari to confirm.

Sessions start clean, so sign in each time. One `browser_run_code_unsafe` call:

```js
async (page) => {
  await page.goto('https://localhost:8090/login');
  await page.fill('input[name=email]', 'mari@roommates.test');
  await page.fill('input[name=password]', 'Roommates-dev-1');
  await page.press('input[name=password]', 'Enter');
  await page.waitForURL((u) => !u.pathname.startsWith('/login'));
  return page.url();
}
```

For two-user flows, sign Jaan in on a second context in the same way. Start with `const ctx = await page.context().browser().newContext({ ignoreHTTPSErrors: true }); const jaan = await ctx.newPage();` and keep everything Jaan does inside a single call, because the other MCP tools don't see that context.

## Check each page

Use `browser_snapshot` for structure and text. Take screenshots only as evidence or for visual checks.

One call covers both viewports and both themes. It reports horizontal scroll and saves screenshots (set `name` to the page):

```js
async (page) => {
  const name = 'home', out = [];
  // without this, theme switches are captured mid-transition and dark mode looks broken
  await page.addStyleTag({ content: '*, *::before, *::after { transition: none !important; }' });
  for (const [w, h] of [[390, 844], [1440, 900]])
    for (const colorScheme of ['light', 'dark']) {
      await page.setViewportSize({ width: w, height: h });
      await page.emulateMedia({ colorScheme });
      const scroll = await page.evaluate(() => document.documentElement.scrollWidth > document.documentElement.clientWidth);
      await page.screenshot({ path: `.playwright-mcp/${name}-${w}-${colorScheme}.png`, fullPage: true });
      out.push({ w, colorScheme, scroll });
    }
  await page.setViewportSize({ width: 1440, height: 900 });
  return out;
}
```

`emulateMedia` only changes the theme while the user's mode setting is "system". Otherwise switch it in settings. Look at the screenshots for clipped text and overlapping controls. Read `browser_console_messages` and look for failed `browser_network_requests`.

Ignore the `/sw.js` certificate error and `Failed to register a ServiceWorker` on every page. Browsers won't load a service worker from a self-signed cert, even with HTTPS errors ignored, so offline and push can't be tested here.

The page language follows the browser. Headless Edge shows English, and the seed's own text is Estonian. Keep text checks tied to the locale files.

### HTMX actions

Confirm the user stays in place, the page doesn't reload, and feedback appears next to the action. To detect a reload, set `window.__stay = 1` before the action and check it survives. `hx-boost` forms can still reload, because `Redirect` sends `HX-Redirect` to HTMX requests.

### Selectors

- The bought toggle has `role="checkbox"`, so select it by `aria-label`.
- Sending an email invite also creates an invite, so the Revoke count goes up by one.
- A browser `type=email` check stops a bad address before it reaches the server, so use a valid address to see server errors.
- Member actions sit behind `summary[aria-label="Manage <name>"]`. Always scope to that person's row, never `.last()`, or you can hit the wrong member.
- Chat renders newest first in the DOM. Find a message by its text before deleting it, never by position.
- Theme and mode radios are hidden, so click their `label`.
- Removing a member revokes all invites, so create a new one before rejoining.
- `task seed` does nothing once the users exist, so it can't restore changed seed data.

## Email

Mailpit's inbox is at `http://127.0.0.1:8025`. Read verify, reset and invite links from there, and open them on `localhost:8090`. It only receives mail when `.env` sets `SMTP_HOST=localhost` and `SMTP_PORT=1025`. Otherwise emails go to the server log.

## Report

```md
**Verified in browser**
- <page or action>: 390 + 1440, light + dark, <what was exercised>
**Not checked, needs manual check**
- PWA install, push notifications, iOS Safari
- <anything skipped and why>
**Issues found**
- <issue, with screenshot path in .playwright-mcp/>
```

Only list something as verified if it was exercised in this run.

## Clean up

- Run `browser_close`.
- Stop `task dev` only if this run started it. Leave the docker services running.
- Delete screenshots in `.playwright-mcp/` that the report doesn't cite.

## Limits

- Auth routes (including `/login`) and email routes allow 10 requests a minute, and invite routes 20. Sign in once per user and reuse the session.
- Never read or print `.env`.
- The browser is headless Edge. Report the PWA, push notifications and iOS behaviour as manual checks.
- Only change seed data, never other rows in the dev database.
