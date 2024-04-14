package postgres

import "errors"

var (
	InsertUpdateBannerErr = errors.New("cannot insert/update banner: " +
		"banner with this uniq data already exist or incorrect data")

	EmptyArrTagIdsErr = errors.New("no one tag_id in request")
	PermissionErr     = errors.New("users have no permission for this data")
)
