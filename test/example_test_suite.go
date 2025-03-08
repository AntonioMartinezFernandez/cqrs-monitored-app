package test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/cmd/di"
	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/test/arrangers"
	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/test/helpers"

	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var routerHandler http.Handler

type ExampleTestSuite struct {
	suite.Suite

	Ctx            context.Context
	CommonServices *di.CommonServices
	SystemServices *di.SystemModuleServices
	BookServices   *di.BookModuleServices
	HttpServices   *di.HttpServices

	SuiteArranger         *arrangers.IdentityOperationsDispatcherSuiteArranger
	IntervalExecutorSuite *helpers.IntervalExecutorTestSuite
}

func (suite *ExampleTestSuite) SetupSuite() {
	if suite.CommonServices == nil {
		suite.CommonServices = di.InitCommonServicesWithEnvFiles("../.env.example", "../.env.example.test")
	}
	if suite.HttpServices == nil {
		suite.HttpServices = di.InitHttpServices(suite.CommonServices)
	}
	if suite.SystemServices == nil {
		suite.SystemServices = di.InitSystemModuleServices(suite.CommonServices, suite.HttpServices)
	}
	if suite.BookServices == nil {
		suite.BookServices = di.InitBookModuleServices(suite.CommonServices, suite.HttpServices)
	}

	suite.Ctx = context.Background()

	suite.SuiteArranger = arrangers.NewIdentityOperationsDispatcherSuiteArranger(suite.CommonServices)
}

func (suite *ExampleTestSuite) TearDownSuite() {
}

func (suite *ExampleTestSuite) SetupTest() {
	suite.SuiteArranger.Arrange(suite.Ctx)
}

func (suite *ExampleTestSuite) TearDownTest() {
}

func (suite *ExampleTestSuite) ExecuteJsonRequest(verb string, path string, body []byte, headers map[string]string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(verb, path, bytes.NewBuffer(body))
	if len(headers) != 0 {
		for headerName, value := range headers {
			req.Header.Set(headerName, value)
		}
	}

	assert.NoError(suite.T(), err)

	req.Header.Set("Content-Type", "application/json")
	return suite.ExecuteRequest(req)
}

func (suite *ExampleTestSuite) ExecuteRequest(req *http.Request) *httptest.ResponseRecorder {
	resRec := httptest.NewRecorder()

	if routerHandler == nil {
		routerHandler = suite.HttpServices.Router.GetMuxRouter()
	}

	routerHandler.ServeHTTP(resRec, req)

	return resRec
}

func (suite *ExampleTestSuite) CheckResponse(expectedStatusCode int, expectedResponse string, response *httptest.ResponseRecorder, formats ...interface{}) {
	ja := jsonassert.New(suite.T())
	suite.CheckResponseCode(expectedStatusCode, response.Code)

	receivedResponse := response.Body.String()
	if receivedResponse == "" {
		assert.Equal(suite.T(), expectedResponse, receivedResponse)
		return
	}
	if formats != nil {
		ja.Assertf(receivedResponse, expectedResponse, formats)
	} else {
		ja.Assertf(receivedResponse, expectedResponse)
	}
}

func (suite *ExampleTestSuite) CheckResponseCode(expected, actual int) {
	if expected != actual {
		suite.T().Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
