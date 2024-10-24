#!/bin/bash

../bin/my-cli <<EOF
match (n) detach delete n
exit
EOF

../bin/graph-generator <<EOF
200
3
10000
movie
50
user
50
done
exit
EOF

../bin/my-cli <<EOF
CALL db.clearQueryCaches()
match (m:movie)-[r]->(u:user) where m.ID=1 return u.ID
EOF