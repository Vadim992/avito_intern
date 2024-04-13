package internal

import (
	"github.com/Vadim992/avito/internal/dto"
	"github.com/Vadim992/avito/internal/mws"
	"github.com/Vadim992/avito/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

type reqBodyUserBanner struct {
	t         *testing.T
	tableTest []testUserBanner
}

type paramsUserBanner struct {
	tagId     int
	featureId int
	role      int
}

type testUserBanner struct {
	name           string
	data           paramsUserBanner
	expectedStatus int
	expectedData   dto.BannerContent
}

func TestGetUsersBanner(t *testing.T) {

	tableTest := []struct {
		name           string
		path           string
		header         string
		headerVal      string
		expectedStatus int
		expectedData   *dto.BannerContent
		expectedErr    *ErrorStruct
	}{
		{
			name:           "Unauthorized",
			path:           "/user_banner?tag_id=1&feature_id=1",
			header:         "token",
			headerVal:      "",
			expectedStatus: http.StatusUnauthorized,
			expectedData:   nil,
			expectedErr:    nil,
		},
		{
			name:           "Forbidden",
			path:           "/user_banner?tag_id=1&feature_id=1",
			header:         "token",
			headerVal:      "forbidden_token",
			expectedStatus: http.StatusForbidden,
			expectedData:   nil,
			expectedErr:    nil,
		},
		{
			name:           "NotFound_user",
			path:           "/user_banner?tag_id=1&feature_id=1000",
			header:         "token",
			headerVal:      "user_token",
			expectedStatus: http.StatusNotFound,
			expectedData:   nil,
			expectedErr:    nil,
		},
		{
			name:           "NotFound_admin",
			path:           "/user_banner?tag_id=1000&feature_id=1",
			header:         "token",
			headerVal:      "user_token",
			expectedStatus: http.StatusNotFound,
			expectedData:   nil,
			expectedErr:    nil,
		},
	}
	//{
	//name: "success_user",SS
	//	path: "/"
	//expectedStatus: http.StatusOK,
	//	expectedData:   dto.NewBannerContent("title3", "text3", "url3"),
	//},
	//	{
	//		name: "success_admin",
	//		data: paramsUserBanner{
	//			tagId:     1,
	//			featureId: 2,
	//			role:      1,
	//		},
	//		expectedStatus: http.StatusOK,
	//		expectedData:   dto.NewBannerContent("title2", "text2", "url2"),
	//	},
	//	{
	//		name: "forbidden_user",
	//		data: paramsUserBanner{
	//			tagId:     1,
	//			featureId: 2,
	//			role:      2,
	//		},
	//		expectedStatus: http.StatusForbidden,
	//		expectedData:   dto.NewBannerContent("title2", "text2", "url2"),
	//	},
	//	{
	//		name: "not_found_user",
	//		data: paramsUserBanner{
	//			tagId:     1,
	//			featureId: 20,
	//			role:      2,
	//		},
	//		expectedStatus: http.StatusNotFound,
	//	},
	//}

	pathRoles := []int{mws.ADMIN, mws.USER}
	tokenMap := map[string]int{
		"admin_toke": mws.ADMIN,
		"user_token": mws.USER,
	}

	mockDb := NewMockDB()
	inMemory := storage.NewStorage()

	app := NewApp(&mockDb, inMemory, tokenMap)

	handler := mws.Auth(pathRoles, app.tokenMap, app.GetUserBanner)

	for _, test := range tableTest {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, test.path, nil)
			r.Header.Set(test.header, test.headerVal)

			w := httptest.NewRecorder()

			handler(w, r)

			resp := w.Result()

			if resp.StatusCode != test.expectedStatus {
				t.Fatalf("wrong")
			}

		})
	}
}
