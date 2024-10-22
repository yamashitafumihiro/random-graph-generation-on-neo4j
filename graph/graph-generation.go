package graph

import (
	"context"
	"fmt"
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

type NodeInfo struct {
	ID     int
	Entity string
}

func NewGenerator(driver neo4j.DriverWithContext) *Generator {
	return &Generator{driver: driver}
}

func (generator *Generator) CreateGraph(ctx context.Context, nodeCount, edgeCount, propertySize int, entities map[string]float64) (Result, error) {
	session := generator.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(session, ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return generator.createGraphTx(ctx, tx, nodeCount, edgeCount, propertySize, entities)
	})

	if err != nil {
		return Result{}, err
	}

	return result.(Result), nil
}

func (generator *Generator) createGraphTx(ctx context.Context, tx neo4j.ManagedTransaction, nodeCount, edgeCount, propertySize int, entities map[string]float64) (Result, error) {
	nodes, err := createNode(ctx, tx, nodeCount, propertySize, entities)
	if err != nil {
		return Result{}, err
	}

	edgesCreated, err := createEdge(ctx, tx, nodes, edgeCount)
	if err != nil {
		return Result{}, err
	}

	return Result{NodesCreated: nodeCount, EdgesCreated: edgesCreated}, nil
}

func generateProperties(size int) map[string]interface{} {
	properties := make(map[string]interface{})
	for i := 1; i <= size; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		properties[key] = value
	}
	return properties
}

func selectEntity(entities map[string]float64) (string, error) {
	r := rand.Float64() * 100
	sum := 0.0

	for entity, probability := range entities {
		sum += probability
		if r <= sum {
			return entity, nil
		}
	}

	return "", fmt.Errorf("failed to select an entity: probabilities may not sum to 100")
}

func createNode(ctx context.Context, tx neo4j.ManagedTransaction, nodeCount, propertySize int, entities map[string]float64) ([]NodeInfo, error) {
	var nodes []NodeInfo
	entityCounters := make(map[string]int)

	for i := 1; i <= nodeCount; i++ {
		entity, err := selectEntity(entities)
		if err != nil {
			return nil, err
		}

		entityCounters[entity]++
		entityID := entityCounters[entity]

		properties := generateProperties(propertySize)
		properties["ID"] = entityID

		_, err = tx.Run(ctx, "CREATE (n:"+entity+") SET n += $props", map[string]interface{}{"props": properties})
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, NodeInfo{ID: entityID, Entity: entity})
	}
	return nodes, nil
}

func edgeExists(usedEdges map[string]struct{}, edgeKey string) bool {
	_, exists := usedEdges[edgeKey]
	return exists
}

func createEdge(ctx context.Context, tx neo4j.ManagedTransaction, nodes []NodeInfo, edgeCount int) (int, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalEdges := 0

	usedEdges := make(map[string]struct{})

	for _, node := range nodes {
		usedTargets := make(map[int]struct{})

		for j := 0; j < edgeCount; j++ {
			var target NodeInfo
			for {
				target = nodes[r.Intn(len(nodes))]

				edgeKey := fmt.Sprintf("%s-%d->%s-%d", node.Entity, node.ID, target.Entity, target.ID)
				if target.ID == node.ID || isTargetUsed(usedTargets, target.ID) || edgeExists(usedEdges, edgeKey) {
					continue
				}
				break
			}

			_, err := tx.Run(ctx, `
                MATCH (a:`+node.Entity+` {ID: $source}), (b:`+target.Entity+` {ID: $target})
                CREATE (a)-[:CONNECTED]->(b)
            `, map[string]interface{}{"source": node.ID, "target": target.ID})
			if err != nil {
				return totalEdges, err
			}

			usedEdges[fmt.Sprintf("%s-%d->%s-%d", node.Entity, node.ID, target.Entity, target.ID)] = struct{}{}
			usedTargets[target.ID] = struct{}{}
			totalEdges++
		}
	}
	return totalEdges, nil
}

func isTargetUsed(usedTargets map[int]struct{}, targetID int) bool {
	_, exists := usedTargets[targetID]
	return exists
}
