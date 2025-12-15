# Database Configuration and Models

## Environment Configuration
The application supports loading configuration from both `configs/config.yaml` and `.env` files.
Environment variables take precedence over the YAML configuration file.

### Setup
1. Copy `.env.example` to `.env`:
   ```bash
   cp .env.example .env
   ```
2. Update `.env` with your local database credentials:
   ```properties
   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=your_password
   DB_NAME=next_ai_gateway
   ```

## Docker Deployment
To run the database in Docker with automatic initialization (schema + seed data):

1. **Start Database**:
   ```bash
   # From project root
   docker-compose -f deployments/docker-compose.yml up -d
   ```
   
2. **Reset Database (Re-seed)**:
   If you want to wipe the database and re-apply init.sql and seed.sql:
   ```bash
   # Windows
   scripts/reset_db.bat
   ```
   
   Or manually:
   ```bash
   docker-compose -f deployments/docker-compose.yml down -v
   docker-compose -f deployments/docker-compose.yml up -d
   ```

## Database Models
The Go models are located in `internal/model/entity/models.go` and correspond to the schema defined in `docs/init.sql`.

### Generation and Updates
Currently, the models are manually maintained based on `docs/init.sql`.
If you modify `docs/init.sql`, please update the corresponding structs in `internal/model/entity/models.go` to ensure consistency.

### Tables
- `departments`: Organization structure
- `users`: User information
- `ai_models`: AI Model routing configuration
- `api_provider_keys`: Provider keys and health status
- `token_quotas`: Usage limits
- `request_logs`: Audit logs (partitioning recommended for production)
- `monthly_costs`: Aggregated cost data

## Version Control
- `.env` files are ignored by git to protect sensitive credentials.
- `bin/` and `logs/` directories are also ignored.
