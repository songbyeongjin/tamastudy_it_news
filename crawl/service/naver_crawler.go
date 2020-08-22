package service

import (
	"github.com/gocolly/colly"
	"log"
	"regexp"
	"strings"
	"sync"
	"tamastudy_news_crawler/crawl/model"
	"time"
)

const (
	naverNewsPrefixUrl      = "news.naver.com/"
	itScienceSection        = "105"
	naverNewsRootUrl        = httpsUrl + naverNewsPrefixUrl + `main/ranking/popularDay.nhn?rankingType=popular_day&sectionId=` + itScienceSection
	naverCssSelectorUrl     = ".ranking_list li .ranking_headline a"
	naverCssSelectorTitle   = "#articleTitle"
	naverCssSelectorContent = "#articleBodyContents"
	naverCssSelectorPress   = ".press_logo img"
	naverCssSelectorDate    = ".t11"
	deleteString            = "// flash 오류를 우회하기 위한 함수 추가 function _flash_removeCallback() {} "
)

type NaverCrawler struct {
}

func (crawler NaverCrawler) CrawlAndSave() error{
	news := crawler.Crawl()
	if err := crawler.Save(news); err != nil{
		return err
	}

	return nil
}

func (crawler NaverCrawler) Crawl() []*model.News{
	naverNewsUrls := crawler.GetNewsUrls(naverNewsRootUrl)
	naverNews := crawler.GetNews(naverNewsUrls)

	return naverNews
}

func (crawler NaverCrawler) Save(news []*model.News) error{
	fileName := `\naver.txt`
	if err := NewsSave(news, fileName); err != nil{
		return err
	}

	return nil
}

//get Naver News url From nate root url
func (crawler NaverCrawler) GetNewsUrls(rootUrl string) []string {
	urls := make([]string, 0, NewsCount)
	c := colly.NewCollector()
	var wg sync.WaitGroup
	wg.Add(NewsCount)

	// Find and visit all links
	c.OnHTML(naverCssSelectorUrl, func(e *colly.HTMLElement) {
		if len(urls) < NewsCount {
			url := e.Attr(Href)
			urls = append(urls, naverNewsPrefixUrl+url)
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

//get Naver News Object from naver urls
func (crawler NaverCrawler) GetNews(newsUrls []string) []*model.News {
	naverNews := make([]*model.News, NewsCount, NewsCount)
	for i := 0; i < len(naverNews); i++ {
		naverNews[i] = &model.News{}
	}

	cSlice := make([]*colly.Collector, NewsCount, NewsCount)
	dateDuplicateCheck := make([]bool, 10)
	var wg sync.WaitGroup

	//Set callback
	for i := 0; i < NewsCount; i++ {
		inIndex := i
		cSlice[i] = colly.NewCollector()
		cSlice[i].OnHTML(naverCssSelectorTitle, func(e *colly.HTMLElement) {
			naverNews[inIndex].Title = e.Text

			wg.Done()
		})

		cSlice[i].OnHTML(naverCssSelectorContent, func(e *colly.HTMLElement) {
			content := e.Text
			space := regexp.MustCompile(`\s+`)
			str := space.ReplaceAllString(content, " ")

			index := strings.Index(str, deleteString)
			startIndex := index + len(deleteString)
			str2 := str[startIndex:]

			if -1 != strings.Index(str2, "moveCall") {
			}

			naverNews[inIndex].Content = MinimizeContent(str2, 200)

			wg.Done()
		})

		cSlice[i].OnHTML(naverCssSelectorPress, func(e *colly.HTMLElement) {
			naverNews[inIndex].Press = e.Attr("title")

			wg.Done()
		})

		cSlice[i].OnHTML(naverCssSelectorDate, func(e *colly.HTMLElement) {
			if dateDuplicateCheck[inIndex] {
				return
			}
			dateDuplicateCheck[inIndex] = true

			dateString := e.Text[:10]
			replacedDateString := strings.ReplaceAll(dateString, ".", "-")
			naverNews[inIndex].Date, _ = time.Parse(layoutYYYYMMDD, replacedDateString)

			wg.Done()
		})
	}

	for i, url := range newsUrls {
		naverNews[i].Url = url
		naverNews[i].Portal = "naver"
		inUrl := url
		inIndex := i

		wg.Add(1 * 4) //4 = field number in colly call back func(title, content, press, date)
		go func(c *colly.Collector) {
			err := c.Visit(httpsUrl + inUrl)
			if err != nil{
				log.Fatal(err)
			}
		}(cSlice[inIndex]) // i+1 is ranking
	}

	wg.Wait()

	return naverNews
}