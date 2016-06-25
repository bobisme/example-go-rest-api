CREATE TABLE states (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    abbrev CHAR(2),

    created_at DATETIME,
    modified_at DATETIME,
    -- let's use deleted_at to render a row "deleted", but keep it around for
    -- statistical use
    deleted_at DATETIME NULL
);

CREATE TABLE cities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    state_id INTEGER,
    lat REAL,
    lon REAL,
    -- sqlite doesn't know trigonometry, so let's help it out some
    lat_sin REAL,
    lat_cos REAL,
    lon_sin REAL,
    lon_cos REAL,

    created_at DATETIME,
    modified_at DATETIME,
    deleted_at DATETIME NULL,

    FOREIGN KEY(state_id) REFERENCES states(id)
);

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT,
    last_name TEXT,

    email TEXT,
    password_hash TEXT,

    created_at DATETIME,
    modified_at DATETIME,
    deleted_at DATETIME NULL
);

CREATE TABLE visits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    city_id INTEGER,
    -- possible use: check in by coordinates
    lat REAL,
    lon REAL,
    lat_sin REAL,
    lat_cos REAL,
    lon_sin REAL,
    lon_cos REAL,
    -- "city" if by city, "coords" if by coordinates
    visit_method TEXT,

    created_at DATETIME,
    -- modified won't be used, but it's here for consitency and future use
    modified_at DATETIME,
    deleted_at DATETIME NULL
);
