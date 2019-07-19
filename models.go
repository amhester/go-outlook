package outlook

// User microsoft user object
type User struct {
	FirstName string `json:"givenName,omitempty"`
	LastName  string `json:"surName,omitempty"`
	Name      string `json:"displayName,omitempty"`
	ID        string `json:"id,omitempty"`
	Email     string `json:"userPrincipalName,omitempty"`
	JobTitle  string `json:"jobTitle,omitempty"`
}

// RefreshTokenRequest microsoft token request object
type RefreshTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
	GrantType    string `json:"grant_type"`
}

// RefreshTokenResponse microsoft token response object
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
}

// FolderListResult struct representing a response from the outlook mailFolders endpoint
type FolderListResult struct {
	Context  string    `json:"@odata.context,omitempty"`
	NextLink string    `json:"@odata.nextLink,omitempty"`
	Total    int64     `json:"@odata.count,omitempty"`
	Value    []*Folder `json:"value,omitempty"`
}

// Folder struct representing an outlook calendar object
type Folder struct {
	ID               string `json:"id,omitempty"`
	DisplayName      string `json:"displayName,omitempty"`
	ParentFolderID   string `json:"parentFolderId,omitempty"`
	ChildFolderCount int    `json:"childFolderCount,omitempty"`
	UnreadItemCount  int    `json:"unreadItemCount,omitempty"`
	TotalItemCount   int    `json:"totalItemCount,omitempty"`
}

// MessageListResult struct representing a response from the outlook messages endpoint
type MessageListResult struct {
	Context  string     `json:"@odata.context,omitempty"`
	NextLink string     `json:"@odata.nextLink,omitempty"`
	Total    int64      `json:"@odata.count,omitempty"`
	Value    []*Message `json:"value,omitempty"`
}

// Message microsoft message object
// TODO: Add all fields from outlook
type Message struct {
	ID             string       `json:"id,omitempty"`
	MessageID      string       `json:"internetMessageId,omitempty"`
	CreatedOn      string       `json:"createdDateTime,omitempty"`
	ReceivedOn     string       `json:"receivedDateTime,omitempty"`
	SentOn         string       `json:"sentDateTime,omitempty"`
	Subject        string       `json:"subject,omitempty"`
	BodyPreview    string       `json:"bodyPreview,omitempty"`
	Importance     string       `json:"importance,omitempty"`
	ConversationID string       `json:"conversationId,omitempty"`
	IsRead         bool         `json:"isread,omitempty"`
	Body           *MessageBody `json:"body,omitempty"`
	Sender         *Recipient   `json:"sender,omitempty"`
	From           *Recipient   `json:"from,omitempty"`
	To             []*Recipient `json:"toRecipients,omitempty"`
	CC             []*Recipient `json:"ccRecipients,omitempty"`
	BCC            []*Recipient `json:"bccRecipients,omitempty"`
	ReplyTo        []*Recipient `json:"replyTo,omitempty"`
}

// BodyContentType enum
const (
	BodyContentTypeText = "TEXT"
	BodyContentTypeHTML = "HTML"
)

// MessageBody microsoft body content object
type MessageBody struct {
	ContentType string `json:"contentType,omitempty"`
	Content     string `json:"content,omitempty"`
}

// Recipient microsoft message recipient object
type Recipient struct {
	EmailAddress *EmailAddress `json:"emailAddress,omitempty"`
}

// EmailAddress microsoft message email object
type EmailAddress struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
}

// CalendarListResult you can tell by the way it is
type CalendarListResult struct {
	Context  string      `json:"@odata.context,omitempty"`
	NextLink string      `json:"@odata.nextLink,omitempty"`
	Total    int64       `json:"@odata.count,omitempty"`
	Value    []*Calendar `json:"value,omitempty"`
}

// Calendar outlook calendar object
type Calendar struct {
	ID                  string        `json:"id,omitempty"`
	Name                string        `json:"name,omitempty"`
	Color               string        `json:"color,omitempty"`
	CanShare            bool          `json:"canShare,omitempty"`
	CanViewPrivateItems bool          `json:"canViewPrivateItems,omitempty"`
	CanEdit             bool          `json:"canEdit,omitempty"`
	Owner               *EmailAddress `json:"owner,omitempty"`
}

// EventListResult you can tell by the way it is
type EventListResult struct {
	Context  string   `json:"@odata.context,omitempty"`
	NextLink string   `json:"@odata.nextLink,omitempty"`
	Total    int64    `json:"@odata.count,omitempty"`
	Value    []*Event `json:"value,omitempty"`
}

// Essentially, enums of possible values for outlook calendar events. Would like to change to iota+custom json serializer/deserializer.
const (
	// EventShowAs
	EventShowAsFree      = "free"
	EventShowAsTentative = "tentative"
	EventShowAsBusy      = "busy"
	EventShowAsOOF       = "oof"
	EventShowAsElsewhere = "workingElsewhere"
	EventShowAsUnknown   = "unknown"

	// EventType
	EventTypeSingleInstance = "singleInstance"
	EventTypeOccurrence     = "occurrence"
	EventTypeException      = "exception"
	EventTypeSeriesMaster   = "seriesMaster"

	// EventSensitivity
	EventSensitivityNormal       = "normal"
	EventSensitivityPersonal     = "personal"
	EventSensitivityPrivate      = "private"
	EventSensitivityConfidential = "confidential"

	// EventImportance
	EventImportanceLow    = "low"
	EventImportanceNormal = "normal"
	EventImportanceHigh   = "high"
)

