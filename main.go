package main

import (
	"fmt"
	"sync"
	"tamastudy_news_crawler/crawl/service"
)
func main()  {
	crawlers := []service.Crawler{
		service.NaverCrawler{},
		service.YahooCrawler{},
	}

	wg := sync.WaitGroup{}
	wg.Add(len(crawlers))

	for _, crawler := range crawlers{
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
