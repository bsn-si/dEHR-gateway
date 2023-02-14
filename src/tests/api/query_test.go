package api_test

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common/fakeData"
)

func (testWrap *testWrap) queryExecSuccess(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		targetUrl := testWrap.serverURL + "/v1/query/aql"
		request, err := http.NewRequest(http.MethodGet, targetUrl, nil)
		if err != nil {
			t.Error(err)
			return
		}

		user := testData.users[0]
		request.URL.Query().Add("ehr_id", user.ehrID)

		query := url.QueryEscape(`SELECT e/ehr_id/value AS ID FROM EHR e [ehr_id/value=$ehr_id]`)
		request.URL.Query().Add("q", query)

		request.URL.Query().Add("fetch", "10")
		request.URL.Query().Add("offset", "0")
		request.URL.Query().Add("ehr_id", user.ehrID)

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected success, received status: %d", response.StatusCode)
		}

		_, err = io.ReadAll(response.Body)
		if err != nil {
			// TODO compare with real result
			t.Errorf("Expected body content")
		}
	}
}

func (testWrap *testWrap) queryExecPostSuccess(testData *TestData) func(t *testing.T) {
	// TODO should be realize after AQL inserts will done
	return func(t *testing.T) {
		if len(testData.users) == 0 || testData.users[0].ehrID == "" {
			t.Fatal("Created EHR required")
		}

		user := testData.users[0]

		url := testWrap.serverURL + "/v1/query/aql"

		request, err := http.NewRequest(http.MethodPost, url, queryExecPostCreateBodyRequest(user.ehrID))
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected success, received status: %d", response.StatusCode)
		}
	}
}

func (testWrap *testWrap) queryExecPostFail(testData *TestData) func(t *testing.T) {
	return func(t *testing.T) {
		if len(testData.users) == 0 {
			t.Fatal("Test user required")
		}

		user := testData.users[0]

		url := testWrap.serverURL + "/v1/query/aql"

		request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte("111qqqEEE")))
		if err != nil {
			t.Error(err)
			return
		}

		request.Header.Set("Content-type", "application/json")
		request.Header.Set("AuthUserId", user.id)
		request.Header.Set("Authorization", "Bearer "+user.accessToken)

		response, err := testWrap.httpClient.Do(request)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected fail, received status: %d", response.StatusCode)
		}
	}
}

func queryExecPostCreateBodyRequest(ehrID string) *bytes.Reader {
	req := fakeData.QueryExecRequest(ehrID)
	return bytes.NewReader(req)
}
