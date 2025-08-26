package api

import (
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"os"
	"reflect"
	"slices"
	"time"
	
	"mangathorg/internal/models/api"
	"mangathorg/internal/utils"
)

// retrieveCacheData
//
//	@Description: retrieves a specific cache's file content and returns it as a models.CacheData.
//	@param info: name of the cache file.
//	@return models.CacheData
func retrieveCacheData(info string) api.CacheData {
	filename := utils.DataPath + info + ".json"
	
	api.RLock(filename)
	defer api.RUnlock(filename)
	
	data, err := os.ReadFile(filename)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	
	if len(data) == 0 {
		return nil
	}
	
	var cacheData api.CacheData
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
func retrieveSingleCacheData(info string, id string, order string, offset int) api.SingleCacheData {
	cacheData := retrieveCacheData(info)
	
	if id == "" && cacheData != nil {
		return cacheData[0]
	} else if cacheData == nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("cache fetching error: CacheData not found")))
		return api.SingleCacheData{Data: nil}
	}
	
	for _, datum := range cacheData {
		if datum.Id == id && datum.Order == order && datum.Offset == offset {
			return datum
		}
	}
	utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("cache fetching error: SingleCacheData not found")))
	return api.SingleCacheData{Data: nil}
}

// checkStatus
//
//	@Description: checks if the item's id is present in the `info` category in status.json.
//	@param info: kind of item.
//	@param id: the item's id.
//	@return bool
func checkStatus(info string, id string) bool {
	filename := utils.DataPath + "status.json"
	
	api.RLock(filename)
	defer api.RUnlock(filename)
	
	data, err := os.ReadFile(filename)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	
	if len(data) == 0 {
		return false
	}
	
	var status api.StatusCache
	err = json.Unmarshal(data, &status)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	
	switch info {
	case api.Status.LastUploaded:
		return status.LastUploaded
	case api.Status.Popular:
		return status.Popular
	case api.Status.Categories:
		return status.Categories != nil && slices.Contains(status.Categories, id)
	case api.Status.ChaptersScan:
		return status.ChaptersScan != nil && slices.Contains(status.ChaptersScan, id)
	case api.Status.MangaFeeds:
		return status.MangaFeeds != nil && slices.Contains(status.MangaFeeds, id)
	case api.Status.MangaStats:
		return status.MangaStats != nil && slices.Contains(status.MangaStats, id)
	case api.Status.Mangas:
		return status.Mangas != nil && slices.Contains(status.Mangas, id)
	case api.Status.Tags:
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
func isCache(r api.MangaRequest) (bool, string, string) {
	switch {
	case reflect.DeepEqual(r, TopPopularRequest):
		return checkStatus(api.Status.Popular, ""), api.Status.Popular, ""
	case reflect.DeepEqual(r, TopLatestUploadedRequest):
		return checkStatus(api.Status.LastUploaded, ""), api.Status.LastUploaded, ""
	default:
		return false, "", ""
	}
}

// cacheRequest
//
//	@Description: retrieves a cached request.
//	@param r
//	@return models.SingleCacheData
func cacheRequest(r api.MangaRequest) api.SingleCacheData {
	switch {
	case reflect.DeepEqual(r, TopPopularRequest):
		return retrieveSingleCacheData(api.Status.Popular, "", "", 0)
	case reflect.DeepEqual(r, TopLatestUploadedRequest):
		return retrieveSingleCacheData(api.Status.LastUploaded, "", "", 0)
	default:
		return api.SingleCacheData{}
	}
}

// updateCacheStatus
//
//	@Description: updates an item's cache status.
//	@param info: type of item.
//	@param id: the item's id.
func updateCacheStatus(info string, id string) {
	filename := utils.DataPath + "status.json"
	
	api.RWLock(filename)
	defer api.RWUnlock(filename)
	
	data, err := os.ReadFile(filename)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	
	var status api.StatusCache
	if len(data) != 0 {
		err = json.Unmarshal(data, &status)
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
	}
	
	switch info {
	case api.Status.LastUploaded:
		status.LastUploaded = true
	case api.Status.Popular:
		status.Popular = true
	case api.Status.Categories:
		if slices.Contains(status.Categories, id) {
			return
		}
		status.Categories = append(status.Categories, id)
	case api.Status.ChaptersScan:
		if slices.Contains(status.ChaptersScan, id) {
			return
		}
		status.ChaptersScan = append(status.ChaptersScan, id)
	case api.Status.MangaFeeds:
		if slices.Contains(status.MangaFeeds, id) {
			return
		}
		status.MangaFeeds = append(status.MangaFeeds, id)
	case api.Status.MangaStats:
		if slices.Contains(status.MangaStats, id) {
			return
		}
		status.MangaStats = append(status.MangaStats, id)
	case api.Status.Mangas:
		if slices.Contains(status.Mangas, id) {
			return
		}
		status.Mangas = append(status.Mangas, id)
	case api.Status.Tags:
		status.Tags = true
	}
	
	data, errJSON := json.MarshalIndent(status, "", "\t")
	if errJSON != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errJSON))
	}
	errWrite := os.WriteFile(filename, data, 0666)
	if errWrite != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errWrite))
	}
}

