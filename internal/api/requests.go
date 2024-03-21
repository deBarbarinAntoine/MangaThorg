package api

import (
	"errors"
	"log"
	"log/slog"
	"mangathorg/internal/models"
	"mangathorg/internal/utils"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

var baseURL string = "https://api.mangadex.org/"

var TopPopularRequest = models.MangaRequest{
	OrderType:    "rating",
	OrderValue:   "desc",
	IncludedTags: nil,
	ExcludedTags: nil,
	Limit:        6,
	Offset:       0,
}

var TopLatestUploadedRequest = models.MangaRequest{
	OrderType:    "latestUploadedChapter",
	OrderValue:   "desc",
	IncludedTags: nil,
	ExcludedTags: nil,
	Limit:        6,
	Offset:       0,
}

func FetchMangaById(id string, order string, offset int) models.MangaUsefullData {
	if id == "" {
		return models.MangaUsefullData{}
	}
	var manga models.MangaUsefullData
	apiManga := MangaRequestById(id)

	manga = apiManga.Data.Format()
	manga.Fill(StatRequest(id), FeedRequest(id, order, offset))

	return manga
}

func FetchMangasById(favorites []models.MangaUser, order string, offset int) []models.MangaUsefullData {
	if favorites == nil {
		return nil
	}
	var mangas []models.MangaUsefullData

	for _, favorite := range favorites {
		mangas = append(mangas, FetchMangaById(favorite.Id, order, offset))
	}

	return mangas
}

func MangaRequestById(id string) models.ApiSingleManga {
	if checkStatus(models.Status.Mangas, id) {
		mangaCache := retrieveSingleCacheData(models.Status.Mangas, id, "", 0)
		manga, err := mangaCache.Manga()
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
		log.Println("retrieving manga from cache") // testing

		// handling missing id in the manga cache data
		if manga.Id == "" {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("manga retrieved from cache has no Id")))
			deleteCacheData(models.Status.Mangas, mangaCache)
		} else {
			return models.ApiSingleManga{Data: manga}
		}
	}
	var apiSingleManga models.ApiSingleManga
	err := apiSingleManga.SendRequest(baseURL, "manga/"+id, nil)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}

	err = apiSingleManga.Data.SingleCacheData(id, "desc", 0).Write(dataPath+models.Status.Mangas+".json", true)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	updateCacheStatus(models.Status.Mangas, id)

	if reflect.DeepEqual(apiSingleManga, models.ApiSingleManga{}) {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("empty apiSingleManga")))
	}

	return apiSingleManga
}

func FetchManga(request models.MangaRequest) models.MangasInBulk {
	apiManga := MangaRequest(request)

	return apiManga.Format()
}

func MangaRequest(request models.MangaRequest) models.ApiManga {
	var exists bool
	var info, id string
	if exists, info, id = isCache(request); exists {
		apiManga, err := cacheRequest(request).ApiManga()
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
		return apiManga
	}

	var apiManga models.ApiManga
	err := apiManga.SendRequest(baseURL, "manga", request.ToQuery())
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}

	if info != "" {
		err = apiManga.SingleCacheData("", request.OrderValue, request.Offset).Write(dataPath+info+".json", id != "")
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
		updateCacheStatus(info, id)
	}
	return apiManga
}

func TagsRequest() models.ApiTags {
	if checkStatus(models.Status.Tags, "") {
		tagCache := retrieveSingleCacheData(models.Status.Tags, "", "", 0)
		apiTags, err := tagCache.ApiTags()
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
		log.Println("retrieving tags from cache") // testing
		return apiTags
	}

	var apiTags models.ApiTags
	err := apiTags.SendRequest(baseURL, "manga/tag", nil)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}

	err = apiTags.SingleCacheData("", "", 0).Write(dataPath+models.Status.Tags+".json", false)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	updateCacheStatus(models.Status.Tags, "")

	return apiTags
}

func TagSelect(id string) models.ApiTag {
	tags := TagsRequest()
	for _, tag := range tags.Data {
		if tag.Id == id {
			return tag
		}
	}
	utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("tag [id:"+id+"] not found")))
	return models.ApiTag{}
}

func FetchSortedTags() models.OrderedTags {
	allTags := TagsRequest().Data
	var orderedTags models.OrderedTags
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
	orderedTags.PublicTags = models.MangaPublic
	orderedTags.StatusTags = models.MangaStatus
	return orderedTags
}

