# Forum Project — Team Roadmap (4 members)

9 phases as plan, order — server, static frontend,
database, then features. Each phase below is split into 4 tasks so everyone
has something to own without waiting on each other. Roles are labeled
**A/B/C/D**.

Phases 1 and 2 can run **at the same time** — the server work doesn't
depend on the HTML existing yet, and the HTML doesn't depend on the server
being finished. Everything from Phase 3 onward is sequential (each phase
needs the previous one's tables/endpoints).

## Target file structure

## Target file structure (grows phase by phase — see notes on each file)
 
```
forum/
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
├── README.md
├── cmd/
│   └── server/
│       └── main.go                # Phase 1, extended every phase
├── internal/
│   ├── router/
│   │   └── router.go              # Phase 1
│   ├── middleware/
│   │   └── middleware.go          # Phase 1 (logging), Phase 4 (auth guard), Phase 9 (errors)
│   ├── db/
│   │   ├── db.go                  # Phase 3
│   │   └── schema.sql             # Phase 3
│   ├── models/
│   │   ├── user.go                # Phase 3/4
│   │   ├── category.go            # Phase 3
│   │   ├── post.go                # Phase 5
│   │   ├── comment.go             # Phase 6
│   │   └── reaction.go            # Phase 7
│   ├── auth/
│   │   ├── password.go            # Phase 4
│   │   ├── session.go             # Phase 4
│   │   └── handlers.go            # Phase 4
│   └── handlers/
│       ├── posts.go               # Phase 5
│       ├── comments.go            # Phase 6
│       ├── reactions.go           # Phase 7
│       └── filters.go             # Phase 8
├── web/
│   ├── static/
│   │   ├── css/style.css          # Phase 2
│   │   └── js/
│   │       ├── main.js            # Phase 2 (nav/shared), Phase 5+ (fetch calls)
│   │       ├── auth.js            # Phase 4
│   │       ├── posts.js           # Phase 5
│   │       ├── comments.js        # Phase 6
│   │       └── reactions.js       # Phase 7
│   └── pages/
│       ├── index.html             # Phase 2
│       ├── login.html             # Phase 2
│       ├── register.html          # Phase 2
│       ├── new-post.html          # Phase 2
│       ├── post.html              # Phase 2
│       ├── category.html          # Phase 2
│       ├── profile.html           # Phase 2
│       └── 404.html               # Phase 1
└── (tests live alongside the code they test: internal/models/user_test.go, etc.)
```
---

## Phase 1 — Static HTTP Server *(runs in parallel with Phase 2)*

**Goal:** Server that serves static files, nothing else yet.

| Owner | Task |
|---|---|
| **A** | `cmd/server/main.go` — server startup, graceful shutdown (SIGINT/SIGTERM). |
| **B** | `internal/router/router.go` — mux setup, static file serving under `/static/`, `/` → `index.html`, custom `404.html` handler. |
| **C** | `internal/middleware/middleware.go` — logging middleware (method, path, status, duration). |
| **D** | `Dockerfile` (early skeleton — just enough to `docker build`/`run` this stage) + root `README.md` with setup/run instructions, so the team has a working environment from day one. |

**Test before merging:** everyone pulls, runs `go run ./cmd/server`, confirms
`curl -i localhost:8080/` → 200, `curl -i localhost:8080/nonsense` → 404 with
the custom page, and log lines appear per request.

---

## Phase 2 — Static Frontend Pages *(runs in parallel with Phase 1)*

**Goal:** Every page exists, styled, navigable, with fake/hardcoded data.

| Owner | Task |
|---|---|
| **A** | Shared layout: nav bar, footer, base `style.css` (colors, typography, spacing — everyone else's pages build on this) + `index.html` feed page. |
| **B** | `login.html` + `register.html` — forms only, no submit logic. |
| **C** | `new-post.html` + `post.html` (post detail, comment list mockup, comment form, like/dislike buttons). |
| **D** | `category.html` (subforum view) + `profile.html` ("My posts" / "My liked posts" tabs). |

**Test before merging:** click through every page from the nav, no dead
links, consistent look (this depends on A's shared CSS landing first — A
should push that early so B/C/D can build on it).

**Sync point:** as a team, review all pages together before Phase 3 — this
is where you finalize what fields each entity needs (does a post need tags
beyond category? do comments need a parent-comment field for replies?).
Whatever you decide here becomes the Phase 3 schema.

---

## Phase 3 — Database Schema & Connection Layer

**Goal:** Schema matching what Phase 2 actually needs.

| Owner | Task |
|---|---|
| **A** | ERD + `internal/db/schema.sql` — all tables, foreign keys, unique indexes (own the schema design, but get sign-off from the team since everyone builds on it). |
| **B** | `internal/db/db.go` — connection setup, pragmas (WAL, foreign_keys, busy_timeout), idempotent migration runner. |
| **C** | Seed data (categories) + `internal/models/category.go`. |
| **D** | `internal/models/user.go` skeleton — struct + `CreateUser`/`GetUserByEmail`, since Phase 4's owner needs this ready first. |

**Test before merging:** run the server, inspect the `.db` file with
`sqlite3` CLI, confirm every table + seed data exists; D's temporary
insert/select route (delete after) proves the round trip works.

---

## Phase 4 — Authentication

**Goal:** Real register/login wired to the Phase 2 forms.

| Owner | Task |
|---|---|
| **A** | `internal/auth/password.go` (bcrypt) + `internal/auth/handlers.go` register endpoint (duplicate email → 409). |
| **B** | `internal/auth/session.go` (UUID sessions, expiry, cookie helpers) + login/logout endpoints. |
| **C** | `internal/middleware/middleware.go` auth-guard (reads cookie, validates session, attaches user to context, 401 on failure). |
| **D** | `web/static/js/auth.js` — wire `login.html`/`register.html` to the API, inline error display, real logged-in/out nav state. |

**Test before merging:** register → cookie set; duplicate email → 409;
wrong password → 401; guarded route rejects an expired/deleted session;
logout clears cookie and nav.

---

## Phase 5 — Posts & Categories

| Owner | Task |
|---|---|
| **A** | `internal/models/post.go` — create (+ `post_categories` insert), get-by-id with categories joined. |
| **B** | `internal/models/post.go` (list/feed query) + `internal/handlers/posts.go` — `GET /api/posts`, `GET /api/posts/{id}` (public). |
| **C** | `internal/handlers/posts.go` — `POST /api/posts` (auth required, category validation). |
| **D** | `web/static/js/posts.js` — wire `new-post.html` submit, `index.html` feed render, `post.html` detail render to the real API. |

**Test before merging:** create a post with 2 categories as a logged-in
user → shows in feed with both tags; logged-out create attempt → 401;
direct post URL shows correct content.

---

## Phase 6 — Comments

| Owner | Task |
|---|---|
| **A** | `internal/models/comment.go` — create, list-by-post. |
| **B** | `internal/handlers/comments.go` — `POST /api/posts/{id}/comments` (auth), `GET /api/posts/{id}/comments` (public). |
| **C** | `web/static/js/comments.js` — wire `post.html` comment form + list to the API. |
| **D** | Comment count consistency: make sure `index.html` feed and `post.html` both show accurate counts; write the model query/join needed for this and coordinate with A if it touches `comment.go`. |

**Test before merging:** comment as logged-in user → appears immediately;
guest sees comments, no post form (or disabled with login prompt); feed
comment counts match actual counts.

---

## Phase 7 — Likes & Dislikes

| Owner | Task |
|---|---|
| **A** | `internal/models/reaction.go` — toggle logic (like ↔ dislike ↔ remove) using the Phase 3 unique index. |
| **B** | `internal/handlers/reactions.go` — `POST /api/posts/{id}/react`, `POST /api/comments/{id}/react` (auth required). |
| **C** | `web/static/js/reactions.js` — like/dislike buttons, optimistic UI update without full reload. |
| **D** | Cross-cutting QA: verify counts are correct and visible to logged-out users on every page they appear (feed, post detail, comments) — this touches multiple owners' work, so D's job this phase is integration testing + fixing small display bugs across A/B/C's pieces. |

**Test before merging:** like → count up, active state shown; dislike
while like active → like drops, dislike rises (not both); click again →
removed; guest action → 401; refresh → state persists.

---

## Phase 8 — Filters

| Owner | Task |
|---|---|
| **A** | `GET /api/posts?category={id}` (public) in `internal/handlers/filters.go` or extended `posts.go`. |
| **B** | `GET /api/posts?mine=true` (auth required, scoped to logged-in user). |
| **C** | `GET /api/posts?liked=true` (auth required, scoped to logged-in user via the reactions table). |
| **D** | Frontend: wire `category.html` to A's filter, `profile.html` tabs to B's and C's filters. |

**Test before merging:** category filter returns only that category's posts;
"my posts"/"my liked posts" as guest → 401, not a silently empty list;
as two different logged-in test accounts, each only sees their own scoped
results.

---

## Phase 9 — Error Handling, Tests, Docker, Deploy

| Owner | Task |
|---|---|
| **A** | Audit every handler from Phases 4–8 for consistent HTTP status codes (400/401/403/404/405/500) — this is a cross-cutting review, expect to touch other people's files with their sign-off. |
| **B** | Unit tests: `internal/models/*_test.go` (repository logic) + `internal/auth/*_test.go` (password/session logic). |
| **C** | Finalize `Dockerfile` (multi-stage, cgo-aware, `debian-slim` + `libsqlite3-0` runtime) + `docker-compose.yml` with a named volume. |
| **D** | Final `README.md` (setup, run, env vars, architecture overview) + manual QA pass hitting every error case across the whole app (bad login, duplicate register, unauthenticated actions, bad IDs, wrong HTTP methods). |

**Test before calling it done:** `go test ./...` passes; `docker build` +
`docker run` serves the whole app; data survives a container restart;
D's manual QA pass finds no unhandled 500s where a clean 4xx should be.

---

## A note on parallelism after Phase 2

From Phase 3 onward, phases are sequential (each depends on tables/endpoints
from the last), but *within* a phase the 4 tasks are designed to touch
different files so people aren't blocked on each other mid-phase. The one
recurring dependency: whoever owns the model/repository layer (usually A)
should land their piece first in each phase, since the handler and frontend
owners build on top of it — flag this in standup so it doesn't become a
bottleneck.