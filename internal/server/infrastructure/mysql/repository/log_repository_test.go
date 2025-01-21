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

type LogRepositorySuite struct {
	suite.Suite
	repo *LogRepository
}

func (suite *LogRepositorySuite) SetupTest() {
	suite.repo = NewLogRepository(dbConnTest)
}

func (suite *LogRepositorySuite) TearDownTest() {
}

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

func (suite *LogRepositorySuite) TestInsertCTRLog() {
	err := suite.repo.CTRSave(context.Background(), &domain.CTRLog{
		CreatedAt: time.Now(),
		Objectid:  "123456",
	})
	assert.NoError(suite.T(), err, "Failed to insert CTR log.")
}

func TestLogRepositorySuite(t *testing.T) {
	suite.Run(t, new(LogRepositorySuite))
}
