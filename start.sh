#!/bin/bash

echo "Starting Expense Management System..."

# Build and start all services
docker compose up --build -d

echo ""
echo "Services starting..."
echo ""
echo "URLs:"
echo "   Frontend: http://localhost:3000"
echo "   Backend:  http://localhost:8080"
echo "   Health:   http://localhost:8080/api/health"
echo ""
echo "Demo Accounts:"
echo "   Employee: employee1@example.com / password123"
echo "   Manager:  manager@example.com / password123"
echo ""
echo "Useful Commands:"
echo "   View logs:     docker-compose logs -f"
echo "   Stop:          docker-compose down"
echo "   Restart:       docker-compose restart"
echo "   Clean & Reset: ./cleanup.sh"
echo ""
echo "Wait ~30 seconds for all services to be ready..."
