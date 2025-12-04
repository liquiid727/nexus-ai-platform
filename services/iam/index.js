const express = require("express");
const cors = require("cors");
const helmet = require("helmet");
const bcrypt = require("bcryptjs");
const jwt = require("jsonwebtoken");
const { authenticator } = require("otplib");
const crypto = require("crypto");
const { v4: uuidv4 } = require("uuid");
require('dotenv').config();
const db = require('./db');
const { createProvider } = require('./sms');

const app = express();
app.use(helmet());
app.use(cors());
app.use(express.json());

const usersByEmail = new Map();
const usersByPhone = new Map();
const verifyTokens = new Map();
const smsCodes = new Map();
const sessionsByUser = new Map();
const rbacCache = new Map();

function getSecret() {
  return process.env.JWT_SECRET || "dev-secret-change";
}

app.post("/auth/register-email", async (req, res) => {
  const { email, password } = req.body || {};
  if (!email || !password) return res.status(400).send({ error: "invalid_input" });
  const hash = await bcrypt.hash(password, 10);
  const userId = uuidv4();
  if (db.pool) {
    try {
      await db.pool.query('INSERT INTO users (id, email, password_hash, verified, created_at) VALUES ($1,$2,$3,$4,$5)', [userId, email, hash, false, Date.now()]);
    } catch (e) {
      if ((e.message || '').includes('duplicate')) return res.status(409).send({ error: 'email_exists' });
      return res.status(500).send({ error: 'server_error' });
    }
    const token = crypto.randomBytes(24).toString('hex');
    await db.pool.query('INSERT INTO email_verify_tokens (token, email, created_at) VALUES ($1,$2,$3)', [token, email, Date.now()]);
    const response = { user_id: userId };
    if ((process.env.NODE_ENV || "development") === "development") response.verification_token = token;
    return res.status(201).send(response);
  }
  if (usersByEmail.has(email)) return res.status(409).send({ error: "email_exists" });
  const user = { id: userId, email, password_hash: hash, verified: false, created_at: Date.now() };
  usersByEmail.set(email, user);
  const token = crypto.randomBytes(24).toString("hex");
  verifyTokens.set(token, email);
  const response = { user_id: userId };
  if ((process.env.NODE_ENV || "development") === "development") response.verification_token = token;
  res.status(201).send(response);
});

app.get("/auth/verify-email", async (req, res) => {
  const token = req.query.token;
  if (!token) return res.status(400).send({ error: "invalid_token" });
  if (db.pool) {
    const r = await db.pool.query('SELECT email FROM email_verify_tokens WHERE token=$1', [token]);
    if (!r.rows.length) return res.status(404).send({ error: 'not_found' });
    const email = r.rows[0].email;
    await db.pool.query('UPDATE users SET verified=true WHERE email=$1', [email]);
    await db.pool.query('DELETE FROM email_verify_tokens WHERE token=$1', [token]);
    return res.send({ success: true });
  }
  const email = verifyTokens.get(token);
  if (!email) return res.status(404).send({ error: "not_found" });
  const user = usersByEmail.get(email);
  if (!user) return res.status(404).send({ error: "not_found" });
  user.verified = true;
  verifyTokens.delete(token);
  res.send({ success: true });
});

app.post("/auth/register-phone", async (req, res) => {
  const { phone } = req.body || {};
  if (!phone) return res.status(400).send({ error: "invalid_input" });
  const code = (Math.floor(Math.random() * 900000) + 100000).toString();
  if (db.pool) {
    const exists = await db.pool.query('SELECT 1 FROM users WHERE phone=$1', [phone]);
    if (exists.rows.length) return res.status(409).send({ error: 'phone_exists' });
    await db.pool.query('INSERT INTO sms_codes (phone, code, created_at) VALUES ($1,$2,$3) ON CONFLICT (phone) DO UPDATE SET code=EXCLUDED.code, created_at=EXCLUDED.created_at', [phone, code, Date.now()]);
    try { await createProvider().sendCode(phone, code); } catch {}
    const response = { sent: true };
    if ((process.env.NODE_ENV || "development") === "development") response.verification_code = code;
    return res.status(201).send(response);
  }
  if (usersByPhone.has(phone)) return res.status(409).send({ error: "phone_exists" });
  smsCodes.set(phone, { code, ts: Date.now() });
  const response = { sent: true };
  if ((process.env.NODE_ENV || "development") === "development") response.verification_code = code;
  res.status(201).send(response);
});

