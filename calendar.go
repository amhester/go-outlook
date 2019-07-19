package outlook

import (
	"context"
	"fmt"
)

// CalendarService manages communication with microsofts graph for calendar resources.
type CalendarService struct {
	session  *Session
	basePath string
}

// NewCalendarService returns a new instance of a CalendarService.
func NewCalendarService(session *Session) *CalendarService {
	return &CalendarService{
		session:  session,
		basePath: "/calendars",
	}
}

// CalendarListCall struct allowing for fluent style configuration of calls to the calendar list endpoint.
type CalendarListCall struct {
	service    *CalendarService
	nextLink   string
	maxResults int64
}

// List returns a CalendarListCall builder struct
func (cs *CalendarService) List() *CalendarListCall {
	return &CalendarListCall{
		service:    cs,
		maxResults: 10,
	}
}

// MaxResults sets the $top query parameter for the calendar list call.
func (clc *CalendarListCall) MaxResults(pageSize int64) *CalendarListCall {
	clc.maxResults = pageSize
	return clc
}

// NextLink uses the link provided to set the $skip query parameter for the calendar list call.
func (clc *CalendarListCall) NextLink(link string) *CalendarListCall {
	clc.nextLink = link
	return clc
}

// Do executes the calendar list call, returning the calendar list result.
func (clc *CalendarListCall) Do(ctx context.Context) (*CalendarListResult, error) {
	params := map[string]interface{}{
		"$top":   clc.maxResults,
		"$count": true,
	}
	if clc.nextLink != "" {
		params["$skip"] = parsePageLink(clc.nextLink, "$skip")
	}

	var result CalendarListResult
	if _, err := clc.service.session.Get(ctx, clc.service.basePath, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CalendarGetCall struct allowing for fluent style configuration of calls to the calendar get endpoint.
type CalendarGetCall struct {
	service    *CalendarService
	calendarID string
}

// Get returns an instance of a CalendarGetCall with the given calendarID.
func (cs *CalendarService) Get(calendarID string) *CalendarGetCall {
	return &CalendarGetCall{
		service:    cs,
		calendarID: calendarID,
	}
}

// Do executes the http get request to microsoft's graph api to get the call's calendar.
func (cgc *CalendarGetCall) Do(ctx context.Context) (*Calendar, error) {
	path := fmt.Sprintf("%s/%s", cgc.service.basePath, cgc.calendarID)
	calendar := Calendar{}
	if _, err := cgc.service.session.Get(ctx, path, nil, &calendar); err != nil {
		return nil, err
	}
	return &calendar, nil
}

// CalendarCreateCall struct allowing for fluent style configuration of calls to the calendar create endpoint.
type CalendarCreateCall struct {
	service  *CalendarService
	calendar *Calendar
}

// Create returns an instance of a CalendarCreateCall.
func (cs *CalendarService) Create() *CalendarCreateCall {
	return &CalendarCreateCall{
		service:  cs,
		calendar: &Calendar{},
	}
}

// Calendar sets the calendar data to be created on the call.
func (ccc *CalendarCreateCall) Calendar(calendar *Calendar) *CalendarCreateCall {
	ccc.calendar = calendar
	return ccc
}

// Do executes the http post request to microsoft's graph api to create the call's calendar.
func (ccc *CalendarCreateCall) Do(ctx context.Context) (*Calendar, error) {
	if _, err := ccc.service.session.Post(ctx, ccc.service.basePath, ccc.calendar, ccc.calendar); err != nil {
		return nil, err
	}
	return ccc.calendar, nil
}

// CalendarUpdateCall struct allowing for fluent style configuration of calls to the calendar update endpoint.
type CalendarUpdateCall struct {
	service    *CalendarService
	calendarID string
	calendar   *Calendar
}

// Update returns an instance of a CalendarUpdateCall with the given calendarID.
func (cs *CalendarService) Update(calendarID string) *CalendarUpdateCall {
	return &CalendarUpdateCall{
		service:    cs,
		calendarID: calendarID,
		calendar:   &Calendar{},
	}
}

// Calendar sets the calendar for the call.
func (cuc *CalendarUpdateCall) Calendar(calendar *Calendar) *CalendarUpdateCall {
	cuc.calendar = calendar
	return cuc
}

// Do executes the http patch request to microsoft's graph api to update the call's calendar.
func (cuc *CalendarUpdateCall) Do(ctx context.Context) (*Calendar, error) {
	path := fmt.Sprintf("%s/%s", cuc.service.basePath, cuc.calendarID)
	if _, err := cuc.service.session.Patch(ctx, path, cuc.calendar, cuc.calendar); err != nil {
		return nil, err
	}
	return cuc.calendar, nil
}

// CalendarDeleteCall struct allowing for fluent style configuration of calls to the calendar delete endpoint.
type CalendarDeleteCall struct {
	service    *CalendarService
	calendarID string
}

// Delete returns an instance of a CalendarDeleteCall with the given calendarID.
func (cs *CalendarService) Delete(calendarID string) *CalendarDeleteCall {
	return &CalendarDeleteCall{
		service:    cs,
		calendarID: calendarID,
	}
}

// Do executes the http delete request to microsoft's graph api to delete the call's calendar.
func (cdc *CalendarDeleteCall) Do(ctx context.Context) error {
	path := fmt.Sprintf("%s/%s", cdc.service.basePath, cdc.calendarID)
	if _, err := cdc.service.session.Delete(ctx, path, nil, nil); err != nil {
		return err
	}
	return nil
}
