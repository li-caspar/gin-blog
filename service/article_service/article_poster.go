package article_service

import (
	"caspar/gin-blog/pkg/file"
	"caspar/gin-blog/pkg/logging"
	"caspar/gin-blog/pkg/qrcode"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
)

type ArticlePoster struct {
	PosterName string
	*Article
	Qr *qrcode.Qrcode
}

func NewArticlePoster(posterName string, article *Article, qr *qrcode.Qrcode) *ArticlePoster {
	return &ArticlePoster{
		PosterName: posterName,
		Article:    article,
		Qr:         qr,
	}
}

func GetPosterFlag() string {
	return "poster"
}

func (a *ArticlePoster) checkMergeImage(path string) bool {
	if file.CheckExist(path + a.PosterName) == true {
		return false
	}
	return true
}

func (a *ArticlePoster) OpenMergedImage(path string) (*os.File, error) {
	f, err := file.MustOpen(a.PosterName, path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

type Rect struct {
	Name string
	X0   int
	Y0   int
	X1   int
	Y1   int
}

type Pt struct {
	X int
	Y int
}

type ArticlePosterBg struct {
	Name string
	*ArticlePoster
	*Rect
	*Pt
}

func NewArticlePosterBg(name string, ap *ArticlePoster, rect *Rect, pt *Pt) *ArticlePosterBg {
	return &ArticlePosterBg{
		Name:          name,
		ArticlePoster: ap,
		Rect:          rect,
		Pt:            pt,
	}
}

func (a *ArticlePosterBg) Generate() (string, string, error) {
	fullPath := qrcode.GetQrCodeFullPath()//获取二维码存储路径
	fileNmae, path, err := a.Qr.Encode(fullPath) //生成二维码图片
	if err != nil {
		return "", "", err
	}
	if !a.checkMergeImage(path) {//检查合并后图像（指的是存放合并后的海报）是否存在
		mergedF, err := a.OpenMergedImage(path) //生成待合并的图像 mergedF
		if err != nil {
			logging.Warn(err)
			return "", "", err
		}
		defer mergedF.Close()
		bgF, err := file.MustOpen(a.Name, path) //打开事先存放的背景图 bgF
		if err != nil {
			logging.Warn(err)
			return "", "", err
		}
		defer bgF.Close()

		qrF, err := file.MustOpen(fileNmae, path) //打开生成的二维码图像 qrF
		if err != nil {
			logging.Warn(err)
			return "", "", err
		}
		defer qrF.Close()

		bgImage, err := jpeg.Decode(bgF)//解码 bgF 和 qrF 返回 image.Image
		if err != nil {
			logging.Warn(err)
			return "", "", err
		}
		qrImage, err := jpeg.Decode(qrF)
		if err != nil {
			logging.Warn(err)
			return "", "", err
		}
		jpg := image.NewRGBA(image.Rect(a.Rect.X0, a.Rect.Y0, a.Rect.X1, a.Rect.Y1))//创建一个新的 RGBA 图像

		draw.Draw(jpg, jpg.Bounds(), bgImage, bgImage.Bounds().Min, draw.Over)//在 RGBA 图像上绘制 背景图（bgF）
		draw.Draw(jpg, jpg.Bounds(), qrImage, qrImage.Bounds().Min.Sub(image.Pt(a.Pt.X, a.Pt.Y)), draw.Over)//在已绘制背景图的 RGBA 图像上，在指定 Point 上绘制二维码图像（qrF）

		jpeg.Encode(mergedF, jpg, nil);//将绘制好的 RGBA 图像以 JPEG 4：2：0 基线格式写入合并后的图像文件（mergedF）

	}
	return fileNmae, path, nil
}
