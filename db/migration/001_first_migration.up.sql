CREATE TABLE documents (
    url                 TEXT    NOT NULL PRIMARY KEY,
    pub_date            BIGINT  NOT NULL,
    fetch_time          BIGINT  NOT NULL,
    text                TEXT    NOT NULL,
    first_fetch_time    BIGINT  NOT NULL
);
