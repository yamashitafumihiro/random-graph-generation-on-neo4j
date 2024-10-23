package main

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"random-graph-generation-on-neo4j/graph"
	"random-graph-generation-on-neo4j/performance"
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

	nodeCounts := []int{100, 500, 1000}
	edgeCounts := []int{2, 4, 8}
	propertySizes := []int{100, 500, 1000}
	entities := map[string]float64{"Entity1": 100}

	queries := []string{
		"MATCH (n) RETURN count(n)",
		"MATCH (n) WHERE n.key1 = $value RETURN n",
		"MATCH p = (n)-[*1..2]-(m) RETURN p LIMIT 100",
	}

	for _, nodeCount := range nodeCounts {
		for _, edgeCount := range edgeCounts {
			for _, propertySize := range propertySizes {

				generator := graph.NewGenerator(driver)
				_, err := generator.CreateGraph(ctx, nodeCount, edgeCount, propertySize, entities)
				if err != nil {
					log.Fatal(err)
				}

				for _, query := range queries {
					avgTime, err := performance.MeasureQueryPerformance(ctx, driver, query, 5)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Printf("Average time for query '%s': %f ms\n", query, avgTime)
				}
				err = performance.ClearCache(ctx, driver)
				if err != nil {
					log.Fatal(err)
				}

				err = performance.ClearDatabase(ctx, driver)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
