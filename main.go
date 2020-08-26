package main

import (
	"fmt"
	"sync"
	"tamastudy_news_crawler/db"
	"tamastudy_news_crawler/db/repository_impl"
	"tamastudy_news_crawler/service/service_impl"
	"tamastudy_news_crawler/service/service_interface"
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

func getCrawler() []service_interface.ICrawlerService {
	// *** DI
	mysqlDb := db.NewDbHandler()

	newsRepository := repository_impl.NewNewsRepository(mysqlDb)

	naverCrawlerService := service_impl.NewNaverCrawlerService(newsRepository)
	yahooCrawlerService := service_impl.NewYahooCrawlerService(newsRepository)

	crawlersService := []service_interface.ICrawlerService{
		naverCrawlerService,
		yahooCrawlerService,
	}
	// DI ***

	return crawlersService
}
