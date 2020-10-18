package impl

import (
	"github.com/gocolly/colly"
	"log"
	"regexp"
	"strings"
	"sync"
	"tamastudy_news_crawler/lib/common"
	"tamastudy_news_crawler/lib/entity/model"
	"tamastudy_news_crawler/lib/entity/repository_inter"
	"tamastudy_news_crawler/lib/usecase/news/inter"
	"time"
)

const (
	naverNewsPrefixUrl      = "news.naver.com/"
	itScienceSection        = "105"
	naverNewsRootUrl        =
		common.HttpsUrl +
			naverNewsPrefixUrl +
		`main/ranking/popularDay.nhn?rankingType=popular_day&sectionId=` +
			itScienceSection
	naverCssSelectorUrl     = ".ranking_list li .ranking_headline a"
	naverCssSelectorTitle   = "#articleTitle"
	naverCssSelectorContent = "#articleBodyContents"
	naverCssSelectorPress   = ".press_logo img"
	naverCssSelectorDate    = ".t11"
	deleteString            = "// flash 오류를 우회하기 위한 함수 추가 function _flash_removeCallback() {} "
)

type naverCrawlerService struct {
	newsRepository repository_inter.INewsRepository
	portal         string
}

func NewNaverCrawlerService(newsRepository repository_inter.INewsRepository) inter.ICrawlerService {
	naverCrawlerService := naverCrawlerService{newsRepository: newsRepository, portal : "naver"}

	return naverCrawlerService
}

func (crawler naverCrawlerService) CrawlAndSave() error{
	news := crawler.Crawl()
	if err := crawler.Save(news); err != nil{
		return err
	}

	return nil
}

func (crawler naverCrawlerService) Crawl() []*model.News {
	naverNewsUrls := crawler.GetNewsUrls(naverNewsRootUrl)
	naverNews := crawler.GetNews(naverNewsUrls)

	return naverNews
}

func (crawler naverCrawlerService) Save(news []*model.News) error{
	/*
	fileName := `\naver.txt`
	//if err := common.NewsSave(news, fileName); err != nil{
		return err
	}
	 */

	if err := crawler.newsRepository.DeleteAllByPortalAndAllCreate(crawler.portal, news); err != nil{
		return err
	}

	return nil
}

//get Naver News url From nate root url
func (crawler naverCrawlerService) GetNewsUrls(rootUrl string) []string {
	urls := make([]string, 0, common.NewsCount)
	c := colly.NewCollector()
	var wg sync.WaitGroup
	wg.Add(common.NewsCount)

	// Find and visit all links
	c.OnHTML(naverCssSelectorUrl, func(e *colly.HTMLElement) {
		if len(urls) < common.NewsCount {
			url := e.Attr(common.Href)
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
func (crawler naverCrawlerService) GetNews(newsUrls []string) []*model.News {
	naverNews := make([]*model.News, common.NewsCount, common.NewsCount)
	for i := 0; i < len(naverNews); i++ {
		naverNews[i] = &model.News{}
	}

	cSlice := make([]*colly.Collector, common.NewsCount, common.NewsCount)
	dateDuplicateCheck := make([]bool, 10)
	var wg sync.WaitGroup
	wg.Add(common.NewsCount * 4) //4 = field number in colly call back func(title, content, press, date)


	//Set callback
	for i := 0; i < common.NewsCount; i++ {
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

			naverNews[inIndex].Content = common.MinimizeContent(str2, 200)

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
			naverNews[inIndex].Date, _ = time.Parse(common.LayoutYYYYMMDD, replacedDateString)

			wg.Done()
		})
	}

	for i, url := range newsUrls {
		naverNews[i].Url = url
		naverNews[i].Portal = "naver"
		naverNews[i].CountryCode = model.KoreaCode
		inUrl := url
		inIndex := i

		go func(c *colly.Collector) {
			err := c.Visit(common.HttpsUrl + inUrl)
			if err != nil{
				log.Fatal(err)
			}
		}(cSlice[inIndex]) // i+1 is ranking
	}

	wg.Wait()

	return naverNews
}