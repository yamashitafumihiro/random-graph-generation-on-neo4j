package main

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"math/rand"
	"random-graph-generation-on-neo4j/io"
	"time"
)

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

	io.Interface(ctx, driver)
}
