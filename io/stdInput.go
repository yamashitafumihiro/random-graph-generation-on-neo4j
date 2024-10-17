package io

import (
	"fmt"
	"log"
)

func Input() (int, int) {
	var nodeCount, edgeCount int

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
	return nodeCount, edgeCount
}
