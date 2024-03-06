package models

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var client = http.Client{
	Timeout: time.Second * 5,
}

func (r MangaRequest) Params() MangaRequestParam {
	return MangaRequestParam{
		Order:              "order[" + r.OrderType + "]",
		IncludedTags:       "includedTags",
		ExcludedTags:       "excludedTags",
		TranslatedLanguage: "availableTranslatedLanguage[]",
		Title:              "title",
		Author:             "author",
		AuthorOrArtist:     "authorOrArtist",
		Status:             "status[]",
		Public:             "publicationDemographic[]",
		ContentRating:      "contentRating[]",
		Limit:              "limit",
		Offset:             "offset",
	}
}

func (r MangaRequest) ToQuery() url.Values {
	var q = make(url.Values)
	params := r.Params()
	q[params.Order] = []string{r.OrderValue}
	if r.IncludedTags != nil {
		q[params.IncludedTags] = r.IncludedTags
	}
	if r.ExcludedTags != nil {
		q[params.ExcludedTags] = r.ExcludedTags
	}
	if r.Title != "" {
		q[params.Title] = []string{r.Title}
	}
	if r.Author != "" {
		q[params.Author] = []string{r.Author}
	}
	if r.AuthorOrArtist != "" {
		q[params.AuthorOrArtist] = []string{r.AuthorOrArtist}
	}
	if r.Status != nil {
		q[params.Status] = r.Status
	}
	if r.Public != nil {
		q[params.Public] = r.Public
	}
	q[params.TranslatedLanguage] = []string{"en"}
	q[params.ContentRating] = []string{"safe"}
	q[params.Limit] = []string{strconv.Itoa(r.Limit)}
	q[params.Offset] = []string{strconv.Itoa(r.Offset)}
	return q
}

func (data *ApiManga) SingleCacheData(order string) SingleCacheData {
	var cache SingleCacheData
	if len(data.Data) == 1 {
		cache.Id = data.Data[0].Id
	}
	cache.UpdatedTime = time.Now()
	cache.Order = order
	cache.Data = data
	return cache
}

func (data *ApiCover) SingleCacheData(order string) SingleCacheData {
	var cache SingleCacheData
	if len(data.Data) == 1 {
		cache.Id = data.Data[0].Id
	}
	cache.UpdatedTime = time.Now()
	cache.Order = order
	cache.Data = data
	return cache
}

func (data *ApiTags) SingleCacheData(order string) SingleCacheData {
	if order != "" {
		log.Println(errors.New("error: ApiTags.SingleCacheData() order not null"))
		return SingleCacheData{}
	}
	var cache SingleCacheData
	cache.UpdatedTime = time.Now()
	cache.Data = data
	return cache
}

func (data *Cover) SingleCacheData() SingleCacheData {
	var cache SingleCacheData
	cache.Id = data.Id
	cache.UpdatedTime = time.Now()
	cache.Data = data
	return cache
}

func (data *Manga) SingleCacheData(id string) SingleCacheData {
	var cache SingleCacheData
	cache.Id = id
	cache.UpdatedTime = time.Now()
	cache.Data = data
	return cache
}

func (data *ApiMangaFeed) SingleCacheData(id string) SingleCacheData {
	var cache SingleCacheData
	cache.Id = id
	cache.UpdatedTime = time.Now()
	cache.Data = data
	return cache
}

func (data *ApiChapterScan) SingleCacheData(id string) SingleCacheData {
	var cache SingleCacheData
	cache.Id = id
	cache.UpdatedTime = time.Now()
	cache.Data = data
	return cache
}

func (data *ApiMangaStats) SingleCacheData(id string) SingleCacheData {
	var cache SingleCacheData
	cache.Id = id
	cache.UpdatedTime = time.Now()
	cache.Data = data
	return cache
}

func (data *ApiManga) SendRequest(baseURL string, endpoint string, query url.Values) error {
	body, err := Request(baseURL+endpoint, query)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return err
	}

	err = data.CheckResponse()
	if err != nil {
		return err
	}

	return nil
}

func (data *ApiSingleManga) SendRequest(baseURL string, endpoint string, query url.Values) error {
	body, err := Request(baseURL+endpoint, query)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return err
	}

	err = data.CheckResponse()
	if err != nil {
		return err
	}

	return nil
}

func (data *ApiCover) SendRequest(baseURL string, endpoint string, query url.Values) error {
	body, err := Request(baseURL+endpoint, query)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return err
	}

	err = data.CheckResponse()
	if err != nil {
		return err
	}

	return nil
}

func (data *ApiTags) SendRequest(baseURL string, endpoint string, query url.Values) error {
	body, err := Request(baseURL+endpoint, query)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return err
	}

	err = data.CheckResponse()
	if err != nil {
		return err
	}

	return nil
}

