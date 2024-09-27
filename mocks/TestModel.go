package mocks

type TestModel struct {
	ID    uint   `gorm:"primaryKey"`
	Email string `gorm:"unique"`
}
