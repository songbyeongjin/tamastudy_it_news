package service_interface

import (
	"tamastudy_news_crawler/domain/model"
)
type ICrawlerService interface {
	CrawlAndSave() error
	Crawl() []*model.News
	Save([]*model.News) error
}