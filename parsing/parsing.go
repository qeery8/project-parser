package parsing

import (
	"fmt"
	"github.com/qeery8/http"
)

type Item struct {
	ID      string  `json:"id"`
	Title   string  `json:"title"`
	Price   float64 `json:"price"`
	WebSlug string  `json:"web_slug"`
	Images  []struct {
		Original string `json:"original"`
		Small    string `json:"small"`
		Medium   string `json:"medium"`
		Large    string `json:"large"`
		XLarge   string `json:"xlarge"`
	} `json:"images"`
}

type SearchResponce struct {
	SearchObjects []struct {
		Content Item `json:"content"`
	} `json:"search_objects"`
}

func ParseWallapop(offset int) ([]string, error) {
	url := fmt.Sprintf("https://api.wallapop.com/api/v3/cars/search?latitude=40.4168&longitude=-3.7038&start=%d&num=%d", offset, offset+50)

	headers := map[string]string{}

	var result SearchResponce
	err := http.GetAPIResponse(url, &result, headers)
	if err != nil {
		return nil, fmt.Errorf("falied to fetch wallapop %w", err)
	}

	if len(result.SearchObjects) == 0 {
		return []string{"Not found"}, nil
	}

	var out []string
	for _, obj := range result.SearchObjects {
		item := obj.Content
		link := fmt.Sprintf("https://es.wallapop.com/item/%s", item.WebSlug)
		image := ""
		text := fmt.Sprintf("title: %s\n price: %.2f\n link: %s\n %s\n", item.Title, item.Price, link, image)

		out = append(out, text)
	}
	return out, nil
}
