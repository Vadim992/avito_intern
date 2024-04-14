package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Vadim992/avito/internal/dto"
	"github.com/Vadim992/avito/internal/mws"
	"github.com/Vadim992/avito/internal/storage"
	"github.com/lib/pq"
	"strings"
)

const (
	numTags     = 1000
	numFeatures = 1000
	numBanners  = 1000
)

func (db *DB) FillDb() error {
	for i := 1; i <= numBanners; i++ {
		values := fmt.Sprintf(`('title%d','text%d', 'url%d', true,date_trunc('seconds',current_timestamp),
date_trunc('seconds',current_timestamp))`, i, i, i)

		stmt := fmt.Sprintf(`INSERT INTO banners_data (title, text, url, is_active, created_at, updated_at)
VALUES %s`, values)
		_, err := db.DB.Exec(stmt)

		if err != nil {
			return err
		}
	}

	for i := 1; i <= numFeatures; i++ {
		for j := 1; j <= numTags; j++ {
			stmt := `INSERT INTO banners (banner_id, tag_id, feature_id) VALUES ($1, $2, $1);`
			_, err := db.DB.Exec(stmt, i, j)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (db *DB) FillStorage(inMemory storage.Storage) error {
	stmt := `SELECT id,title, text, url, is_active, feature_id, ARRAY_AGG(tag_id) AS tag_ids 
FROM banners 
JOIN banners_data
 ON banners_data.id = banners.banner_id
GROUP BY id,title, text, url, is_active, feature_id;`

	rows, err := db.DB.Query(stmt)

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {

		var banner dto.PostPatchBanner
		var content dto.BannerContent
		var bannerId int

		err := rows.Scan(&bannerId, &content.Title, &content.Text,
			&content.Url, &banner.IsActive, &banner.FeatureId,
			pq.Array(&banner.TagIds))

		if err != nil {
			return err
		}

		banner.Content = &content

		err = inMemory.Save(bannerId, &banner)

		if err != nil {
			return err
		}
	}

	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

func (db *DB) GetUserBanner(tagId, featureId, role int) (*dto.GetBanner, error) {
	stmt := `SELECT id, title, text, url, is_active FROM banners 
    JOIN banners_data 
        ON banners.banner_id = banners_data.id 
    WHERE tag_id = $1 AND feature_id = $2 %s LIMIT 1;`

	row := db.DB.QueryRow(stmt, tagId, featureId)

	var content dto.BannerContent
	var banner dto.GetBanner
	err := row.Scan(&banner.BannerId, &content.Title, &content.Text, &content.Url, &banner.IsActive)

	if err != nil {
		return nil, err
	}

	if !*banner.IsActive && role == mws.USER {
		return nil, PermissionErr
	}

	banner.Content = &content

	return &banner, nil
}

func (db *DB) GetBanners(ctx context.Context, whereStmt, limitOffsetStmt string) ([]dto.GetBanner, error) {
	stmt := fmt.Sprintf(`SELECT banners_data.*, feature_id, ARRAY_AGG(tag_id) AS tag_ids FROM banners_data
	JOIN banners
	 ON banners_data.id = banners.banner_id %s
	GROUP BY id, feature_id, title, text, url, is_active, created_at, updated_at %s;`, whereStmt, limitOffsetStmt)

	rows, err := db.DB.Query(stmt)

	if err != nil {
		return nil, err
	}

	result := make([]dto.GetBanner, 0)
	for rows.Next() {
		var banner dto.GetBanner
		var content dto.BannerContent

		err := rows.Scan(&banner.BannerId, &content.Title, &content.Text,
			&content.Url, &banner.IsActive, &banner.CreatedAt,
			&banner.UpdatedAt, &banner.FeatureId, pq.Array(&banner.TagIds))

		if err != nil {
			return nil, err
		}

		banner.Content = &content

		result = append(result, banner)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, sql.ErrNoRows
	}

	return result, nil
}

func (db *DB) FillArr() error {
	var b strings.Builder
	b.WriteString("{")
	for i := 1; i <= numTags; i++ {
		if i != numTags {
			b.WriteString(fmt.Sprintf("%d,", i))
			continue
		}
		b.WriteString(fmt.Sprintf("%d}", i))
	}

	stmt := `UPDATE banners_data SET tag_ids = $1;`

	_, err := db.DB.Exec(stmt, b.String())
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) InsertBanner(banner dto.PostPatchBanner) (int, error) {
	stmt := `INSERT INTO banners_data (title, text, url,is_active, created_at, updated_at)
	VALUES($1, $2, $3, $4, date_trunc('seconds',current_timestamp),
	    date_trunc('seconds',current_timestamp))
	RETURNING id;`

	var bannerId int
	content := banner.Content
	row := db.DB.QueryRow(stmt, content.Title, content.Text, content.Url,
		banner.IsActive)

	if err := row.Scan(&bannerId); err != nil {
		return 0, err
	}

	stmt = `INSERT INTO banners (banner_id, feature_id, tag_id)
	VALUES($1, $2, $3);`

	for _, tagId := range banner.TagIds {
		_, err := db.DB.Exec(stmt, bannerId, banner.FeatureId, tagId)

		if err != nil {
			return 0, err
		}
	}

	return bannerId, nil
}

func (db *DB) checkId(tx *sql.Tx, id int) error {
	stmt := fmt.Sprintf(`SELECT id FROM banners_data
                WHERE id = $1 LIMIT 1;`)

	var bannerId int
	err := tx.QueryRow(stmt, id).Scan(&bannerId)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) getFeatureAndTagIds(tx *sql.Tx, id int) (int, []int64, error) {
	stmt := `SELECT feature_id, ARRAY_AGG(tag_id) AS tag_ids FROM banners
                WHERE banner_id = $1 GROUP BY feature_id;`

	var featureId int
	var tagIds = make([]int64, 0)
	err := tx.QueryRow(stmt, id).Scan(&featureId, pq.Array(&tagIds))

	if err != nil {
		return 0, nil, err
	}

	return featureId, tagIds, nil
}

func (db *DB) UpdateBannerId(id int, banner dto.PostPatchBanner) (*storage.UpdateDeleteFromDB, error) {
	tx, err := db.DB.Begin()

	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = db.checkId(tx, id)

	if err != nil {
		return nil, err
	}

	reqStr, err := validateBannersDataPatch(banner)

	if err != nil {
		return nil, err
	}

	oldFeature, oldTags, err := db.getFeatureAndTagIds(tx, id)

	if err != nil {
		return nil, err
	}

	switch {
	case banner.FeatureId != nil && banner.TagIds != nil:
		err := db.updateOldTagsAndFeatureOrOldTags(tx, id, oldFeature, *banner.FeatureId, banner.TagIds)

		if err != nil {
			return nil, err
		}
	case banner.FeatureId != nil:
		err := db.updateOldFeature(tx, id, *banner.FeatureId)

		if err != nil {
			return nil, err
		}
	case banner.TagIds != nil:
		err := db.updateOldTagsAndFeatureOrOldTags(tx, id, oldFeature, oldFeature, banner.TagIds)

		if err != nil {
			return nil, err
		}
	}

	var resBanner *dto.PostPatchBanner

	if oldFeature != 0 || oldTags != nil || reqStr != "" {
		resBanner, err = db.updateContent(tx, id, reqStr)
		if err != nil {
			return nil, err
		}

	}

	storageStruct := storage.NewUpdateDeleteFromDB(id, &oldFeature, oldTags,
		nil, resBanner)

	err = tx.Commit()

	return storageStruct, err
}

func (db *DB) updateOldTagsAndFeatureOrOldTags(tx *sql.Tx, bannerId, oldFeatureId, featureId int, tagIds []int64) error {
	stmt := `DELETE FROM banners WHERE banner_id = $1 AND feature_id = $2;`

	_, err := tx.Exec(stmt, bannerId, oldFeatureId)

	if err != nil {
		return err
	}

	stmt = `INSERT INTO banners (banner_id, tag_id, feature_id) VALUES ($1, $2, $3)`

	for _, tagId := range tagIds {
		_, err := tx.Exec(stmt, bannerId, tagId, featureId)

		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) updateOldFeature(tx *sql.Tx, bannerId, featureId int) error {
	stmt := `UPDATE banners SET feature_id = $1 WHERE banner_id = $2;`

	_, err := tx.Exec(stmt, featureId, bannerId)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) updateContent(tx *sql.Tx, bannerId int, setStmt string) (*dto.PostPatchBanner, error) {
	stmt := fmt.Sprintf(`UPDATE banners_data SET %s updated_at=date_trunc('seconds',current_timestamp)
               WHERE id = $1 
RETURNING title, text, url, is_active;`, setStmt)

	var result dto.PostPatchBanner
	var content dto.BannerContent

	err := tx.QueryRow(stmt, bannerId).Scan(&content.Title,
		&content.Text, &content.Url, &result.IsActive)

	if err != nil {
		return nil, err
	}

	result.Content = &content

	return &result, nil
}

func (db *DB) DeleteBanner(id int) (*storage.UpdateDeleteFromDB, error) {
	tx, err := db.DB.Begin()

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	stmt := `SELECT id FROM banners_data WHERE id = $1;`

	var bannerId int
	err = tx.QueryRow(stmt, id).Scan(&bannerId)

	if err != nil {
		return nil, err
	}

	stmt = `SELECT feature_id, ARRAY_AGG(tag_id) AS tag_ids FROM banners WHERE banner_id = $1
GROUP BY feature_id;`

	var featureId int
	tagIds := make([]int64, 0)

	err = tx.QueryRow(stmt, id).Scan(&featureId, pq.Array(&tagIds))

	if err != nil {
		return nil, err
	}

	stmt = `DELETE FROM banners_data WHERE id = $1;`

	_, err = tx.Exec(stmt, id)

	if err != nil {
		return nil, err
	}

	storageStruct := storage.NewUpdateDeleteFromDB(id, &featureId, tagIds, nil, nil)

	err = tx.Commit()

	return storageStruct, err
}
