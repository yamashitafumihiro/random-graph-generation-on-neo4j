package io

import (
	"bufio"
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"os"
	"strings"
	"time"
)

func Interface(ctx context.Context, driver neo4j.DriverWithContext) {
	fmt.Println("Hello neo4j. enter 'exit' to terminate this session.")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("neo4j> ")
		if !scanner.Scan() {
			break
		}
		query := scanner.Text()
		if strings.TrimSpace(query) == "exit" {
			break
		}
		if err := executeAndPrintQuery(ctx, driver, query); err != nil {
			log.Printf("Error executing query: %v\n", err)
		}
	}
}

func executeAndPrintQuery(ctx context.Context, driver neo4j.DriverWithContext, query string) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() {
		if err := session.Close(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	startTime := time.Now()
	result, err := session.Run(ctx, query, nil)
	elapsedTime := time.Since(startTime)

	if err != nil {
		return err
	}

	for result.Next(ctx) {
		record := result.Record()
		for _, value := range record.Values {
			fmt.Printf("%v ", value)
		}
		fmt.Println()
	}
	fmt.Printf("Query executed in %v\n", elapsedTime)

	if err = result.Err(); err != nil {
		return err
	}
	return nil
}
