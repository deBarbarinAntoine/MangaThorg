package api

import (
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"mangathorg/internal/models"
	"mangathorg/internal/utils"
	"os"
	"reflect"
	"slices"
	"time"
)

// dataPath is the absolute path to the cache directory.
var dataPath string = utils.Path + "cache/"

// retrieveCacheData
//
//	@Description: retrieves a specific cache's file content and returns it as a models.CacheData.
//	@param info: name of the cache file.
//	@return models.CacheData
func retrieveCacheData(info string) models.CacheData {
	data, err := os.ReadFile(dataPath + info + ".json")
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}

	if len(data) == 0 {
		return nil
	}

	var cacheData models.CacheData
	err = json.Unmarshal(data, &cacheData)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	return cacheData
}

// retrieveSingleCacheData
//
//	@Description: retrieves a specific entry in a specific cache file and returns it as a models.SingleCacheData.
//	@param info: name of the cache file.
//	@param id: id of the item searched.
//	@param order
//	@param offset
//	@return models.SingleCacheData
func retrieveSingleCacheData(info string, id string, order string, offset int) models.SingleCacheData {
	cacheData := retrieveCacheData(info)

	if id == "" && cacheData != nil {
		return cacheData[0]
	} else if cacheData == nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("cache fetching error: CacheData not found")))
		return models.SingleCacheData{Data: nil}
	}

	for _, datum := range cacheData {
		if datum.Id == id && datum.Order == order && datum.Offset == offset {
			return datum
		}
	}
	utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("cache fetching error: SingleCacheData not found")))
	return models.SingleCacheData{Data: nil}
}

// checkStatus
//
//	@Description: checks if the item's id is present in the `info` category in status.json.
//	@param info: kind of item.
//	@param id: the item's id.
//	@return bool
func checkStatus(info string, id string) bool {
	data, err := os.ReadFile(dataPath + "status.json")
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}

	if len(data) == 0 {
		return false
	}

	var status models.StatusCache
	err = json.Unmarshal(data, &status)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}

	switch info {
	case models.Status.LastUploaded:
		return status.LastUploaded
	case models.Status.Popular:
		return status.Popular
	case models.Status.Categories:
		return status.Categories != nil && slices.Contains(status.Categories, id)
	case models.Status.ChaptersScan:
		return status.ChaptersScan != nil && slices.Contains(status.ChaptersScan, id)
	case models.Status.MangaFeeds:
		return status.MangaFeeds != nil && slices.Contains(status.MangaFeeds, id)
	case models.Status.MangaStats:
		return status.MangaStats != nil && slices.Contains(status.MangaStats, id)
	case models.Status.Mangas:
		return status.Mangas != nil && slices.Contains(status.Mangas, id)
	case models.Status.Tags:
		return status.Tags
	default:
		return false
	}
}

// isCache
//
//	@Description: checks if the request has already been cached.
//	@param r
//	@return bool
//	@return string
//	@return string
func isCache(r models.MangaRequest) (bool, string, string) {
	switch {
	case reflect.DeepEqual(r, TopPopularRequest):
		return checkStatus(models.Status.Popular, ""), models.Status.Popular, ""
	case reflect.DeepEqual(r, TopLatestUploadedRequest):
		return checkStatus(models.Status.LastUploaded, ""), models.Status.LastUploaded, ""
	default:
		return false, "", ""
	}
}

// cacheRequest
//
//	@Description: retrieves a cached request.
//	@param r
//	@return models.SingleCacheData
func cacheRequest(r models.MangaRequest) models.SingleCacheData {
	switch {
	case reflect.DeepEqual(r, TopPopularRequest):
		return retrieveSingleCacheData(models.Status.Popular, "", "", 0)
	case reflect.DeepEqual(r, TopLatestUploadedRequest):
		return retrieveSingleCacheData(models.Status.LastUploaded, "", "", 0)
	default:
		return models.SingleCacheData{}
	}
}

// updateCacheStatus
//
//	@Description: updates an item's cache status.
//	@param info: type of item.
//	@param id: the item's id.
func updateCacheStatus(info string, id string) {
	data, err := os.ReadFile(dataPath + "status.json")
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}

	var status models.StatusCache
	if len(data) != 0 {
		err = json.Unmarshal(data, &status)
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
	}

	switch info {
	case models.Status.LastUploaded:
		status.LastUploaded = true
	case models.Status.Popular:
		status.Popular = true
	case models.Status.Categories:
		if slices.Contains(status.Categories, id) {
			return
		}
		status.Categories = append(status.Categories, id)
	case models.Status.ChaptersScan:
		if slices.Contains(status.ChaptersScan, id) {
			return
		}
		status.ChaptersScan = append(status.ChaptersScan, id)
	case models.Status.MangaFeeds:
		if slices.Contains(status.MangaFeeds, id) {
			return
		}
		status.MangaFeeds = append(status.MangaFeeds, id)
	case models.Status.MangaStats:
		if slices.Contains(status.MangaStats, id) {
			return
		}
		status.MangaStats = append(status.MangaStats, id)
	case models.Status.Mangas:
		if slices.Contains(status.Mangas, id) {
			return
		}
		status.Mangas = append(status.Mangas, id)
	case models.Status.Tags:
		status.Tags = true
	}

	data, errJSON := json.MarshalIndent(status, "", "\t")
	if errJSON != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errJSON))
	}
	errWrite := os.WriteFile(dataPath+"status.json", data, 0666)
	if errWrite != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errWrite))
	}
}

