package article_service

import (
	"caspar/gin-blog/models"
	"caspar/gin-blog/pkg/gredis"
	"caspar/gin-blog/pkg/logging"
	"caspar/gin-blog/service/cache_service"
	"encoding/json"
)

type Article struct {
	ID            int
	TagId         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	State         int
	CreateBy      string
	ModifiedBy    string

	PageNum  int
	PageSize int
}

func (a *Article) Add() error {
	article := map[string]interface{}{
		"tag_id":          a.TagId,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"cover_image_url": a.CoverImageUrl,
		"created_by":      a.CreateBy,
		"state":           a.State,
	}
	err := models.AddArticle(article)
	if err != nil {
		logging.Info(err)
	}
	return err
}

func (a *Article) Get() (*models.Article, error) {
	var cacheArticle *models.Article

	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticle)
			return cacheArticle, nil
		}
	}

	article, err := models.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}

	gredis.Set(key, article, 3600)
	return article, nil
}

func (a *Article) GetAll() ([]*models.Article, error) {
	var (
		articles, cacheArticles []*models.Article
	)
	cache := cache_service.Article{ID: a.ID, TagID: a.TagId, State: a.State, PageNum: a.PageNum, PageSize: a.PageSize}
	key := cache.GetArticlesKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticles)
			return cacheArticles, nil
		}
	}
	articles, err := models.GetArticles(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}
	gredis.Set(key, articles, 3600)
	return articles, nil
}

func (a *Article) ExistByID() (bool, error) {
	return models.ExitArticleByID(a.ID)
}

func (a *Article) Count() (int, error) {
	count, err := models.GetArticleTotal(a.getMaps())
	if err != nil {
		logging.Info(err)
		return 0, err
	}
	return count, nil
}

func (a *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if a.State != -1 {
		maps["state"] = a.State
	}
	if a.TagId != -1 {
		maps["tag_id"] = a.TagId
	}
	return maps
}
