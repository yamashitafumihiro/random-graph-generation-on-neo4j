package graph

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
)

type Generator struct {
	driver neo4j.DriverWithContext
}

func NewGenerator(driver neo4j.DriverWithContext) *Generator {
	return &Generator{driver: driver}
}

func (generator *Generator) CreateGraph(ctx context.Context, nodeCount, edgeCount int) (int, error) {
	session := generator.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(session, ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return generator.createGraphTx(ctx, tx, nodeCount, edgeCount)
	})

	if err != nil {
		return 0, err
	}

	return result.(int), nil
}

func (generator *Generator) createGraphTx(ctx context.Context, tx neo4j.ManagedTransaction, nodeCount, edgeCount int) (int, error) {
	for i := 1; i <= nodeCount; i++ {
		_, err := tx.Run(ctx, "CREATE (n:Node {ID: $id})", map[string]interface{}{"id": i})
		if err != nil {
			return 0, err
		}
	}

	return nodeCount, nil
}
