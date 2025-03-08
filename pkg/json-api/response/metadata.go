package json_api_response

type Metadata []MetadataItem

func NewMetadata(items ...MetadataItem) Metadata {
	return items
}

func (m Metadata) MetadataMap() map[string]any {
	metadata := map[string]any{}
	for _, item := range m {
		metadata[item.Key] = item.Value
	}

	return metadata
}

type MetadataItem struct {
	Key   string
	Value any
}

func NewMetadataItem(key string, value any) MetadataItem {
	return MetadataItem{
		Key:   key,
		Value: value,
	}
}
