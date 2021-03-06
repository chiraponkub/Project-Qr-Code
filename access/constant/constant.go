package constant

import "errors"

const (
	LocalsKeyControl string = "CTRL"
	SecretKey = "T-DEV Co., Ltd."
	Http = "http://www.localhost:8080/qr"
)



type UserRole string

const (
	Owner    UserRole = "owner"
	Admin    UserRole = "admin"
	Operator UserRole = "operator"
)

var UserRoleData = []UserRole{
	Owner,
	Admin,
	Operator,
}

func (userRole UserRole) Role() (result *string, Errors error) {
	switch userRole {
	case Owner:
		fullName := "เจ้าของ"
		result = &fullName
	case Admin:
		fullName := "ผู้ดูแลระบบ"
		result = &fullName
	case Operator:
		fullName := "ช่างซ่อมบำรุง"
		result = &fullName
	default:
		Errors = errors.New("unimplemented")
	}
	return
}

type CategoryFile string

const (
	QRCode   CategoryFile = "qr_code"
)

var CategoryFileData = []CategoryFile{
	QRCode,
}

type QrCode string

const (
	SaveFileLocationQrCode  QrCode = "fileqrcode"
	SaveFileLocationZipFile QrCode = "zipfile"
	SaveFileLocationExposed QrCode = "Exposed File QrCode"
)



