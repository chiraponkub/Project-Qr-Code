package control

import (
	"archive/zip"
	"encoding/json"
	"errors"
	uuid2 "github.com/gofrs/uuid"
	"github.com/yeqown/go-qrcode"
	"gorm.io/datatypes"
	"io"
	"os"
	"qrcode/access/constant"
	rdbmsstructure "qrcode/access/rdbms/structure"
	"qrcode/present/structure"
	"qrcode/utility"
	"strconv"
	"strings"
)

func (ctrl *APIControl) GetAllQrCode() (response []structure.GetQrCode, Error error) {
	var getQrCodeArray []structure.GetQrCode

	data, err := ctrl.access.RDBMS.GetAllQrCode()
	if err != nil {
		Error = errors.New("ไม่มีข้อมูล")
		return
	}
	for _, res := range data {
		owner, _ := ctrl.access.RDBMS.GetAccount(int(res.OwnerId))
		resGetQrCode := structure.GetQrCode{
			OwnerId:       res.OwnerId,
			OwnerUsername: owner.Username,
			CreatedAt:     res.CreatedAt,
			UpdatedAt:     res.UpdatedAt,
			TemplateName:  res.TemplateName,
			QrCodeId:      res.QrCodeUUID.String(),
			CodeName:      res.Code + "-" + res.Count,
			URL:           ctrl.access.ENV.URLQRCode + res.QrCodeUUID.String(),
			Active:        res.Active,
		}
		getQrCodeArray = append(getQrCodeArray, resGetQrCode)
	}
	response = getQrCodeArray
	return
}

func (ctrl *APIControl) GetQrCodeById(OwnerId int) (response []structure.GetQrCode, Error error) {
	var getQrCodeArray []structure.GetQrCode
	check, err := ctrl.access.RDBMS.CheckAccountId(uint(OwnerId))
	if err != nil {
		Error = errors.New("ไม่มีผู้ใช้คนนี้ในระบบ")
		return
	}
	if !(check.Role == string(constant.Owner)) {
		Error = errors.New("สิทธิ์ผู้ใช้งานไม่ถูกต้อง")
		return
	}
	data, err := ctrl.access.RDBMS.GetQrCodeByOwnerId(OwnerId)
	if err != nil {
		Error = err
		return
	}
	for _, res := range data {
		owner, _ := ctrl.access.RDBMS.GetAccount(int(res.OwnerId))
		resGetQrCode := structure.GetQrCode{
			OwnerId:       res.OwnerId,
			OwnerUsername: owner.Username,
			CreatedAt:     res.CreatedAt,
			UpdatedAt:     res.UpdatedAt,
			TemplateName:  res.TemplateName,
			QrCodeId:      res.QrCodeUUID.String(),
			CodeName:      res.Code + "-" + res.Count,
			URL:           ctrl.access.ENV.URLQRCode + res.QrCodeUUID.String(),
			Active:        res.Active,
		}
		getQrCodeArray = append(getQrCodeArray, resGetQrCode)
	}
	response = getQrCodeArray
	return
}

func (ctrl *APIControl) InsertDataQrCode(req *structure.InsertDataQrCode) (Error error) {
	check, err := ctrl.access.RDBMS.GetQrCodeByQrCodeId(int(req.OwnerId), req.QrCodeId.String())
	if err != nil {
		Error = errors.New("ไม่พบ QrCode นี้อยู่ในระบบ")
		return
	}
	if check.TemplateName != "" {
		Error = errors.New("QrCode ได้ถูกตั้งค่า Template แล้ว")
		return
	}

	b, err := json.Marshal(req.Info)
	if err != nil {
		Error = err
		return
	}

	data := rdbmsstructure.QrCode{
		TemplateName: req.TemplateName,
		Info:         datatypes.JSON(b),
		QrCodeUUID:   req.QrCodeId,
		First:        true,
	}

	err = ctrl.access.RDBMS.UpdateQrCodeById(data)
	if err != nil {
		Error = err
		return
	}

	return
}

func (ctrl *APIControl) GetDataQrCode(QrCodeId string) (response structure.GetDataQrCode, Error error) {
	data, err := ctrl.access.RDBMS.GetDataQrCode(QrCodeId)
	if err != nil {
		Error = errors.New("ไม่มี QrCode นี้อยู่ในระบบ")
		return
	}
	var HistoryArray []structure.GetHistory
	History, err := ctrl.access.RDBMS.GetHistory(QrCodeId)

	for _, h := range History {
		dataHistory := structure.GetHistory{
			HistoryInfo: h.HistoryInfo,
			UserId:      h.UserId,
			UpdatedAt:   h.UpdatedAt,
		}
		HistoryArray = append(HistoryArray, dataHistory)
	}

	response = structure.GetDataQrCode{
		QrCodeId:     data.QrCodeUUID.String(),
		Info:         data.Info,
		Ops:          data.Ops,
		OwnerId:      int(data.OwnerId),
		TemplateName: data.TemplateName,
		CodeName:     data.Code + "-" + data.Count,
		HistoryInfo:  HistoryArray,
	}
	return
}

