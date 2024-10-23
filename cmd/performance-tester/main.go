package main

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"math/rand"
	"random-graph-generation-on-neo4j/graph"
	"random-graph-generation-on-neo4j/performance"
	"time"
)

type TestResult struct {
	NodeCount    int
	EdgeCount    int
	PropertySize int
	Query        string
	AverageTime  float64
}

func main() {
	ctx := context.Background()
	rand.NewSource(time.Now().UnixNano())
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

	nodeCounts := []int{100}
	edgeCounts := []int{4}
	propertySizes := []int{10, 100, 1000, 10000}
	entities := map[string]float64{
		"user":  50,
		"movie": 50,
	}

	queries := []string{
		"MATCH (n:user)-[r:CONNECTED]->(m:movie) WHERE n.key1 = $value RETURN n, r, m",
	}

	var results []TestResult

	for _, nodeCount := range nodeCounts {
		for _, edgeCount := range edgeCounts {
			for _, propertySize := range propertySizes {
				generator := graph.NewGenerator(driver)
				_, err := generator.CreateGraph(ctx, nodeCount, edgeCount, propertySize, entities)
				if err != nil {
					log.Fatal(err)
				}

				for _, query := range queries {
					value := fmt.Sprintf("value%d", rand.Intn(propertySize)+1)
					params := map[string]interface{}{"value": value}

					avgTime, err := performance.MeasureQueryPerformance(ctx, driver, query, params, 5)
					if err != nil {
						log.Fatal(err)
					}
					result := TestResult{
						NodeCount:    nodeCount,
						EdgeCount:    edgeCount,
						PropertySize: propertySize,
						Query:        query,
						AverageTime:  avgTime,
					}
					results = append(results, result)
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

	fmt.Printf("%-10s %-10s %-15s %-50s %-15s\n", "Nodes", "Edges", "Properties", "Query", "AvgTime(ms)")
	for _, res := range results {
		fmt.Printf("%-10d %-10d %-15d %-50s %-15f\n", res.NodeCount, res.EdgeCount, res.PropertySize, res.Query, res.AverageTime)
	}
}
