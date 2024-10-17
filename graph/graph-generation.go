package graph

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"math/rand"
	"time"
)

type Generator struct {
	driver neo4j.DriverWithContext
}

type Result struct {
	NodesCreated int
	EdgesCreated int
}

func NewGenerator(driver neo4j.DriverWithContext) *Generator {
	return &Generator{driver: driver}
}

func (generator *Generator) CreateGraph(ctx context.Context, nodeCount, edgeCount int) (Result, error) {
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
		return Result{}, err
	}

	return result.(Result), nil
}

func (generator *Generator) createGraphTx(ctx context.Context, tx neo4j.ManagedTransaction, nodeCount, edgeCount int) (Result, error) {
	nodesCreated, err := createNode(ctx, tx, nodeCount)
	if err != nil {
		return Result{}, err
	}

	edgesCreated, err := createEdge(ctx, tx, nodeCount, edgeCount)
	if err != nil {
		return Result{}, err
	}

	return Result{NodesCreated: nodesCreated, EdgesCreated: edgesCreated}, nil

}

func createNode(ctx context.Context, tx neo4j.ManagedTransaction, nodeCount int) (int, error) {
	for i := 1; i <= nodeCount; i++ {
		_, err := tx.Run(ctx, "CREATE (n:Node {ID: $id})", map[string]interface{}{"id": i})
		if err != nil {
			return 0, err
		}
	}
	return nodeCount, nil
}

func createEdge(ctx context.Context, tx neo4j.ManagedTransaction, nodeCount, edgeCount int) (int, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalEdges := 0

	for i := 1; i <= nodeCount; i++ {
		usedTargets := make(map[int]struct{})

		for j := 0; j < edgeCount; j++ {
			var targetID int
			for {
				targetID = r.Intn(nodeCount) + 1
				if targetID != i && !isTargetUsed(usedTargets, targetID) {
					break // 適切なターゲットが見つかった場合はループを抜ける
				}
			}

			_, err := tx.Run(ctx, `
				MATCH (a:Node {ID: $source}), (b:Node {ID: $target})
				CREATE (a)-[:CONNECTED]->(b)
			`, map[string]interface{}{"source": i, "target": targetID})
			if err != nil {
				return totalEdges, err
			}

			usedTargets[targetID] = struct{}{}
			totalEdges++
		}
	}
	return totalEdges, nil
}

func isTargetUsed(usedTargets map[int]struct{}, targetID int) bool {
	_, exists := usedTargets[targetID]
	return exists
}
