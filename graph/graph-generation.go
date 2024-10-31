package graph

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type Generator struct {
	driver         neo4j.DriverWithContext
	entityCounters map[string]int
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
	return &Generator{
		driver:         driver,
		entityCounters: make(map[string]int),
	}
}

func (generator *Generator) CreateGraph(ctx context.Context, nodeCount, edgeCount, propertySize, batchSize int, entities map[string]float64) (Result, error) {
	var totalNodesCreated int
	var totalEdgesCreated int
	var allNodes []NodeInfo

	for i := 0; i < nodeCount; i += batchSize {
		currentBatchSize := batchSize
		if i+batchSize > nodeCount {
			currentBatchSize = nodeCount - i
		}

		nodes, err := generator.createNodesBatch(ctx, currentBatchSize, propertySize, entities)
		if err != nil {
			return Result{}, err
		}
		totalNodesCreated += len(nodes)
		allNodes = append(allNodes, nodes...)
	}

	for i := 0; i < len(allNodes); i += batchSize {
		end := i + batchSize
		if end > len(allNodes) {
			end = len(allNodes)
		}
		nodesBatch := allNodes[i:end]

		edgesCreated, err := generator.createEdgesBatch(ctx, nodesBatch, allNodes, edgeCount)
		if err != nil {
			return Result{}, err
		}
		totalEdgesCreated += edgesCreated
	}

	return Result{NodesCreated: totalNodesCreated, EdgesCreated: totalEdgesCreated}, nil
}

func (generator *Generator) createNodesBatch(ctx context.Context, batchSize, propertySize int, entities map[string]float64) ([]NodeInfo, error) {
	session := generator.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(session, ctx)

	nodes := []NodeInfo{}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		for i := 0; i < batchSize; i++ {
			entity, err := selectEntity(entities)
			if err != nil {
				return nil, err
			}

			generator.entityCounters[entity]++
			entityID := generator.entityCounters[entity]

			properties := generateProperties(propertySize)
			properties["ID"] = entityID

			_, err = tx.Run(ctx, "CREATE (n:"+entity+") SET n += $props", map[string]interface{}{"props": properties})
			if err != nil {
				return nil, err
			}

			nodes = append(nodes, NodeInfo{ID: entityID, Entity: entity})
		}
		return nil, nil
	})

	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (generator *Generator) createEdgesBatch(ctx context.Context, nodesBatch []NodeInfo, allNodes []NodeInfo, edgeCount int) (int, error) {
	session := generator.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(session, ctx)

	totalEdges := 0

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		for _, node := range nodesBatch {
			usedTargets := make(map[int]struct{})

			for j := 0; j < edgeCount; j++ {
				var target NodeInfo
				for {
					target = allNodes[r.Intn(len(allNodes))]

					if target.ID == node.ID {
						continue
					}
					if _, exists := usedTargets[target.ID]; exists {
						continue
					}
					break
				}

				_, err := tx.Run(ctx, `
                    MATCH (a:`+node.Entity+` {ID: $source}), (b:`+target.Entity+` {ID: $target})
                    CREATE (a)-[:CONNECTED]->(b)
                `, map[string]interface{}{"source": node.ID, "target": target.ID})
				if err != nil {
					return nil, err
				}

				usedTargets[target.ID] = struct{}{}
				totalEdges++
			}
		}
		return nil, nil
	})

	if err != nil {
		return 0, err
	}

	return totalEdges, nil
}

func generateProperties(size int) map[string]interface{} {
	properties := make(map[string]interface{})

	for i := 1; i <= size; i++ {
		key := fmt.Sprintf("key%d", i)
		value := rand.Intn(10) + 1
		properties[key] = strconv.Itoa(value)
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
