package service

import (
	"github.com/gocolly/colly"
	"log"
	"regexp"
	"sync"
	"tamastudy_news_crawler/crawl/model"
	"time"
)

type YahooCrawler struct {
}

const (
	yahooNewsRootUrl         = httpsUrl  + `news.yahoo.co.jp/topics/it`
	yahooCssSelectorFirstUrl = ".newsFeed_item_link"
	yahooCssSelectorSecondUrl = ".pickupMain_detailLink a"
	yahooCssSelectorTitle    = ".sc-eXEjpC"
	yahooCssSelectorContent  = ".article_body"
	yahooCssSelectorPress    = ".pickupMain_media"
	yahooCssSelectorDate     = ".sc-bwCtUz time"
)

func (crawler YahooCrawler) CrawlAndSave() error{
	news := crawler.Crawl()
	if err := crawler.Save(news); err != nil{
		return err
	}

	return nil
}

func (crawler YahooCrawler) Crawl() []*model.News {
	firstUrls  := crawler.GetFirstNewsUrls(yahooNewsRootUrl)
	secondUrls, press := crawler.GetSecondNewsUrls(firstUrls)
	news := crawler.GetNews(secondUrls, press)

	return news
}

func (crawler YahooCrawler) Save(news []*model.News) error{
	fileName := `\yahoo.txt`
	if err := NewsSave(news, fileName); err != nil{
		return err
	}

	return nil
}

//get Naver News url From nate root url
func (crawler YahooCrawler) GetFirstNewsUrls(rootUrl string) []string {
	urls := make([]string, 0, NewsCount)
	c := colly.NewCollector()
	var wg sync.WaitGroup
	wg.Add(NewsCount)

	// Find and visit all links
	c.OnHTML(yahooCssSelectorFirstUrl, func(e *colly.HTMLElement) {
		if len(urls) < NewsCount {
			url := e.Attr(Href)
			urls = append(urls, url)
			wg.Done()
		}
	})

	err := c.Visit(rootUrl)
	if err != nil{
		log.Fatal(err)
	}

	wg.Wait()

	return urls
}

func (crawler YahooCrawler) GetSecondNewsUrls(firstUrls []string) ([]string, []string){
	urls := make([]string, 0, NewsCount)
	press := make([]string, 0, NewsCount)
	cSlice := make([]*colly.Collector, NewsCount, NewsCount)
	var wg sync.WaitGroup
	wg.Add(NewsCount * 2)//2 ==field count(url, press)

	for i := 0; i < NewsCount; i++ {
		inIndex := i

		cSlice[inIndex] = colly.NewCollector()

		cSlice[inIndex].OnHTML(yahooCssSelectorSecondUrl, func(e *colly.HTMLElement) {
			defer wg.Done()
			url := e.Attr("href")
			urls = append(urls, url)
		})

		cSlice[inIndex].OnHTML(yahooCssSelectorPress, func(e *colly.HTMLElement) {
			defer wg.Done()
			p := e.Text
			press = append(press, p)
		})
	}

	for i, firstUrl := range firstUrls {
		inUrl := firstUrl
		inIndex := i

		go func(c *colly.Collector) {
			if err := c.Visit(inUrl); err != nil{
				log.Fatal(err)
			}

		}(cSlice[inIndex]) // i+1 is ranking
	}

	wg.Wait()

	return urls, press
}


//get Yahoo News Object from Yahoo urls
func (crawler YahooCrawler) GetNews(newsUrls []string, newsPress []string) []*model.News {
	yahooNews := make([]*model.News, NewsCount, NewsCount)
	for i := 0; i < len(yahooNews); i++ {
		yahooNews[i] = &model.News{}
	}

	cSlice := make([]*colly.Collector, NewsCount, NewsCount)
	var wg sync.WaitGroup

	//Set callback
	for i := 0; i < NewsCount; i++ {
		inIndex := i
		cSlice[i] = colly.NewCollector()

		cSlice[i].OnHTML(yahooCssSelectorTitle, func(e *colly.HTMLElement) {
			yahooNews[inIndex].Title = e.Text

			wg.Done()
		})

		cSlice[i].OnHTML(yahooCssSelectorContent, func(e *colly.HTMLElement) {
			content := e.Text
			space := regexp.MustCompile(`\s+`)
			trimmedContent := space.ReplaceAllString(content, " ")

			yahooNews[inIndex].Content = MinimizeContent(trimmedContent, 200)

			wg.Done()
		})

		cSlice[i].OnHTML(yahooCssSelectorDate, func(e *colly.HTMLElement) {
			yahooNews[inIndex].Date, _ = time.Parse(layoutYYYYMMDD, time.Now().String())

			wg.Done()
		})
	}

	for i := 0; i< NewsCount; i++{
		yahooNews[i].Url = newsUrls[i]
		yahooNews[i].Portal = "yahoo"
		yahooNews[i].Press = newsPress[i]

		inUrl := newsUrls[i]
		inIndex := i

		wg.Add(1 * 3)//4 = field number in colly call back func(title, content, date)
		go func(c *colly.Collector) {
			if err := c.Visit(inUrl); err != nil{
				log.Fatal(err)
			}
		}(cSlice[inIndex]) // i+1 is ranking
	}

	wg.Wait()

	return yahooNews
}
