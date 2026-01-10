#!/bin/bash

echo "leaning up Expense Management System..."

# Stop all containers
echo "Stopping containers..."
docker compose down -v

# Remove all containers, networks, volumes
echo "Removing volumes..."
docker volume rm expense-management-system_postgres_data 2>/dev/null || true

# Remove node_modules and build artifacts
echo "Cleaning frontend..."
rm -rf frontend/node_modules
rm -rf frontend/.nuxt
rm -rf frontend/.output
rm -rf frontend/dist

# Remove Go build cache
echo "Cleaning backend..."
cd backend
go clean -cache -modcache -i -r 2>/dev/null || true
rm -f api
cd ..

echo "✅ Cleanup complete!"
echo ""
echo "To start fresh, run: ./start.sh"
