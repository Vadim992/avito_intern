package internal

import (
	"encoding/json"
	"github.com/Vadim992/avito/internal/dto"
	"github.com/Vadim992/avito/internal/mws"
	"github.com/Vadim992/avito/internal/req"
	"github.com/Vadim992/avito/internal/storage"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getJSONResp(data interface{}) []byte {
	res, _ := json.Marshal(data)

	return res
}

func TestGetUsersBanner(t *testing.T) {

	tableTest := []struct {
		name           string
		path           string
		header         string
		headerVal      string
		expectedStatus int
		expectedData   []byte
	}{
		{
			name:           "Unauthorized",
			path:           "/user_banner?tag_id=1&feature_id=1",
			header:         "token",
			headerVal:      "",
			expectedStatus: http.StatusUnauthorized,
			expectedData:   []byte(""),
		},
		{
			name:           "Forbidden",
			path:           "/user_banner?tag_id=1&feature_id=1",
			header:         "token",
			headerVal:      "forbidden_token",
			expectedStatus: http.StatusForbidden,
			expectedData:   []byte(""),
		},
		{
			name:           "Forbidden_user(is_active=false)_map",
			path:           "/user_banner?tag_id=2&feature_id=2",
			header:         "token",
			headerVal:      "user_token",
			expectedStatus: http.StatusForbidden,
			expectedData:   []byte(""),
		},
		{
			name:           "Forbidden_user(is_active=false)_db",
			path:           "/user_banner?tag_id=4&feature_id=4",
			header:         "token",
			headerVal:      "user_token",
			expectedStatus: http.StatusForbidden,
			expectedData:   []byte(""),
		},
		{
			name:           "NotFound_user",
			path:           "/user_banner?tag_id=1&feature_id=1000",
			header:         "token",
			headerVal:      "user_token",
			expectedStatus: http.StatusNotFound,
			expectedData:   []byte(""),
		},
		{
			name:           "NotFound_admin",
			path:           "/user_banner?tag_id=1000&feature_id=1",
			header:         "token",
			headerVal:      "admin_token",
			expectedStatus: http.StatusNotFound,
			expectedData:   []byte(""),
		},

		{
			name:           "BadReq_no_tag_user",
			path:           "/user_banner?feature_id=3",
			header:         "token",
			headerVal:      "user_token",
			expectedStatus: http.StatusBadRequest,
			expectedData:   getJSONResp(ErrorStruct{req.QueryDataErr.Error()}),
		},
		{
			name:           "BadReq_no_tag_admin",
			path:           "/user_banner?feature_id=3",
			header:         "token",
			headerVal:      "admin_token",
			expectedStatus: http.StatusBadRequest,
			expectedData:   getJSONResp(ErrorStruct{req.QueryDataErr.Error()}),
		},
		{
			name:           "BadReq_no_Feature_user",
			path:           "/user_banner?tag_id=2",
			header:         "token",
			headerVal:      "user_token",
			expectedStatus: http.StatusBadRequest,
			expectedData:   getJSONResp(ErrorStruct{req.QueryDataErr.Error()}),
		},
		{
			name:           "BadReq_no_feature_admin",
			path:           "/user_banner?tag_id=2",
			header:         "token",
			headerVal:      "user_token",
			expectedStatus: http.StatusBadRequest,
			expectedData:   getJSONResp(ErrorStruct{req.QueryDataErr.Error()}),
		},
		{
			name:           "OK_user_map",
			path:           "/user_banner?tag_id=2&feature_id=5",
			header:         "token",
			headerVal:      "user_token",
			expectedStatus: http.StatusOK,
			expectedData: getJSONResp(dto.NewBannerContent("title5", "text5",
				"url5")),
		},
		{
			name:           "OK_admin_map",
			path:           "/user_banner?tag_id=3&feature_id=1",
			header:         "token",
			headerVal:      "admin_token",
			expectedStatus: http.StatusOK,
			expectedData: getJSONResp(dto.NewBannerContent("title1", "text1",
				"url1")),
		},
		{
			name:           "OK_user_db",
			path:           "/user_banner?tag_id=3&feature_id=3&use_last_revision=true",
			header:         "token",
			headerVal:      "user_token",
			expectedStatus: http.StatusOK,
			expectedData: getJSONResp(dto.NewBannerContent("title3", "text3",
				"url3")),
		},
		{
			name:           "OK_admin_db",
			path:           "/user_banner?tag_id=3&feature_id=1&use_last_revision=true",
			header:         "token",
			headerVal:      "admin_token",
			expectedStatus: http.StatusOK,
			expectedData: getJSONResp(dto.NewBannerContent("title1", "text1",
				"url1")),
		},
	}

	pathRoles := []int{mws.ADMIN, mws.USER}
	tokenMap := map[string]int{
		"admin_token": mws.ADMIN,
		"user_token":  mws.USER,
	}

	mockDb := NewMockDB()
	inMemory := storage.NewInMemoryStorage()
	mockDb.FillInMemory(inMemory)

	app := NewApp(&mockDb, inMemory, tokenMap)

	handler := mws.Auth(pathRoles, app.tokenMap, app.GetUserBanner)

	for _, test := range tableTest {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, test.path, nil)
			r.Header.Set(test.header, test.headerVal)

			w := httptest.NewRecorder()

			handler(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			require.Equal(t, test.expectedStatus, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, test.expectedData, body)

		})
	}
}
