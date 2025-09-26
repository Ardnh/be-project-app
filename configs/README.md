# configs/ - Configuration Files

## Purpose
Contains configuration files for different environments.

## Files
- **config.yaml**: Default development configuration
- **config.prod.yaml**: Production configuration
- **config.test.yaml**: Test environment configuration

## Environment Variables
Configuration can be overridden using environment variables:
- `DB_HOST`, `DB_PORT`, `DB_USERNAME`, `DB_PASSWORD`
- `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`
- `JWT_SECRET`

## Usage
Application automatically loads the appropriate config based on `APP_ENV` environment variable.
