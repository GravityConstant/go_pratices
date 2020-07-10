package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

var (
	url      = "http://gwy.rst.fujian.gov.cn/signupcount"
	data     [3004][7]string
	dataStr  string
	isInTime bool
)

func ExampleScrape() {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var n int

	dataStr, _ = doc.Find("table").Html()
	doc.Find("thead").Each(func(i int, thead *goquery.Selection) {
		var d [7]string
		thead.Find("th").Each(func(j int, th *goquery.Selection) {
			d[j] = th.Text()
		})
		if len(d[0]) > 0 {
			data[n] = d
			n++
		}
	})
	doc.Find("tbody").Each(func(i int, tbody *goquery.Selection) {
		tbody.Find("tr").Each(func(j int, tr *goquery.Selection) {
			var d [7]string
			tr.Find("td").Each(func(k int, td *goquery.Selection) {
				d[k] = td.Text()
			})
			if len(d[0]) > 0 {
				data[n] = d
				n++
			}
		})
	})

}

func main() {
	ExampleScrape()
	// for i, d := range data {
	// 	if i > 2 {
	// 		break
	// 	}
	// 	fmt.Println(d)
	// }

	logger := log.New(os.Stdout, "INFO: ", log.Lshortfile)

	tick := time.Tick(300 * time.Second)
	go func() {
		for range tick {
			isInTime = true
		}
	}()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Static("/assets", "./statics")
	r.LoadHTMLGlob("./views/*/*")
	r.GET("/", func(c *gin.Context) {
		if isInTime {
			logger.Println("fresh data")
			ExampleScrape()
			isInTime = false
		}
		c.HTML(http.StatusOK, "index/index.html", gin.H{
			"title":   "gov exam",
			"data":    data,
			"dataStr": dataStr,
		})
	})
	logger.Println("Listening on port: 10003")
	r.Run(":10003")
}
