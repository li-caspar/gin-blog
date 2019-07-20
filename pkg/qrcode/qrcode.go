package qrcode

import (
	"caspar/gin-blog/pkg/file"
	"caspar/gin-blog/pkg/setting"
	"caspar/gin-blog/pkg/util"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/jpeg"
)

type Qrcode struct {
	Url string
	Width int
	Height int
	Ext string
	Level qr.ErrorCorrectionLevel
	Mode qr.Encoding
}

const EXT_JPG = ".jpg"

func NewQrcode(url string, width, height int, level qr.ErrorCorrectionLevel, mode qr.Encoding) *Qrcode{
	return &Qrcode{
		Url:url,
		Width:width,
		Height:height,
		Ext:EXT_JPG,
		Level:level,
		Mode:mode,
	}
}

func GetQrCodePath() string {
	return setting.AppSetting.QrCodeSavePath
}

func GetQrCodeFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetQrCodePath()
}

func GetQrCodeFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetQrCodePath() + name
}

func GetQrCodeFileName(value string) string {
	return util.EncodeMD5(value)
}

func (q *Qrcode) GetQrCodeExt() string {
	return q.Ext
}

func (q *Qrcode) CheckEncode(path string) bool {
	src := path + GetQrCodeFileName(q.Url) + q.GetQrCodeExt()
	if file.CheckExist(src) == true {
		return false
	}
	return true
}

func (q *Qrcode) Encode(path string)(string, string, error){
	name := GetQrCodeFileName(q.Url) + q.GetQrCodeExt()
    src := path + name
    if file.CheckExist(src) == true {
    	code, err := qr.Encode(q.Url, q.Level, q.Mode)
    	if err != nil {
    		return "", "", err
		}
    	code, err = barcode.Scale(code, q.Width, q.Height)
    	if err != nil {
    		return "", "", err
		}
    	f, err := file.MustOpen(name, path)
    	if err != nil {
    		return "", "", nil
		}
    	defer f.Close()
    	err = jpeg.Encode(f, code, nil)
    	if err != nil {
    		return "", "", err
		}
	}
    return name, path, nil
}
