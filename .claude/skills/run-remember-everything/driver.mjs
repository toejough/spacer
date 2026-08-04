#!/usr/bin/env node
// Minimal headless-Chromium REPL for driving the Remember Everything web app.
// Usage: node driver.mjs <script-file>   (or pipe commands via stdin)
//
// Commands (one per line):
//   nav <url>
//   wait-for <selector>            (waits for element visible; text=... also supported)
//   click <selector>
//   fill <selector> <text...>
//   press <key>
//   screenshot <name>
//   console-errors
//   eval <js>
//   sleep <ms>
//
// Screenshots land in ./screenshots/<name>.png

import { chromium } from 'playwright';
import fs from 'fs';
import readline from 'readline';

const screenshotsDir = new URL('./screenshots/', import.meta.url).pathname;
fs.mkdirSync(screenshotsDir, { recursive: true });

const browser = await chromium.launch({ args: ['--no-sandbox'] });
const page = await browser.newPage();
const consoleErrors = [];
page.on('console', (msg) => { if (msg.type() === 'error') consoleErrors.push(msg.text()); });
page.on('pageerror', (err) => consoleErrors.push(String(err)));

function resolveSelector(sel) {
  if (sel.startsWith('text=')) return `text=${sel.slice(5)}`;
  return sel;
}

async function runLine(line) {
  const trimmed = line.trim();
  if (!trimmed || trimmed.startsWith('#')) return;
  const [cmd, ...rest] = trimmed.split(' ');
  const arg = rest.join(' ');
  switch (cmd) {
    case 'nav':
      await page.goto(arg, { waitUntil: 'load' });
      console.log(`[nav] ${arg}`);
      break;
    case 'wait-for':
      await page.waitForSelector(resolveSelector(arg), { timeout: 10000, state: 'visible' });
      console.log(`[wait-for] ${arg} visible`);
      break;
    case 'click':
      await page.click(resolveSelector(arg));
      console.log(`[click] ${arg}`);
      break;
    case 'fill': {
      const [sel, ...text] = rest;
      await page.fill(resolveSelector(sel), text.join(' '));
      console.log(`[fill] ${sel} = ${text.join(' ')}`);
      break;
    }
    case 'press':
      await page.keyboard.press(arg);
      console.log(`[press] ${arg}`);
      break;
    case 'screenshot': {
      const name = arg || `shot-${Date.now()}`;
      const path = `${screenshotsDir}${name}.png`;
      await page.screenshot({ path });
      console.log(`[screenshot] ${path}`);
      break;
    }
    case 'console-errors':
      console.log('[console-errors]', JSON.stringify(consoleErrors));
      break;
    case 'eval': {
      const result = await page.evaluate(arg);
      console.log('[eval]', JSON.stringify(result));
      break;
    }
    case 'sleep':
      await new Promise(r => setTimeout(r, parseInt(arg, 10)));
      break;
    // Custom high-level command for this app's drag-and-drop, since it's
    // pointer-event-based (not HTML5 DnD) — Playwright's built-in drag
    // helpers don't fire the right events. Simulates a manual
    // pointerdown -> pointermove -> pointerup sequence between two
    // elements' centers.
    case 'drag': {
      const [fromSel, toSel] = rest;
      const from = await page.$(resolveSelector(fromSel));
      const to = await page.$(resolveSelector(toSel));
      if (!from || !to) throw new Error(`drag: selector not found (from=${fromSel} to=${toSel})`);
      const fb = await from.boundingBox();
      const tb = await to.boundingBox();
      const fx = fb.x + fb.width / 2, fy = fb.y + fb.height / 2;
      const tx = tb.x + tb.width / 2, ty = tb.y + tb.height / 2;
      await page.mouse.move(fx, fy);
      await page.mouse.down();
      await page.mouse.move(fx + 10, fy + 10, { steps: 3 });
      await page.mouse.move(tx, ty, { steps: 8 });
      await page.mouse.up();
      console.log(`[drag] ${fromSel} -> ${toSel}`);
      break;
    }
    default:
      console.log(`[?] unknown command: ${cmd}`);
  }
}

const scriptFile = process.argv[2];
const lines = scriptFile
  ? fs.readFileSync(scriptFile, 'utf8').split('\n')
  : fs.readFileSync(0, 'utf8').split('\n');

for (const line of lines) {
  try {
    await runLine(line);
  } catch (err) {
    console.error(`[error] on line "${line}":`, err.message);
  }
}

await browser.close();
