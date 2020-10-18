package inter

import (
	"tamastudy_news_crawler/lib/entity/model"
)

type ICrawlerService interface {
	CrawlAndSave() error
	Crawl() ([]*model.Employment, error)
	Save([]*model.Employment) error
}