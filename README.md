# random graph generator

You can create a random graph on neo4j consisting of an arbitrary number of nodes and edges by this generator.

## Requirement

- Go 1.20+ (https://go.dev/dl/)
- Neo4j

### Recommended: Running Neo4j with Docker

You can run Neo4j with Docker:

```bash
docker run -p 7474:7474 -p 7687:7687 -d \
  --env NEO4J_AUTH=none \
  neo4j:5.x
```

After starting, access the Neo4j Browser at: http://localhost:7474

## Getting started

```bash
$ git clone git@github.com:yamashitafumihiro/random-graph-generation-on-neo4j.git
$ go run main.go
Connection established.
Enter the number of nodes (nodeCount): 10
Enter the number of edges (edgeCount): 2
Created graph with 10 nodes and 20 edges
```

- If the connection to the neo4j server is successful, the message `Connection established` is displayed.
  If the connection fails, an error message is displayed. Please review the port settings, NEO4J_AUTH value, etc.
- Once connected, you will be asked to enter the number of nodes and the number of edges.

