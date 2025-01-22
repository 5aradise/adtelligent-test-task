-- +goose Up
CREATE TABLE IF NOT EXISTS sources (
    id SERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS campaigns (
    id SERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS campaign_sources (
    source_id INT NOT NULL,
    campaign_id INT NOT NULL,
    PRIMARY KEY (source_id, campaign_id),
    FOREIGN KEY (source_id) REFERENCES sources(id) ON DELETE CASCADE,
    FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS creatives (
    id SERIAL PRIMARY KEY,
    campaign_id INT NOT NULL,
    cents_per_view INT NOT NULL,
    duration_in_sec INT NOT NULL,
    hls_playlist_path VARCHAR(256) NOT NULL,
    FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS creatives;

DROP TABLE IF EXISTS campaign_sources;

DROP TABLE IF EXISTS campaigns;

DROP TABLE IF EXISTS sources;
