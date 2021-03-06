package impl

import (
	"github.com/gocolly/colly"
	"log"
	"regexp"
	"sync"
	"tamastudy_news_crawler/lib/common"
	"tamastudy_news_crawler/lib/entity/model"
	"tamastudy_news_crawler/lib/entity/repository_inter"
	"tamastudy_news_crawler/lib/usecase/news/inter"
	"time"
)

const (
	yahooNewsRootUrl         = common.HttpsUrl + `news.yahoo.co.jp/topics/it`
	yahooCssSelectorFirstUrl = ".newsFeed_item_link"
	yahooCssSelectorTitle    = ".sc-epnACN"
	yahooCssSelectorContent  = ".article_body"
	yahooCssSelectorPress    = ".pickupMain_media"
	yahooCssSelectorSecondUrlPress     = ".pickupMain_articleInfo"
)

var essayCounter = 0

type yahooCrawlerService struct {
	newsRepository repository_inter.INewsRepository
	portal         string
}

func NewYahooCrawlerService(newsRepository repository_inter.INewsRepository) inter.ICrawlerService {
	yahooCrawlerService := yahooCrawlerService{
		newsRepository: newsRepository,
		portal : "yahoo"}

	return yahooCrawlerService
}


func (crawler yahooCrawlerService) CrawlAndSave() error{
	news := crawler.Crawl()
	if err := crawler.Save(news); err != nil{
		return err
	}

	return nil
}

func (crawler yahooCrawlerService) Crawl() []*model.News {
	firstUrls  := crawler.GetFirstNewsUrls(yahooNewsRootUrl)
	secondUrls, press := crawler.GetSecondNewsUrlsAndPress(firstUrls)
	news := crawler.GetNews(secondUrls, press)

	return news
}

func (crawler yahooCrawlerService) Save(news []*model.News) error{
	if err := crawler.newsRepository.DeleteAllByPortalAndAllCreate(crawler.portal, news); err != nil{
		return err
	}

	return nil
}

//get Naver News url From nate root url
func (crawler yahooCrawlerService) GetFirstNewsUrls(rootUrl string) []string {
	urls := make([]string, 0, common.NewsCount)
	c := colly.NewCollector()
	var wg sync.WaitGroup
	wg.Add(common.NewsCount)

	// Find and visit all links
	c.OnHTML(yahooCssSelectorFirstUrl, func(e *colly.HTMLElement) {
		if len(urls) < common.NewsCount {
			url := e.Attr(common.Href)
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

func (crawler yahooCrawlerService) GetSecondNewsUrlsAndPress(firstUrls []string) ([]string, []string){
	urls := make([]string, common.NewsCount, common.NewsCount)
	press := make([]string, common.NewsCount, common.NewsCount)
	cSlice := make([]*colly.Collector, common.NewsCount, common.NewsCount)
	var wg sync.WaitGroup
	wg.Add(common.NewsCount * 1)

	for i := 0; i < common.NewsCount; i++ {
		inIndex := i

		cSlice[inIndex] = colly.NewCollector()
		cSlice[inIndex].OnHTML(yahooCssSelectorSecondUrlPress, func(e *colly.HTMLElement) {
			defer wg.Done()

			if p := e.ChildText(yahooCssSelectorPress); p != ""{
				url := e.ChildAttr(" a", "href")
				urls[inIndex] = url
				press[inIndex] = p


			}else{
				//sometime news don't have press then that news item is essay
				press[inIndex] = "*"
				urls[inIndex] = "*"
				essayCounter++
			}
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

	//delete item that not found press(=essay item)
	for i, url := range urls{
		if url == "*" {
			urls = append(urls[0:i], urls[i+1:]...)
			press = append(press[0:i], press[i+1:]...)
		}
	}

	return urls, press
}


//get Yahoo News Object from Yahoo urls
func (crawler yahooCrawlerService) GetNews(newsUrls []string, newsPress []string) []*model.News {
	newsCount := len(newsUrls)
	yahooNews := make([]*model.News, newsCount, newsCount)
	for i := 0; i < newsCount; i++ {
		yahooNews[i] = &model.News{}
	}

	cSlice := make([]*colly.Collector, newsCount, newsCount)
	var wg sync.WaitGroup
	wg.Add( newsCount * 2)//3 = field number in colly call back func(title, content)

	//Set callback
	for i := 0; i < newsCount; i++ {
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

			yahooNews[inIndex].Content = common.MinimizeContent(trimmedContent, 200)
			wg.Done()
		})
	}

	for i := 0; i < newsCount; i++{
		yahooNews[i].Url = newsUrls[i]
		yahooNews[i].Portal = "yahoo"
		yahooNews[i].Press = newsPress[i]
		yahooNews[i].CountryCode = model.JapanCode
		yahooNews[i].Date = time.Now()

		inUrl := newsUrls[i]
		inIndex := i

		go func(c *colly.Collector) {
			if err := c.Visit(inUrl); err != nil{
				log.Fatal(err)
			}
		}(cSlice[inIndex]) // i+1 is ranking
	}


	wg.Wait()

	return yahooNews
}
