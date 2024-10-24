package io

import (
	"bufio"
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"os"
	"strings"
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
		result, err := executeQuery(ctx, driver, query)
		if err != nil {
			log.Fatal()
		}
		fmt.Println(result)
	}
}

func executeQuery(ctx context.Context, driver neo4j.DriverWithContext, query string) (neo4j.ResultWithContext, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(session, ctx)

	result, err := session.Run(ctx, query, nil)
	if err != nil {
		return result, err
	}
	return result, nil
}
