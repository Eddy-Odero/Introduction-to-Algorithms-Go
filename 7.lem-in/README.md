# Lem-in (Ant Farm Simulation)

Lem-in is a digital ant farm simulation written entirely in Go. The program reads a structural description of an ant colony (rooms, coordinates, tunnels, and a population of ants) from a file, discovers the optimal path combination to avoid traffic jams, and simulates the quickest way to move all ants from the `##start` room to the `##end` room.

## Table of Contents
- [Project Structure](#project-structure)
- [Installation & Cloning](#installation--cloning)
- [How It Works](#how-it-works)
- [Input & Output Format](#input--output-format)
- [Testing](#testing)

---

## Project Structure

The project follows a clean, modular design separating parsing, graph theory, distribution strategy, and execution tracing:

```text
lem-in/
├── cmd/lem-in/        # Main package, entry point, wires everything together
├── internal/parser/    # Phase 2-3: Reads files -> builds validated in-memory models
├── internal/graph/     # Phase 5-6: Graph representation + path-finding algorithms (BFS, multi-path)
├── internal/solver/    # Phase 7: Coordinates assignment logic to minimize total operational turns
├── internal/simulator/ # Phase 8: Turn-by-turn movement simulation + formatted standard output
├── testdata/           # Sample colony maps (valid and invalid) for testing validation states
├── go.mod
└── README.md
```

---

## Installation & Cloning

To download and run this project locally, execute the following steps in your terminal:

1. **Clone the Repository**
   ```bash
   git clone https://learn.zone01kisumu.ke/git/edwaodero/lem-in
   cd lem-in
   ```

2. **Run the Program**
   Pass any valid colony map path from the `testdata/` directory or your local files as an argument:
   ```bash
   go run ./cmd/lem-in/ testdata/valid_colony.txt
   ```

---

## How It Works

The simulation engine coordinates execution through 4 internal layers:
1. **`parser`**: Validates structural safety rules (e.g., checks duplicate room names, coordinates, structural loops, presence of `##start`/`##end`).
2. **`graph`**: Generates a network graph mapping node rooms. Employs multi-path discovery algorithms to gather combinations of non-overlapping (vertex-disjoint) path configurations.
3. **`solver`**: Evaluates queue depths dynamically to determine exactly which routes optimize the traffic flow for N total ants.
4. **`simulator`**: Tracks real-time turn stepping, outputting standard movement string records while enforcing room capacities.

---

## Input & Output Format

### Example Input File
```text
3
##start
0 1 2
##end
1 9 2
2 5 0
3 5 4
0-2
0-3
2-1
3-1
```

### Example Simulation Output
```text
[Initial map configuration printed here...]

L1-2 L2-3
L1-1 L2-1 L3-2
L3-1
```

---

## Testing

Unit tests are co-located directly inside each subsystem folder within `internal/` to independently verify data parsing accuracy, routing choices, and safety edge cases.

### Run All Tests Simultaneously
To scan the entire workspace recursively and execute every test file found across all subfolders:
```bash
go test -v ./...
```

### Test Specific Subsystems
If you are iterating on a single package layer and want focused feedback, run your tests locally inside that target directory:

* **Test the Parser Layer:**
  ```bash
  go test -v ./internal/parser/...
  ```
* **Test Path-finding Algorithms:**
  ```bash
  go test -v ./internal/graph/...
  ```
* **Test Optimization & Distribution Logic:**
  ```bash
  go test -v ./internal/solver/...
  ```
* **Test Simulation Calculations:**
  ```bash
  go test -v ./internal/simulator/...
  ```

### Verify Code Coverage
To analyze how extensively your unit tests cover your code statements within the core logic domains:
```bash
go test -cover ./internal/...
```
## Read file build time
# 1. Compile your program cleanly
```bash
go build -o lem-in ./cmd/lem-in/
```
# 2. Time the compiled binary against the map file
```bash
time ./lem-in testdata/example06.txt
```