# Advent of Code 2024 - Progress Report

- **Language:** Go 1.24
- **Days Completed:** 1-15

## Algorithm Implementations & Learning

### Day 1: Historian Hysteria

- *Algorithm:** Sorting + Linear Search
- *Data Structures:** Slices, Maps
- *Key Learning:** Go's `sort.Ints()` and frequency counting patterns
- *Time Complexity:** O(n log n)

### Day 2: Red-Nosed Reports

- *Algorithm:** Brute Force with Single Element Removal
- *Data Structures:** Slices
- *Key Learning:** Problem Dampener - testing all single-element removals
- *Time Complexity:** O(n²) for dampened version, O(n) for basic validation

### Day 3: Mull It Over

- *Algorithm:** Regular Expression Parsing
- *Data Structures:** Regex captures
- *Key Learning:** Go's `regexp` package for pattern extraction
- *Time Complexity:** O(n)

### Day 4: Ceres Search

- *Algorithm:** 2D Grid Search + Direction Vectors
- *Data Structures:** 2D byte slices
- *Key Learning:** Multi-directional pattern matching in grids
- *Time Complexity:** O(n*m*8\*k) where k is pattern length

### Day 5: Print Queue

- *Algorithm:** Topological Sort + Dependency Validation
- *Data Structures:** Adjacency lists, Maps
- *Key Learning:** Cycle detection and ordering constraints
- *Time Complexity:** O(V + E)

### Day 6: Guard Gallivant

- *Algorithm:** Path Simulation + Cycle Detection
- *Data Structures:** 2D grids, Set for visited states
- *Key Learning:** State tracking with position+direction tuples
- *Time Complexity:** O(n*m*4)

### Day 7: Bridge Repair

- *Algorithm:** Iterative Enumeration + Base Conversion
- *Data Structures:** Slices, bit manipulation
- *Key Learning:** Binary/ternary enumeration for operator combinations
- *Time Complexity:** O(2^n) for part 1, O(3^n) for part 2

### Day 8: Resonant Collinearity

- *Algorithm:** Coordinate Geometry + Antinode Calculation
- *Data Structures:** Maps for frequency grouping
- *Key Learning:** Mathematical relationships between points and line extensions
- *Time Complexity:** O(n²)

### Day 9: Disk Fragmenter

- *Algorithm:** Iterative Search + Block Movement
- *Data Structures:** Dynamic arrays, file mapping
- *Key Learning:** Memory defragmentation with whole-file movement constraints
- *Time Complexity:** O(n²) for part 1, O(n²-n³) for part 2

### Day 10: Hoof It

- *Algorithm:** DFS for Path Finding + Trail Validation
- *Data Structures:** Recursion stack, visited sets
- *Key Learning:** Multi-source pathfinding on elevation constraints
- *Time Complexity:** O(V + E)

### Day 11: Plutonian Pebbles

- *Algorithm:** Dynamic Programming + Memoization
- *Data Structures:** Nested maps for caching (`map[string]map[int]*big.Int`)
- *Key Learning:** Exponential growth optimization through recursive memoization with arbitrary precision arithmetic
- *Time Complexity:** O(n\*k) with memoization vs O(2^n) naive
- *Implementation Notes:** Uses `math/big` for overflow protection, two-level memoization cache

### Day 12: Garden Groups

- *Algorithm:** Connected Components + Boundary Analysis
- *Data Structures:** BFS traversal, region sets
- *Key Learning:** Dual approaches for side counting - segment merging and corner detection
- *Time Complexity:** O(n\*m)
- *Implementation Notes:**

- Part 1: Perimeter calculation during BFS
- Part 2: Corner counting using convex/concave analysis (corners = sides principle)

### Day 13: Claw Contraption

- *Algorithm:** Linear Algebra + Gaussian Elimination
- *Data Structures:** Coefficient variables for 2x2 systems
- *Key Learning:** Integer solution validation with elimination method for constraint satisfaction
- *Time Complexity:** O(1) per system
- *Implementation Notes:** Uses elimination method equivalent to Cramer's rule, handles large coordinate offsets (10^13)

### Day 14: Claw Contraption

- *Algorithm:** Linear Algebra + Gaussian Elimination
- *Data Structures:** Coefficient variables for 2x2 systems
- *Key Learning:** Integer solution validation with elimination method for constraint satisfaction
- *Time Complexity:** O(1) per system
- *Implementation Notes:** Uses elimination method equivalent to Cramer's rule, handles large coordinate offsets (10^13)

### Day 15: Warehouse Wanderer

- *Algorithm:** Grid Simulation + Component Movement
- *Data Structures:** 2D grids, Sets (maps) for connected box tracking
- *Key Learning:** Simulating complex movement rules on a dynamic grid; chain pushing of connected components; transforming - roblem representations (normal vs. wide warehouse)
- *Time Complexity:** O(n \* m \* k) where n,m = grid dimensions, k = number of moves
- *Implementation Notes:**
-
- Part 1: Simulates robot movement and box pushing in a standard grid.
- Part 2: Transforms grid to "wide" representation with multi-cell boxes and pushes entire components vertically/horizontally.
- Uses BFS-like traversal to collect movable box groups.

### Day 16: Reindeer Maze

- Algorithm: Dijkstra's Algorithm (Forward + Backward) + State Space Search
- Data Structures: Priority Queue (min-heap), 3D state tracking (x, y, direction)
- Key Learning: Bidirectional pathfinding for optimal path counting, state includes orientation for turn costs
- Time Complexity: O(V log V) where V = cells × 4 directions
- Implementation Notes: Forward dijkstra finds minimum cost backward dijkstra from end calculates reverse distances, optimal tiles identified where forward_cost + backward_cost = minimum_total

### Day 17: Chronospatial Computer

(I will be honest here. I couldn't solve the part 2 on my own. Watched tons of
video and claude to solved this part 2)

- Algorithm: 3-bit Computer Simulation + Reverse Engineering (Backtracking)
- Data Structures: Virtual Machine State (registers A,B,C + instruction pointer), Candidate Queue
- Key Learning: Quine generation through base-8 digit reconstruction, working backwards from target output
- Time Complexity: O(8^n) where n = program length, but pruned heavily by early validation
- Implementation Notes: Forward simulation executes 8 opcodes with combo/literal operands, reverse search builds register A by testing 3-bit chunks from program end to start, each iteration shifts left 3 bits and tests 8 possibilities (0-7)

### Day 18: RAM Run

I solved this one step by step - Part 1 was straightforward pathfinding, but Part 2 required optimizing the search strategy to avoid checking every single byte position.

- Algorithm: BFS Pathfinding + Binary Search Optimization
Data Structures: Grid representation via corruption map (hash set), BFS queue with state tracking
- Key Learning: Binary search transforms O(n²) sequential checking into O(n log n), critical for large datasets where naive approaches timeout
- Time Complexity: O(n log n × grid_size) where n = number of falling bytes, each binary search iteration runs BFS in O(grid_size)
- Implementation Notes: BFS explores 4-directional movement with boundary/corruption checking, binary search bisects the byte sequence to find the exact transition point from "path exists" to "path blocked", corruption tracking uses hash map for O(1) lookup of blocked coordinates

---
