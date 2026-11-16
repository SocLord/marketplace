const { spawn } = require('child_process');
const fs = require('fs');
const path = require('path');

const edgePaths = [
  'C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe',
  'C:\\Program Files\\Microsoft\\Edge\\Application\\msedge.exe'
];
const edgeExe = edgePaths.find(p => fs.existsSync(p));
if (!edgeExe) {
  console.error('Edge not found');
  process.exit(1);
}

const outDir = 'C:\\Users\\user\\Documents\\GitHub\\marketplace\\docs\\screens';
if (!fs.existsSync(outDir)) {
  fs.mkdirSync(outDir, { recursive: true });
}

const tempProfile = 'C:\\Users\\user\\AppData\\Local\\Temp\\edge_cdp_screens';
if (!fs.existsSync(tempProfile)) {
  fs.mkdirSync(tempProfile, { recursive: true });
}

const edgeProc = spawn(edgeExe, [
  '--headless=new',
  '--disable-gpu',
  '--remote-debugging-port=9222',
  `--user-data-dir=${tempProfile}`,
  '--window-size=1280,950',
  'about:blank'
]);

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

async function run() {
  await sleep(2000);
  const versionRes = await fetch('http://localhost:9222/json/version');
  const versionData = await versionRes.json();
  const browserWsUrl = versionData.webSocketDebuggerUrl;

  const listRes = await fetch('http://localhost:9222/json/list');
  const pages = await listRes.json();
  const pageWsUrl = pages[0].webSocketDebuggerUrl;

  const ws = new WebSocket(pageWsUrl);
  let id = 1;
  const pending = new Map();

  ws.onmessage = (event) => {
    const msg = JSON.parse(event.data);
    if (msg.id && pending.has(msg.id)) {
      const { resolve, reject } = pending.get(msg.id);
      pending.delete(msg.id);
      if (msg.error) reject(msg.error);
      else resolve(msg.result);
    }
  };

  await new Promise(resolve => ws.onopen = resolve);

  function send(method, params = {}) {
    return new Promise((resolve, reject) => {
      const reqId = id++;
      pending.set(reqId, { resolve, reject });
      ws.send(JSON.stringify({ id: reqId, method, params }));
    });
  }

  await send('Page.enable');
  await send('Runtime.enable');
  await send('Emulation.setDeviceMetricsOverride', {
    width: 1280,
    height: 950,
    deviceScaleFactor: 1,
    mobile: false
  });

  // Login tokens
  const buyerLoginRes = await fetch('http://localhost:8089/api/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email: 'buyer@market.tj', password: '123456' })
  });
  const buyerAuth = (await buyerLoginRes.json()).data;

  const sellerLoginRes = await fetch('http://localhost:8089/api/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email: 'seller1@market.tj', password: '123456' })
  });
  const sellerAuth = (await sellerLoginRes.json()).data;

  const modLoginRes = await fetch('http://localhost:8089/api/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email: 'moderator@market.tj', password: '123456' })
  });
  const modAuth = (await modLoginRes.json()).data;

  async function capture(url, filename, authUser = null, preScript = null) {
    console.log(`Capturing ${filename} from ${url}...`);
    await send('Page.navigate', { url: 'http://localhost:3001/' });
    await sleep(500);

    if (authUser) {
      await send('Runtime.evaluate', {
        expression: `
          localStorage.setItem('market_token', ${JSON.stringify(authUser.token)});
          localStorage.setItem('market_user', ${JSON.stringify(JSON.stringify(authUser.user))});
        `
      });
    } else {
      await send('Runtime.evaluate', {
        expression: `
          localStorage.removeItem('market_token');
          localStorage.removeItem('market_user');
        `
      });
    }

    await send('Page.navigate', { url });
    await sleep(1500);

    if (preScript) {
      await send('Runtime.evaluate', { expression: preScript });
      await sleep(500);
    }

    const { data } = await send('Page.captureScreenshot', { format: 'png' });
    const buffer = Buffer.from(data, 'base64');
    fs.writeFileSync(path.join(outDir, filename), buffer);
    console.log(`Saved ${filename} (${buffer.length} bytes)`);
  }

  // 1. Home
  await capture('http://localhost:3001/', '01_home.png');

  // 2. Catalog
  await capture('http://localhost:3001/catalog', '02_catalog.png');

  // 3. Product detail
  await capture('http://localhost:3001/product/1', '03_product_detail.png');

  // 4. Cart (add items to cart first)
  const cartScript = `
    localStorage.setItem('market_cart', JSON.stringify([
      { product_id: 3, name: 'Гӯшмонакҳои бесими Sony WH-1000XM5', price: 3450, old_price: 3900, image_url: 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=600', shop_id: 1, shop_name: 'TechStore', stock_quantity: 28, quantity: 1 },
      { product_id: 7, name: 'Кросссовкаҳои варзишии Nike Air Max Pulse', price: 980, old_price: 1250, image_url: 'https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=600', shop_id: 2, shop_name: 'ModnoStyle', stock_quantity: 30, quantity: 1 }
    ]));
  `;
  await capture('http://localhost:3001/cart', '04_cart.png', buyerAuth, cartScript);

  // 5. Checkout
  await capture('http://localhost:3001/checkout', '05_checkout.png', buyerAuth, cartScript);

  // 6. Orders
  await capture('http://localhost:3001/orders', '06_buyer_orders.png', buyerAuth);

  // 7. Shops list
  await capture('http://localhost:3001/shops', '07_shops_list.png');

  // 8. Single Shop
  await capture('http://localhost:3001/shop/techstore', '08_shop_page.png');

  // 9. Seller Dashboard
  await capture('http://localhost:3001/seller/dashboard', '09_seller_dashboard.png', sellerAuth);

  // 10. Moderator Panel
  await capture('http://localhost:3001/moderator', '10_moderator_panel.png', modAuth);

  ws.close();
  edgeProc.kill();
  console.log('All screenshots captured successfully!');
  process.exit(0);
}

run().catch(err => {
  console.error(err);
  edgeProc.kill();
  process.exit(1);
});