func (ctrl *APIControl) DeleteQrCode(req structure.DelQrCode) (Error error) {
	if len(req.QrCodeId) == 0 {
		Error = errors.New("ไม่มี Qr-Code ถูกส่งมา")
	}
	for _, del := range req.QrCodeId {
		_, err := ctrl.access.RDBMS.GetDataQrCode(del)
		if err != nil {
			Error = errors.New("Qr-Code ที่จะลบไม่มีอยู่ในระบบ")
			return
		}
		err = ctrl.access.RDBMS.DeleteQrCode(del)
		if err != nil {
			Error = err
			return
		}
	}
	return
}

func (ctrl *APIControl) UpdateStatusQrCode(QrCodeId string, req structure.StatusQrCode) (Error error) {
	res, err := ctrl.access.RDBMS.GetDataQrCode(QrCodeId)
	if err != nil {
		Error = errors.New("Qr-Code ที่จะลบไม่มีอยู่ในระบบ")
		return
	}
	data := rdbmsstructure.QrCode{
		QrCodeUUID: res.QrCodeUUID,
		Active:     *req.Active,
	}
	err = ctrl.access.RDBMS.UpdateQrCodeActive(data)
	if err != nil {
		Error = err
		return
	}
	return
}

func (ctrl *APIControl) CreateQrCode(req structure.GenQrCode) (Error error) {
	req.CodeName = strings.Trim(req.CodeName, "\t \n")
	req.TemplateName = strings.Trim(req.TemplateName, "\t \n")
	if !(len(req.CodeName) <= 20) {
		Error = errors.New("CodeName ต้องไม่เกิน 20 ตัว")
		return
	}
	if req.CodeName == "" {
		Error = errors.New("CodeName ต้องไม่ว่าง")
		return
	}

	ownerId := int(req.OwnerId)
	data, err := ctrl.access.RDBMS.GetAccount(ownerId)
	if err != nil {
		Error = errors.New("ไม่มีผู้ใช้คนนี้อยู่ในระบบ")
		return
	}
	if data.Role != string(constant.Owner) {
		Error = errors.New("ผู้ใช้คนนี้ไม่มีสิทธิ์ในการสร้าง QR-Code")
		return
	}
	structureInfo, err := utility.CheckTemplate(req.TemplateName)
	if err != nil {
		Error = err
		return
	}
	byteInfo, err := json.Marshal(structureInfo)
	if err != nil {
		Error = err
		return
	}
	// สร้าง QR-Code
	count, err := ctrl.access.RDBMS.CountCode(req.OwnerId, req.TemplateName, req.CodeName)
	if err != nil {
		Error = err
		return
	}

	var counts = 0
	check, err := ctrl.access.RDBMS.CheckCode(req.OwnerId, req.TemplateName, req.CodeName)
	if err != nil {
		Error = err
		return
	}
	if check.Code == req.CodeName {
		counts = len(count)
	}

	var arrayString []string
	for i := 0 + 1; i <= req.Amount; i++ {
		uuid, err := uuid2.NewV4()
		if err != nil {
			Error = err
			return
		}
		number := strconv.Itoa(counts + i)

		save := rdbmsstructure.QrCode{
			OwnerId:      req.OwnerId,
			TemplateName: req.TemplateName,
			Info:         datatypes.JSON(byteInfo),
			Ops:          datatypes.JSON(""),
			QrCodeUUID:   uuid,
			Code:         req.CodeName,
			Count:        number,
			First:        false,
			Active:       true,
		}

		err = ctrl.access.RDBMS.CreateQrCode(save)
		if err != nil {
			Error = err
			return
		}
		arrayString = append(arrayString, uuid.String())
	}

	return
}

