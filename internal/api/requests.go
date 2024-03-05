package api

import (
	"errors"
	"log"
	"log/slog"
	"mangathorg/internal/models"
	"mangathorg/internal/utils"
	"net/url"
	"strconv"
	"sync"
)

var baseURL string = "https://api.mangadex.org/"

var TopPopularRequest = models.MangaRequest{
	OrderType:    "rating",
	OrderValue:   "desc",
	IncludedTags: nil,
	ExcludedTags: nil,
	Limit:        10,
	Offset:       0,
}

var TopLatestUploadedRequest = models.MangaRequest{
	OrderType:    "latestUploadedChapter",
	OrderValue:   "desc",
	IncludedTags: nil,
	ExcludedTags: nil,
	Limit:        10,
	Offset:       0,
}

func FetchManga(request models.MangaRequest) []models.MangaWhole {
	apiManga := MangaRequest(request)
	coversId, err := apiManga.CoversId()
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	covers := CoverRequest(coversId)
	var mangas []models.MangaWhole
	for i, manga := range apiManga.Data {
		var mangaWhole = models.MangaWhole{
			Manga: manga,
			Cover: covers[i],
		}
		mangas = append(mangas, mangaWhole)
	}

	return mangas
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
		err = apiManga.SingleCacheData(request.OrderValue).Write(dataPath+info+".json", id != "")
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
		updateCacheStatus(info, id)
	}

	return apiManga
}

func CoverRequest(ids []string) []models.Cover {
	var allCovers = make([]models.Cover, len(ids))
	var wg sync.WaitGroup
	var toRequest []int
	for i, id := range ids {
		wg.Add(1)
		go cacheCover(id, i, &allCovers, &wg, &toRequest)
	}
	wg.Wait()

	var apiCoverResponse models.ApiCover
	var query = make(url.Values)
	for _, i := range toRequest {
		query["ids[]"] = append(query["ids[]"], ids[i])
	}
	err := apiCoverResponse.SendRequest(baseURL, "cover", query)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}

	covers, errDiv := apiCoverResponse.Divide()
	if errDiv != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}

	for _, cover := range covers {
		cache := cover.SingleCacheData()
		err = cache.Write(dataPath+models.Status.Covers+".json", true)
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
		updateCacheStatus(models.Status.Covers, cache.Id)
	}
	for i, j := range toRequest {
		if i >= len(covers) {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("index ["+strconv.Itoa(i)+"] of ["+strconv.Itoa(len(covers))+"] out of range")))
			break
		}
		allCovers[j] = covers[i]
	}

	return allCovers
}

func TagsRequest() models.ApiTags {
	if checkStatus(models.Status.Tags, "") {
		tagCache := retrieveSingleCacheData(models.Status.Tags, "")
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

	err = apiTags.SingleCacheData("").Write(dataPath+models.Status.Tags+".json", false)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	updateCacheStatus(models.Status.Tags, "")

	return apiTags
}
