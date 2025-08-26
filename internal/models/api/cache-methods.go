package api

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	
	"mangathorg/internal/utils"
)

// cacheMutexes encapsulates the map to control access.
type cacheMutexes struct {
	mu sync.RWMutex
	m  map[string]*sync.RWMutex
}

// newCacheMutexes initializes the map and wraps it in the struct.
func newCacheMutexes() *cacheMutexes {
	return &cacheMutexes{
		m: map[string]*sync.RWMutex{
			filepath.Join(utils.DataPath, "status.json"):               new(sync.RWMutex),
			filepath.Join(utils.DataPath, Status.LastUploaded+".json"): new(sync.RWMutex),
			filepath.Join(utils.DataPath, Status.Popular+".json"):      new(sync.RWMutex),
			filepath.Join(utils.DataPath, Status.Tags+".json"):         new(sync.RWMutex),
			filepath.Join(utils.DataPath, Status.Categories+".json"):   new(sync.RWMutex),
			filepath.Join(utils.DataPath, Status.Mangas+".json"):       new(sync.RWMutex),
			filepath.Join(utils.DataPath, Status.MangaFeeds+".json"):   new(sync.RWMutex),
			filepath.Join(utils.DataPath, Status.ChaptersScan+".json"): new(sync.RWMutex),
			filepath.Join(utils.DataPath, Status.MangaStats+".json"):   new(sync.RWMutex),
		},
	}
}

// Get safely retrieves a mutex from the map.
func (c *cacheMutexes) Get(filename string) *sync.RWMutex {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.m[filename]
}

// CacheMutexes is the exported instance.
var CacheMutexes = newCacheMutexes()

// RLock locks a mutex with read access for a given filename.
func RLock(filename string) {
	if mutex := CacheMutexes.Get(filename); mutex != nil {
		mutex.RLock()
	}
}

// RUnlock unlocks a mutex with read access for a given filename.
func RUnlock(filename string) {
	if mutex := CacheMutexes.Get(filename); mutex != nil {
		mutex.RUnlock()
	}
}

// RWLock locks a mutex with read and write access for a given filename.
func RWLock(filename string) {
	if mutex := CacheMutexes.Get(filename); mutex != nil {
		mutex.Lock()
	}
}

// RWUnlock unlocks a mutex with read and write access for a given filename.
func RWUnlock(filename string) {
	if mutex := CacheMutexes.Get(filename); mutex != nil {
		mutex.Unlock()
	}
}

// Exists
//
//	@Description: checks if a singleCacheData exists in a CacheData.
//	@receiver cache
//	@param id
//	@param order
//	@param offset
//	@return bool
func (cache CacheData) Exists(id string, order string, offset int) bool {
	for _, data := range cache {
		if data.Id == id && data.Order == order && data.Offset == offset {
			return true
		}
	}
	return false
}

// Update
//
//	@Description: updates a SingleCacheData present in a CacheData.
//	@receiver cache
//	@param newDatum
//	@return CacheData
//	@return error
func (cache CacheData) Update(newDatum SingleCacheData) (CacheData, error) {
	for ind, data := range cache {
		if data.Id == newDatum.Id && data.Order == newDatum.Order && data.Offset == newDatum.Offset {
			cache[ind] = newDatum
			return cache, nil
		}
	}
	return nil, errors.New("singleCacheData not found")
}

// Delete
//
//	@Description: removes a SingleCacheData from a CacheData.
//	@receiver cache
//	@param id
//	@param order
//	@param offset
//	@return CacheData
//	@return error
func (cache CacheData) Delete(id string, order string, offset int) (CacheData, error) {
	var idx int
	var found bool
	for i, single := range cache {
		if single.Id == id && single.Order == order && single.Offset == offset {
			idx = i
			found = true
			break
		}
	}
	if found {
		return append(cache[:idx], cache[idx+1:]...), nil
	}
	return cache, errors.New("singleCacheData not found")
}

