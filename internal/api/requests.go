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
		log.Printf(utils.GetCurrentFuncName()+" info: %#v\n", info)
		log.Printf(utils.GetCurrentFuncName()+" id: %#v\n", id)
		log.Printf(utils.GetCurrentFuncName()+" cache exists: %#v\n", exists)
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

	err = response.SingleCacheData(request.OrderValue).Write(dataPath+info+".json", id != "")
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	updateCacheStatus(info, id)
	//log.Printf(utils.GetCurrentFuncName()+" response: %#v\n", response) // testing
	log.Printf(utils.GetCurrentFuncName()+" info: %#v\n", info)
	log.Printf(utils.GetCurrentFuncName()+" id: %#v\n", id)
	log.Printf(utils.GetCurrentFuncName()+" cache exists: %#v\n", exists)

	return response
}
