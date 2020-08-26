package repository_interface

import(
	"tamastudy_news_crawler/domain/model"
)
type INewsRepository interface {
	DeleteAllByPortal(portal string) error
	Create(*model.News) error
	DeleteAllByPortalAndAllCreate(portal string, news []*model.News)  error
}