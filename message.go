package outlook

import (
	"context"
	"fmt"
	"time"
)

// MessageService manages communication with microsofts graph for message resources.
type MessageService struct {
	session  *Session
	basePath string
}

// NewMessageService returns a new instance of a MessageService.
func NewMessageService(session *Session) *MessageService {
	return &MessageService{
		session:  session,
		basePath: "/messages",
	}
}

// MessageListCall struct allowing for fluent style configuration of calls to the message list endpoint.
type MessageListCall struct {
	service    *MessageService
	folderID   string
	nextLink   string
	maxResults int64
	startTime  time.Time
	endTime    time.Time
}

// List returns a MessageListCall builder struct
func (ms *MessageService) List(folderID string) *MessageListCall {
	return &MessageListCall{
		service:    ms,
		maxResults: 10,
		folderID:   folderID,
	}
}

// MaxResults sets the $top query parameter for the message list call.
func (mlc *MessageListCall) MaxResults(pageSize int64) *MessageListCall {
	mlc.maxResults = pageSize
	return mlc
}

// NextLink uses the link provided to set the $skip query parameter for the message list call.
func (mlc *MessageListCall) NextLink(link string) *MessageListCall {
	mlc.nextLink = link
	return mlc
}

// StartTime sets the startDateTime query parameter for the message list call.
func (mlc *MessageListCall) StartTime(start time.Time) *MessageListCall {
	mlc.startTime = start
	return mlc
}

// EndTime sets the endDateTime query parameter for the message list call.
func (mlc *MessageListCall) EndTime(end time.Time) *MessageListCall {
	mlc.endTime = end
	return mlc
}

// Do executes the message list call, returning the message list result.
func (mlc *MessageListCall) Do(ctx context.Context) (*MessageListResult, error) {
	params := map[string]interface{}{
		"$top":          mlc.maxResults,
		"$count":        true,
		"startDateTime": mlc.startTime.Format(DefaultQueryDateTimeFormat),
		"endDateTime":   mlc.endTime.Format(DefaultQueryDateTimeFormat),
	}
	if mlc.nextLink != "" {
		params["$skip"] = parsePageLink(mlc.nextLink, "$skip")
	}

	path := fmt.Sprintf("/mailFolders/%s%s", mlc.folderID, mlc.service.basePath)

	var result MessageListResult
	if _, err := mlc.service.session.Get(ctx, path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
