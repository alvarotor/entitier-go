package controllers

// // MockLogger is a mock implementation of the logger.Logger interface
// type MockLogger struct {
// 	mock.Mock
// }

// func (m *MockLogger) Error(msg string) {
// 	m.Called(msg)
// }

// func (m *MockLogger) Debug(msg string) {
// 	m.Called(msg)
// }

// func (m *MockLogger) Info(msg string) {
// 	m.Called(msg)
// }

// func (m *MockLogger) Warn(msg string) {
// 	m.Called(msg)
// }

// // Implement any other methods required by the logger.Logger interface

// // MockGenericService is a mock implementation of the services.GenericService interface
// type MockGenericService[T any, X string | uint] struct {
// 	mock.Mock
// }

// func (m *MockGenericService[T, X]) Create(model T) (T, error) {
// 	args := m.Called(model)
// 	return args.Get(0).(T), args.Error(1)
// }

// func (m *MockGenericService[T, X]) Get(id X, preload string) (*T, error) {
// 	args := m.Called(id, preload)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*T), args.Error(1)
// }

// func (m *MockGenericService[T, X]) GetAll() ([]*T, error) {
// 	args := m.Called()
// 	return args.Get(0).([]*T), args.Error(1)
// }

// func (m *MockGenericService[T, X]) Delete(id X, permanently bool) error {
// 	args := m.Called(id, permanently)
// 	return args.Error(0)
// }

// func (m *MockGenericService[T, X]) Update(id X, amended T) error {
// 	args := m.Called(id, amended)
// 	return args.Error(0)
// }

// // TestModel is a sample struct for testing
// type TestModel struct {
// 	ID   uint   `json:"id"`
// 	Name string `json:"name"`
// }

// var _ services.GenericService[TestModel, uint] = (*MockGenericService[TestModel, uint])(nil)

// func TestNewGenericController(t *testing.T) {
// 	mockLogger := new(MockLogger)
// 	mockDB := &gorm.DB{}

// 	controller := NewGenericController[TestModel, uint](mockLogger, mockDB)

// 	assert.NotNil(t, controller)
// 	assert.IsType(t, &controllerGeneric[TestModel, uint]{}, controller)
// }

// func TestCreate(t *testing.T) {
// 	mockLogger := new(MockLogger)
// 	mockService := new(MockGenericService[TestModel, uint])

// 	controller := &controllerGeneric[TestModel, uint]{
// 		svcT: mockService,
// 		log:  mockLogger,
// 	}

// 	testModel := TestModel{Name: "Test"}
// 	expectedModel := TestModel{ID: 1, Name: "Test"}

// 	mockService.On("Create", testModel).Return(expectedModel, nil)

// 	result, err := controller.Create(testModel)

// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedModel, result)
// 	mockService.AssertExpectations(t)
// }

// func TestGet(t *testing.T) {
// 	mockLogger := new(MockLogger)
// 	mockService := new(MockGenericService[TestModel, uint])

// 	controller := &controllerGeneric[TestModel, uint]{
// 		svcT: mockService,
// 		log:  mockLogger,
// 	}

// 	gin.SetMode(gin.TestMode)
// 	w := httptest.NewRecorder()
// 	c, _ := gin.CreateTestContext(w)

// 	c.Params = gin.Params{{Key: "id", Value: "1"}}

// 	expectedModel := TestModel{ID: 1, Name: "Test"}

// 	mockService.On("Get", uint(1), "User").Return(&expectedModel, nil)

// 	controller.Get(c)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var response map[string]TestModel
// 	err := json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedModel, response["item"])

// 	mockService.AssertExpectations(t)
// }

// // ... (rest of the test cases remain the same, just update the mock service type)

// func TestGetInvalidID(t *testing.T) {
// 	mockLogger := new(MockLogger)
// 	mockService := new(MockGenericService[TestModel, uint])

// 	controller := &controllerGeneric[TestModel, uint]{
// 		svcT: mockService,
// 		log:  mockLogger,
// 	}

// 	gin.SetMode(gin.TestMode)
// 	w := httptest.NewRecorder()
// 	c, _ := gin.CreateTestContext(w)

// 	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

// 	mockLogger.On("Error", "must provide valid id").Return()

// 	controller.Get(c)

// 	assert.Equal(t, http.StatusBadRequest, w.Code)

// 	var response map[string]string
// 	err := json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "must provide valid id", response["err"])

// 	mockLogger.AssertExpectations(t)
// }

// func TestGetServiceError(t *testing.T) {
// 	mockLogger := new(MockLogger)
// 	mockService := new(MockGenericService[TestModel, uint])

// 	controller := &controllerGeneric[TestModel, uint]{
// 		svcT: mockService,
// 		log:  mockLogger,
// 	}

// 	gin.SetMode(gin.TestMode)
// 	w := httptest.NewRecorder()
// 	c, _ := gin.CreateTestContext(w)

// 	c.Params = gin.Params{{Key: "id", Value: "1"}}

// 	mockService.On("Get", uint(1), "User").Return((*TestModel)(nil), errors.New("service error"))
// 	mockLogger.On("Error", "service error").Return()

// 	controller.Get(c)

// 	assert.Equal(t, http.StatusInternalServerError, w.Code)

// 	var response map[string]string
// 	err := json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "service error", response["err"])

// 	mockService.AssertExpectations(t)
// 	mockLogger.AssertExpectations(t)
// }
