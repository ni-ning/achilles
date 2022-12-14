package blog

import (
	"achilles/global"
	"achilles/internal/model"
	"achilles/pkg/app"
	"fmt"
)

type TagListCountRequest struct {
	Name  string `form:"name" binding:"max=100"`
	State *uint8 `form:"state"`
}

// 与 TagListRequest 中 State取值不同
type TagCreateRequest struct {
	Name  string `form:"name" binding:"required,min=1,max=100"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

// 全量更新
type TagUpdateRequest struct {
	ID    uint64 `form:"id" binding:"required,gte=1"`
	Name  string `form:"name" binding:"min=1,max=100"`
	State *uint8 `form:"state,default=1" binding:"oneof=0 1"`
}

type TagDeleteRequest struct {
	ID uint64 `form:"id" binding:"required,gte=1"`
}

func GetTagList(req TagListCountRequest, page, pageSize int) ([]*model.Tag, error) {
	var tags []*model.Tag
	db := global.DBEngine.Model(&model.Tag{})
	if req.Name != "" {
		db = db.Where("name = ?", req.Name)
	}
	if req.State != nil {
		db = db.Where("state = ?", req.State)
	}

	fmt.Printf("\nGetTagList Req:%#v\n", req)

	pageOffset := app.GetPageOffset(page, pageSize)
	if pageOffset >= 0 && pageSize > 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}

	if err := db.Debug().Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func GetTagCount(req TagListCountRequest) (int, error) {
	var count int64
	db := global.DBEngine.Model(&model.Tag{})
	if req.Name != "" {
		db = db.Where("name = ?", req.Name)
	}
	if req.State != nil {
		db = db.Where("state = ?", req.State)
	}

	fmt.Printf("\nGetTagCount Req:%#v\n", req)

	if err := db.Debug().Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func CreateTag(req TagCreateRequest) error {
	fmt.Printf("\nCreateTag Req:%#v\n", req)

	return global.DBEngine.Debug().Create(&model.Tag{Name: req.Name, State: req.State}).Error
}

func UpdateTag(req TagUpdateRequest) error {
	fmt.Printf("\nUpdateTag Req:%#v\n", req)

	db := global.DBEngine.Model(&model.Tag{})
	db = db.Debug().Where("id = ?", req.ID)
	return db.Updates(map[string]interface{}{"name": req.Name, "state": req.State}).Error
}

func DeleteTag(req TagDeleteRequest) error {
	fmt.Printf("\nDeleteTag Req:%#v\n", req)

	return global.DBEngine.Delete(&model.Tag{}, req.ID).Error
}
