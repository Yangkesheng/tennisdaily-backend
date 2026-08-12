package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: dbclient <dsn> <list|query|migrate> [sql|sqlFile]")
		os.Exit(1)
	}
	dsn := os.Args[1]
	op := os.Args[2]

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Fprintln(os.Stderr, "open failed:", err)
		os.Exit(1)
	}
	defer db.Close()

	if op == "list" {
		rows, err := db.Query("SHOW DATABASES")
		if err != nil {
			fmt.Fprintln(os.Stderr, "query failed:", err)
			os.Exit(1)
		}
		defer rows.Close()
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				fmt.Fprintln(os.Stderr, "scan failed:", err)
				os.Exit(1)
			}
			fmt.Println(name)
		}
		return
	}

	if op == "migrate" {
		if len(os.Args) < 4 {
			fmt.Fprintln(os.Stderr, "missing sqlFile")
			os.Exit(1)
		}
		content, err := os.ReadFile(os.Args[3])
		if err != nil {
			fmt.Fprintln(os.Stderr, "read failed:", err)
			os.Exit(1)
		}
		if _, err := db.Exec(string(content)); err != nil {
			fmt.Fprintln(os.Stderr, "exec failed:", err)
			os.Exit(1)
		}
		fmt.Println("migrate ok:", os.Args[3])
		return
	}

	if op == "query" {
		if len(os.Args) < 4 {
			fmt.Fprintln(os.Stderr, "missing sql")
			os.Exit(1)
		}
		rows, err := db.Query(os.Args[3])
		if err != nil {
			fmt.Fprintln(os.Stderr, "query failed:", err)
			os.Exit(1)
		}
		defer rows.Close()
		columns, _ := rows.Columns()
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		for rows.Next() {
			if err := rows.Scan(pointers...); err != nil {
				fmt.Fprintln(os.Stderr, "scan failed:", err)
				os.Exit(1)
			}
			for i, column := range columns {
				fmt.Printf("%s=%v ", column, values[i])
			}
			fmt.Println()
		}
		return
	}

	fmt.Fprintln(os.Stderr, "unknown op:", op)
	os.Exit(1)
}