func (ctrl *APIControl) AddFileZipById(req structure.FileZip) (file string, Error error) {
	var arrayFileName []structure.ArrayFileName
	if len(req.QrCodeId) == 0 {
		Error = errors.New("QrCodeId ต้องไม่ว่าง")
		return
	}
	check, err := ctrl.access.RDBMS.GetAccount(req.OwnerId)
	if err != nil {
		Error = errors.New("ไม่มีผู้ใช้คนนี้อยู่ในระบบ")
		return
	}
	if check.Role != string(constant.Owner) {
		Error = errors.New("ผู้ใช้คนนี้ไม่มีสิทธิ์ในการสร้าง QR-Code")
		return
	}

	err = os.RemoveAll(string(constant.SaveFileLocationZipFile))
	if err != nil {
		Error = err
		return
	}
	err = os.Mkdir(string(constant.SaveFileLocationQrCode), 0755)
	err = os.Mkdir(string(constant.SaveFileLocationZipFile), 0755)
	err = os.Mkdir(string(constant.SaveFileLocationResQR), 0755)
	if err != nil {
		Error = err
		return
	}
	for _, QrCodeId := range req.QrCodeId {
		data, err := ctrl.access.RDBMS.GetQrCodeByQrCodeId(req.OwnerId, QrCodeId)
		if err != nil {
			Error = errors.New("ไม่มี Qr-Code ที่นี้อยู่ในระบบ")
			err = os.RemoveAll(string(constant.SaveFileLocationZipFile))
			if err != nil {
				Error = err
				return
			}
			err = os.RemoveAll(string(constant.SaveFileLocationQrCode))
			if err != nil {
				Error = err
				return
			}
			err = os.RemoveAll(string(constant.SaveFileLocationResQR))
			if err != nil {
				Error = err
				return
			}
			return
		}

		filename := data.Code + "-" + data.Count
		//todo สร้าง QrCode
		qrc, err := qrcode.New(constant.Http + "/" + data.QrCodeUUID.String(),
			//qrcode.WithLogoImageFilePNG(text),
			//qrcode.WithBgColor(col),
		)
		if err != nil {
			Error = err
			return
		}
		// save file
		pathQr := string(constant.SaveFileLocationQrCode) + "/" + filename + ".jpg"
		if err = qrc.Save(pathQr); err != nil {
			Error = err
			return
		}

		// todo เขียนข้อความ
		ResQr := string(constant.SaveFileLocationResQR) + "/" + filename + ".png"
		RequestTextOnImg := utility.Request{
			BgImgPath: pathQr,
			Text:      filename,
			PathSave:  ResQr,
		}
		err = utility.TextOnImg(RequestTextOnImg)
		if err != nil {
			Error = err
			return
		}

		files := structure.ArrayFileName{
			FileName: filename,
		}
		arrayFileName = append(arrayFileName, files)
	}
	output := string(constant.SaveFileLocationZipFile) + "/" + "FileZip.zip"
	if err := ZipFilesById(output, arrayFileName); err != nil {
		Error = err
		return
	}
	err = os.RemoveAll(string(constant.SaveFileLocationQrCode))
	if err != nil {
		Error = err
		return
	}
	err = os.RemoveAll(string(constant.SaveFileLocationResQR))
	if err != nil {
		Error = err
		return
	}
	file = output
	return
}

func (ctrl *APIControl) AddFileZipByOwner(req structure.FileZipByOwner) (file string, Error error) {
	ownerId := int(req.OwnerId)
	dataOwnerId, err := ctrl.access.RDBMS.GetAccount(ownerId)
	if err != nil {
		Error = errors.New("ไม่มีผู้ใช้คนนี้อยู่ในระบบ")
		return
	}
	if dataOwnerId.Role != string(constant.Owner) {
		Error = errors.New("ผู้ใช้คนนี้ไม่มีสิทธิ์ในการสร้าง QR-Code")
		return
	}
	data, err := ctrl.access.RDBMS.GetQrCodeByOwnerId(ownerId)
	if err != nil {
		Error = err
		return
	}
	if len(data) == 0 {
		Error = errors.New("ยังไม่สร้าง Qr-Code ใน template ที่ถูกเลือก")
		return
	}
	err = os.RemoveAll(string(constant.SaveFileLocationZipFile))
	if err != nil {
		Error = err
		return
	}
	err = os.Mkdir(string(constant.SaveFileLocationQrCode), 0755)
	err = os.Mkdir(string(constant.SaveFileLocationZipFile), 0755)
	err = os.Mkdir(string(constant.SaveFileLocationResQR), 0755)
	if err != nil {
		Error = err
		return
	}
	var pathQr string
	var arrayFileName []structure.ArrayFileName
	for _, res := range data {
		qrc, err := qrcode.New(constant.Http + "/" + res.QrCodeUUID.String())
		if err != nil {
			Error = err
			return
		}
		filename := res.Code + "-" + res.Count
		pathQr = string(constant.SaveFileLocationQrCode) + "/" + filename + ".jpg"
		// save file
		if err = qrc.Save(pathQr); err != nil {
			Error = err
			return
		}
		ResQr := string(constant.SaveFileLocationResQR) + "/" + filename + ".png"
		RequestTextOnImg := utility.Request{
			BgImgPath: pathQr,
			Text:      filename,
			PathSave:  ResQr,
		}
		err = utility.TextOnImg(RequestTextOnImg)
		if err != nil {
			Error = err
			return
		}
		files := structure.ArrayFileName{
			FileName: filename,
		}
		arrayFileName = append(arrayFileName, files)
	}
	output := string(constant.SaveFileLocationZipFile) + "/" + "FileZip.zip"
	if err := ZipFilesByTemplateName(output, arrayFileName); err != nil {
		Error = err
		return
	}

	file = output
	err = os.RemoveAll(string(constant.SaveFileLocationQrCode))
	if err != nil {
		Error = err
		return
	}
	err = os.RemoveAll(string(constant.SaveFileLocationResQR))
	if err != nil {
		Error = err
		return
	}
	return
}