// ApiManga
//
//	@Description: converts a SingleCacheData to an ApiManga.
//	@receiver datum
//	@return ApiManga
//	@return error
func (datum SingleCacheData) ApiManga() (ApiManga, error) {
	var apiManga ApiManga
	data, err := json.Marshal(datum.Data)
	if err != nil {
		return ApiManga{}, err
	}
	err = json.Unmarshal(data, &apiManga)
	if err != nil {
		return ApiManga{}, err
	}
	return apiManga, nil
}

// Manga
//
//	@Description: converts a SingleCacheData to a Manga.
//	@receiver datum
//	@return Manga
//	@return error
func (datum SingleCacheData) Manga() (Manga, error) {
	var manga Manga
	data, err := json.Marshal(datum.Data)
	if err != nil {
		return Manga{}, err
	}
	err = json.Unmarshal(data, &manga)
	if err != nil {
		return Manga{}, err
	}
	return manga, nil
}

// ApiTags
//
//	@Description: converts a SingleCacheData to an ApiTags.
//	@receiver datum
//	@return ApiTags
//	@return error
func (datum SingleCacheData) ApiTags() (ApiTags, error) {
	var apiTags ApiTags
	data, err := json.Marshal(datum.Data)
	if err != nil {
		return ApiTags{}, err
	}
	err = json.Unmarshal(data, &apiTags)
	if err != nil {
		return ApiTags{}, err
	}
	return apiTags, nil
}

// ApiMangaFeed
//
//	@Description: converts a SingleCacheData to an ApiMangaFeed.
//	@receiver datum
//	@return ApiMangaFeed
//	@return error
func (datum SingleCacheData) ApiMangaFeed() (ApiMangaFeed, error) {
	var apiMangaFeed ApiMangaFeed
	data, err := json.Marshal(datum.Data)
	if err != nil {
		return ApiMangaFeed{}, err
	}
	err = json.Unmarshal(data, &apiMangaFeed)
	if err != nil {
		return ApiMangaFeed{}, err
	}
	return apiMangaFeed, nil
}

// ApiChapterScan
//
//	@Description: converts a SingleCacheData to an ApiChapterScan.
//	@receiver datum
//	@return ApiChapterScan
//	@return error
func (datum SingleCacheData) ApiChapterScan() (ApiChapterScan, error) {
	var apiChapterScan ApiChapterScan
	data, err := json.Marshal(datum.Data)
	if err != nil {
		return ApiChapterScan{}, err
	}
	err = json.Unmarshal(data, &apiChapterScan)
	if err != nil {
		return ApiChapterScan{}, err
	}
	return apiChapterScan, nil
}

// ApiMangaStats
//
//	@Description: converts a SingleCacheData to an ApiMangaStats.
//	@receiver datum
//	@return ApiMangaStats
//	@return error
func (datum SingleCacheData) ApiMangaStats() (ApiMangaStats, error) {
	var apiMangaStats ApiMangaStats
	data, err := json.Marshal(datum.Data)
	if err != nil {
		return ApiMangaStats{}, err
	}
	err = json.Unmarshal(data, &apiMangaStats)
	if err != nil {
		return ApiMangaStats{}, err
	}
	return apiMangaStats, nil
}

// Write
//
//	@Description: writes a SingleCacheData in a specific cache file appending it or not.
//	@receiver datum
//	@param filePath
//	@param Append
//	@return error
func (datum SingleCacheData) Write(filePath string, Append bool) error {
	
	RWLock(filePath)
	defer RWUnlock(filePath)
	
	var cacheData CacheData
	if Append {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		if len(data) != 0 {
			err = json.Unmarshal(data, &cacheData)
			if err != nil {
				return err
			}
		}
	}
	if cacheData.Exists(datum.Id, datum.Order, datum.Offset) {
		var err error
		cacheData, err = cacheData.Update(datum)
		if err != nil {
			return err
		}
	} else {
		cacheData = append(cacheData, datum)
	}
	data, errJSON := json.MarshalIndent(cacheData, "", "\t")
	if errJSON != nil {
		return errJSON
	}
	errWrite := os.WriteFile(filePath, data, 0666)
	if errWrite != nil {
		return errWrite
	}
	return nil
}
