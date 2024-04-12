package postgres

import "errors"

var (
	InsertBannerErr = errors.New("banner with this one of your tag_id and feature_id already exist")

	EmptyArrTagIdsErr = errors.New("no one tag_id in request")

	UpdateBannerErr = errors.New("cannot update banner: " +
		"banner with this uniq data already exist or incorrect data")
)