// Event microsoft event object
// TODO: Add all fields from outlook
type Event struct {
	ID                         string               `json:"id,omitempty"`
	CreatedOn                  string               `json:"createdDateTime,omitempty"`
	UpdatedOn                  string               `json:"lastModifiedDateTime,omitempty"`
	ICalUID                    string               `json:"iCalUId,omitempty"`
	Categories                 []string             `json:"categories,omitempty"`
	Subject                    string               `json:"subject,omitempty"`
	BodyPreview                string               `json:"bodyPreview,omitempty"`
	Importance                 string               `json:"importance,omitempty"`
	IsOrganizer                bool                 `json:"isOrganizer,omitempty"`
	IsCancelled                bool                 `json:"isCancelled,omitempty"`
	SeriesID                   string               `json:"seriesMasterId,omitempty"`
	Type                       string               `json:"type,omitempty"`
	Body                       *MessageBody         `json:"body,omitempty"`
	Start                      *DateTimeTimeZone    `json:"start,omitempty"`
	OriginalStart              string               `json:"originalStart,omitempty"` // YYYY-mm-ddT00:00:00Z
	OriginalStartTimezone      string               `json:"originalStartTimeZone,omitempty"`
	End                        *DateTimeTimeZone    `json:"end,omitempty"`
	AllDay                     bool                 `json:"isAllDay,omitempty"`
	Location                   *Location            `json:"location,omitempty"`
	Locations                  []*Location          `json:"locations,omitempty"`
	Attendees                  []*Attendee          `json:"attendees,omitempty"`
	Organizer                  *Recipient           `json:"organizer,omitempty"`
	ResponseStatus             *ResponseStatus      `json:"responseStatus,omitempty"`
	WebLink                    string               `json:"webLink,omitempty"`
	OnlineMeetingURL           string               `json:"onlineMeetingUrl,omitempty"`
	ShowAs                     string               `json:"showAs,omitempty"`
	Sensitivity                string               `json:"sensitivity,omitempty"`
	ResponseRequested          bool                 `json:"responseRequested,omitempty"`
	ReminderMinutesBeforeStart int                  `json:"reminderMinutesBeforeStart,omitempty"`
	Recurrence                 *PatternedRecurrence `json:"recurrence,omitempty"`
	ReminderOn                 bool                 `json:"isReminderOn,omitempty"`
	HasAttachments             bool                 `json:"hasAttachments,omitempty"`
}

// ResponseStatus something
type ResponseStatus struct {
	Response string `json:"response,omitempty"`
	Time     string `json:"time,omitempty"`
}

// DateTimeTimeZone microsoft event datetime-timezone object
type DateTimeTimeZone struct {
	DateTime string `json:"dateTime,omitempty"`
	Timezone string `json:"timeZone,omitempty"`
}

// Location microsoft event location object
type Location struct {
	DisplayName string   `json:"displayName,omitempty"`
	Address     *Address `json:"address,omitempty"`
	Type        string   `json:"locationType,omitempty"`
}

// Address microsoft event location address object
type Address struct {
	Street  string `json:"street,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Country string `json:"countryOrRegion,omitempty"`
	Postal  string `json:"postalCode,omitempty"`
}

// Attendee microsoft event attendee object
type Attendee struct {
	Type         string          `json:"type,omitempty"`
	Status       *ResponseStatus `json:"status,omitempty"`
	EmailAddress *EmailAddress   `json:"emailAddress,omitempty"`
}

// PatternedRecurrence microsoft recurrence definition.
type PatternedRecurrence struct {
	Pattern *RecurrencePattern `json:"pattern,omitempty"`
	Range   *RecurrenceRange   `json:"range,omitempty"`
}

// RecurrencePattern enums for recurrence definition
const (
	// RecurrencePatternType
	RecurrencePatternTypeDaily           = "daily"
	RecurrencePatternTypeWeekly          = "weekly"
	RecurrencePatternTypeAbsoluteMonthly = "absoluteMonthly"
	RecurrencePatternTypeRelativeMonthly = "relativeMonthly"
	RecurrencePatternTypeAbsoluteYearly  = "absoluteYearly"
	RecurrencePatternTypeRelativeYearly  = "relativeYearly"

	// RecurrencePatternIndex
	RecurrencePatternIndexFirst  = "first"
	RecurrencePatternIndexSecond = "second"
	RecurrencePatternIndexThird  = "third"
	RecurrencePatternIndexFourth = "fourth"
	RecurrencePatternIndexLast   = "last"
)

// RecurrencePattern microsoft event recurrence pattern definition.
type RecurrencePattern struct {
	DayOfMonth     int      `json:"dayOfMonth,omitempty"`
	DaysOfWeek     []string `json:"daysOfWeek,omitempty"`
	FirstDayOfWeek string   `json:"firstDayOfWeek,omitempty"`
	Index          string   `json:"index,omitempty"`
	Interval       int      `json:"interval,omitempty"`
	Month          int      `json:"month,omitempty"`
	Type           string   `json:"type,omitempty"`
}

// RecurrenceRangeType enum
const (
	RecurrenceRangeTypeEndDate  = "endDate"
	RecurrenceRangeTypeNoEnd    = "noEnd"
	RecurrenceRangeTypeNumbered = "numbered"
)

// RecurrenceRange microsoft event recurrence range definition.
type RecurrenceRange struct {
	EndDate             string `json:"endDate,omitempty"` // FORMAT - YYYY-mm-dd
	NumberOfOccurrences int    `json:"numberOfOccurrences,omitempty"`
	RecurrenceTimezone  string `json:"recurrenceTimeZone,omitempty"`
	StartDate           string `json:"startDate,omitempty"`
	Type                string `json:"type,omitempty"`
}
