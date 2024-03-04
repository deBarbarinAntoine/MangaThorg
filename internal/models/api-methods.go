package models

import (
	"net/url"
	"strconv"
	"time"
)

func (r MangaRequest) Params() MangaRequestParam {
	return MangaRequestParam{
		Order:          "order[" + r.OrderType + "]",
		IncludedTags:   "includedTags",
		ExcludedTags:   "excludedTags",
		Title:          "title",
		Author:         "author",
		AuthorOrArtist: "authorOrArtist",
		Status:         "status[]",
		Public:         "publicationDemographic[]",
		ContentRating:  "contentRating[]",
		Limit:          "limit",
		Offset:         "offset",
	}
}

func (r MangaRequest) ToQuery(q *url.Values) {
	params := r.Params()
	(*q)[params.Order] = []string{r.OrderValue}
	if r.IncludedTags != nil {
		(*q)[params.IncludedTags] = r.IncludedTags
	}
	if r.ExcludedTags != nil {
		(*q)[params.ExcludedTags] = r.ExcludedTags
	}
	if r.Title != "" {
		(*q)[params.Title] = []string{r.Title}
	}
	if r.Author != "" {
		(*q)[params.Author] = []string{r.Author}
	}
	if r.AuthorOrArtist != "" {
		(*q)[params.AuthorOrArtist] = []string{r.AuthorOrArtist}
	}
	if r.Status != nil {
		(*q)[params.Status] = r.Status
	}
	if r.Public != nil {
		(*q)[params.Public] = r.Public
	}
	(*q)[params.ContentRating] = []string{"safe"}
	(*q)[params.Limit] = []string{strconv.Itoa(r.Limit)}
	(*q)[params.Offset] = []string{strconv.Itoa(r.Offset)}
}

func (res ApiManga) SingleCacheData(order string) SingleCacheData {
	var data SingleCacheData
	if len(res.Data) == 1 {
		data.Id = res.Data[0].Id
	}
	data.UpdatedTime = time.Now()
	data.Order = order
	data.Data = res
	//log.Printf("mangathorg/internal/methods/ApiManga.SingleCacheData() data: %#v\n", data) // testing
	return data
}
