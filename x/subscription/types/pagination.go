package types

const DefaultLimit = 100

func InitPagintionRequestDefaults(pageRequest *PageRequest) *PageRequest {
	if pageRequest == nil {
		pageRequest = &PageRequest{}
	}

	pageRequestCopy := *pageRequest
	if pageRequestCopy.Limit == 0 {
		pageRequestCopy.Limit = DefaultLimit
	}

	return &pageRequestCopy
}
