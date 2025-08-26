package api

import (
	"errors"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"sync"
	
	"mangathorg/internal/models/api"
	"mangathorg/internal/models/server"
	"mangathorg/internal/utils"
)

// BaseApiURL is the common URL used for all MangaDex API requests.
var BaseApiURL string = "https://api.mangadex.org/"

// TopPopularRequest is the exact request done to retrieve the six most popular
// mangas for the principal page.
var TopPopularRequest = api.MangaRequest{
	OrderType:    "rating",
	OrderValue:   "desc",
	IncludedTags: nil,
	ExcludedTags: nil,
	Limit:        6,
	Offset:       0,
}

// TopLatestUploadedRequest is the exact request done to retrieve the six latest
// uploaded mangas for the principal page.
var TopLatestUploadedRequest = api.MangaRequest{
	OrderType:    "latestUploadedChapter",
	OrderValue:   "desc",
	IncludedTags: nil,
	ExcludedTags: nil,
	Limit:        6,
	Offset:       0,
}

// FetchMangaById
//
//	@Description: fetches a specific manga according to its id,
//	with a set of chapters according to the `order` and `offset`.
//	@param id
//	@param order
//	@param offset
//	@return models.MangaUsefullData
func FetchMangaById(id string, order string, offset int) api.MangaUsefullData {
	if id == "" {
		return api.MangaUsefullData{}
	}
	var manga api.MangaUsefullData
	apiManga := MangaRequestById(id)
	
	manga = apiManga.Data.Format()
	feed := FeedRequest(id, order, offset)
	manga.Fill(StatRequest(id), feed)
	if order == "asc" {
		for i := range manga.Chapters {
			manga.Chapters[i].Offset = offset + i
		}
	} else {
		for i := range manga.Chapters {
			manga.Chapters[i].Offset = (feed.Total - 1) - (offset + i)
		}
	}
	
	return manga
}

// fillMangaListById
//
//	@Description: fetches a specific manga according to its id and add it to
//	`mangaList` (meant to be called as a goroutine to optimize timing).
//	@param id
//	@param order
//	@param offset
//	@param mangaList
//	@param wg
func fillMangaListById(id, order string, offset int, mangaList *[]api.MangaUsefullData, wg *sync.WaitGroup) {
	defer wg.Done()
	*mangaList = append(*mangaList, FetchMangaById(id, order, offset))
}

// FetchMangasById
//
//	@Description: fetches a list of mangas according to their ids.
//	@param favorites
//	@param order
//	@param offset
//	@return []models.MangaUsefullData
func FetchMangasById(favorites []server.MangaUser, order string, offset int) []api.MangaUsefullData {
	if favorites == nil {
		return nil
	}
	var mangas []api.MangaUsefullData
	var wg sync.WaitGroup
	
	for _, favorite := range favorites {
		wg.Add(1)
		go fillMangaListById(favorite.Id, order, offset, &mangas, &wg)
	}
	wg.Wait()
	
	var sortedMangas []api.MangaUsefullData
	
	for _, favorite := range favorites {
		for _, manga := range mangas {
			if favorite.Id == manga.Id {
				sortedMangas = append(sortedMangas, manga)
			}
		}
	}
	
	return sortedMangas
}

// MangaRequestById
//
//	@Description: requests a single manga according to its id.
//	@param id
//	@return models.ApiSingleManga
func MangaRequestById(id string) api.ApiSingleManga {
	if checkStatus(api.Status.Mangas, id) {
		mangaCache := retrieveSingleCacheData(api.Status.Mangas, id, "", 0)
		manga, err := mangaCache.Manga()
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
		log.Println("retrieving manga from cache") // testing
		
		// handling missing id in the manga cache data
		if manga.Id == "" {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("manga retrieved from cache has no Id")))
			deleteCacheData(api.Status.Mangas, mangaCache)
		} else {
			return api.ApiSingleManga{Data: manga}
		}
	}
	var apiSingleManga api.ApiSingleManga
	err := apiSingleManga.SendRequest(BaseApiURL, "manga/"+id, nil)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	
	err = apiSingleManga.Data.SingleCacheData(id, "desc", 0).Write(utils.DataPath+api.Status.Mangas+".json", true)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	updateCacheStatus(api.Status.Mangas, id)
	
	if reflect.DeepEqual(apiSingleManga, api.ApiSingleManga{}) {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("empty apiSingleManga")))
	}
	
	return apiSingleManga
}

// FetchManga
//
//	@Description: fetches mangas according to a request.
//	@param request
//	@return models.MangasInBulk
func FetchManga(request api.MangaRequest) api.MangasInBulk {
	apiManga := MangaRequest(request)
	
	return apiManga.Format()
}

