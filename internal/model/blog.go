package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Model struct {
	ID        uint64    `gorm:"primarykey;" json:"id"` // 主键
	CreatedAt time.Time `json:"-"`                     // 创建时间
	CreatedBy string    `json:"-"`                     // 创建人
	UpdatedAt time.Time `json:"-"`                     // 修改时间
	UpdatedBy string    `json:"-"`                     // 修改人
	// DeletedAt gorm.DeletedAt `json:"-"`                     // 删除时间 库表中对应字段 deleted_at datetime类型
	IsDeleted soft_delete.DeletedAt `gorm:"softDelete:flag" json:"-" ` // 删除时间 库表中对应字段 is_deleted bool
}

type Article struct {
	*Model
	Title         string `json:"title"`           // 文章标题
	Desc          string `json:"desc"`            // 文章简述
	Content       string `json:"content"`         // 文章内容
	CoverImageUrl string `json:"cover_image_url"` // 封面图片地址
	State         uint8  `json:"state"`           // 状态 0 为禁用、1 为启用
}

func (a Article) TableName() string {
	return "blog_article"
}

type Tag struct {
	*Model
	Name  string `json:"name"`  // 标签名称
	State uint8  `json:"state"` // 状态 0 为禁用、1 为启用
}

func (t Tag) TableName() string {
	return "blog_tag"
}

type ArticleTag struct {
	*Model
	TagID     uint64 `json:"tag_id"`     // 标签 ID
	ArticleID uint64 `json:"article_id"` // 文章 ID
}

func (a ArticleTag) TableName() string {
	return "blog_article_tag"
}
