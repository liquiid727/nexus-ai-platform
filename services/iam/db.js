const { Pool } = require('pg');
const fs = require('fs');
const path = require('path');

let pool = null;

function getDatabaseUrl() {
  return process.env.DATABASE_URL || 'postgres://postgres:example@localhost:5432/postgres';
}

async function connect() {
  try {
    pool = new Pool({ connectionString: getDatabaseUrl() });
    await pool.query('SELECT 1');
    return pool;
  } catch (e) {
    console.error('DB connect failed, fallback to memory:', e.message);
    pool = null;
    return null;
  }
}

async function migrate() {
  if (!pool) return;
  const sql = fs.readFileSync(path.join(__dirname, '../../db/migrations/001_init.sql'), 'utf8');
  await pool.query(sql);
}

module.exports = { connect, migrate, get pool() { return pool; } };

