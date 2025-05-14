package parsing

import (
	"fmt"
	"github.com/qeery8/http"
)

type Item struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Price struct {
		Amount float64 `json:"amount"`
	} `json:"price"`
	WebSlug string `json:"web_slug"`
	Images  []struct {
		Original string `json:"original"`
	} `json:"images"`
}

type SearchResponce struct {
	SearchObjects []struct {
		Content Item `json:"content"`
	} `json:"search_objects"`
}

func ParseWallapop() (string, error) {
	url := "https://api.wallapop.com/api/v3/general/search?keywords=iphone&latitude=40.4168&longitude=-3.7038&order_by=most_relevance"

	var result SearchResponce
	err := http.GetAPIResponse(url, &result, nil)
	if err != nil {
		return "", fmt.Errorf("falied to fetch wallapop", err)
	}

	if len(result.SearchObjects) == 0 {
		return "Not found", nil
	}

	var out string
	for i, obj := range result.SearchObjects {
		item := obj.Content
		link := fmt.Sprintf("https://es.wallapop.com/item/%s", item.WebSlug)
		out += fmt.Sprintf("title: %s\n, price: %.2f\n, link: %s\n", item.Title, item.Price.Amount, link)

		if i >= 4 {
			break
		}
	}
	return out, nil
}
