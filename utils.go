package outlook

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func checkResponse(res *http.Response) error {
	status := res.StatusCode
	if status >= 200 && status < 300 {
		return nil
	}

	statusErr := &ErrStatusCode{Code: status}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(data) > 0 {
		statusErr.Message = string(data)
	}
	if statusErr.Code == 429 {
		rawRetrySecs := res.Header.Get("Retry-After")
		if rawRetrySecs != "" {
			retrySecs, _ := strconv.ParseInt(rawRetrySecs, 10, 64)
			statusErr.SuggestedRetryDuration = time.Duration(retrySecs) * time.Second
		}
	}

	return statusErr
}

func createQueryString(params map[string]interface{}) string {
	query := url.Values{}
	for key, val := range params {
		query.Set(key, fmt.Sprintf("%v", val))
	}
	finalQuery := query.Encode()
	if finalQuery == "" {
		return ""
	}
	return fmt.Sprintf("?%s", finalQuery)
}

func parsePageLink(link, key string) string {
	parsed, _ := url.Parse(link)
	q := parsed.Query()
	return q.Get(key)
}