func (data *ApiMangaFeed) SendRequest(baseURL string, endpoint string, query url.Values) error {
	if query == nil {
		query = make(url.Values)
	}
	query.Add("translatedLanguage[]", "en")
	body, err := Request(baseURL+endpoint, query)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return err
	}

	err = data.CheckResponse()
	if err != nil {
		return err
	}

	return nil
}

func (data *ApiChapterScan) SendRequest(baseURL string, endpoint string, query url.Values) error {
	body, err := Request(baseURL+endpoint, query)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return err
	}

	err = data.CheckResponse()
	if err != nil {
		return err
	}

	return nil
}

func (data *ApiMangaStats) SendRequest(baseURL string, endpoint string, query url.Values) error {
	body, err := Request(baseURL+endpoint, query)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return err
	}

	err = data.CheckResponse()
	if err != nil {
		return err
	}

	return nil
}

func Request(url string, query url.Values) ([]byte, error) {
	req, errReq := http.NewRequest(http.MethodGet, url, nil)
	if errReq != nil {
		return nil, errReq
	}

	if query != nil {
		req.URL.RawQuery = query.Encode()
	}

	res, errRes := client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	} else {
		return nil, errRes
	}

	body, errBody := io.ReadAll(res.Body)
	if errBody != nil {
		return nil, errBody
	}
	return body, nil
}

func (data *ApiMangaStats) Stats(id string) Statistics {
	m := data.Statistics.(map[string]interface{})
	if value, ok := m[id].(interface{}); ok {
		var stats Statistics
		jsonData, err := json.Marshal(value)
		if err != nil {
			return Statistics{}
		}
		err = json.Unmarshal(jsonData, &stats)
		if err != nil {
			return Statistics{}
		}
		return stats
	}
	return Statistics{}
}

func (data *ApiCover) Divide() ([]Cover, error) {
	if data.Response == "collection" {
		var apiCovers []Cover
		for _, cover := range data.Data {
			apiCovers = append(apiCovers, cover)
		}
		return apiCovers, nil
	}
	return nil, errors.New("error: data is not a collection")
}

func (data *ApiSingleManga) CoverId() string {
	for _, relationship := range data.Data.Relationships {
		if relationship.Type == "cover_art" {
			return relationship.Id
		}
	}
	log.Println("cover_id not found")
	return ""
}

func (data *ApiManga) CoversId() ([]string, error) {
	var ids []string

	for _, manga := range data.Data {
		for _, relationship := range manga.Relationships {
			if relationship.Type == "cover_art" {
				ids = append(ids, relationship.Id)
			}
		}
	}
	var err error
	if len(ids) != len(data.Data) {
		err = errors.New("apiManga.CoversId: ids and mangas number doesn't match")
	}

	return ids, err
}

func (data *ApiTags) CheckResponse() error {
	if len(data.Errors) > 0 {
		var msg string
		for _, err := range data.Errors {
			msg += "error " + strconv.Itoa(err.Status) + ": " + err.Title + " -> " + err.Detail
		}
		return errors.New(msg)
	}
	return nil
}

func (data *ApiCover) CheckResponse() error {
	if len(data.Errors) > 0 {
		var msg string
		for _, err := range data.Errors {
			msg += "error " + strconv.Itoa(err.Status) + ": " + err.Title + " -> " + err.Detail
		}
		return errors.New(msg)
	}
	return nil
}

func (data *ApiManga) CheckResponse() error {
	if len(data.Errors) > 0 {
		var msg string
		for _, err := range data.Errors {
			msg += "error " + strconv.Itoa(err.Status) + ": " + err.Title + " -> " + err.Detail
		}
		return errors.New(msg)
	}
	return nil
}

func (data *ApiSingleManga) CheckResponse() error {
	if len(data.Errors) > 0 {
		var msg string
		for _, err := range data.Errors {
			msg += "error " + strconv.Itoa(err.Status) + ": " + err.Title + " -> " + err.Detail
		}
		return errors.New(msg)
	}
	return nil
}

func (data *ApiMangaFeed) CheckResponse() error {
	if len(data.Errors) > 0 {
		var msg string
		for _, err := range data.Errors {
			msg += "error " + strconv.Itoa(err.Status) + ": " + err.Title + " -> " + err.Detail
		}
		return errors.New(msg)
	}
	return nil
}

func (data *ApiChapterScan) CheckResponse() error {
	if len(data.Errors) > 0 {
		var msg string
		for _, err := range data.Errors {
			msg += "error " + strconv.Itoa(err.Status) + ": " + err.Title + " -> " + err.Detail
		}
		return errors.New(msg)
	}
	return nil
}

func (data *ApiMangaStats) CheckResponse() error {
	if len(data.Errors) > 0 {
		var msg string
		for _, err := range data.Errors {
			msg += "error " + strconv.Itoa(err.Status) + ": " + err.Title + " -> " + err.Detail
		}
		return errors.New(msg)
	}
	return nil
}
