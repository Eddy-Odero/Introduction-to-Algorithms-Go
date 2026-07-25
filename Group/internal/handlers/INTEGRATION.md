# Phase 6 Integration Guide — Comments

## Overview
This document describes how to integrate the comment system into the existing forum application.

## Files Added
- `internal/models/comment.go` — Database operations
- `internal/models/comment_test.go` — Model unit tests
- `internal/handlers/comments.go` — HTTP handler
- `internal/handlers/comments_test.go` — Handler tests

## Integration Steps

### 1. Initialize Database in main.go
```go
import "forum/internal/db"

func main() {
    database, err := db.Init("forum.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Pass database to router
    mux := router.SetupRouter(uiDir, database)
    // ...
}
