
DROP TABLE IF EXISTS banners;
DROP TABLE IF EXISTS banners_data CASCADE;

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


INSERT INTO banners_data (title, text, url, is_active, created_at, updated_at)
 VALUES
     ('title', 'text', 'url', true, date_trunc('seconds',current_timestamp), date_trunc('seconds',current_timestamp)),
     ('title', 'text', 'url', true, date_trunc('seconds',current_timestamp), date_trunc('seconds',current_timestamp)),
     ('title', 'text', 'url', true, date_trunc('seconds',current_timestamp), date_trunc('seconds',current_timestamp));

INSERT INTO banners (banner_id, feature_id, tag_id)
    VALUES
    (1, 1, 1),
    (1, 1, 2),
    (1, 1, 3),
    (2, 1, 4),
    (2, 1, 5),
    (2, 1, 6),
    (3, 2, 7),
    (3, 2, 8),
    (3, 2, 9);