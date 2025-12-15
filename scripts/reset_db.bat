@echo off
echo Stopping existing containers...
docker-compose --env-file .env -f deployments/docker-compose.yml down -v

echo Starting database container...
docker-compose --env-file .env -f deployments/docker-compose.yml up -d

echo Waiting for database to initialize...
timeout /t 10

echo Database restarted and initialized with seed data.
