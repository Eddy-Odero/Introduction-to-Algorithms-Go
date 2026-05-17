#!/bin/sh
echo "Stopping old containers..."
docker rm -f ascii-art-container 2>/dev/null
echo "Building production multi-stage image..."
docker build -t ascii-art-web:1.0 .
echo "Pruning dangling build cache objects..."
docker image prune -f
echo "Launching isolated web container..."
docker run -d -p 8080:8080 --name ascii-art-container ascii-art-web:1.0
echo "Deployment complete! Application active at http://localhost:8080"
