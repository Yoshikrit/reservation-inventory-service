package entity

type ProductStockHistory struct {
	ID        int64  `gorm:"primaryKey;autoIncrement;type:bigserial"`
	ProductID string `gorm:"size:14;not null;index"`
	OldQty    uint   `gorm:"not null"`
	NewQty    uint   `gorm:"not null"`
	Delta     int    `gorm:"not null"` // negative = deduct, positive = replenish/create
	Reason    string `gorm:"size:50;not null"`
	AuditModel
}
