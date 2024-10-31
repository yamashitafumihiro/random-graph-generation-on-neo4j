package io

import (
	"fmt"
	"log"
)

func Input() (int, int, int, map[string]float64, int) {
	var nodeCount, edgeCount, propertySize, batchSize int

	fmt.Print("Enter the number of nodes (nodeCount): ")
	_, err := fmt.Scan(&nodeCount)
	if err != nil {
		log.Fatal("Invalid input for nodeCount:", err)
	}

	fmt.Print("Enter the number of edges (edgeCount): ")
	_, err = fmt.Scan(&edgeCount)
	if err != nil {
		log.Fatal("Invalid input for edgeCount:", err)
	}

	fmt.Print("Enter the size of properties (propertySize): ")
	_, err = fmt.Scan(&propertySize)
	if err != nil {
		log.Fatal("Invalid input for propertySize:", err)
	}

	fmt.Print("Enter the batch size (batchSize): ")
	_, err = fmt.Scan(&batchSize)
	if err != nil {
		log.Fatal("Invalid input for batchSize:", err)
	}

	entities := make(map[string]float64)
	var entityName string
	var probability float64

	for {
		fmt.Print("Enter entity name (or 'done' to finish): ")
		_, err = fmt.Scan(&entityName)
		if err != nil {
			log.Fatal("Invalid input for entity name:", err)
		}
		if entityName == "done" {
			break
		}

		fmt.Printf("Enter probability for %s: ", entityName)
		_, err = fmt.Scan(&probability)
		if err != nil {
			log.Fatal("Invalid input for probability:", err)
		}

		entities[entityName] = probability
	}

	return nodeCount, edgeCount, propertySize, entities, batchSize
}
