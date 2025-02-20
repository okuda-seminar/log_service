// Package repository provides implementations for handling database operations
// related to logging and CTR logs in the application.
package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"log_service/internal/server/domain"
)

// LogRepositorySuite is a test suite for testing the LogRepository.
// It includes setup and teardown functions, along with test cases for
// log insertion and retrieval.
type LogRepositorySuite struct {
	suite.Suite
	repo *LogRepository // The LogRepository instance used in tests.
}

// SetupTest initializes the repository for each test in the suite.
func (suite *LogRepositorySuite) SetupTest() {
	suite.repo = NewLogRepository(dbConnTest)
}

// TearDownTest is called after each test in the suite to clean up resources.
// Currently, no teardown logic is implemented.
func (suite *LogRepositorySuite) TearDownTest() {
	// Add teardown logic if needed in the future.
}

// TestInsertLog tests the insertion of a log entry into the database.
func (suite *LogRepositorySuite) TestInsertLog() {
	err := suite.repo.Save(context.Background(), &domain.Log{
		LogLevel:           "INFO",
		Date:               time.Now(),
		DestinationService: "UserService",
		SourceService:      "AuthService",
		RequestType:        "POST",
		Content:            "User created successfully.",
	})

	require.NoError(suite.T(), err, "Failed to insert log.")
}

// TestList tests the retrieval of log entries from the database.
func (suite *LogRepositorySuite) TestList() {
	err := suite.repo.Save(context.Background(), &domain.Log{
		LogLevel:           "INFO",
		Date:               time.Now(),
		DestinationService: "UserService",
		SourceService:      "AuthService",
		RequestType:        "POST",
		Content:            "Test Get Log.",
	})
	require.NoError(suite.T(), err)

	results, err := suite.repo.List(context.Background())
	require.NoError(suite.T(), err, "Failed to get logs.")
	assert.GreaterOrEqual(suite.T(), len(results), 1, "Want 1 logs but got %d", len(results))
}

// TestInsertCTRLog tests the insertion of a CTR log entry into the database.
func (suite *LogRepositorySuite) TestInsertCTRLog() {
	err := suite.repo.CTRSave(context.Background(), &domain.CTRLog{
		EventType: "tap",
		CreatedAt: time.Now(),
		ObjectID:  "123456",
	})
	assert.NoError(suite.T(), err, "Failed to insert CTR log.")
}

// TestListCTRLogs tests the retrieval of CTR log entries from the database.
func (suite *LogRepositorySuite) TestListCTRLogs() {
	err := suite.repo.CTRSave(context.Background(), &domain.CTRLog{
		EventType: "tap",
		CreatedAt: time.Now(),
		ObjectID:  "123456",
	})
	require.NoError(suite.T(), err)

	results, err := suite.repo.CTRList(context.Background())
	require.NoError(suite.T(), err, "Failed to get CTR logs.")
	assert.GreaterOrEqual(suite.T(), len(results), 1, "Want 1 CTR logs but got %d", len(results))
}

// TestLogRepositorySuite runs the LogRepositorySuite tests using testify's suite package.
func TestLogRepositorySuite(t *testing.T) {
	suite.Run(t, new(LogRepositorySuite))
}
