package api

import (
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"mangathorg/internal/models"
	"mangathorg/internal/utils"
	"net/http"
	"time"
)

var baseURL string = "https://api.mangadex.org/"
var client = http.Client{
	Timeout: time.Second * 5,
}

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

	endpoint := "manga"
	req, errReq := http.NewRequest(http.MethodGet, baseURL+endpoint, nil)
	if errReq != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errReq))
	}
	q := req.URL.Query()
	request.ToQuery(&q)
	req.URL.RawQuery = q.Encode()

	res, errRes := client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	} else {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errRes))
	}

	body, errBody := io.ReadAll(res.Body)
	if errBody != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errBody))
	}

	var response models.ApiManga
	err := json.Unmarshal(body, &response)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}

	if info != "" {
		err = response.SingleCacheData(request.OrderValue).Write(dataPath+info+".json", id != "")
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
		updateCacheStatus(info, id)
	}

	return response
}

func CoverRequest(id string) models.ApiCover {
	if checkStatus(models.Status.Covers, id) {
		coverCache := retrieveSingleCacheData(models.Status.Covers, id)
		apiCover, err := coverCache.ApiCover()
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
		log.Println("retrieving cover from cache") // testing
		return apiCover
	}

	endpoint := "cover/"
	req, errReq := http.NewRequest(http.MethodGet, baseURL+endpoint+id, nil)
	if errReq != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errReq))
	}

	res, errRes := client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	} else {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errRes))
	}

	body, errBody := io.ReadAll(res.Body)
	if errBody != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errBody))
	}

	var response models.ApiCover
	err := json.Unmarshal(body, &response)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}

	err = response.SingleCacheData("").Write(dataPath+models.Status.Covers+".json", id != "")
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	updateCacheStatus(models.Status.Covers, id)

	return response
}
