package parsing

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

func ParseWallapop() (string, error) {
	resp, err := http.Get("https://es.wallapop.com/")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", nil
	}

	var result strings.Builder

	doc.Find(".item").Each(func(index int, item *goquery.Selection) {
		name := strings.TrimSpace(item.Find(".title").Text())
		price := strings.TrimSpace(item.Find(".price").Text())

		if name != "" {
			result.WriteString(fmt.Sprintf("%d. %s - %s\n", index+1, name, price))
		}
	})

	if result.Len() == 0 {
		return "no items found on Wallapop.", nil
	}

	return result.String(), nil
}
