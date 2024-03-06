package models

import (
	"encoding/json"
	"errors"
	"os"
)

func (cache CacheData) Exists(id string) bool {
	for _, data := range cache {
		if data.Id == id {
			return true
		}
	}
	return false
}

func (cache CacheData) Select(id string) SingleCacheData {
	for _, data := range cache {
		if data.Id == id {
			return data
		}
	}
	return SingleCacheData{}
}

func (cache CacheData) Update(newDatum SingleCacheData) (CacheData, error) {
	for ind, data := range cache {
		if data.Id == newDatum.Id {
			cache[ind] = newDatum
			return cache, nil
		}
	}
	return nil, errors.New("SingleCacheData not found")
}

func (cache CacheData) Delete(id string) (CacheData, error) {
	var idx int
	var found bool
	for i, single := range cache {
		if single.Id == id {
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

func (datum SingleCacheData) Cover() (Cover, error) {
	var cover Cover
	data, err := json.Marshal(datum.Data)
	if err != nil {
		return Cover{}, err
	}
	err = json.Unmarshal(data, &cover)
	if err != nil {
		return Cover{}, err
	}
	return cover, nil
}

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

func (datum SingleCacheData) Write(filePath string, Append bool) error {
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
	if cacheData.Exists(datum.Id) {
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
