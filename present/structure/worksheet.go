package structure

import (
	"github.com/gofrs/uuid"
	"qrcode/access/constant"
	"time"
)

type TypeWorksheet struct {
	Data []TypeWorksheets `json:"data"`
}
type TypeWorksheets struct {
	TypeWorksheet constant.Worksheets `json:"type_worksheet"`
}

type GetWorksheet struct {
	Info      interface{} `json:"info,omitempty"`
	Worksheet []Worksheet `json:"worksheet,omitempty"`
}

type Worksheet struct {
	ID              uint              `json:"id,omitempty"`
	QrCodeID        uuid.UUID         `json:"qr_code_id,omitempty"`
	Info            interface{}       `json:"info,omitempty"`
	Text            string            `json:"text,omitempty"`
	Type            string            `json:"type,omitempty"`
	Ops             *string           `json:"ops,omitempty"`
	OwnerId         uint              `json:"owner_id,omitempty"`
	StatusWorksheet []StatusWorksheet `json:"status_worksheet"`
}

type StatusWorksheet struct {
	Status     string      `json:"status"`
	UpdateAt   *time.Time  `json:"update_at"`
	Text       *string      `json:"text,omitempty"`
	Equipments []Equipment `json:"equipments,omitempty"`
}

type Equipment struct {
	NameEquipment string `json:"name_equipment,omitempty"`
}

//type StatusWorksheet struct {
//	StatusWorksheet1 StatusWorksheet1 `json:"status_worksheet_1"`
//	StatusWorksheet2 StatusWorksheet2 `json:"status_worksheet_2"`
//	StatusWorksheet3 StatusWorksheet3 `json:"status_worksheet_3"`
//}

//type StatusWorksheet1 struct {
//	Status   string    `json:"status"`
//	UpdateAt time.Time `json:"update_at"`
//}
//
//type StatusWorksheet2 struct {
//	Status   string     `json:"status"`
//	UpdateAt *time.Time `json:"update_at"`
//}

//type StatusWorksheet3 struct {
//	Status   string     `json:"status"`
//	UpdateAt *time.Time `json:"update_at"`
//	Text     string     `json:"text"`
//}

type ReportID struct {
	ReportID   uint   `json:"report_id" query:"report_id"`
	LineUserId string `json:"line_user_id" query:"line_user_id"`
}

type InsertWorksheet struct {
	QrCodeID uuid.UUID `json:"qr_code_id,omitempty" validate:"required"`
	Text     string    `json:"text,omitempty" validate:"required"`
	Type     string    `json:"type,omitempty" validate:"required"`
}

type UpdateWorksheet struct {
	LineUserId string `json:"line_user_id" validate:"required"`
	Text       string `json:"text,omitempty" validate:"required"`
	Equipments []Equipment `json:"equipments,omitempty"`
}
