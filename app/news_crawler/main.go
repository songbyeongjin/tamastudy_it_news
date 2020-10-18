package main

import (
	"fmt"
	"sync"
	"tamastudy_news_crawler/lib/external/db"
	"tamastudy_news_crawler/lib/external/db/repository_impl"
	"tamastudy_news_crawler/lib/usecase/news/impl"
	"tamastudy_news_crawler/lib/usecase/news/inter"
)

func main()  {
	crawlers, err := getCrawlers()
	if err != nil{
		fmt.Printf("main fail : %s\n", err)
		return
	}

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

func getCrawlers() ([]inter.ICrawlerService, error) {
	// *** DI
	mysqlDb, err := db.NewDbHandler()
	if err != nil{
		return nil, err
	}

	newsRepository := repository_impl.NewNewsRepository(mysqlDb)

	naverCrawlerService := impl.NewNaverCrawlerService(newsRepository)
	yahooCrawlerService := impl.NewYahooCrawlerService(newsRepository)

	crawlersService := []inter.ICrawlerService{
		naverCrawlerService,
		yahooCrawlerService,
	}
	// DI ***

	return crawlersService, nil
}
