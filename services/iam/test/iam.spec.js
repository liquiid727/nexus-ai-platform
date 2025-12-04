const request = require('supertest');
const express = require('express');
const appModule = require('../index');

// Since index.js starts the server immediately, we will target the running instance at localhost
const base = 'http://localhost:3001';

describe('IAM minimal flows', () => {
  it('registers email and verifies', async () => {
    const email = `user_${Date.now()}@example.com`;
    const res = await request(base).post('/auth/register-email').send({ email, password: 'Passw0rd!' });
    expect(res.status).toBe(201);
    expect(res.body.user_id).toBeTruthy();
  });
});

