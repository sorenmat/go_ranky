#!/bin/bash
 n=$1
shift
while [ $(( n -= 1 )) -ge 0 ]
do
    curl -s -X PUT -H "Content-Type: application/json" -d '{"Name": "SorenMat"}' http://localhost:8080/players -o /dev/null 
done
