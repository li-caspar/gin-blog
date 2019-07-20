package tag_service

import (
	"caspar/gin-blog/models"
	"caspar/gin-blog/pkg/gredis"
	"caspar/gin-blog/pkg/logging"
	"caspar/gin-blog/service/cache_service"
	"encoding/json"
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

