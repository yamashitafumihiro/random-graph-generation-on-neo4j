package io

import (
	"fmt"
	"log"
)

func Input() (int, int, int) {
	var nodeCount, edgeCount, propertySize int

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

	return nodeCount, edgeCount, propertySize
}
