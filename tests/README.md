# 🧪 Transfer Service Tests

This directory contains comprehensive tests for the transfer service application, ensuring reliability and correctness of all business logic.

## 📁 Test Structure

```
tests/
├── service/
│   ├── account_service_test.go    # Account service unit tests
│   └── transaction_service_test.go # Transaction service unit tests
├── run_tests.sh                   # Test runner script
└── README.md                      # This file
```

## 🎯 Test Coverage

### Account Service Tests (`tests/service/account_service_test.go`)

| Test Case | Description | Status |
|-----------|-------------|--------|
| `TestCreateAccount_Success` | ✅ Create account with valid data | ✅ |
| `TestCreateAccount_DuplicateID` | ❌ Attempt to create account with existing ID | ✅ |
| `TestCreateAccount_InvalidRequest` | ⚠️ Handle invalid request data | ✅ |
| `TestGetAccount_Success` | ✅ Retrieve existing account | ✅ |
| `TestGetAccount_NotFound` | ❌ Attempt to get non-existent account | ✅ |

### Transaction Service Tests (`tests/service/transaction_service_test.go`)

| Test Case | Description | Status |
|-----------|-------------|--------|
| `TestTransferValidation_SameAccounts` | ❌ Validate same source/destination accounts are rejected | ✅ |
| `TestTransferValidation_ValidAccounts` | ✅ Validate different accounts are accepted | ✅ |
| `TestTransferValidation_AmountValidation` | ⚠️ Validate transfer amounts (positive, zero, negative) | ✅ |
| `TestTransferValidation_AccountIDValidation` | ⚠️ Validate account ID combinations | ✅ |

## 🚀 Running Tests

### Run All Tests
```bash
# Using the test runner script
./tests/run_tests.sh

# Or directly with Go
go test ./tests/service -v
```

### Run Specific Test Categories
```bash
# Run only account service tests
go test ./tests/service -v -run "Test.*Account.*"

# Run only transaction service tests
go test ./tests/service -v -run "Test.*Validation.*"

# Run only success cases
go test ./tests/service -v -run "Test.*Success"

# Run only error cases
go test ./tests/service -v -run "Test.*Error"
```

### Run Individual Tests
```bash
# Run specific test
go test ./tests/service -v -run "TestCreateAccount_Success"

# Run all tests in service package with coverage
go test ./tests/service -v -cover

# Run tests with race detection
go test ./tests/service -v -race
```

## 📊 Test Features

- **🔍 Mock Dependencies**: Uses mock repositories to isolate business logic
- **🎯 Edge Cases**: Tests both success and failure scenarios
- **💰 Decimal Precision**: Validates 5-decimal precision requirements
- **🔒 Business Rules**: Ensures all business validation rules are enforced
- **📝 Structured Logging**: Tests include proper logging verification

## 🛠️ Test Dependencies

The tests use the following Go testing packages:
- `testing` - Standard Go testing framework
- `github.com/stretchr/testify` - Enhanced assertions and mocking
- `github.com/shopspring/decimal` - Decimal arithmetic for monetary values

## 📈 Coverage Goals

- **Account Service**: 100% business logic coverage
- **Transaction Service**: 100% validation logic coverage
- **Error Handling**: All error paths tested
- **Edge Cases**: Boundary conditions validated

## 🔧 Test Configuration

Tests are configured to run with:
- Verbose output (`-v` flag)
- Coverage reporting (`-cover` flag)
- Race detection (`-race` flag when needed)
- Timeout handling for long-running operations

## 📝 Adding New Tests

When adding new tests:

1. **Follow Naming Convention**: `Test[FunctionName]_[Scenario]`
2. **Use Descriptive Names**: Make test purpose clear from the name
3. **Test Both Success and Failure**: Cover all code paths
4. **Mock External Dependencies**: Don't rely on real database connections
5. **Validate Business Rules**: Ensure all validation logic is tested
6. **Include Edge Cases**: Test boundary conditions and error scenarios

## 🎉 Test Results

All tests should pass before merging code changes. The test suite ensures:
- ✅ Account creation and retrieval work correctly
- ✅ Transaction validation prevents invalid transfers
- ✅ Error handling provides meaningful responses
- ✅ Business rules are enforced consistently
- ✅ Decimal precision is maintained throughout

---

**Happy Testing! 🧪✨**
