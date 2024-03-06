package api

import (
	"encoding/json"
	"log"
	"log/slog"
	"mangathorg/internal/models"
	"mangathorg/internal/utils"
	"os"
	"reflect"
	"slices"
	"time"
)

var dataPath string = utils.Path + "cache/"

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

func retrieveSingleCacheData(info string, id string) models.SingleCacheData {
	cacheData := retrieveCacheData(info)

	if id == "" && cacheData != nil {
		return cacheData[0]
	}

	for _, datum := range cacheData {
		if datum.Id == id {
			return datum
		}
	}
	return models.SingleCacheData{}
}

func checkCache(info string, id string) bool {
	if !checkStatus(info, id) {
		return false
	}

	cacheData := retrieveCacheData(info)

	return cacheData.Exists(id)
}

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
	case models.Status.Covers:
		return status.Covers != nil && slices.Contains(status.Covers, id)
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

func cacheRequest(r models.MangaRequest) models.SingleCacheData {
	switch {
	case reflect.DeepEqual(r, TopPopularRequest):
		return retrieveSingleCacheData(models.Status.Popular, "")
	case reflect.DeepEqual(r, TopLatestUploadedRequest):
		return retrieveSingleCacheData(models.Status.LastUploaded, "")
	default:
		return models.SingleCacheData{}
	}
}

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
	case models.Status.Covers:
		if slices.Contains(status.Covers, id) {
			return
		}
		status.Covers = append(status.Covers, id)
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

func cacheCover(id string, index int, covers *[]models.Cover) bool {
	if checkStatus(models.Status.Covers, id) {
		coverCache := retrieveSingleCacheData(models.Status.Covers, id)
		apiCover, err := coverCache.Cover()
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
		log.Println("retrieving cover from cache") // testing
		(*covers)[index] = apiCover
		return true
	}
	return false
}

func isOldCache(data models.SingleCacheData) bool {
	// possible evolution: to customize time limit, add a time parameter, or add the info parameter
	// to set a different time according to the data type (info).
	if time.Since(data.UpdatedTime) > time.Hour*24 {
		return true
	}
	return false
}

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
	case models.Status.Covers:
		index := slices.Index(status.Covers, id)
		if index != -1 {
			status.Covers = append(status.Covers[:index], status.Covers[index+1:]...)
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
		cacheData, err = cacheData.Delete(data.Id)
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
	}
	deleteCacheStatus(info, data.Id)
}

func clearCache(info string) {
	cache := retrieveCacheData(info)
	for _, data := range cache {
		if isOldCache(data) {
			deleteCacheData(info, data)
		}
	}
}

func CacheMonitor() {
	time.Sleep(time.Second * 10)
	hour := 3
	var duration time.Duration
	var infos = []string{models.Status.Popular, models.Status.LastUploaded, models.Status.Tags, models.Status.Mangas, models.Status.Covers, models.Status.MangaFeeds, models.Status.MangaStats, models.Status.ChaptersScan, models.Status.Categories}
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
