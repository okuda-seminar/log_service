package repository

import (
	"database/sql"
	"testing"

	dbTest "log_service/app/infrastructure/mysql/db/db_test"
)

var (
	dbConnTest *sql.DB
)

func TestMain(m *testing.M) {
	resource, pool := dbTest.CreateContainer()
	defer dbTest.CloseContainer(resource, pool)

	dbConnTest = dbTest.ConnectDB(resource, pool)
	defer dbConnTest.Close()

	dbTest.SetupTestDB("../db/schema/000001_log.up.sql")

	m.Run()
}