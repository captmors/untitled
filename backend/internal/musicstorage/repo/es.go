package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	. "untitled/internal/musicstorage/mdl"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticRepo struct {
	client *elasticsearch.Client
	index  string
}

func NewElasticRepo(client *elasticsearch.Client, index string) *ElasticRepo {
	return &ElasticRepo{
		client: client,
		index:  index,
	}
}

func (r *ElasticRepo) IndexTrack(track Track_MONGO) error {
	body := QueryMap{
		"title":  track.Title,
		"artist": track.Artist,
		"album":  track.Album,
		"genre":  track.Genre,
		"format": track.Format,
		"ptr":    track.Ptr,
		"uuid":   track.ID.String(),
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	res, err := r.client.Index(
		r.index,
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

// aliases
type QueryMap map[string]interface{}
type FilterMap []map[string]interface{}

func addFilter(query QueryMap, filter QueryMap) {
	query["bool"].(QueryMap)["filter"] = append(query["bool"].(QueryMap)["filter"].(FilterMap), filter)
}

// SearchTracks performs a combined search with full and partial match and filters
func (r *ElasticRepo) SearchTracks(req TrackSearchRequest) ([]QueryMap, error) {
	var buf bytes.Buffer

	query := QueryMap{
		"bool": QueryMap{
			"must":   FilterMap{},
			"filter": FilterMap{},
		},
	}
	mustClauses := FilterMap{}

	if len(req.GroupSearch.Fields) > 0 {
		for _, field := range req.GroupSearch.Fields {
			if req.GroupSearch.Refine {
				mustClauses = append(mustClauses, QueryMap{
					"fuzzy": QueryMap{
						field: QueryMap{
							"value":     req.FieldSearch[field].Query,
							"fuzziness": "AUTO",
						},
					},
				})
			} else {
				mustClauses = append(mustClauses, QueryMap{
					"match": QueryMap{
						field: req.FieldSearch[field].Query,
					},
				})
			}
		}
	}

	for field, fieldInfo := range req.FieldSearch {
		if fieldInfo.Query != "" {
			if fieldInfo.Refine {
				mustClauses = append(mustClauses, QueryMap{
					"fuzzy": QueryMap{
						field: QueryMap{
							"value":     fieldInfo.Query,
							"fuzziness": "AUTO",
						},
					},
				})
			} else {
				mustClauses = append(mustClauses, QueryMap{
					"match": QueryMap{
						field: fieldInfo.Query,
					},
				})
			}
		}
	}

	query["bool"].(QueryMap)["must"] = mustClauses

	// Добавляем фильтры по жанру и формату
	if req.Genre != nil {
		addFilter(query, QueryMap{
			"term": QueryMap{
				"genre": *req.Genre,
			},
		})
	}

	if req.Format != nil {
		addFilter(query, QueryMap{
			"term": QueryMap{
				"format": *req.Format,
			},
		})
	}

	// Сортировка по продолжительности
	sortOrder := "desc"
	if req.SortByDurationAsc {
		sortOrder = "asc"
	}

	query["sort"] = []QueryMap{
		{
			"duration": QueryMap{
				"order": sortOrder,
			},
		},
	}

	// Сериализация запроса в JSON
	if err := json.NewEncoder(&buf).Encode(QueryMap{
		"query": query,
	}); err != nil {
		return nil, fmt.Errorf("error encoding query: %w", err)
	}

	// Выполняем запрос
	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex(r.index),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting response from ES: %w", err)
	}
	defer res.Body.Close()

	var result QueryMap
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing ES response: %w", err)
	}

	// Сбор результатов
	hits := result["hits"].(QueryMap)["hits"].([]interface{})
	tracks := make([]QueryMap, 0, len(hits))
	for _, hit := range hits {
		source, _ := ExtractInterface(hit.(QueryMap), "_source")
		tracks = append(tracks, source)
	}

	return tracks, nil
}

func (r *ElasticRepo) AutocompleteTitle(prefix string) ([]string, error) {
	var buf bytes.Buffer

	query := QueryMap{
		"suggest": QueryMap{
			"track-title-suggest": QueryMap{
				"prefix": prefix,
				"completion": QueryMap{
					"field": "title_suggest",
					"fuzzy": QueryMap{"fuzziness": 2},
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error encoding suggest query: %w", err)
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex(r.index),
		r.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting response from ES: %w", err)
	}
	defer res.Body.Close()

	var result QueryMap
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing ES response: %w", err)
	}

	suggestions := result["suggest"].(QueryMap)["track-title-suggest"].([]interface{})
	titles := []string{}
	for _, suggest := range suggestions {
		options := suggest.(QueryMap)["options"].([]interface{})
		for _, option := range options {
			titles = append(titles, option.(QueryMap)["text"].(string))
		}
	}

	return titles, nil
}

// utils
func ExtractInterface(m map[string]interface{}, key string) (map[string]interface{}, bool) {
	if val, ok := m[key]; ok {
		if mapVal, ok := val.(map[string]interface{}); ok {
			return mapVal, true
		}
	}
	return nil, false
}
