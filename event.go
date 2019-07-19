package outlook

import (
	"context"
	"fmt"
	"strings"
	"time"
)

var (
	// DefaultEventFields the default set of fields that will be requested from microsoft's graph api when fecthing events
	DefaultEventFields = strings.Join([]string{
		"id",
		"start",
		"end",
		"createdDateTime",
		"lastModifiedDateTime",
		"iCalUId",
		"subject",
		"isAllDay",
		"isCancelled",
		"isOrganizer",
		"showAs",
		"onlineMeetingUrl",
		"recurrence",
		"responseStatus",
		"location",
		"attendees",
		"organizer",
		"categories",
		"seriesMasterId",
	}, ",")
)

// EventService manages communication with microsofts graph for event resources.
type EventService struct {
	session  *Session
	basePath string
}

// NewEventService returns a new instance of a EventService.
func NewEventService(session *Session) *EventService {
	return &EventService{
		session:  session,
		basePath: "/events",
	}
}

// EventListCall struct allowing for fluent style configuration of calls to the event list endpoint.
type EventListCall struct {
	service    *EventService
	calendarID string
	nextLink   string
	maxResults int64
	startTime  time.Time
	endTime    time.Time
}

// List returns a EventListCall struct
func (es *EventService) List(calendarID string) *EventListCall {
	return &EventListCall{
		service:    es,
		maxResults: 10,
		calendarID: calendarID,
	}
}

// MaxResults sets the $top query parameter for the event list call.
func (elc *EventListCall) MaxResults(pageSize int64) *EventListCall {
	elc.maxResults = pageSize
	return elc
}

// NextLink uses the link provided to set the $skip query parameter for the event list call.
func (elc *EventListCall) NextLink(link string) *EventListCall {
	elc.nextLink = link
	return elc
}

// StartTime sets the startDateTime query parameter for the event list call.
func (elc *EventListCall) StartTime(start time.Time) *EventListCall {
	elc.startTime = start
	return elc
}

// EndTime sets the endDateTime query parameter for the event list call.
func (elc *EventListCall) EndTime(end time.Time) *EventListCall {
	elc.endTime = end
	return elc
}

// Do executes the event list call, returning the event list result.
func (elc *EventListCall) Do(ctx context.Context) (*EventListResult, error) {
	params := map[string]interface{}{
		"$top":          elc.maxResults,
		"$count":        true,
		"startDateTime": elc.startTime.Format(DefaultQueryDateTimeFormat),
		"endDateTime":   elc.endTime.Format(DefaultQueryDateTimeFormat),
		"$select":       DefaultEventFields,
	}
	if elc.nextLink != "" {
		params["$skip"] = parsePageLink(elc.nextLink, "$skip")
	}

	var path string
	if elc.calendarID == "primary" {
		path = "/calendarView"
	} else {
		path = fmt.Sprintf("/calendars/%s%s", elc.calendarID, "/calendarView")
	}

	var result EventListResult
	if _, err := elc.service.session.Get(ctx, path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// EventGetCall struct allowing for fluent style configuration of calls to the event get endpoint.
type EventGetCall struct {
	service    *EventService
	calendarID string
	eventID    string
}

// Get returns an instance of an EventGetCall with the given calendarID and eventID.
func (es *EventService) Get(calendarID string, eventID string) *EventGetCall {
	return &EventGetCall{
		service:    es,
		calendarID: calendarID,
		eventID:    eventID,
	}
}

// Do executes the http get to microsoft's graph api to get the call's event.
func (egc *EventGetCall) Do(ctx context.Context) (*Event, error) {
	var path string
	if egc.calendarID == "primary" {
		path = fmt.Sprintf("/events/%s", egc.eventID)
	} else {
		path = fmt.Sprintf("/calendars/%s%s/%s", egc.calendarID, egc.service.basePath, egc.eventID)
	}
	event := Event{}
	if _, err := egc.service.session.Get(ctx, path, nil, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

// EventCreateCall struct allowing for fluent style configuration of calls to the event create endpoint.
type EventCreateCall struct {
	service    *EventService
	calendarID string
	event      *Event
}

// Create returns an instance of en EventCreateCall with the given calendarID.
func (es *EventService) Create(calendarID string) *EventCreateCall {
	return &EventCreateCall{
		service:    es,
		calendarID: calendarID,
		event:      &Event{},
	}
}

// Event sets the event on the EventCreateCall.
func (ecc *EventCreateCall) Event(event *Event) *EventCreateCall {
	ecc.event = event
	return ecc
}

// Do executes the http post to microsoft's graph api to create the call's event.
func (ecc *EventCreateCall) Do(ctx context.Context) (*Event, error) {
	path := fmt.Sprintf("/calendars/%s%s", ecc.calendarID, ecc.service.basePath)
	if _, err := ecc.service.session.Post(ctx, path, ecc.event, ecc.event); err != nil {
		return nil, err
	}
	return ecc.event, nil
}

// EventUpdateCall struct allowing for fluent style configuration of calls to the event update endpoint.
type EventUpdateCall struct {
	service    *EventService
	calendarID string
	event      *Event
}

// Update returns an instance of an EventUpdateCall with the given calendarID.
func (es *EventService) Update(calendarID string) *EventUpdateCall {
	return &EventUpdateCall{
		service:    es,
		calendarID: calendarID,
		event:      &Event{},
	}
}

// Event sets the event on the EventUpdateCall.
func (euc *EventUpdateCall) Event(event *Event) *EventUpdateCall {
	euc.event = event
	return euc
}

// Do executes the http patch to microsoft's graph api to update the call's event.
func (euc *EventUpdateCall) Do(ctx context.Context) (*Event, error) {
	var path string
	if euc.calendarID == "primary" {
		path = fmt.Sprintf("/events/%s", euc.event.ID)
	} else {
		path = fmt.Sprintf("/calendars/%s%s/%s", euc.calendarID, euc.service.basePath, euc.event.ID)
	}
	if _, err := euc.service.session.Patch(ctx, path, euc.event, euc.event); err != nil {
		return nil, err
	}
	return euc.event, nil
}

// EventDeleteCall struct allowing for fluent style configuration of calls to the event delete endpoint.
type EventDeleteCall struct {
	service    *EventService
	calendarID string
	eventID    string
}

// Delete returns an instance of an EventDeleteCall with the given calendarID and eventID.
func (es *EventService) Delete(calendarID, eventID string) *EventDeleteCall {
	return &EventDeleteCall{
		service:    es,
		calendarID: calendarID,
		eventID:    eventID,
	}
}

// Do executes the http delete to microsoft's graph api to delete the call's event.
func (edc *EventDeleteCall) Do(ctx context.Context) error {
	var path string
	if edc.calendarID == "primary" {
		path = fmt.Sprintf("/events/%s", edc.eventID)
	} else {
		path = fmt.Sprintf("/calendars/%s%s/%s", edc.calendarID, edc.service.basePath, edc.eventID)
	}
	if _, err := edc.service.session.Delete(ctx, path, nil, nil); err != nil {
		return err
	}
	return nil
}