// isOldCache
//
//	@Description: checks if a cached item is outdated or not.
//	@param data
//	@return bool
func isOldCache(data api.SingleCacheData) bool {
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
	filename := utils.DataPath + "status.json"
	
	api.RWLock(filename)
	defer api.RWUnlock(filename)
	
	data, err := os.ReadFile(filename)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
	
	var status api.StatusCache
	if len(data) != 0 {
		err = json.Unmarshal(data, &status)
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
	}
	
	switch info {
	case api.Status.LastUploaded:
		status.LastUploaded = false
	case api.Status.Popular:
		status.Popular = false
	case api.Status.Categories:
		index := slices.Index(status.Categories, id)
		if index != -1 {
			status.Categories = append(status.Categories[:index], status.Categories[index+1:]...)
		}
	case api.Status.ChaptersScan:
		index := slices.Index(status.ChaptersScan, id)
		if index != -1 {
			status.ChaptersScan = append(status.ChaptersScan[:index], status.ChaptersScan[index+1:]...)
		}
	case api.Status.MangaFeeds:
		index := slices.Index(status.MangaFeeds, id)
		if index != -1 {
			status.MangaFeeds = append(status.MangaFeeds[:index], status.MangaFeeds[index+1:]...)
		}
	case api.Status.MangaStats:
		index := slices.Index(status.MangaStats, id)
		if index != -1 {
			status.MangaStats = append(status.MangaStats[:index], status.MangaStats[index+1:]...)
		}
	case api.Status.Mangas:
		index := slices.Index(status.Mangas, id)
		if index != -1 {
			status.Mangas = append(status.Mangas[:index], status.Mangas[index+1:]...)
		}
	case api.Status.Tags:
		status.Tags = false
	}
	
	data, errJSON := json.MarshalIndent(status, "", "\t")
	if errJSON != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errJSON))
	}
	errWrite := os.WriteFile(filename, data, 0666)
	if errWrite != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errWrite))
	}
}

// deleteCacheData
//
//	@Description: removes a single item cached.
//	@param info: item's type.
//	@param data: item.
func deleteCacheData(info string, data api.SingleCacheData) {
	filename := utils.DataPath + info + ".json"
	
	api.RWLock(filename)
	defer api.RWUnlock(filename)
	
	switch info {
	case api.Status.LastUploaded, api.Status.Popular, api.Status.Tags:
		err := os.WriteFile(filename, []byte("[]"), 0666)
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
	hour := 1
	var duration time.Duration
	var infos = []string{api.Status.Popular, api.Status.LastUploaded, api.Status.Tags, api.Status.Mangas, api.Status.MangaFeeds, api.Status.MangaStats, api.Status.ChaptersScan, api.Status.Categories}
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
	filename := utils.DataPath + status + ".json"
	
	api.RWLock(filename)
	defer api.RWUnlock(filename)
	
	err := os.WriteFile(filename, []byte("[]"), 0666)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		return
	}
}

// emptyCacheStatus
//
//	@Description: empties status.json.
func emptyCacheStatus() {
	filename := utils.DataPath + "status.json"
	
	api.RWLock(filename)
	defer api.RWUnlock(filename)
	
	data, err := json.MarshalIndent(api.StatusCache{}, "", "\t")
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		return
	}
	err = os.WriteFile(filename, data, 0666)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
}

// EmptyCache
//
//	@Description: empties the whole cache (all specific files and status.json).
func EmptyCache() {
	log.Println("Emptying cache...")
	var infos = []string{api.Status.Popular, api.Status.LastUploaded, api.Status.Tags, api.Status.Mangas, api.Status.MangaFeeds, api.Status.MangaStats, api.Status.ChaptersScan, api.Status.Categories}
	for _, status := range infos {
		emptyFile(status)
	}
	emptyCacheStatus()
}
