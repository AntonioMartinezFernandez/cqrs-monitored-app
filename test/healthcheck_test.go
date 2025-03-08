package test

import (
	"net/http"
	"testing"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/test/helpers"

	"github.com/stretchr/testify/suite"
)

type HealthcheckSuite struct {
	ExampleTestSuite
}

func (suite *HealthcheckSuite) SetupSuite() {
	suite.ExampleTestSuite.SetupSuite()
}

func (suite *HealthcheckSuite) SetupTest() {
	suite.ExampleTestSuite.SetupTest()
}

func TestHealthcheckSuite(t *testing.T) {
	suite.Run(t, new(HealthcheckSuite))
}

func (suite *HealthcheckSuite) TestHandleGetHealthcheckRequest() {
	// Make http request
	response := suite.ExecuteJsonRequest(
		http.MethodGet,
		"/system/healthcheck",
		nil,
		helpers.EmptyHeaders(),
	)

	suite.CheckResponse(
		http.StatusOK,
		`{"data":{"type":"healthcheck","id":"<<PRESENCE>>","attributes":{"service_name":"cqrs-monitored-app","status":{"system":"ok"}}}}`,
		response,
	)
}
