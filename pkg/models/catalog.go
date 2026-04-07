package models

type CatalogType struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	TypeName string `json:"type_name"`
}

type CatalogTypeResponse struct {
	CatalogTypes []CatalogType `json:"catalog_types"`
}

type CatalogEntry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CatalogEntriesResponse struct {
	CatalogEntries []CatalogEntry `json:"catalog_entries"`
	PaginationMeta struct {
		PageSize         int `json:"page_size"`
		TotalRecordCount int `json:"total_record_count"`
	} `json:"pagination_meta"`
}
