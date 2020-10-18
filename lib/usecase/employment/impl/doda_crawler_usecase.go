package impl

import (
	"tamastudy_news_crawler/lib/entity/model"
	"tamastudy_news_crawler/lib/entity/repository_inter"
	"tamastudy_news_crawler/lib/usecase/employment/inter"
)

var dodaStr = "doda"

type dodaCrawlerService struct {
	newsRepository repository_inter.IEmploymentRepository
	site string
}

func NewDodaCrawlerService(newsRepository repository_inter.IEmploymentRepository) inter.ICrawlerService {
	dodaCrawlerService := dodaCrawlerService{newsRepository: newsRepository, site : dodaStr}

	return dodaCrawlerService
}

func (crawler dodaCrawlerService) CrawlAndSave() error{
	if _, err := crawler.Crawl(); err != nil{
		return err
	}

	if err := crawler.Save(nil); err != nil{
		return err
	}

	return nil
}

func (crawler dodaCrawlerService) Crawl() ([]*model.Employment,error) {

	return nil, nil
}

func (crawler dodaCrawlerService) Save(news []*model.Employment) error{

	return nil
}