package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type Token struct {
    ID    int
    Usage int
}

type TokenPool struct {
    tokens    []Token
    mutex     sync.RWMutex
    lastReset time.Time
}

func NewTokenPool(size int) *TokenPool { 
    tokens := make([]Token, size)
    for i := range tokens {
        tokens[i] = Token{ID: i + 1, Usage: 0}
    }
    return &TokenPool{
        tokens:    tokens,
        lastReset: time.Now(),
    }
}

func (tp *TokenPool) checkAndResetIfNeeded() {
    if time.Since(tp.lastReset) >= 24*time.Hour {
        tp.mutex.Lock()
        defer tp.mutex.Unlock()
        
        for i := range tp.tokens {
            tp.tokens[i].Usage = 0
        }
        tp.lastReset = time.Now()
    }
}

// Modified SelectToken to better distribute usage
func (tp *TokenPool) SelectToken() *Token {
    tp.checkAndResetIfNeeded()
    
    tp.mutex.Lock()
    defer tp.mutex.Unlock()

    // Find tokens with minimum usage
    minUsage := tp.tokens[0].Usage
    for _, t := range tp.tokens {
        if t.Usage < minUsage {
            minUsage = t.Usage
        }
    }

    // Get all tokens with minimum usage
    var candidates []int
    for i, t := range tp.tokens {
        if t.Usage == minUsage {
            candidates = append(candidates, i)
        }
    }

    // Randomly select from candidates with weighted probability
    var selectedIndex int
    if len(candidates) > 1 && rand.Float32() < 0.3 { // 30% chance to reuse a token
        // Select from tokens that have been used before
        usedTokens := make([]int, 0)
        for i, t := range tp.tokens {
            if t.Usage > 0 {
                usedTokens = append(usedTokens, i)
            }
        }
        if len(usedTokens) > 0 {
            selectedIndex = usedTokens[rand.Intn(len(usedTokens))]
        } else {
            selectedIndex = candidates[rand.Intn(len(candidates))]
        }
    } else {
        selectedIndex = candidates[rand.Intn(len(candidates))]
    }

    tp.tokens[selectedIndex].Usage++
    return &tp.tokens[selectedIndex]
}

func (tp *TokenPool) GetUsageStats() (map[int]int, []Token) {
    tp.mutex.RLock()
    defer tp.mutex.RUnlock()

    usageMap := make(map[int]int)
    for _, t := range tp.tokens {
        usageMap[t.ID] = t.Usage
    }

    // Find least used tokens
    minUsage := tp.tokens[0].Usage
    for _, t := range tp.tokens {
        if t.Usage < minUsage {
            minUsage = t.Usage
        }
    }

    leastUsed := make([]Token, 0)
    for _, t := range tp.tokens {
        if t.Usage == minUsage {
            leastUsed = append(leastUsed, t)
        }
    }

    sort.Slice(leastUsed, func(i, j int) bool {
        return leastUsed[i].ID < leastUsed[j].ID
    })

    return usageMap, leastUsed
}

func (tp *TokenPool) SimulateOperations(count int) {
    for i := 0; i < count; i++ {
        tp.SelectToken()
    }
}

func main() {
    rand.Seed(time.Now().UnixNano())
    pool := NewTokenPool(1000)

    var operations int
    fmt.Print("Enter simulation operations count: ")
    fmt.Scan(&operations)

    startTime := time.Now()
    pool.SimulateOperations(operations)
    duration := time.Since(startTime)

    usageMap, leastUsed := pool.GetUsageStats()

    fmt.Printf("\nSimulation Time: %d operations (completed in %v)\n", operations, duration)
    
    // Print all token usages
    maxUsage := 0
    totalUsage := 0
    usedTokens := 0
    
    // Get all used tokens
    type TokenUsage struct {
        ID    int
        Usage int
    }
    
    usedTokensList := make([]TokenUsage, 0)
    for id, usage := range usageMap {
        if usage > 0 {
            usedTokensList = append(usedTokensList, TokenUsage{ID: id, Usage: usage})
            usedTokens++
        }
        totalUsage += usage
        if usage > maxUsage {
            maxUsage = usage
        }
    }

    // Sort by ID for consistent output
    sort.Slice(usedTokensList, func(i, j int) bool {
        return usedTokensList[i].ID < usedTokensList[j].ID
    })

    // Print first few used tokens
    fmt.Println("\nToken Usage:")
    for i := 0; i < min(5, len(usedTokensList)); i++ {
        fmt.Printf("Token %d: %d use(s)\n", usedTokensList[i].ID, usedTokensList[i].Usage)
    }
    if len(usedTokensList) > 5 {
        fmt.Println("...")

    fmt.Printf("\nLeast Used Token(s):\n")
    for i, token := range leastUsed {
        if i < 5 { // Show only first 5 least used tokens
            fmt.Printf("Token %d: %d use(s)\n", token.ID, token.Usage)
        } else {
            fmt.Printf("... and %d more tokens with %d use(s)\n", 
                len(leastUsed)-5, token.Usage)
            break
        }
    }

    fmt.Printf("\nSummary:\n")
    fmt.Printf("Total operations: %d\n", operations)
    fmt.Printf("Tokens used: %d out of 1000\n", usedTokens)
    fmt.Printf("Maximum usage: %d\n", maxUsage)
    fmt.Printf("Average usage per used token: %.2f\n", 
        float64(totalUsage)/float64(max(1, usedTokens)))
}
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}