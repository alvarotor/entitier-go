package service_test

// var ctx = context.Background()

// func TestGenericService_GetAll_Success(t *testing.T) {
// 	mockRepo := new(mocks.IGenericRepo[mocks.TestModel, uint])

// 	mockRepo.On("GetAll", ctx).Return([]*mocks.TestModel{
// 		{Email: "test1@example.com"},
// 		{Email: "test2@example.com"},
// 	}, nil)

// 	service := service.NewGenericService[mocks.TestModel, uint](mockRepo)
// 	result, err := service.GetAll(ctx)

// 	assert.NoError(t, err)
// 	assert.Equal(t, 2, len(result))
// 	assert.Equal(t, "test1@example.com", result[0].Email)
// }

// func TestGenericService_GetAll_Error(t *testing.T) {
// 	mockRepo := new(mocks.IGenericRepo[mocks.TestModel, uint])

// 	mockRepo.On("GetAll", ctx).Return([]*mocks.TestModel{}, errors.New("DB error"))

// 	service := service.NewGenericService[mocks.TestModel, uint](mockRepo)
// 	result, err := service.GetAll(ctx)

// 	assert.Error(t, err)
// 	assert.Nil(t, result)
// 	assert.Equal(t, "DB error", err.Error())
// }

// func TestGenericService_Get_Success(t *testing.T) {
// 	mockRepo := new(mocks.IGenericRepo[mocks.TestModel, uint])

// 	mockRepo.On("Get", ctx, uint(1), "").Return(&mocks.TestModel{ID: 1, Email: "test1@example.com"}, nil)

// 	service := service.NewGenericService[mocks.TestModel, uint](mockRepo)
// 	result, err := service.Get(ctx, 1, "")

// 	assert.NoError(t, err)
// 	assert.Equal(t, "test1@example.com", result.Email)
// }

// func TestGenericService_Get_Error(t *testing.T) {
// 	mockRepo := new(mocks.IGenericRepo[mocks.TestModel, uint])

// 	mockRepo.On("Get", ctx, uint(1), "").Return(nil, gorm.ErrRecordNotFound)

// 	service := service.NewGenericService[mocks.TestModel, uint](mockRepo)
// 	result, err := service.Get(ctx, 1, "")

// 	assert.Error(t, err)
// 	assert.Nil(t, result)
// 	assert.Equal(t, gorm.ErrRecordNotFound, err)
// }

// func TestGenericService_Create_Success(t *testing.T) {
// 	mockRepo := new(mocks.IGenericRepo[mocks.TestModel, uint])

// 	inputModel := mocks.TestModel{Email: "test@example.com"}
// 	createdModel := mocks.TestModel{ID: 1, Email: "test@example.com"}

// 	mockRepo.On("Create", ctx, inputModel).Return(createdModel, nil)

// 	service := service.NewGenericService[mocks.TestModel, uint](mockRepo)

// 	result, err := service.Create(ctx, inputModel)

// 	assert.NoError(t, err)
// 	assert.Equal(t, "test@example.com", result.Email)
// 	assert.Equal(t, createdModel, result)
// }

// func TestGenericService_Create_Error(t *testing.T) {
// 	mockRepo := new(mocks.IGenericRepo[mocks.TestModel, uint])

// 	inputModel := mocks.TestModel{Email: "test@example.com"}
// 	mockRepo.On("Create", ctx, inputModel).Return(mocks.TestModel{}, errors.New("DB error"))

// 	service := service.NewGenericService[mocks.TestModel, uint](mockRepo)

// 	result, err := service.Create(ctx, inputModel)

// 	assert.Error(t, err)
// 	assert.Equal(t, inputModel, result)
// 	assert.Equal(t, "DB error", err.Error())
// }

// func TestGenericService_Delete_Success(t *testing.T) {
// 	mockRepo := new(mocks.IGenericRepo[mocks.TestModel, uint])

// 	mockRepo.On("Delete", ctx, uint(1), true).Return(nil)

// 	service := service.NewGenericService[mocks.TestModel, uint](mockRepo)
// 	err := service.Delete(ctx, 1, true)

// 	assert.NoError(t, err)
// }

// func TestGenericService_Delete_Error(t *testing.T) {
// 	mockRepo := new(mocks.IGenericRepo[mocks.TestModel, uint])

// 	mockRepo.On("Delete", ctx, uint(1), true).Return(gorm.ErrRecordNotFound)

// 	service := service.NewGenericService[mocks.TestModel, uint](mockRepo)
// 	err := service.Delete(ctx, 1, true)

// 	assert.Error(t, err)
// 	assert.Equal(t, gorm.ErrRecordNotFound, err)
// }

// func TestGenericService_Update_Success(t *testing.T) {
// 	mockRepo := new(mocks.IGenericRepo[mocks.TestModel, uint])

// 	updatedModel := mocks.TestModel{Email: "updated@example.com"}

// 	mockRepo.On("Update", ctx, uint(1), updatedModel).Return(nil)

// 	service := service.NewGenericService[mocks.TestModel, uint](mockRepo)

// 	err := service.Update(ctx, 1, updatedModel)

// 	assert.NoError(t, err)
// 	assert.Equal(t, "updated@example.com", updatedModel.Email)
// }

// func TestGenericService_Update_Error(t *testing.T) {
// 	mockRepo := new(mocks.IGenericRepo[mocks.TestModel, uint])

// 	updatedModel := mocks.TestModel{ID: 1, Email: "updated@example.com"}

// 	mockRepo.On("Update", ctx, uint(1), updatedModel).Return(gorm.ErrRecordNotFound)

// 	service := service.NewGenericService[mocks.TestModel, uint](mockRepo)

// 	err := service.Update(ctx, 1, updatedModel)

// 	assert.Error(t, err)
// 	assert.Equal(t, gorm.ErrRecordNotFound, err)
// }
