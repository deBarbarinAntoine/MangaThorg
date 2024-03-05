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
	"sync"
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

func cacheCover(id string, index int, covers *[]models.Cover, wg *sync.WaitGroup, coversNotFound *[]int) {
	defer wg.Done()
	if checkStatus(models.Status.Covers, id) {
		coverCache := retrieveSingleCacheData(models.Status.Covers, id)
		apiCover, err := coverCache.Cover()
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		}
		log.Println("retrieving cover from cache") // testing
		(*covers)[index] = apiCover
		return
	}
	*coversNotFound = append(*coversNotFound, index)
}
