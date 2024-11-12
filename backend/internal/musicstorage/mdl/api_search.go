package mdl

type SearchFieldInfo struct {
	Query      string `json:"query,omitempty"`
	Refine     bool   `json:"refine,omitempty"`
}

type GroupSearchInfo struct {
	Fields     []string `json:"fields"`
	Refine     bool     `json:"refine"`
}

type TrackSearchRequest struct {
	GroupSearch GroupSearchInfo          `json:"group_search"`
	FieldSearch map[string]SearchFieldInfo `json:"field_search"`
	Genre       *string                  `json:"genre,omitempty"`
	Format      *string                  `json:"format,omitempty"`
	SortByDurationAsc bool               `json:"sort_by_duration_asc"`
}
