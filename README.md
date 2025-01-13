# Token Pool Manager

A highly concurrent, efficient token pool management system implemented in Go. This system manages a pool of 1000 tokens, handling their distribution and usage tracking with automatic 24-hour resets.

## Project Overview

This system simulates a token pool manager that could be used in real-world scenarios like:
- Load balancers distributing traffic across servers
- Connection pool managers
- Resource allocation systems
- Rate limiters

## Key Go Concepts Used

### 1. Concurrent Programming
- **sync.RWMutex**: Used for thread-safe operations when accessing and modifying token states
  ```go
  type TokenPool struct {
      mutex     sync.RWMutex
      // ...
  }
  ```
- **Lock Patterns**: Implementation of read/write locks to optimize concurrent access
  - Read locks for statistics gathering
  - Write locks for token updates

### 2. Structs and Methods
- **Custom Types**: 
  ```go
  type Token struct {
      ID    int
      Usage int
  }
  ```
- **Method Receivers**: Implementation of methods on TokenPool type
  ```go
  func (tp *TokenPool) SelectToken() *Token
  ```

### 3. Slices and Maps
- **Dynamic Arrays**: Using slices for token storage
  ```go
  tokens    []Token
  ```
- **Hash Maps**: For efficient usage statistics tracking
  ```go
  usageMap := make(map[int]int)
  ```

### 4. Random Number Generation
- **math/rand**: Implementing weighted random selection for fair token distribution
  ```go
  selectedIndex := candidates[rand.Intn(len(candidates))]
  ```

### 5. Time Management
- **time.Time**: Handling 24-hour reset functionality
  ```go
  lastReset time.Time
  ```

### 6. Error Handling
- **Defensive Programming**: Proper handling of edge cases and concurrent scenarios

## Why This Approach?

1. **Concurrency Safety**: 
   - RWMutex ensures thread-safe operations
   - Separate locks for read/write operations optimize performance

2. **Memory Efficiency**:
   - Pre-allocated slices minimize memory reallocations
   - Efficient data structures for token tracking

3. **Fair Distribution**:
   - Weighted random selection ensures balanced token usage
   - 30% probability for token reuse prevents starvation

4. **Performance Optimization**:
   - O(1) token access using slices
   - Efficient statistics collection using maps

## Getting Started

### Prerequisites
- Go 1.16 or higher
- Git (for cloning the repository)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd token-manager
```

2. Initialize Go module:
```bash
go mod init token-manager
```

3. Run the application:
```bash
go run main.go
```

### Usage

1. When prompted, enter the number of operations to simulate:
```bash
Enter simulation operations count: 1000
```

2. The program will display:
   - Operation statistics
   - Token usage distribution
   - Performance metrics
   - Least used tokens

### Output

![image](https://github.com/user-attachments/assets/933bff6c-0586-46b6-85ee-6631e0d6bb65)

![image](https://github.com/user-attachments/assets/60667694-adcb-4344-9b67-046727fd15bd)



## Design Decisions

1. **Token Selection Strategy**:
   - Uses weighted random selection
   - Balances between new and previously used tokens
   - Prevents both clustering and starvation

2. **Reset Mechanism**:
   - Automatic 24-hour reset
   - Thread-safe implementation
   - Minimal performance impact

3. **Statistics Collection**:
   - Real-time usage tracking
   - Efficient data structures
   - Comprehensive reporting



