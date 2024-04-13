package req

import (
	"encoding/json"
	"github.com/Vadim992/avito/internal/postgres"
	"github.com/Vadim992/avito/pkg/logger"
	"github.com/lib/pq"
	"net/http"
	"net/url"
	"strconv"
)

const (
	tagIdQuery           = "tag_id"
	featureIdQuery       = "feature_id"
	offsetQuery          = "offset"
	limitQuery           = "limit"
	useLastRevisionQuery = "use_last_revision"
)

var (
	mapStructFieldToDB = map[string]string{
		"BannerId":  "banner_id",
		"TagIds":    "tag_ids",
		"FeatureId": "feature_id",
		"IsActive":  "is_active",
	}
)

func hasRequiredQuery(q url.Values, param string) (string, error) {
	if !q.Has(param) {
		return "", QueryDataErr
	}

	value := q.Get(param)
	if value == "" {
		return "", QueryDataErr
	}

	return value, nil
}

func hasQuery(q url.Values, param string) (string, bool) {
	if !q.Has(param) {
		return "", false
	}

	value := q.Get(param)
	if value == "" {
		return "", false
	}

	return value, true
}

func convertStrToInt(str string) (int, error) {
	num, err := strconv.Atoi(str)

	if err != nil {
		logger.ErrLog.Printf("cannot convert string to int: %s", err)

		return 0, QueryDataErr
	}

	if num < 0 {
		return 0, QueryDataErr
	}

	return num, nil
}

func hasRequiredIntQuery(q url.Values, queryName string) (int, error) {
	numStr, err := hasRequiredQuery(q, queryName)

	if err != nil {
		return 0, err
	}

	return convertStrToInt(numStr)
}

func hasIntQuery(q url.Values, queryName string) (int, bool, error) {
	numStr, ok := hasQuery(q, queryName)

	if !ok {
		return 0, false, nil
	}

	num, err := convertStrToInt(numStr)

	return num, ok, err
}

func convertStrToBool(str string) (bool, error) {
	boolVal, err := strconv.ParseBool(str)

	if err != nil {
		logger.ErrLog.Printf("cannot convert string to bool: %s", err)

		return false, QueryDataErr
	}

	return boolVal, nil
}

func hasBoolQuery(q url.Values, queryName string) (bool, error) {
	boolStr, ok := hasQuery(q, queryName)

	if !ok {
		return false, nil
	}

	boolVal, err := convertStrToBool(boolStr)

	return boolVal, err
}

func sentDataToFront(data interface{}, w http.ResponseWriter, code int) error {
	res, err := json.Marshal(data)

	if err != nil {
		return err
	}

	w.WriteHeader(code)

	w.Write(res)

	return nil
}

func ValidateDriverErrors(err error) error {
	pqErr, ok := err.(*pq.Error)

	if !ok {
		return err
	}

	switch pqErr.Code {
	case pq.ErrorCode("23505"):
		return postgres.InsertUpdateBannerErr
	}

	return pqErr
}

func deleteDuplicateFromTagIds(tadIds []int64) []int64 {
	m := make(map[int64]struct{}, len(tadIds))
	newTagIds := make([]int64, 0, len(tadIds))

	for _, val := range tadIds {
		if _, ok := m[val]; !ok {
			newTagIds = append(newTagIds, val)

			m[val] = struct{}{}
		}
	}

	return newTagIds
}
