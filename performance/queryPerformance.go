package performance

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"time"
)

func MeasureQueryPerformance(ctx context.Context, driver neo4j.DriverWithContext, query string, params map[string]interface{}, runs int) (float64, error) {
	var totalDuration time.Duration

	for i := 0; i < runs; i++ {
		err := ClearCache(ctx, driver)
		if err != nil {
			return 0, err
		}

		session := driver.NewSession(ctx, neo4j.SessionConfig{})
		defer func(session neo4j.SessionWithContext, ctx context.Context) {
			err := session.Close(ctx)
			if err != nil {
				log.Fatal(err)
			}
		}(session, ctx)

		start := time.Now()

		_, err = session.Run(ctx, query, params)
		if err != nil {
			return 0, err
		}

		duration := time.Since(start)
		totalDuration += duration
	}
	avgDuration := totalDuration.Seconds() * 1000 / float64(runs)
	return avgDuration, nil
}

func ClearCache(ctx context.Context, driver neo4j.DriverWithContext) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(session, ctx)

	_, err := session.Run(ctx, "CALL db.clearQueryCaches()", nil)
	if err != nil {
		return err
	}
	return nil
}

func ClearDatabase(ctx context.Context, driver neo4j.DriverWithContext) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(session, ctx)

	_, err := session.Run(ctx, "MATCH (n) DETACH DELETE n", nil)
	if err != nil {
		return err
	}
	return nil
}
