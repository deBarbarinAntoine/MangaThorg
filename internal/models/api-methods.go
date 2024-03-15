package models

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var client = http.Client{
	Timeout: time.Second * 5,
}

var ApiErrorStatus bool = false

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

func (data *ApiManga) SingleCacheData(id string, order string, pag int) SingleCacheData {
	var cache SingleCacheData
	if len(data.Data) == 1 {
		cache.Id = data.Data[0].Id
	}
	cache.UpdatedTime = time.Now()
	cache.Order = order
	cache.Data = data
	cache.Page = pag
	return cache
}

func (data *ApiCover) SingleCacheData(id string, order string, pag int) SingleCacheData {
	var cache SingleCacheData
	if len(data.Data) == 1 {
		cache.Id = data.Data[0].Id
	}
	cache.UpdatedTime = time.Now()
	cache.Order = order
	cache.Page = pag
	cache.Data = data
	return cache
}

func (data *ApiTags) SingleCacheData(id string, order string, pag int) SingleCacheData {
	if order != "" {
		log.Println(errors.New("error: ApiTags.SingleCacheData() order not null"))
		return SingleCacheData{}
	}
	var cache SingleCacheData
	cache.UpdatedTime = time.Now()
	cache.Data = data
	cache.Page = pag
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

func (data *ApiMangaFeed) SingleCacheData(id string, order string, pag int) SingleCacheData {
	var cache SingleCacheData
	cache.Id = id
	cache.Order = order
	cache.Page = pag
	cache.UpdatedTime = time.Now()
	cache.Data = data
	return cache
}

func (data *ApiChapterScan) SingleCacheData(id string, order string, pag int) SingleCacheData {
	var cache SingleCacheData
	cache.Id = id
	cache.Page = pag
	cache.UpdatedTime = time.Now()
	cache.Data = data
	return cache
}

func (data *ApiMangaStats) SingleCacheData(id string, order string, pag int) SingleCacheData {
	var cache SingleCacheData
	cache.Id = id
	cache.Page = pag
	cache.UpdatedTime = time.Now()
	cache.Data = data
	return cache
}

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

func (data *ApiCover) SendRequest(baseURL string, endpoint string, query url.Values) error {
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

	var err error
	if res.StatusCode != 200 {
		err = errors.New("error " + res.Status)
	}

	return body, err
}

func (data *ApiMangaStats) Stats(id string) Statistics { // fixme
	if data.Statistics == nil {
		log.Println("ApiMangaStats.Stats(): data.Statistics is null") // testing
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

func (data *ApiManga) Format() []MangaUsefullData {
	var formattedMangas []MangaUsefullData
	for _, datum := range data.Data {
		formattedMangas = append(formattedMangas, datum.Format())
	}
	return formattedMangas
}

func (data *ApiSingleManga) Format() MangaUsefullData {
	return data.Data.Format()
}

func (data *Manga) Format() MangaUsefullData {

	var feed ApiMangaFeed
	var query = make(url.Values)
	query.Add("order[chapter]", "asc")
	query.Add("contentRating[]", "safe")
	query.Add("includes[]", "scanlation_group")
	query.Add("limit", "1")

	// fixme: optimization needed (takes too much time to load)
	//err := feed.SendRequest("https://api.mangadex.org/", "manga/"+data.Id+"/feed", query)
	//if err != nil {
	//	log.Println("request error:", err)
	//}
	//var firstChapterId string
	//if feed.Data != nil {
	//	firstChapterId = feed.Data[0].Id
	//}

	var manga = MangaUsefullData{
		Id:                     data.Id,
		Title:                  data.Attributes.Title.En,
		Author:                 "",
		Description:            data.Attributes.Description.En,
		FirstChapterId:         data.Attributes.LatestUploadedChapter, // fixme: change to FirstChapterId
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

func (manga *MangaUsefullData) Fill(stats Statistics, feed ApiMangaFeed) {
	manga.Rating = math.Round(stats.Rating.Bayesian*10) / 10
	manga.Chapters = feed.Format()
	manga.NbChapter = feed.Total
	manga.FirstChapterId = manga.Chapters[0].Id
}

func (data *ApiMangaFeed) Format() []ChapterUsefullData {
	var chapters []ChapterUsefullData
	for _, chapter := range data.Data {
		chapters = append(chapters, chapter.Format())
	}
	return chapters
}

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

func (data *ApiTags) CheckResponse() error {
	if len(data.Errors) > 0 {
		var msg string
		for _, err := range data.Errors {
			if err.Status >= 500 {
				ApiErrorStatus = true
			}
			msg += "error " + strconv.Itoa(err.Status) + ": " + err.Title + " -> " + err.Detail
		}
		return errors.New(msg)
	}
	ApiErrorStatus = false
	return nil
}

func (data *ApiCover) CheckResponse() error {
	if len(data.Errors) > 0 {
		var msg string
		for _, err := range data.Errors {
			if err.Status >= 500 {
				ApiErrorStatus = true
			}
			msg += "error " + strconv.Itoa(err.Status) + ": " + err.Title + " -> " + err.Detail
		}
		return errors.New(msg)
	}
	ApiErrorStatus = false
	return nil
}

func (data *ApiManga) CheckResponse() error {
	if len(data.Errors) > 0 {
		var msg string
		for _, err := range data.Errors {
			if err.Status >= 500 {
				ApiErrorStatus = true
			}
			msg += "error " + strconv.Itoa(err.Status) + ": " + err.Title + " -> " + err.Detail
		}
		return errors.New(msg)
	}
	ApiErrorStatus = false
	return nil
}

func (data *ApiSingleManga) CheckResponse() error {
	if len(data.Errors) > 0 {
		var msg string
		for _, err := range data.Errors {
			if err.Status >= 500 {
				ApiErrorStatus = true
			}
			msg += "error " + strconv.Itoa(err.Status) + ": " + err.Title + " -> " + err.Detail
		}
		return errors.New(msg)
	}
	ApiErrorStatus = false
	return nil
}

func (data *ApiMangaFeed) CheckResponse() error {
	if len(data.Errors) > 0 {
		var msg string
		for _, err := range data.Errors {
			if err.Status >= 500 {
				ApiErrorStatus = true
			}
			msg += "error " + strconv.Itoa(err.Status) + ": " + err.Title + " -> " + err.Detail
		}
		return errors.New(msg)
	}
	ApiErrorStatus = false
	return nil
}

func (data *ApiChapterScan) CheckResponse() error {
	if len(data.Errors) > 0 {
		var msg string
		for _, err := range data.Errors {
			if err.Status >= 500 {
				ApiErrorStatus = true
			}
			msg += "error " + strconv.Itoa(err.Status) + ": " + err.Title + " -> " + err.Detail
		}
		return errors.New(msg)
	}
	ApiErrorStatus = false
	return nil
}

func (data *ApiMangaStats) CheckResponse() error {
	if len(data.Errors) > 0 {
		var msg string
		for _, err := range data.Errors {
			if err.Status >= 500 {
				ApiErrorStatus = true
			}
			msg += "error " + strconv.Itoa(err.Status) + ": " + err.Title + " -> " + err.Detail
		}
		return errors.New(msg)
	}
	ApiErrorStatus = false
	return nil
}
