package repository_inter

import(
	"tamastudy_news_crawler/lib/entity/model"
)
type INewsRepository interface {
	DeleteAllByPortal(portal string) error
	Create(*model.News) error
	DeleteAllByPortalAndAllCreate(portal string, news []*model.News)  error
}