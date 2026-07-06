# lem-in

A digital ant farm: read a colony description (rooms + tunnels + ant count)
from a file, compute the fastest way to move all ants from `##start` to
`##end`, and print each turn's moves.

Built step-by-step as a learning project. Each phase below gets implemented,
tested, then documented with findings/issues before moving to the next.

## Goal

Given a file describing:
- number of ants
- rooms (`name x y`)
- links between rooms (`name1-name2`)
- one `##start` and one `##end` room

Output:
```
number_of_ants
the_rooms
the_links

Lx-y Lz-w Lr-o ...   <- turn 1
Lx-y Lz-w ...         <- turn 2
...
```
Where each `Lx-y` means "ant x moves to room y" this turn.

If the input is malformed in any way (bad ant count, no start/end, duplicate
rooms, links to unknown rooms, cycles that trap parsing, etc.) print:
```
ERROR: invalid data format
```
(optionally with a more specific suffix).

## Project structure

```
lem-in/
├── cmd/lem-in/        # main package, entry point, wires everything together
├── internal/parser/    # Phase 2-3: read file -> validated in-memory model
├── internal/graph/     # Phase 5-6: graph representation + path-finding (BFS, multi-path)
├── internal/solver/    # Phase 7: decide which ants take which path, minimize turns
├── internal/simulator/ # Phase 8: turn-by-turn movement simulation + formatted output
├── testdata/           # sample colony files, valid and invalid, for manual/automated testing
├── go.mod
└── README.md           # this file
```

## Phases

- [x] **Phase 1** — Project skeleton: read filename from args, read file contents
- [x] **Phase 2** — Data structures: `Room`, graph adjacency representation
- [ ] **Phase 3** — Parsing: lines -> rooms, links, start/end, ant count, skip comments
- [ ] **Phase 4** — Validation: all the ways input can be invalid -> `ERROR: invalid data format`
- [ ] **Phase 5** — BFS: find a single shortest path start -> end
- [ ] **Phase 6** — Multiple shortest/non-overlapping paths (max-flow style thinking)
- [ ] **Phase 7** — Assign ants to paths to minimize total turns
- [ ] **Phase 8** — Turn-by-turn simulation + output formatting matching spec

## Findings & Issues Log

Notes, bugs, and "aha" moments get appended here after each phase, in order.

### Phase 1
- `os.Args` check must be `!= 2` (exactly one arg expected), not `< 2` — otherwise
  extra stray arguments are silently ignored.
- Every error path must `return` (or `os.Exit`) immediately — printing an error
  and letting execution fall through to the next line still runs that next line
  (e.g. printing an empty string after a failed read).
- Error output should use the spec's required prefix `ERROR: invalid data format`
  rather than ad-hoc Go error strings — raw Go errors (`open x: no such file...`)
  leak internal detail and don't match spec.
- Confirmed: reading + `strings.Split` on `"\n"` correctly preserves comment
  lines (`#comment`) and section markers (`##start`, `##end`) as plain lines —
  parsing/filtering these is deferred to Phase 3, not handled here.
- Verified against `testdata/sample.txt` (spec's example colony): all 24 lines
  print correctly, including both comment lines and both markers.

### Phase 2
- `Room.Links []*Room` stores pointers, not copies — two rooms linking to the
  same neighbor share the exact same `*Room` in memory. Important once ant
  simulation needs to check/mark room occupancy.
- Deliberately did NOT put a `Visited` flag on `Room`. Rooms get reused across
  multiple path searches (Phase 6), so per-search state must live in a
  separate map created fresh each search, not as a permanent struct field.
- `Colony.Rooms` is a `map[string]*Room` for O(1) name lookup while linking
  rooms together during parsing.
- `AddRoom`/`AddLink` return `bool` (not error) to signal duplicate-room /
  unknown-room-in-link conditions — parser (Phase 3/4) decides how to turn
  that into an `ERROR: invalid data format` message.
