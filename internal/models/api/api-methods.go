package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

// Client is the http.Client used for all API requests.
var Client = http.Client{
	Timeout: time.Second * 5,
}

// Params
//
//	@Description: generates all parameters names for a specific manga request.
//	@receiver r
//	@return MangaRequestParam
func (r MangaRequest) Params() MangaRequestParam {
	if r.OrderType == "" {
		r.OrderType = "rating"
	}
	return MangaRequestParam{
		Order:              "order[" + r.OrderType + "]",
		IncludedTags:       "includedTags[]",
		ExcludedTags:       "excludedTags[]",
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

// ToQuery
//
//	@Description: generates a query from a manga request.
//	@receiver r
//	@return url.Values
func (r MangaRequest) ToQuery() url.Values {
	var q = make(url.Values)
	params := r.Params()
	if r.OrderValue == "" {
		r.OrderValue = "desc"
	}
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

// SingleCacheData
//
//	@Description: converts an ApiManga to a SingleCacheData.
//	@receiver data
//	@param id
//	@param order
//	@param offset
//	@return SingleCacheData
func (data *ApiManga) SingleCacheData(id string, order string, offset int) SingleCacheData {
	var cache SingleCacheData
	if len(data.Data) == 1 {
		cache.Id = data.Data[0].Id
	}
	cache.UpdatedTime = time.Now()
	cache.Order = order
	cache.Data = data
	cache.Offset = offset
	return cache
}

// SingleCacheData
//
//	@Description: converts an ApiTags to a SingleCacheData.
//	@receiver data
//	@param id
//	@param order
//	@param offset
//	@return SingleCacheData
func (data *ApiTags) SingleCacheData(id string, order string, offset int) SingleCacheData {
	if order != "" {
		log.Println(errors.New("error: ApiTags.SingleCacheData() order not null"))
		return SingleCacheData{}
	}
	var cache SingleCacheData
	cache.UpdatedTime = time.Now()
	cache.Data = data
	cache.Offset = offset
	return cache
}

// SingleCacheData
//
//	@Description: converts a Manga to a SingleCacheData.
//	@receiver data
//	@param id
//	@param order
//	@param offset
//	@return SingleCacheData
func (data *Manga) SingleCacheData(id string, order string, offset int) SingleCacheData {
	var cache SingleCacheData
	cache.Id = id
	cache.UpdatedTime = time.Now()
	cache.Data = data
	return cache
}

// SingleCacheData
//
//	@Description: converts an ApiMangaFeed to a SingleCacheData.
//	@receiver data
//	@param id
//	@param order
//	@param offset
//	@return SingleCacheData
func (data *ApiMangaFeed) SingleCacheData(id string, order string, offset int) SingleCacheData {
	var cache SingleCacheData
	cache.Id = id
	cache.Order = order
	cache.Offset = offset
	cache.UpdatedTime = time.Now()
	cache.Data = data
	return cache
}

// SingleCacheData
//
//	@Description: converts a ApiChapterScan to a SingleCacheData.
//	@receiver data
//	@param id
//	@param order
//	@param offset
//	@return SingleCacheData
func (data *ApiChapterScan) SingleCacheData(id string, order string, offset int) SingleCacheData {
	var cache SingleCacheData
	cache.Id = id
	cache.Offset = offset
	cache.UpdatedTime = time.Now()
	cache.Data = data
	return cache
}

// SingleCacheData
//
//	@Description: converts an ApiMangaStats to a SingleCacheData.
//	@receiver data
//	@param id
//	@param order
//	@param offset
//	@return SingleCacheData
func (data *ApiMangaStats) SingleCacheData(id string, order string, offset int) SingleCacheData {
	var cache SingleCacheData
	cache.Id = id
	cache.Offset = offset
	cache.UpdatedTime = time.Now()
	cache.Data = data
	return cache
}

// SendRequest
//
//	@Description: sends an ApiManga request.
//	@receiver data
//	@param baseURL
//	@param endpoint
//	@param query
//	@return error
func (data *ApiManga) SendRequest(baseURL string, endpoint string, query url.Values) error {
	if query == nil {
		query = make(url.Values)
	}
	
	query.Add("includes[]", "cover_art")
	query.Add("includes[]", "author")
	
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

// SendRequest
//
//	@Description: sends an ApiSingleManga request.
//	@receiver data
//	@param baseURL
//	@param endpoint
//	@param query
//	@return error
func (data *ApiSingleManga) SendRequest(baseURL string, endpoint string, query url.Values) error {
	if query == nil {
		query = make(url.Values)
	}
	
	query.Add("includes[]", "cover_art")
	query.Add("includes[]", "author")
	
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

// SendRequest
//
//	@Description: sends an ApiTags request.
//	@receiver data
//	@param baseURL
//	@param endpoint
//	@param query
//	@return error
func (data *ApiTags) SendRequest(baseURL string, endpoint string, query url.Values) error {
	if query == nil {
		query = make(url.Values)
	}
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

// SendRequest
//
//	@Description: sends an ApiMangaFeed request.
//	@receiver data
//	@param baseURL
//	@param endpoint
//	@param query
//	@return error
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
		log.Println("ApiMangaFeed.SendRequest: unmarshal error!")
		return err
	}
	
	err = data.CheckResponse()
	if err != nil {
		return err
	}
	
	return nil
}

// SendRequest
//
//	@Description: sends an ApiChapterScan request.
//	@receiver data
//	@param baseURL
//	@param endpoint
//	@param query
//	@return error
func (data *ApiChapterScan) SendRequest(baseURL string, endpoint string, query url.Values) error {
	if query == nil {
		query = make(url.Values)
	}
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

// SendRequest
//
//	@Description: sends an ApiMangaStats request.
//	@receiver data
//	@param baseURL
//	@param endpoint
//	@param query
//	@return error
func (data *ApiMangaStats) SendRequest(baseURL string, endpoint string, query url.Values) error {
	if query == nil {
		query = make(url.Values)
	}
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

// Request
//
//	@Description: sends a request to an `url` with a `query`.
//	@param url
//	@param query
//	@return []byte
//	@return error
func Request(url string, query url.Values) ([]byte, error) {
	req, errReq := http.NewRequest(http.MethodGet, url, nil)
	if errReq != nil {
		return nil, errReq
	}
	
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	
	res, errRes := Client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	} else {
		return nil, errRes
	}
	
	body, errBody := io.ReadAll(res.Body)
	if errBody != nil {
		return nil, errBody
	}
	
	var err error
	if res.StatusCode != 200 {
		err = errors.New("error " + res.Status)
	}
	
	return body, err
}

// Stats
//
//	@Description: converts an ApiMangaStats to a Statistics structure.
//	@receiver data
//	@param id
//	@return Statistics
func (data *ApiMangaStats) Stats(id string) Statistics {
	if data.Statistics == nil {
		log.Println("ApiMangaStats.Stats(): data.Statistics is null") // testing
		return Statistics{}
	}
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

// Format
//
//	@Description: converts an ApiManga to a MangasInBulk with all needed data.
//	@receiver data
//	@return MangasInBulk
func (data *ApiManga) Format() MangasInBulk {
	var formattedMangas MangasInBulk
	var wg sync.WaitGroup
	for _, datum := range data.Data {
		wg.Add(1)
		go func(data *Manga, mangas *[]MangaUsefullData, wg *sync.WaitGroup) {
			defer wg.Done()
			*mangas = append(*mangas, data.Format())
		}(&datum, &formattedMangas.Mangas, &wg)
	}
	wg.Wait()
	formattedMangas.NbMangas = data.Total
	return formattedMangas
}

// Format
//
//	@Description: converts a Manga to a MangaUsefullData with all needed data.
//	@receiver data
//	@return MangaUsefullData
func (data *Manga) Format() MangaUsefullData {
	
	var feed ApiMangaFeed
	var query = make(url.Values)
	query.Add("order[chapter]", "asc")
	query.Add("contentRating[]", "safe")
	query.Add("includes[]", "scanlation_group")
	query.Add("limit", "1")
	
	err := feed.SendRequest("https://api.mangadex.org/", "manga/"+data.Id+"/feed", query)
	if err != nil {
		log.Println("request error:", err)
	}
	var firstChapterId string
	if feed.Data != nil && len(feed.Data) > 0 {
		firstChapterId = feed.Data[0].Id
	}
	
	var manga = MangaUsefullData{
		Id:                     data.Id,
		Title:                  data.Attributes.Title.En,
		Author:                 "",
		Description:            data.Attributes.Description.En,
		FirstChapterId:         firstChapterId,
		LastChapterId:          data.Attributes.LatestUploadedChapter,
		LastChapterNb:          data.Attributes.LastChapter,
		OriginalLanguage:       data.Attributes.OriginalLanguage,
		PublicationDemographic: data.Attributes.PublicationDemographic,
		Status:                 data.Attributes.Status,
		Year:                   data.Attributes.Year,
		Tags:                   data.Attributes.Tags,
		CoverId:                "",
		CoverImg:               "",
		Rating:                 0,
		Chapters:               nil,
		NbChapter:              feed.Total,
	}
	var isCover, isAuthor bool
	for _, relationship := range data.Relationships {
		if relationship.Type == "cover_art" && !isCover {
			manga.CoverId = relationship.Id
			manga.CoverImg = relationship.Attributes.FileName
			isCover = true
		}
		if relationship.Type == "author" && !isAuthor {
			manga.Author = relationship.Attributes.Name
			isAuthor = true
		}
	}
	return manga
}

// Fill
//
//	@Description: adds the statistics and chapter feed to a MangaUsefullData.
//	@receiver manga
//	@param stats
//	@param feed
func (manga *MangaUsefullData) Fill(stats Statistics, feed ApiMangaFeed) {
	manga.Rating = math.Round(stats.Rating.Bayesian*10) / 10
	manga.Chapters = feed.Format()
	manga.NbChapter = feed.Total
}

// Format
//
//	@Description: converts an ApiMangaFeed to a list of ChapterUsefullData.
//	@receiver data
//	@return []ChapterUsefullData
func (data *ApiMangaFeed) Format() []ChapterUsefullData {
	var chapters []ChapterUsefullData
	for _, chapter := range data.Data {
		chapters = append(chapters, chapter.Format())
	}
	return chapters
}

// Format
//
//	@Description: converts a Chapter to a ChapterUsefullData.
//	@receiver data
//	@return ChapterUsefullData
func (data *Chapter) Format() ChapterUsefullData {
	var chapter = ChapterUsefullData{
		Id:                 data.Id,
		Title:              data.Attributes.Title,
		Volume:             data.Attributes.Volume,
		Chapter:            data.Attributes.Chapter,
		Pages:              data.Attributes.Pages,
		TranslatedLanguage: data.Attributes.TranslatedLanguage,
		Uploader:           data.Attributes.Uploader,
		UpdatedAt:          data.Attributes.UpdatedAt,
		ScanlationGroupId:  "",
		ScanlationGroup:    "",
	}
	for _, relationship := range data.Relationships {
		if relationship.Type == "scanlation_group" {
			chapter.ScanlationGroupId = relationship.Id
			chapter.ScanlationGroup = relationship.Attributes.Name
			break
		}
	}
	return chapter
}

// CheckResponse
//
//	@Description: checks the MangaDex API response, looking for any error.
//	@receiver data
//	@return error
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

// CheckResponse
//
//	@Description: checks the MangaDex API response, looking for any error.
//	@receiver data
//	@return error
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

// CheckResponse
//
//	@Description: checks the MangaDex API response, looking for any error.
//	@receiver data
//	@return error
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

// CheckResponse
//
//	@Description: checks the MangaDex API response, looking for any error.
//	@receiver data
//	@return error
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

// CheckResponse
//
//	@Description: checks the MangaDex API response, looking for any error.
//	@receiver data
//	@return error
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

// CheckResponse
//
//	@Description: checks the MangaDex API response, looking for any error.
//	@receiver data
//	@return error
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