// MangaRequest
//
//	@Description: requests a list of mangas.
//	@param request
//	@return models.ApiManga
func MangaRequest(request api.MangaRequest) api.ApiManga {
	var exists bool
	var info, id string
	if exists, info, id = isCache(request); exists {
		apiManga, err := cacheRequest(request).ApiManga()
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
		return apiManga
	}
	
	var apiManga api.ApiManga
	err := apiManga.SendRequest(BaseApiURL, "manga", request.ToQuery())
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	
	if info != "" {
		err = apiManga.SingleCacheData("", request.OrderValue, request.Offset).Write(utils.DataPath+info+".json", id != "")
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
		updateCacheStatus(info, id)
	}
	return apiManga
}

// TagsRequest
//
//	@Description: requests all tags from MangaDex API.
//	@return models.ApiTags
func TagsRequest() api.ApiTags {
	if checkStatus(api.Status.Tags, "") {
		tagCache := retrieveSingleCacheData(api.Status.Tags, "", "", 0)
		apiTags, err := tagCache.ApiTags()
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
		log.Println("retrieving tags from cache") // testing
		return apiTags
	}
	
	var apiTags api.ApiTags
	err := apiTags.SendRequest(BaseApiURL, "manga/tag", nil)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	
	err = apiTags.SingleCacheData("", "", 0).Write(utils.DataPath+api.Status.Tags+".json", false)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	updateCacheStatus(api.Status.Tags, "")
	
	return apiTags
}

// TagSelect
//
//	@Description: selects a tag according to its id.
//	@param id
//	@return models.ApiTag
func TagSelect(id string) api.ApiTag {
	tags := TagsRequest()
	for _, tag := range tags.Data {
		if tag.Id == id {
			return tag
		}
	}
	utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("tag [id:"+id+"] not found")))
	return api.ApiTag{}
}

// FetchSortedTags
//
//	@Description: fetches all tags (public and status included) and sort them by
//	type.
//	@return models.OrderedTags
func FetchSortedTags() api.OrderedTags {
	allTags := TagsRequest().Data
	var orderedTags api.OrderedTags
	for _, tag := range allTags {
		switch tag.Attributes.Group {
		case "format":
			orderedTags.FormatTags = append(orderedTags.FormatTags, tag)
		case "genre":
			orderedTags.GenreTags = append(orderedTags.GenreTags, tag)
		case "theme":
			orderedTags.ThemeTags = append(orderedTags.ThemeTags, tag)
		}
	}
	orderedTags.PublicTags = api.MangaPublic
	orderedTags.StatusTags = api.MangaStatus
	return orderedTags
}

// FeedRequest
//
//	@Description: requests a specific list of chapters according to the manga's
//	`id`, the `order` and the `offset`.
//	@param id
//	@param order
//	@param offset
//	@return models.ApiMangaFeed
func FeedRequest(id, order string, offset int) api.ApiMangaFeed {
	
	// retrieving the total number of chapters
	var total int
	if checkStatus(api.Status.MangaFeeds, id) {
		feedCache := retrieveSingleCacheData(api.Status.MangaFeeds, id, "desc", 0)
		apiMangaFeed, err := feedCache.ApiMangaFeed()
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
		total = apiMangaFeed.Total
	} else {
		var apiMangaFeed api.ApiMangaFeed
		
		var query = make(url.Values)
		query.Add("order[chapter]", order)
		query.Add("translatedLanguage[]", "en")
		query.Add("contentRating[]", "safe")
		query.Add("includes[]", "scanlation_group")
		
		err := apiMangaFeed.SendRequest(BaseApiURL, "manga/"+id+"/feed", query)
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
		total = apiMangaFeed.Total
	}
	
	// checking the offset value
	if total != 0 && total <= offset {
		offset = (total / 15) - 1
	}
	
	if checkStatus(api.Status.MangaFeeds, id) {
		feedCache := retrieveSingleCacheData(api.Status.MangaFeeds, id, order, offset)
		if feedCache.Data != nil {
			apiMangaFeed, err := feedCache.ApiMangaFeed()
			if err != nil {
				utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
			}
			log.Println("retrieving feed from cache") // testing
			return apiMangaFeed
		}
	}
	
	var apiMangaFeed api.ApiMangaFeed
	
	var query = make(url.Values)
	query.Add("order[chapter]", order)
	query.Add("translatedLanguage[]", "en")
	query.Add("contentRating[]", "safe")
	query.Add("includes[]", "scanlation_group")
	query.Add("limit", "15")
	query.Add("offset", strconv.Itoa(offset))
	
	err := apiMangaFeed.SendRequest(BaseApiURL, "manga/"+id+"/feed", query)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	
	err = apiMangaFeed.SingleCacheData(id, order, offset).Write(utils.DataPath+api.Status.MangaFeeds+".json", true)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	updateCacheStatus(api.Status.MangaFeeds, id)
	
	return apiMangaFeed
}

