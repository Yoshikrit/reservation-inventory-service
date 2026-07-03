package entity

func (p Product) TableName() string {
	return "products"
}

type Product struct {
	ProductID   string  `gorm:"primaryKey;type:varchar(14);not null"`
	Name        string  `gorm:"type:varchar(100);index;not null"`
	Description string  `gorm:"type:text"`
	Price       float64 `gorm:"type:decimal(12,2);not null;check:price >= 0"`
	Quantity    uint    `gorm:"type:int;not null;default:0;check:quantity >= 0"`
	AuditModel  `gorm:"embedded"`
}
