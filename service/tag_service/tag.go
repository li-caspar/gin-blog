package tag_service

import (
	"caspar/gin-blog/models"
	"caspar/gin-blog/pkg/export"
	mfile "caspar/gin-blog/pkg/file"
	"caspar/gin-blog/pkg/gredis"
	"caspar/gin-blog/pkg/logging"
	"caspar/gin-blog/service/cache_service"
	"encoding/json"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/tealeg/xlsx"
	"io"
	"strconv"
	"time"
)

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

/*func (t *Tag) Get() (tag *models.Tag, err error){
    cache := cache_service.Tag{ID:t.ID}
    key := cache.GetTagKey()
    if gredis.Exists(key){
    	data, err := gredis.Get(key)
    	if err != nil{
    		logging.Info(err)
		}else{
			json.Unmarshal(data, &tag)
			return
		}
	}
    tag, err := models.get
}*/

func (t *Tag) GetAll() ([]*models.Tag, error) {
	var tags []*models.Tag
	cache := cache_service.Tag{
		Name:     t.Name,
		State:    t.State,
		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}
	key := cache.GetTagsKey()
	if gredis.Exists(key){
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		}else{
			json.Unmarshal(data, &tags)
			return tags, nil
		}
	}
	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (t *Tag) Count() (int, error) {
	return models.GetTagTotal(t.getMaps())
}

func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["stauts"] = t.State
	}
	return maps
}

func (t *Tag) ExistByName() (bool, error){
    return models.ExistTagByName(t.Name)
}

func (t *Tag) ExistByID() (bool, error){
    return models.ExistTagByID(t.ID)
}

func (t *Tag) Add() error{
    if err:= models.AddTag(t.Name, t.State, t.CreatedBy); err != nil {
    	logging.Info(err)
    	return err
	}
    return nil
}

func (t *Tag) Edit() error{
	data := make(map[string]interface{})
	data["modified_by"] = t.ModifiedBy
	data["name"] = t.Name
	data["state"] = t.State
    if err := models.EditTag(t.ID, data); err != nil {
		logging.Info(err)
    	return err
	}
    return nil
}

func (t *Tag) Delete() error {
	if err := models.DeleteTag(t.ID); err != nil {
		logging.Info(err)
		return err
	}
	return nil
}

func (t *Tag) Export() (string, error){
	tags, err := t.GetAll()
	if err != nil {
		return "", err
	}
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("标签信息")
	if err != nil {
		return "", err
	}
	titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	row := sheet.AddRow()
	var  cell *xlsx.Cell
	for _, title := range titles {
		cell = row.AddCell()
		cell.Value = title
	}

	for _,v := range tags {
		values := []string{
			strconv.Itoa(v.ID),
			v.Name,
			v.CreatedBy,
			strconv.Itoa(v.CreatedOn),
			v.ModifiedBy,
			strconv.Itoa(v.ModifiedOn),
		}
		row := sheet.AddRow()
		for _,value := range values{
			cell = row.AddCell()
			cell.Value = value
		}
	}
	ts := strconv.Itoa(int(time.Now().Unix()))
	filename := "tags-"+ts+".xlsx"

	fullPath := export.GetExcelFullPath() + filename
	err = mfile.IsNotExisMkDir(export.GetExcelFullPath())
	err = file.Save(fullPath)
	if err != nil {
		return "", nil
	}

	return filename, nil
}

func (t *Tag) Import(r io.Reader) error {
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}
	rows,_ := xlsx.GetRows("标签信息")
	for index, row := range rows {
		if index >0 {
			var data []string
			for _, cell := range row{
               data = append(data, cell)
			}
			models.AddTag(data[0], 1, data[1])
		}

	}
	return nil

}

