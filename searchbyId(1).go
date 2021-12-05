// cumt lib search book project main.go
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/chromedp/chromedp"
)

func main() {

	var bookID string
	fmt.Printf("请输入bookID：")
	//fmt.Scanf(&bookID)
	bookID = "373428"

	input := bufio.NewScanner(os.Stdin)

	// 逐行扫描
	if input.Scan() {
		bookID = input.Text()
	}

	searchBookById(bookID)
}

func searchBookById(bookID string) {

	var htmlContent string

	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	c, _ := chromedp.NewExecAllocator(context.Background(), options...)

	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	defer cancel()
	//chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)
	// 给每个页面的爬取设置超时时间
	ctx, cancel := context.WithTimeout(chromeCtx, 20*time.Second)
	defer cancel()

	bUrl := "https://findcumt.libsp.com/#/searchList/bookDetails/" + bookID

	err := chromedp.Run(ctx,
		chromedp.Navigate(bUrl),
		chromedp.WaitVisible(`#container > div.table___1Mn5Z > div.collectionInfo___1e6a2 > div > div > div > div > div > div > div > div > div.ant-table-body > table > tbody`),
		chromedp.OuterHTML(`#container > div.table___1Mn5Z > div.collectionInfo___1e6a2 > div > div > div > div > div > div > div > div > div.ant-table-body > table`, &htmlContent, chromedp.BySearch),
	)

	ctx.Done()

	if err != nil {
		log.Fatal(err)
	}

	var table [][]string = bookLocaltion(htmlContent)

	for trkey, _ := range table {
		fmt.Printf("|%12s\t|%12s\t|%12s\t", table[trkey][3], table[trkey][1], table[trkey][5])
		fmt.Println("|\n------------------------------------------------------------------")
	}
}

func bookLocaltion(htmlContent string) (table [][]string) {
	var tableCache [][]string

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Fatalln(err)
	}

	trNum := dom.Find("tr").Length()
	tdNum := dom.Find("tr").Find("td").Length()

	tableCache = make([][]string, trNum, trNum)
	for i := 0; i < trNum; i++ {
		tableCache[i] = make([]string, tdNum/trNum, tdNum/trNum)
		for j := 0; j < (tdNum / trNum); j++ {
			tableCache[i][j] = strings.Trim(dom.Find("tr").Find("td").Eq(i*(tdNum/trNum)+j).Text(), " ")
		}
	}

	return tableCache
}
