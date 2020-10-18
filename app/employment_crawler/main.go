package main

import (
	"fmt"
	"sync"
	"tamastudy_news_crawler/lib/external/db"
	"tamastudy_news_crawler/lib/external/db/repository_impl"
	"tamastudy_news_crawler/lib/usecase/employment/impl"
	"tamastudy_news_crawler/lib/usecase/employment/inter"
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

	employmentRepository := repository_impl.NewEmploymentRepository(mysqlDb)

	dodaCrawlerService := impl.NewDodaCrawlerService(employmentRepository)

	crawlersService := []inter.ICrawlerService{
		dodaCrawlerService,
	}
	// DI ***

	return crawlersService, nil
}
