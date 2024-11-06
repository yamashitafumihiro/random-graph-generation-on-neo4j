package main

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"random-graph-generation-on-neo4j/graph"
	"random-graph-generation-on-neo4j/io"
)

func main() {
	ctx := context.Background()

	dbUri := "neo4j://localhost"

	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.NoAuth())
	if err != nil {
		log.Fatal(err)
	}
	defer func(driver neo4j.DriverWithContext, ctx context.Context) {
		err := driver.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(driver, ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection established.")

	generator := graph.NewGenerator(driver)

	nodeCount, edgeCount, nodePropertySize, edgePropertySize, entities, batchSize := io.Input()

	result, err := generator.CreateGraph(ctx, nodeCount, edgeCount, nodePropertySize, edgePropertySize, batchSize, entities)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created graph with %d nodes and %d edges\n", result.NodesCreated, result.EdgesCreated)
}
