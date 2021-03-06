package present

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt"
	"os"
	"qrcode/access/constant"
	"qrcode/control"
	"time"
)

type ContextApi struct {
	fiber.Ctx
	apiControl *control.APIControl
}

func APICreate(ctrl *control.APIControl) {
	app := fiber.New(fiber.Config{})

	app.Use(logger.New(logger.Config{
		Next:         nil,
		Format:       "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat:   "15:04:05",
		TimeZone:     "Local",
		TimeInterval: 500 * time.Millisecond,
		Output:       os.Stderr,
	}))
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",                              // todo รอคีย์ domain จะมาจาก env
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH", //
	}))

	app.Use(func(context *fiber.Ctx) error {
		context.Locals(constant.LocalsKeyControl, ctrl)
		return context.Next()
	})

	api := app.Group("/api")
	api.Post("admin", admin) // todo สำหรับ สมัคร admin เท่านั้น

	qr := app.Group("/qr")
	qr.Get("*", getDataQrCode) // ตอน ScanQrCode

	qrApi := api.Group("/qr-api")
	qrApi.Get("getDataQrCodeJson/:id", getDataQrCodeJson)

	// -- Todo Owner
	owner := api.Group("/owner")
	owner.Post("login", LoginOwner)
	owner.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(constant.SecretKey),
		SuccessHandler: func(context *fiber.Ctx) error {
			user := context.Locals("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			var userRole = claims["role"]
			if userRole == string(constant.Owner) {
				return context.Next()
			} else {
				return context.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Unauthorized",
				})
			}
			return context.Next()
		},
		ErrorHandler: AuthError,
		AuthScheme:   "Bearer",
	}))
	// -- API Owner Account
	owner.Get("getAccount", getAccount)
	owner.Post("register_operator", registerOperatorOwner)
	owner.Get("getOperator", getOperator) // todo ดูข้อมูลทั่งหมดของ Operator ById Owner
	owner.Get("getOperatorById/:id", getOperatorById)
	owner.Put("updateProfile/:id", updateProfile)
	owner.Put("changePasswordOwner", ChangePasswordOwner)
	owner.Put("changePasswordOperator", ChangePasswordOperatorByOwner)
	owner.Delete("deleteAccount/:id", deleteAccountOperator)

	// QrCode
	owner.Put("updateStatusQrCode/:id", updateStatusQrCodeOwner)
	owner.Get("getQrCode", getQrCodeOwnerById) // Id >>> OwnerId
	owner.Put("updateDataQrCode", updateDataQrCode)
	owner.Post("insertDataQrCode", insertDataQrCode)
	owner.Post("updateHistoryInfoDataQrCode", updateHistoryInfoDataQrCode)
	owner.Post("updateOpsDataQrCode", updateOpsDataQrCode)
	owner.Get("getTemplate", getTemplateList)


	// officeEquipment
	owner.Get("getOperatorLine", getOperatorLine)
	owner.Put("updateOption", UpdateOption)
	owner.Put("worksheet/:id", worksheet)
	owner.Get("reportAll/:id", ownerGetWorksheet)
	owner.Get("report/:id", ownerGetWorksheetById)
	owner.Put("worksheet/:id", ownerWorksheet)
	owner.Get("getDataUpdate/:id", ownerGetUpdateWorksheet)
	owner.Put("report/:id", ownerUpdateWorksheet)
	owner.Delete("report/:id", ownerDeleteWorksheet)


	// ops
	ops := api.Group("/ops")
	ops.Post("login", LoginOperator)
	ops.Get("getAccount/:id", getAccountByLineId)

	// insertDataQrCode And UpdateDataQrCode
	ops.Post("insertDataQrCode", insertDataQrCodeOps)
	ops.Put("updateDataQrCode", updateDataQrCodeOps)

	// officeEquipment
	ops.Post("getWorksheet/:id", getWorksheet)
	ops.Get("getTemplate", getTemplateList)
	ops.Get("getTemplate/:id", getTemplate)
	ops.Get("typeReport", getTypeWorksheet)
	ops.Get("report/:id", getWorksheet)
	ops.Get("report", getWorksheetById)
	ops.Post("report", insertWorksheet)
	ops.Put("worksheet/:id", worksheet)
	ops.Get("getDataUpdate/:id", getUpdateWorksheet)
	ops.Put("report/:id", updateWorksheet)
	ops.Delete("report/:id", deleteWorksheet)

	//ops.Use(jwtware.New(jwtware.Config{
	//	SigningKey: []byte(constant.SecretKey),
	//	SuccessHandler: func(context *fiber.Ctx) error {
	//		user := context.Locals("user").(*jwt.Token)
	//		claims := user.Claims.(jwt.MapClaims)
	//		var userRole = claims["role"]
	//		if userRole == string(constant.Operator) {
	//			return context.Next()
	//		} else {
	//			return context.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	//				"error": "Unauthorized",
	//			})
	//		}
	//		return context.Next()
	//	},
	//	ErrorHandler: AuthError,
	//	AuthScheme:   "Bearer",
	//}))

	// -- Todo Admin
	admin := api.Group("/admin")
	admin.Post("login", LoginAdmin)
	admin.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(constant.SecretKey),
		SuccessHandler: func(context *fiber.Ctx) error {
			user := context.Locals("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			var userRole = claims["role"]
			if userRole == string(constant.Admin) {
				return context.Next()
			} else {
				return context.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Unauthorized",
				})
			}
			return context.Next()
		},
		ErrorHandler: AuthError,
		AuthScheme:   "Bearer",
	}))

	// -- Account
	admin.Post("register_owner", registerOwner)
	admin.Post("register_operator", registerOperator)
	admin.Get("getAccount", getAccount)
	admin.Get("getAllAccountOwner", getAllAccountOwner)
	admin.Get("getSubOwner/:id", getSubOwner) // todo ดูข้อมูลทั่งหมดของ Operator ById Owner
	admin.Get("getAllAccountOperator", getAllAccountOperator)
	admin.Get("getOwnerByIdOps/:id", getOwnerByIdOps) // todo ดูข้อมูล Owner ById Ops ยังต้องแก้ SQL ยังไม่ได้ join
	admin.Get("getAccountById/:id", getAccountById)
	admin.Put("updateProfile/:id", updateProfile)
	admin.Put("changePassword/:id", changePassword)
	admin.Delete("deleteAccount/:id", deleteAccount)

	// -- QrCode
	admin.Get("getTemplate", getTemplateList)
	admin.Get("getDateQrCodeById/:id", getDataQrCodeJson)
	admin.Post("createQrCode", createQrCode)
	admin.Post("genQrCodeToFileZipByTemplateName", genQrCodeToFileZipByTemplateName)
	admin.Post("genQrCodeToFileZipByOwner", genQrCodeToFileZipByOwner)
	admin.Post("genQrCodeToFileZipByQrCodeId", genQrCodeToFileZipByQrCodeId)
	admin.Get("getAllQrCodeByOwnerId/:id", getQrCodeById) // Id >>> OwnerId
	admin.Get("getAllQrCode", getAllQrCode)
	admin.Post("insertDataQrCode", insertDataQrCode)
	admin.Post("updateHistoryInfoDataQrCode", updateHistoryInfoDataQrCode)
	admin.Post("updateOpsDataQrCode", updateOpsDataQrCode)
	admin.Delete("delQrCode", deleteQrCode) // todo ลบ QrCode
	admin.Put("updateStatusQrCode/:id", updateStatusQrCode)


	//officeEquipment


	// -- TeamPage
	//admin.Get("TestQrCode", TestQrCode)
	_ = app.Listen(":8000")
}

func ValidateStruct(dataStruct interface{}) error {
	validate := validator.New()
	err := validate.Struct(dataStruct)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(fmt.Sprintf("%s: %s", err.StructField(), err.Tag()))
		}
	} else {
		return nil
	}
	return err
}
