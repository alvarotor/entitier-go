package repositories

type MockModel struct {
	ID    uint   `gorm:"primaryKey"`
	Email string `gorm:"unique"`
	Name  string
}
