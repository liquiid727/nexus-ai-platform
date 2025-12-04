CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  email TEXT UNIQUE,
  phone TEXT UNIQUE,
  password_hash TEXT,
  verified BOOLEAN DEFAULT FALSE,
  phone_verified BOOLEAN DEFAULT FALSE,
  created_at BIGINT
);

CREATE TABLE IF NOT EXISTS email_verify_tokens (
  token TEXT PRIMARY KEY,
  email TEXT,
  created_at BIGINT
);

CREATE TABLE IF NOT EXISTS sms_codes (
  phone TEXT PRIMARY KEY,
  code TEXT,
  created_at BIGINT
);

CREATE TABLE IF NOT EXISTS sessions (
  id UUID PRIMARY KEY,
  user_id UUID,
  device TEXT,
  ip TEXT,
  login_time BIGINT
);

CREATE TABLE IF NOT EXISTS devices (
  id UUID PRIMARY KEY,
  user_id UUID,
  device TEXT,
  first_seen BIGINT,
  last_seen BIGINT
);

CREATE TABLE IF NOT EXISTS roles (
  id SERIAL PRIMARY KEY,
  name TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS permissions (
  id SERIAL PRIMARY KEY,
  name TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS user_roles (
  user_id UUID,
  role_id INT
);

CREATE TABLE IF NOT EXISTS role_permissions (
  role_id INT,
  permission_id INT
);

