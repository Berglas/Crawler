package main

import (
	// import standard libraries
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	// import third party libraries
	"./sqlutil"
	"github.com/PuerkitoBio/goquery"
)

var gConn *sql.DB

func main() {
	gConn = sqlutil.Conn("host=127.0.0.1 dbname=Resumes user=postgres password=1234 sslmode=disable")
	if gConn == nil {
		return
	}
	defer gConn.Close()

	postScrape()
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func postScrape() {
	doc, err := goquery.NewDocument("http://rate.bot.com.tw/xrt?Lang=zh-TW")
	if err != nil {
		log.Fatal(err)
	}

	// use CSS selector found with the browser inspector
	// for each, use index and item
	doc.Find(".table.table-striped.table-bordered.table-condensed.table-hover tbody tr").Each(func(index int, item *goquery.Selection) {
		currency := item.Find(".hidden-phone.print_show").Text()
		buying_rate := item.Find(".rate-content-cash.text-right.print_hide").Eq(0).Text()
		if buying_rate == "-" {
			buying_rate = "NULL"
		}
		selling_rate := item.Find(".rate-content-cash.text-right.print_hide").Eq(1).Text()
		if selling_rate == "-" {
			selling_rate = "NULL"
		}
		quoted_date := doc.Find(".text-info span.time").Text()

		//寫入資料庫
		tx, _ := gConn.Begin()
		_, err = tx.Exec(fmt.Sprintf("INSERT INTO currency_rate VALUES ('%v',%v,%v,'%v')", removeBlank(currency), removeBlank(buying_rate), removeBlank(selling_rate), quoted_date))
		if err != nil {
			log.Print(err)
			tx.Rollback()
			fmt.Println("新增失敗")
		} else {
			tx.Commit()
			fmt.Println("新增成功")
		}
		//fmt.Printf("%s買:%s,賣:%s，時間：%s\n", removeBlank(currency), removeBlank(buying_rate), removeBlank(selling_rate), quoted_date)
	})
}

func removeBlank(r string) (result string) {
	result = strings.Replace(strings.Replace(r, " ", "", -1), "\n", "", -1)
	return
}
