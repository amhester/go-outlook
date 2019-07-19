package outlook

import "context"

// FolderService manages communication with microsofts graph for folder resources.
type FolderService struct {
	session  *Session
	basePath string
}

// NewFolderService returns a new instance of a FolderService.
func NewFolderService(session *Session) *FolderService {
	return &FolderService{
		session:  session,
		basePath: "/mailFolders",
	}
}

// FolderListCall struct allowing for fluent style configuration of calls to the mailFolder list endpoint.
type FolderListCall struct {
	service    *FolderService
	nextLink   string
	maxResults int64
}

// List returns a FolderListCall builder struct
func (fs *FolderService) List() *FolderListCall {
	return &FolderListCall{
		service:    fs,
		maxResults: 10,
	}
}

// MaxResults sets the $top query parameter for the folder list call.
func (flc *FolderListCall) MaxResults(pageSize int64) *FolderListCall {
	flc.maxResults = pageSize
	return flc
}

// NextLink uses the link provided to set the $skip query parameter for the folder list call.
func (flc *FolderListCall) NextLink(link string) *FolderListCall {
	flc.nextLink = link
	return flc
}

// Do executes the folder list call, returning the folder list result.
func (flc *FolderListCall) Do(ctx context.Context) (*FolderListResult, error) {
	params := map[string]interface{}{
		"$top":   flc.maxResults,
		"$count": true,
	}
	if flc.nextLink != "" {
		params["$skip"] = parsePageLink(flc.nextLink, "$skip")
	}

	var result FolderListResult
	if _, err := flc.service.session.Get(ctx, flc.service.basePath, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