func FeedRequest(id, order string, offset int) models.ApiMangaFeed {

	// retrieving the total number of chapters
	var total int
	if checkStatus(models.Status.MangaFeeds, id) {
		feedCache := retrieveSingleCacheData(models.Status.MangaFeeds, id, "desc", 0)
		apiMangaFeed, err := feedCache.ApiMangaFeed()
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
		total = apiMangaFeed.Total
	} else {
		var apiMangaFeed models.ApiMangaFeed

		var query = make(url.Values)
		query.Add("order[chapter]", order)
		query.Add("translatedLanguage[]", "en")
		query.Add("contentRating[]", "safe")
		query.Add("includes[]", "scanlation_group")

		err := apiMangaFeed.SendRequest(baseURL, "manga/"+id+"/feed", query)
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
		total = apiMangaFeed.Total
	}

	// checking the offset value
	if total != 0 && total <= offset {
		offset = (total / 15) - 1
	}

	if checkStatus(models.Status.MangaFeeds, id) {
		feedCache := retrieveSingleCacheData(models.Status.MangaFeeds, id, order, offset)
		if feedCache.Data != nil {
			apiMangaFeed, err := feedCache.ApiMangaFeed()
			if err != nil {
				utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
			}
			log.Println("retrieving feed from cache") // testing
			return apiMangaFeed
		}
	}

	var apiMangaFeed models.ApiMangaFeed

	var query = make(url.Values)
	query.Add("order[chapter]", order)
	query.Add("translatedLanguage[]", "en")
	query.Add("contentRating[]", "safe")
	query.Add("includes[]", "scanlation_group")
	query.Add("limit", "15")
	query.Add("offset", strconv.Itoa(offset))

	err := apiMangaFeed.SendRequest(baseURL, "manga/"+id+"/feed", query)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}

	err = apiMangaFeed.SingleCacheData(id, order, offset).Write(dataPath+models.Status.MangaFeeds+".json", true)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	updateCacheStatus(models.Status.MangaFeeds, id)

	return apiMangaFeed
}

func ScanRequest(id string) models.ApiChapterScan {
	if checkStatus(models.Status.ChaptersScan, id) {
		scanCache := retrieveSingleCacheData(models.Status.ChaptersScan, id, "", 0)
		apiChapterScan, err := scanCache.ApiChapterScan()
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
		log.Println("retrieving chapterScan from cache") // testing
		return apiChapterScan
	}
	var apiChapterScan models.ApiChapterScan
	err := apiChapterScan.SendRequest(baseURL, "at-home/server/"+id, nil)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}

	err = apiChapterScan.SingleCacheData(id, "", 0).Write(dataPath+models.Status.ChaptersScan+".json", true)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	updateCacheStatus(models.Status.ChaptersScan, id)

	return apiChapterScan
}

func StatRequest(id string) models.Statistics {
	if checkStatus(models.Status.MangaStats, id) {
		statCache := retrieveSingleCacheData(models.Status.MangaStats, id, "", 0)
		apiMangaStats, err := statCache.ApiMangaStats()
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
		log.Println("retrieving chapterStats from cache") // testing

		mangaStats := apiMangaStats.Stats(id)
		if reflect.DeepEqual(mangaStats, models.Statistics{}) {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("unable to extract statistics from interface")))
		}
		return mangaStats
	}
	var apiMangaStats models.ApiMangaStats
	err := apiMangaStats.SendRequest(baseURL, "statistics/manga/"+id, nil)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	err = apiMangaStats.SingleCacheData(id, "", 0).Write(dataPath+models.Status.MangaStats+".json", true)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	updateCacheStatus(models.Status.MangaStats, id)

	mangaStats := apiMangaStats.Stats(id)
	if reflect.DeepEqual(mangaStats, models.Statistics{}) {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("unable to extract statistics from interface")))
	}
	return mangaStats
}

func ImageProxy(mangaId, pictureName string) []byte {
	reqUrl := "https://uploads.mangadex.org/covers/" + mangaId + "/" + pictureName
	data, err := models.Request(reqUrl, nil)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}

	return data
}

func ScanProxy(chapterId, quality, hash, img string) []byte {
	chapter := ScanRequest(chapterId)
	if chapter.Chapter.Hash != hash && chapter.Chapter.Hash != "" {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("invalid hash (chapterScan)")))
	}
	reqUrl := chapter.BaseUrl + "/" + quality + "/" + hash + "/" + img
	data, err := models.Request(reqUrl, nil)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		for i, datum := range chapter.Chapter.Data {
			if datum == img {
				img = chapter.Chapter.DataSaver[i]
			}
		}
		reqUrl = chapter.BaseUrl + "/dataSaver/" + hash + "/" + img
		data, err = models.Request(reqUrl, nil)
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
	}

	return data
}

func AddFavoriteInfo(r *http.Request, mangas *[]models.MangaUsefullData) bool {
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

func AddSingleFavoriteInfo(r *http.Request, manga *models.MangaUsefullData) bool {
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
