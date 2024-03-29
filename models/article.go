package models

import "github.com/jinzhu/gorm"

type Article struct {
	Model
	TagID         int    `json:"tag_id" gom:"index"`
	Tag           Tag    `json:"tag"`
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	CreatedBy     string `json:"created_by"`
	ModifiedBy    string `json:"modified_by"`
	State         int    `json:"state"`
}

func ExitArticleByID(id int) (bool, error) {
	var article Article
	err := db.Select("id").Where("id=? AND deleted_on=?", id, 0).First(&article).Error
	if err != nil {
		return false, err
	}
	if article.ID > 0 {
		return true, nil
	}
	return false, nil
}

func GetArticleTotal(maps interface{}) (int, error) {
	var count int
	err := db.Model(&Article{}).Where(maps).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetArticles(pageNum int, pageSize int, maps interface{}) ([]*Article, error) {
	var articles []*Article
	err := db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return articles, nil
}

func GetArticle(id int) (*Article, error) {
	var article Article
	err := db.Where("id=?", id).First(&article).Error
	db.Model(&article).Related(&article.Tag)
	return &article, err
}

func AddArticle(data map[string]interface{}) error {
	return db.Create(&Article{
		TagID:         data["tag_id"].(int),
		Title:         data["title"].(string),
		Desc:          data["desc"].(string),
		Content:       data["content"].(string),
		CoverImageUrl: data["cover_image_url"].(string),
		CreatedBy:     data["created_by"].(string),
		State:         data["state"].(int),
	}).Error
}

func EditArticle(id int, data map[string]interface{}) error {
	if err := db.Model(&Article{}).Where("id=?", id).Update(data).Error; err != nil {
		return err
	}
	return nil
}

func DeleteArticle(id int) error {
	return db.Where("id=?", id).Delete(Article{}).Error
}

/*func (article *Article) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())
	return nil
}

func (article *Article) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedOn", time.Now().Unix())
	return nil
}*/
