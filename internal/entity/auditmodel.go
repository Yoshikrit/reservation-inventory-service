package entity

import (
	"time"

	"gorm.io/gorm"
)

type contextKey string

const (
	ContextKeyEvent   contextKey = "event"
	ContextKeyTraceID contextKey = "trace_id"
)

type AuditModel struct {
	CreatedAt        time.Time `gorm:"autoCreateTime;not null"`
	CreatedByEvent   string    `gorm:"size:100;not null;default:''"`
	CreatedByTraceID string    `gorm:"size:64;not null;default:''"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime;not null"`
	UpdatedByEvent   string    `gorm:"size:100;not null;default:''"`
	UpdatedByTraceID string    `gorm:"size:64;not null;default:''"`
}

func (a *AuditModel) BeforeCreate(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	if ctx == nil {
		return nil
	}
	event, _ := ctx.Value(ContextKeyEvent).(string)
	traceID, _ := ctx.Value(ContextKeyTraceID).(string)
	a.CreatedByEvent = event
	a.CreatedByTraceID = traceID
	a.UpdatedByEvent = event
	a.UpdatedByTraceID = traceID
	tx.Statement.SetColumn("created_by_event", event)
	tx.Statement.SetColumn("created_by_trace_id", traceID)
	tx.Statement.SetColumn("updated_by_event", event)
	tx.Statement.SetColumn("updated_by_trace_id", traceID)
	return nil
}

func (a *AuditModel) BeforeUpdate(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	if ctx == nil {
		return nil
	}
	event, _ := ctx.Value(ContextKeyEvent).(string)
	traceID, _ := ctx.Value(ContextKeyTraceID).(string)
	a.UpdatedByEvent = event
	a.UpdatedByTraceID = traceID
	tx.Statement.SetColumn("updated_by_event", event)
	tx.Statement.SetColumn("updated_by_trace_id", traceID)
	return nil
}
