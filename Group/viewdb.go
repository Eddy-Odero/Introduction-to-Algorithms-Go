package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run viewdb.go <database.db>")
		os.Exit(1)
	}

	dbPath := os.Args[1]

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// List all tables
	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Println("║         TABLES IN DATABASE           ║")
	fmt.Println("╚══════════════════════════════════════╝")
	showTables(db)

	// Show contents of each table
	tables := []string{"categories", "users", "posts", "comments", "reactions", "sessions", "post_categories"}

	for _, table := range tables {
		fmt.Println()
		fmt.Printf("╔══════════════════════════════════════╗\n")
		fmt.Printf("║  TABLE: %-28s ║\n", table)
		fmt.Printf("╚══════════════════════════════════════╝\n")
		showTableContents(db, table)
	}
}

func showTables(db *sql.DB) {
	rows, err := db.Query(`
		SELECT name FROM sqlite_master 
		WHERE type='table' AND name NOT LIKE 'sqlite_%'
		ORDER BY name
	`)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		rows.Scan(&name)
		fmt.Printf("  • %s\n", name)
	}
}

func showTableContents(db *sql.DB, table string) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", table))
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		fmt.Printf("  Error getting columns: %v\n", err)
		return
	}

	if len(columns) == 0 {
		fmt.Println("  (no columns)")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	// Header
	for _, col := range columns {
		fmt.Fprintf(w, "%s\t", col)
	}
	fmt.Fprintln(w)

	// Separator
	for range columns {
		fmt.Fprintf(w, "----------\t")
	}
	fmt.Fprintln(w)

	// Rows
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	count := 0
	for rows.Next() {
		rows.Scan(scanArgs...)
		for _, v := range values {
			s := fmt.Sprintf("%v", v)
			if len(s) > 30 {
				s = s[:27] + "..."
			}
			fmt.Fprintf(w, "%s\t", s)
		}
		fmt.Fprintln(w)
		count++
	}
	w.Flush()

	if count == 0 {
		fmt.Println("  (empty table — no rows yet)")
	} else {
		fmt.Printf("\n  %d row(s)\n", count)
	}
}