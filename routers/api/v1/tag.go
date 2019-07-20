package v1

import (
	"caspar/gin-blog/pkg/app"
	"caspar/gin-blog/pkg/e"
	"caspar/gin-blog/pkg/export"
	"caspar/gin-blog/pkg/logging"
	"caspar/gin-blog/pkg/setting"
	"caspar/gin-blog/pkg/util"
	"caspar/gin-blog/service/tag_service"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetTags(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.Query("name")
	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}
	tagServcie := tag_service.Tag{
		Name:     name,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}
	tags, err := tagServcie.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}
	count, err := tagServcie.Count()

	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"data":  tags,
		"total": count,
	})

}

//@Summary 新增标签
//@Produce json
//@Param name query string true "Name"
//@Param state query int false "State"
//@Param created_by query int false "CreateBy"
//@Success 200 {string} json "{"code":200, "data":{}, "msg":"ok"}"
//@Router /api/v1/tags [post]
func AddTags(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")
	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "create_by").Message("创建人最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{
		Name:      name,
		CreatedBy: createdBy,
		State:     state,
	}

	exists, err := tagService.ExistByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG, nil)
		return
	}

	err = tagService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func EditTag(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")

	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}
	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	tagServcie := tag_service.Tag{
		ID:         id,
		Name:       name,
		ModifiedBy: modifiedBy,
		State:      state,
	}

	exists, err := tagServcie.ExistByID()

	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	if err = tagServcie.Edit(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func DeleteTag(c *gin.Context) {
	appG := app.Gin{C: c}

	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 0, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	tagService := tag_service.Tag{
		ID: id,
	}
	exists, err := tagService.ExistByID()

	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG, nil)
		return
	}

	if err = tagService.Delete(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func ExportTag(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.PostForm("name")
	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}
	tagService := tag_service.Tag{
		Name:  name,
		State: state,

	}
	filename, err := tagService.Export()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXPORT_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"export_url":export.GetExcelFullUrl(filename),
		"export_save_url":export.GetExcelPath() + filename,
	})
}

func ImportTag(c *gin.Context){
	appG := app.Gin{C:c}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	tagService := tag_service.Tag{}
	err = tagService.Import(file)
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, e.ERROR_IMPORT_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)

}
