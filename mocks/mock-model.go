package mocks

type MockModel struct {
	ID    uint   `gorm:"primaryKey"`
	Email string `gorm:"unique"`
	Name  string
}