// ScanRequest
//
//	@Description: requests a chapter's scans according to its `id`.
//	@param id
//	@return models.ApiChapterScan
func ScanRequest(id string) api.ApiChapterScan {
	if checkStatus(api.Status.ChaptersScan, id) {
		scanCache := retrieveSingleCacheData(api.Status.ChaptersScan, id, "", 0)
		apiChapterScan, err := scanCache.ApiChapterScan()
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
		log.Println("retrieving chapterScan from cache") // testing
		return apiChapterScan
	}
	var apiChapterScan api.ApiChapterScan
	err := apiChapterScan.SendRequest(BaseApiURL, "at-home/server/"+id, nil)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	
	err = apiChapterScan.SingleCacheData(id, "", 0).Write(utils.DataPath+api.Status.ChaptersScan+".json", true)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	updateCacheStatus(api.Status.ChaptersScan, id)
	
	return apiChapterScan
}

// StatRequest
//
//	@Description: requests a manga's statistics according to its `id`.
//	@param id
//	@return models.Statistics
func StatRequest(id string) api.Statistics {
	if checkStatus(api.Status.MangaStats, id) {
		statCache := retrieveSingleCacheData(api.Status.MangaStats, id, "", 0)
		apiMangaStats, err := statCache.ApiMangaStats()
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
		log.Println("retrieving chapterStats from cache") // testing
		
		mangaStats := apiMangaStats.Stats(id)
		if reflect.DeepEqual(mangaStats, api.Statistics{}) {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("unable to extract statistics from interface")))
		}
		return mangaStats
	}
	var apiMangaStats api.ApiMangaStats
	err := apiMangaStats.SendRequest(BaseApiURL, "statistics/manga/"+id, nil)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	err = apiMangaStats.SingleCacheData(id, "", 0).Write(utils.DataPath+api.Status.MangaStats+".json", true)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	updateCacheStatus(api.Status.MangaStats, id)
	
	mangaStats := apiMangaStats.Stats(id)
	if reflect.DeepEqual(mangaStats, api.Statistics{}) {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("unable to extract statistics from interface")))
	}
	return mangaStats
}

// ImageProxy
//
//	@Description: requests a single cover image.
//	@param mangaId
//	@param pictureName
//	@return []byte
func ImageProxy(mangaId, pictureName string) []byte {
	reqUrl := "https://uploads.mangadex.org/covers/" + mangaId + "/" + pictureName
	data, err := api.Request(reqUrl, nil)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	
	return data
}

// ScanProxy
//
//	@Description: requests a single scan image.
//	@param chapterId
//	@param quality
//	@param hash
//	@param img
//	@return []byte
func ScanProxy(chapterId, quality, hash, img string) []byte {
	chapter := ScanRequest(chapterId)
	if chapter.Chapter.Hash != hash && chapter.Chapter.Hash != "" {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("invalid hash (chapterScan)")))
	}
	reqUrl := chapter.BaseUrl + "/" + quality + "/" + hash + "/" + img
	data, err := api.Request(reqUrl, nil)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		for i, datum := range chapter.Chapter.Data {
			if datum == img {
				img = chapter.Chapter.DataSaver[i]
			}
		}
		reqUrl = chapter.BaseUrl + "/dataSaver/" + hash + "/" + img
		data, err = api.Request(reqUrl, nil)
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
	}
	
	return data
}

// AddFavoriteInfo
//
//	@Description: adds some user related data to every manga of the list `mangas`.
//	@param r
//	@param mangas
//	@return bool
func AddFavoriteInfo(r *http.Request, mangas *[]api.MangaUsefullData) bool {
	session, sessionId := utils.GetSession(r)
	if sessionId == "" {
		return false
	}
	user, ok := utils.SelectUser(session.Username)
	if !ok {
		return false
	}
	for i, manga := range *mangas {
		for _, favorite := range user.Favorites {
			if favorite.Id == manga.Id {
				(*mangas)[i].IsFavorite = true
				(*mangas)[i].LastChapterRead = favorite.LastChapterRead
			}
		}
	}
	return true
}

// AddSingleFavoriteInfo
//
//	@Description: adds some user related data to a `manga`.
//	@param r
//	@param manga
//	@return bool
func AddSingleFavoriteInfo(r *http.Request, manga *api.MangaUsefullData) bool {
	session, sessionId := utils.GetSession(r)
	if sessionId == "" {
		return false
	}
	user, ok := utils.SelectUser(session.Username)
	if !ok {
		return false
	}
	for _, favorite := range user.Favorites {
		if favorite.Id == manga.Id {
			manga.IsFavorite = true
			manga.LastChapterRead = favorite.LastChapterRead
		}
	}
	return true
}
