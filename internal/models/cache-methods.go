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
