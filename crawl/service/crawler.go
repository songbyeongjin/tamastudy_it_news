package service

import (
	"tamastudy_news_crawler/crawl/model"
)
type Crawler interface {
	CrawlAndSave() error
	Crawl() []*model.News
	Save([]*model.News) error
}