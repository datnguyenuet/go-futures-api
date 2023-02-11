package models

type Position struct {
	PositionId int64 `json:"position_id,omitempty" binding:"required"`
}
