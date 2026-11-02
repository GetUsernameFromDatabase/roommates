---
name: browser-check
description: Run Roommates and check it in a real browser with Playwright MCP. Use to run, screenshot or verify the app, and before calling any UI work verified.
---

# Browser check

## Start

Docker Desktop must be running. From the repo root:

```sh
docker compose up -d --wait db rustfs mailpit
```

If port 8090 is already listening, the user's `task dev` is probably running, so reuse it and don't stop it. Otherwise run `task dev` in the background from `app/` and wait for the port. Then run `task seed`. Re-running the seed is safe.

Always use `https://localhost:8090`. The session cookie is `Secure` and `SameSite=Strict`, so don't mix hosts.

## Sign in

| User | Email | Role |
| --- | --- | --- |
| Mari | `mari@roommates.test` | owner |
| Jaan | `jaan@roommates.test` | member |

Both use the password `Roommates-dev-1`. They share the household "Tähe 3 korter", and Jaan has a payment waiting for Mari to confirm.

Sessions start clean, so sign in each time. For two-user flows, open a second context for Jaan with `browser_run_code_unsafe` (`browser.newContext({ ignoreHTTPSErrors: true })`).

## Check each page

1. Use both viewports: 390x844 and 1440x900 (`browser_resize`).
2. Use both themes (`browser_emulate_media`, or the setting the page offers).
3. Read `browser_console_messages` and look for failed `browser_network_requests`.
4. Look for horizontal scroll, clipped text and overlapping controls.

For HTMX actions, confirm the user stays in place, the page doesn't reload, and feedback appears next to the action.

Ignore the `/sw.js` certificate error and `Failed to register a ServiceWorker` on every page. Browsers won't load a service worker from a self-signed cert, even with HTTPS errors ignored, so offline and push can't be tested here.

The page language follows the browser. Headless Edge shows English, and the seed's own text is Estonian. Keep text checks tied to the locale files.

To detect a reload, set `window.__stay = 1` before the action and check it survives. `hx-boost` forms can still reload, because `Redirect` sends `HX-Redirect` to HTMX requests.

Selector notes:

- The bought toggle has `role="checkbox"`, so select it by `aria-label`.
- Sending an email invite also creates an invite, so the Revoke count goes up by one.
- A browser `type=email` check stops a bad address before it reaches the server, so use a valid address to see server errors.
- Member actions sit behind `summary[aria-label="Manage <name>"]`. Always scope to that person's row, never `.last()`, or you can hit the wrong member.
- Chat renders newest first in the DOM. Find a message by its text before deleting it, never by position.
- Theme and mode radios are hidden, so click their `label`.
- Removing a member revokes all invites, so create a new one before rejoining.
- `task seed` does nothing once the users exist, so it can't restore changed seed data.

Save evidence in `.playwright-mcp/`.

## Email

Mailpit's inbox is at `http://127.0.0.1:8025`. Read verify, reset and invite links from there, and open them on `localhost:8090`. It only receives mail when `.env` sets `SMTP_HOST=localhost` and `SMTP_PORT=1025`. Otherwise emails go to the server log.

## Limits

- Auth and email routes allow 10 requests a minute, and invite routes 20.
- Never read or print `.env`.
- The browser is headless Edge. Report the PWA, push notifications and iOS behaviour as manual checks.
- Only change seed data, never other rows in the dev database.