app.post("/auth/verify-sms", async (req, res) => {
  const { phone, code } = req.body || {};
  if (!phone || !code) return res.status(400).send({ error: "invalid_input" });
  if (db.pool) {
    const r = await db.pool.query('SELECT code FROM sms_codes WHERE phone=$1', [phone]);
    if (!r.rows.length || r.rows[0].code !== code) return res.status(401).send({ error: 'invalid_code' });
    await db.pool.query('DELETE FROM sms_codes WHERE phone=$1', [phone]);
    const existing = await db.pool.query('SELECT id FROM users WHERE phone=$1', [phone]);
    if (!existing.rows.length) {
      const userId = uuidv4();
      await db.pool.query('INSERT INTO users (id, phone, verified, phone_verified, created_at) VALUES ($1,$2,$3,$4,$5)', [userId, phone, true, true, Date.now()]);
    } else {
      await db.pool.query('UPDATE users SET phone_verified=true, verified=true WHERE phone=$1', [phone]);
    }
    return res.send({ success: true });
  }
  const entry = smsCodes.get(phone);
  if (!entry || entry.code !== code) return res.status(401).send({ error: "invalid_code" });
  smsCodes.delete(phone);
  let user = usersByPhone.get(phone);
  if (!user) {
    const userId = uuidv4();
    user = { id: userId, phone, email: null, password_hash: null, verified: true, phone_verified: true, created_at: Date.now(), twofa_enabled: false, totp_secret: null };
    usersByPhone.set(phone, user);
  } else {
    user.phone_verified = true;
    user.verified = user.verified || true;
  }
  res.send({ success: true });
});

app.post("/auth/issue-token", async (req, res) => {
  const { email, password, otp } = req.body || {};
  if (!email || !password) return res.status(400).send({ error: "invalid_input" });
  let user = null;
  if (db.pool) {
    const r = await db.pool.query('SELECT * FROM users WHERE email=$1', [email]);
    if (!r.rows.length) return res.status(404).send({ error: 'not_found' });
    user = r.rows[0];
  } else {
    user = usersByEmail.get(email);
    if (!user) return res.status(404).send({ error: "not_found" });
  }
  if (!user.verified) return res.status(403).send({ error: "email_unverified" });
  const ok = await bcrypt.compare(password, user.password_hash);
  if (!ok) return res.status(401).send({ error: "invalid_credentials" });
  if (user.twofa_enabled) {
    if (!otp || !authenticator.verify({ token: otp, secret: user.totp_secret })) {
      return res.status(401).send({ error: "invalid_otp" });
    }
  }
  const accessToken = jwt.sign({ sub: user.id, email }, getSecret(), { expiresIn: "15m" });
  const refreshToken = jwt.sign({ sub: user.id, type: "refresh" }, getSecret(), { expiresIn: "7d" });
  const device = (req.headers['user-agent'] || 'unknown').slice(0, 128);
  const ip = req.headers['x-forwarded-for'] || req.socket.remoteAddress || '';
  const sessionId = uuidv4();
  if (db.pool) {
    await db.pool.query('INSERT INTO sessions (id, user_id, device, ip, login_time) VALUES ($1,$2,$3,$4,$5)', [sessionId, user.id, device, ip, Date.now()]);
    await db.pool.query('INSERT INTO devices (id, user_id, device, first_seen, last_seen) VALUES ($1,$2,$3,$4,$5)', [uuidv4(), user.id, device, Date.now(), Date.now()]);
  } else {
    const s = sessionsByUser.get(user.id) || [];
    s.push({ id: sessionId, device, login_time: Date.now() });
    sessionsByUser.set(user.id, s);
  }
  res.send({ access_token: accessToken, refresh_token: refreshToken, expires_in: 900 });
});

app.get("/auth/sessions", async (req, res) => {
  const userId = req.query.user_id;
  if (!userId) return res.status(400).send({ error: "invalid_input" });
  if (db.pool) {
    const r = await db.pool.query('SELECT id, device, login_time FROM sessions WHERE user_id=$1 ORDER BY login_time DESC LIMIT 50', [userId]);
    return res.send({ sessions: r.rows });
  }
  const s = sessionsByUser.get(userId) || [];
  res.send({ sessions: s });
});

app.post("/auth/enable-2fa", (req, res) => {
  const { user_id } = req.body || {};
  if (!user_id) return res.status(400).send({ error: "invalid_input" });
  let user = null;
  for (const u of usersByEmail.values()) if (u.id === user_id) user = u;
  if (!user) for (const u of usersByPhone.values()) if (u.id === user_id) user = u;
  if (!user) return res.status(404).send({ error: "not_found" });
  const secret = authenticator.generateSecret();
  user.totp_secret = secret;
  user.twofa_enabled = false;
  const uri = authenticator.keyuri(user.email || user.phone || user.id, "AI_GATEWAY", secret);
  const response = { otpauth_uri: uri };
  if ((process.env.NODE_ENV || "development") === "development") response.secret = secret;
  res.send(response);
});

app.post("/auth/verify-2fa", (req, res) => {
  const { user_id, token } = req.body || {};
  if (!user_id || !token) return res.status(400).send({ error: "invalid_input" });
  let user = null;
  for (const u of usersByEmail.values()) if (u.id === user_id) user = u;
  if (!user) for (const u of usersByPhone.values()) if (u.id === user_id) user = u;
  if (!user || !user.totp_secret) return res.status(404).send({ error: "not_found" });
  const ok = authenticator.verify({ token, secret: user.totp_secret });
  if (!ok) return res.status(401).send({ error: "invalid_otp" });
  user.twofa_enabled = true;
  res.send({ success: true });
});

const port = process.env.PORT || 3001;
(async () => {
  await db.connect();
  await db.migrate();
  app.listen(port, () => {
    console.log(`IAM listening at http://localhost:${port}`);
  });
})();
