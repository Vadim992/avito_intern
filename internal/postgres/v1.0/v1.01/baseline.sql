
DROP TABLE IF EXISTS banners;
DROP TABLE IF EXISTS banners_data CASCADE;

DROP INDEX IF EXISTS banners_data_pkey;
DROP INDEX IF EXISTS banners_feature_id_tag_id_key;
DROP INDEX IF EXISTS banner_id_index;
DROP INDEX IF EXISTS title_index;
DROP INDEX IF EXISTS text_index;
DROP INDEX IF EXISTS url_index;
DROP INDEX IF EXISTS is_active_index;


CREATE TABLE banners_data (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    text TEXT NOT NULL,
    url TEXT NOT NULL,
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE banners (
    banner_id INTEGER NOT NULL,
    feature_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,

    UNIQUE(feature_id, tag_id),
    FOREIGN KEY (banner_id) REFERENCES banners_data (id) ON DELETE CASCADE
);

CREATE INDEX banners_id_index ON banners (banner_id);
CREATE INDEX title_index ON banners_data (title);
CREATE INDEX text_index ON banners_data (text);
CREATE  INDEX url_index ON banners_data (url);
CREATE  INDEX is_active_index ON banners_data (is_active);


-- INSERT INTO banners_data (title, text, url, is_active, created_at, updated_at)
--  VALUES
--      ('title', 'text', 'url', true, date_trunc('seconds',current_timestamp), date_trunc('seconds',current_timestamp)),
--      ('title', 'text', 'url', true, date_trunc('seconds',current_timestamp), date_trunc('seconds',current_timestamp)),
--      ('title', 'text', 'url', true, date_trunc('seconds',current_timestamp), date_trunc('seconds',current_timestamp));
--
-- INSERT INTO banners (banner_id, feature_id, tag_id)
--     VALUES
--     (1, 1, 1),
--     (1, 1, 2),
--     (1, 1, 3),
--     (2, 1, 4),
--     (2, 1, 5),
--     (2, 1, 6),
--     (3, 2, 7),
--     (3, 2, 8),
--     (3, 2, 9);