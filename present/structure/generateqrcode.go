package structure

import (
	"github.com/gofrs/uuid"
	"gorm.io/datatypes"
	"time"
)

type FileZipByTemplateName struct {
	OwnerId      uint   `json:"owner_id" validate:"required"`
	TemplateName string `json:"template_name"`
}

type FileZipByOwner struct {
	OwnerId uint `json:"owner_id" validate:"required"`
}

type FileZip struct {
	OwnerId  int      `json:"owner_id" validate:"required"`
	QrCodeId []string `json:"qr_code_id" validate:"required"`
}

type GetQrCode struct {
	OwnerId       uint      `json:"owner_id"`
	OwnerUsername string    `json:"owner_username"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	TemplateName  string    `json:"template_name"`
	QrCodeId      string    `json:"qr_code_id"`
	CodeName      string    `json:"code_name"`
	URL           string    `json:"url"`
	Active        bool      `json:"active"`
}

type GetDataQrCode struct {
	QrCodeId     string         `json:"qr_code_id"`
	Info         datatypes.JSON `json:"info"`
	HistoryInfo  []GetHistory   `json:"history_info"`
	Ops          []GetOps       `json:"ops"`
	OwnerId      int            `json:"owner_id"`
	OwnerName    string         `json:"owner_name"`
	TemplateName string         `json:"template_name"`
	CodeName     string         `json:"code_name"`
}

type GetOps struct {
	Ops       datatypes.JSON
	User      string
	Role      string
}

type GetHistory struct {
	HistoryInfo datatypes.JSON
	User        string
	UpdatedAt   time.Time
	Role        string
}

type ArrayFileName struct {
	FileName string `json:"file_name"`
}

type GetQrCodeImage struct {
	FileName string `json:"file_name"`
}

type GenQrCode struct {
	OwnerId      uint   `json:"owner_id" validate:"required"`
	CodeName     string `json:"code_name" validate:"required"`
	TemplateName string `json:"template_name"`
	Amount       int    `json:"amount" validate:"required"`
}

type LineQrCodeId struct {
	QrCodeId uuid.UUID   `json:"qr_code_id" query:"qr_id"`
}

type UpdateDataQrCode struct {
	OwnerId  uint        `json:"owner_id" validate:"required"`
	QrCodeId uuid.UUID   `json:"qr_code_id" validate:"required"`
	Info     interface{} `json:"info" validate:"required"`
	LineUserId string `json:"line_user_id"`
	//HistoryInfo interface{} `json:"history_info" validate:"required"`
}

type InsertDataQrCode struct {
	OwnerId      uint        `json:"owner_id" validate:"required"`
	QrCodeId     uuid.UUID   `json:"qr_code_id" validate:"required"`
	TemplateName string      `json:"template_name" validate:"required"`
	LineUserId   string      `json:"line_user_id" validate:"required"`
	Info         interface{} `json:"info" validate:"required"`
}

type UpdateHistoryInfoDataQrCode struct {
	UserId      uint        `json:"user_id" validate:"required"`
	OwnerId     uint        `json:"owner_id" validate:"required"`
	QrCodeId    uuid.UUID   `json:"qr_code_id" validate:"required"`
	HistoryInfo interface{} `json:"history_info" validate:"required"`
}

type UpdateOpsDataQrCode struct {
	UserId   uint        `json:"user_id" validate:"required"`
	OwnerId  uint        `json:"owner_id" validate:"required"`
	QrCodeId uuid.UUID   `json:"qr_code_id" validate:"required"`
	Ops      interface{} `json:"ops" validate:"required"`
}

type DelQrCode struct {
	QrCodeId []string `json:"qr_code_id" validate:"required"`
}

type StatusQrCode struct {
	Active *bool `json:"active"`
}
