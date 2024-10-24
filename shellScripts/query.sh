#!/bin/bash

total_time=0

for i in {1..5}
do
  random_id=$(( RANDOM % 100 + 1 ))

  result=$(../bin/my-cli <<EOF
CALL db.clearQueryCaches()
match (m:movie) where m.key1="value4070" AND m.key10="value8574" return m
exit
EOF
)

  execution_time=$(echo "$result" | grep 'Query executed in' | tail -n 1 | awk '{print $4}' | sed 's/ms//')

  total_time=$(awk "BEGIN {print $total_time + $execution_time}")

  echo "Run $i with ID=$random_id: Executed in ${execution_time}ms"
done

average_time=$(awk "BEGIN {print $total_time / 5}")

echo "Average Query Execution Time: ${average_time}ms"