// isOldCache
//
//	@Description: checks if a cached item is outdated or not.
//	@param data
//	@return bool
func isOldCache(data models.SingleCacheData) bool {
	// possible evolution: to customize time limit, add a time parameter, or add the info parameter
	// to set a different time according to the data type (info).
	if time.Since(data.UpdatedTime) > time.Hour*24 {
		return true
	}
	return false
}

// deleteCacheStatus
//
//	@Description: removes the item's id in status.json.
//	@param info: item's type.
//	@param id: item's id.
func deleteCacheStatus(info string, id string) {
	data, err := os.ReadFile(dataPath + "status.json")
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}

	var status models.StatusCache
	if len(data) != 0 {
		err = json.Unmarshal(data, &status)
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
	}

	switch info {
	case models.Status.LastUploaded:
		status.LastUploaded = false
	case models.Status.Popular:
		status.Popular = false
	case models.Status.Categories:
		index := slices.Index(status.Categories, id)
		if index != -1 {
			status.Categories = append(status.Categories[:index], status.Categories[index+1:]...)
		}
	case models.Status.ChaptersScan:
		index := slices.Index(status.ChaptersScan, id)
		if index != -1 {
			status.ChaptersScan = append(status.ChaptersScan[:index], status.ChaptersScan[index+1:]...)
		}
	case models.Status.MangaFeeds:
		index := slices.Index(status.MangaFeeds, id)
		if index != -1 {
			status.MangaFeeds = append(status.MangaFeeds[:index], status.MangaFeeds[index+1:]...)
		}
	case models.Status.MangaStats:
		index := slices.Index(status.MangaStats, id)
		if index != -1 {
			status.MangaStats = append(status.MangaStats[:index], status.MangaStats[index+1:]...)
		}
	case models.Status.Mangas:
		index := slices.Index(status.Mangas, id)
		if index != -1 {
			status.Mangas = append(status.Mangas[:index], status.Mangas[index+1:]...)
		}
	case models.Status.Tags:
		status.Tags = false
	}

	data, errJSON := json.MarshalIndent(status, "", "\t")
	if errJSON != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errJSON))
	}
	errWrite := os.WriteFile(dataPath+"status.json", data, 0666)
	if errWrite != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errWrite))
	}
}

// deleteCacheData
//
//	@Description: removes a single item cached.
//	@param info: item's type.
//	@param data: item.
func deleteCacheData(info string, data models.SingleCacheData) {
	switch info {
	case models.Status.LastUploaded, models.Status.Popular, models.Status.Tags:
		err := os.WriteFile(dataPath+info+".json", []byte("[]"), 0666)
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
	default:
		cacheData := retrieveCacheData(info)
		var err error
		cacheData, err = cacheData.Delete(data.Id, data.Order, data.Offset)
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
	}
	deleteCacheStatus(info, data.Id)
}

// clearCache
//
//	@Description: clears all outdated caches of a kind.
//	@param info: kind of cache.
func clearCache(info string) {
	cache := retrieveCacheData(info)
	for _, data := range cache {
		if isOldCache(data) {
			deleteCacheData(info, data)
		}
	}
}

// CacheMonitor
//
//	@Description: checks and clears the cache periodically according to its
//	validity (meant to be a goroutine).
func CacheMonitor() {
	time.Sleep(time.Second * 10)
	hour := 3
	var duration time.Duration
	var infos = []string{models.Status.Popular, models.Status.LastUploaded, models.Status.Tags, models.Status.Mangas, models.Status.MangaFeeds, models.Status.MangaStats, models.Status.ChaptersScan, models.Status.Categories}
	for {
		utils.Logger.Info(utils.GetCurrentFuncName(), slog.String("goroutine", "CacheMonitor"))
		for _, status := range infos {
			info := status
			go clearCache(info)
		}
		if time.Now().Hour() != hour || time.Now().Minute() > 15 {
			duration = utils.SetDailyTimer(hour)
		} else {
			duration = time.Hour * 24
		}
		time.Sleep(duration)
	}
}

// emptyFile
//
//	@Description: empties a cache file.
//	@param status: cache type.
func emptyFile(status string) {
	err := os.WriteFile(dataPath+status+".json", []byte("[]"), 0666)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		return
	}
}

// emptyCacheStatus
//
//	@Description: empties status.json.
func emptyCacheStatus() {
	data, err := json.MarshalIndent(models.StatusCache{}, "", "\t")
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		return
	}
	err = os.WriteFile(dataPath+"status.json", data, 0666)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
}

// EmptyCache
//
//	@Description: empties the whole cache (all specific files and status.json).
func EmptyCache() {
	log.Println("Emptying cache...")
	var infos = []string{models.Status.Popular, models.Status.LastUploaded, models.Status.Tags, models.Status.Mangas, models.Status.MangaFeeds, models.Status.MangaStats, models.Status.ChaptersScan, models.Status.Categories}
	for _, status := range infos {
		emptyFile(status)
	}
	emptyCacheStatus()
}
