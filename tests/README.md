# tests/ - Test Files

## Purpose
Contains all test files organized by test type and layer.

## Structure
```
tests/
├── unit/                # Unit tests
├── integration/         # Integration tests
└── fixtures/           # Test data fixtures
```

## Test Types

### Unit Tests
- Test individual functions/methods in isolation
- Mock external dependencies
- Fast execution
- High code coverage

### Integration Tests
- Test component interactions
- Use real databases (test containers)
- Test API endpoints end-to-end
- Slower but more comprehensive

### Fixtures
- Sample data for tests
- JSON/YAML test data files
- Reusable test scenarios

## Running Tests
```bash
# All tests
make test

# Unit tests only
go test -short ./...

# Integration tests
go test -run Integration ./...

# With coverage
make test-coverage
```
