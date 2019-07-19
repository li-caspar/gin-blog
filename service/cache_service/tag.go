package cache_service

import (
	"caspar/gin-blog/pkg/e"
	"strconv"
	"strings"
)

type Tag struct {
	ID       string
	Name     string
	State    int
	PageNum  int
	PageSize int
}

func (t *Tag) GetTagKey() string {
	return e.CACHE_TAG + "_" + t.ID
}

func (t *Tag) CetTagsKey() string {
	keys := []string{
		e.CACHE_TAG,
		"LIST",
	}
	if t.Name != "" {
		keys = append(keys, t.Name)
	}
	if t.State >= 0 {
		keys = append(keys, strconv.Itoa(t.State))
	}
	if t.PageNum > 0 {
		keys = append(keys, strconv.Itoa(t.PageNum))
	}
	if t.PageSize > 0 {
		keys = append(keys, strconv.Itoa(t.PageSize))
	}

	return strings.Join(keys, "_")
}
