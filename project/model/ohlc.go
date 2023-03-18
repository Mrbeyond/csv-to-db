package model

import "mime/multipart"

type Ohcl struct {
	ID     uint64  `json:"-" gorm:"primaryKey;autoIncrement"`
	UNIX   uint64  `json:"unix" binding:"required" gorm:"not null"`
	SYMBOL string  `json:"symbol" binding:"required" gorm:"not null"`
	OPEN   float32 `json:"open" binding:"required" gorm:"not null"`
	HIGH   float32 `json:"high" binding:"required" gorm:"not null"`
	LOW    float32 `json:"low" binding:"required" gorm:"not null"`
	CLOSE  float32 `json:"close" binding:"required" gorm:"not null"`
}

type CreatePayload struct {
	CSVFile multipart.FileHeader `json:"scv_file" form:"scv_file" bindging:"required"`
}
