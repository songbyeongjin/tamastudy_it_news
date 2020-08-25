package main

import (
	"fmt"
	"sync"
	"tamastudy_news_crawler/db"
	"tamastudy_news_crawler/db/repository_impl"
	"tamastudy_news_crawler/service"
	"tamastudy_news_crawler/service/impl"
)

func main()  {
	crawlers := getCrawler()
	wg := sync.WaitGroup{}
	wg.Add(len(crawlers))

	for _, crawler := range crawlers {
		inCrawler := crawler
		go func(){
			defer wg.Done()
			if err := inCrawler.CrawlAndSave(); err != nil{
				fmt.Println(err)
			}
		}()
	}

	wg.Wait()
}

func getCrawler() []service.ICrawlerService{
	// *** DI
	mysqlDb := db.NewDbHandler()

	newsRepository := repository_impl.NewNewsRepository(mysqlDb)

	naverCrawlerService := impl.NewNaverCrawlerService(newsRepository)
	yahooCrawlerService := impl.NewYahooCrawlerService(newsRepository)

	crawlersService := []service.ICrawlerService{
		naverCrawlerService,
		yahooCrawlerService,
	}
	// DI ***

	return crawlersService
}