func (ctrl *APIControl) AddFileZipByTemplateName(req structure.FileZipByTemplateName) (file string, Error error) {
	req.TemplateName = strings.Trim(req.TemplateName, "\t \n")
	if req.TemplateName == "" {
		Error = errors.New("TemplateName ต้องไม่ว่าง")
		return
	}
	_, err := utility.CheckTemplate(req.TemplateName)
	if err != nil {
		Error = err
		return
	}
	ownerId := int(req.OwnerId)
	dataOwnerId, err := ctrl.access.RDBMS.GetAccount(ownerId)
	if err != nil {
		Error = errors.New("ไม่มีผู้ใช้คนนี้อยู่ในระบบ")
		return
	}
	if dataOwnerId.Role != string(constant.Owner) {
		Error = errors.New("ผู้ใช้คนนี้ไม่มีสิทธิ์ในการสร้าง QR-Code")
		return
	}
	data, err := ctrl.access.RDBMS.GetQrCode(req.OwnerId, req.TemplateName)
	if err != nil {
		Error = err
		return
	}
	if len(data) == 0 {
		Error = errors.New("ยังไม่สร้าง Qr-Code ใน template ที่ถูกเลือก")
		return
	}
	err = os.RemoveAll(string(constant.SaveFileLocationZipFile))
	if err != nil {
		Error = err
		return
	}
	err = os.Mkdir(string(constant.SaveFileLocationQrCode), 0755)
	err = os.Mkdir(string(constant.SaveFileLocationZipFile), 0755)
	err = os.Mkdir(string(constant.SaveFileLocationResQR), 0755)
	if err != nil {
		Error = err
		return
	}
	var pathQr string
	var arrayFileName []structure.ArrayFileName
	for _, res := range data {
		qrc, err := qrcode.New(constant.Http + "/" + res.QrCodeUUID.String())
		if err != nil {
			Error = err
			return
		}
		filename := res.Code + "-" + res.Count
		pathQr = string(constant.SaveFileLocationQrCode) + "/" + filename + ".jpg"
		// save file
		if err = qrc.Save(pathQr); err != nil {
			Error = err
			return
		}
		ResQr := string(constant.SaveFileLocationResQR) + "/" + filename + ".png"
		RequestTextOnImg := utility.Request{
			BgImgPath: pathQr,
			Text:      filename,
			PathSave:  ResQr,
		}
		err = utility.TextOnImg(RequestTextOnImg)
		if err != nil {
			Error = err
			return
		}

		files := structure.ArrayFileName{
			FileName: filename,
		}
		arrayFileName = append(arrayFileName, files)
	}
	output := string(constant.SaveFileLocationZipFile) + "/" + "FileZip.zip"
	if err := ZipFilesByTemplateName(output, arrayFileName); err != nil {
		Error = err
		return
	}

	file = output
	err = os.RemoveAll(string(constant.SaveFileLocationQrCode))
	if err != nil {
		Error = err
		return
	}
	err = os.RemoveAll(string(constant.SaveFileLocationResQR))
	if err != nil {
		Error = err
		return
	}
	return
}

func ZipFilesById(filename string, files []structure.ArrayFileName) error {
	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		pathFile := string(constant.SaveFileLocationResQR) + "/" + file.FileName + ".png"
		if err = AddFileToZip(zipWriter, pathFile); err != nil {
			return err
		}
	}
	return nil
}

func ZipFilesByTemplateName(filename string, files []structure.ArrayFileName) error {
	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		pathFile := string(constant.SaveFileLocationResQR) + "/" + file.FileName + ".png"
		if err = AddFileToZip(zipWriter, pathFile); err != nil {
			return err
		}
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}
