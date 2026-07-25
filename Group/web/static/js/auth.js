/* ==========================================================
   EDU-FLIX — auth.js
   Talks to your Go backend for signup/login/password-reset.
   Include this on index.html (login), signup.html, and
   reset-password.html. Other pages can use authGuard() below
   to require login.
========================================================== */

// Now that Go serves both the API and the frontend on the same origin,
// a relative path works everywhere — locally and once deployed.
// No CORS issues, nothing to change when you move to Render.
const API_BASE = '/api';

// ── TOKEN STORAGE ───────────────────────────────────────
function saveSession(token, user) {
  localStorage.setItem('eduflix_token', token);
  localStorage.setItem('eduflix_user', JSON.stringify(user));
}

function getToken() {
  return localStorage.getItem('eduflix_token');
}

function getUser() {
  const raw = localStorage.getItem('eduflix_user');
  return raw ? JSON.parse(raw) : null;
}

function clearSession() {
  localStorage.removeItem('eduflix_token');
  localStorage.removeItem('eduflix_user');
}

// ── API CALLS ───────────────────────────────────────────

async function signup(name, email, password) {
  const res = await fetch(`${API_BASE}/signup`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name, email, password })
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || 'Signup failed');
  saveSession(data.token, data.user);
  return data.user;
}

async function login(email, password) {
  const res = await fetch(`${API_BASE}/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || 'Login failed');
  saveSession(data.token, data.user);
  return data.user;
}

async function forgotPassword(email) {
  const res = await fetch(`${API_BASE}/forgot-password`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email })
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || 'Request failed');
  return data.message;
}

async function resetPassword(token, newPassword) {
  const res = await fetch(`${API_BASE}/reset-password`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ token, new_password: newPassword })
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || 'Reset failed');
  return data.message;
}

async function fetchMe() {
  const token = getToken();
  if (!token) return null;

  const res = await fetch(`${API_BASE}/me`, {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  if (!res.ok) {
    clearSession();
    return null;
  }
  return await res.json();
}

function logout() {
  clearSession();
  window.location.href = 'index.html';
}

// ── PAGE GUARD ──────────────────────────────────────────
// Call this at the top of any protected page (e.g. main.html)
// to redirect unauthenticated users back to login.
async function authGuard() {
  const token = getToken();
  if (!token) {
    window.location.href = 'index.html';
    return null;
  }
  const user = await fetchMe();
  if (!user) {
    window.location.href = 'index.html';
    return null;
  }
  return user;
}

// ── ADMIN GUARD ─────────────────────────────────────────
// Call this on admin.html to block non-admins.
async function adminGuard() {
  const user = await authGuard();
  if (!user) return null;
  if (user.role !== 'admin') {
    alert('Admin access required.');
    window.location.href = 'main.html';
    return null;
  }
  return user;
}

// ── AUTHENTICATED FETCH HELPER ──────────────────────────
// Use this for any call that needs the Bearer token attached.
async function authFetch(url, options = {}) {
  const token = getToken();
  const headers = {
    ...(options.headers || {}),
    'Authorization': `Bearer ${token}`
  };
  return fetch(url, { ...options, headers });
}
