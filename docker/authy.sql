CREATE TABLE organizations
(
    id         BIGSERIAL PRIMARY KEY NOT NULL,
    name       VARCHAR(255)          NOT NULL,
    created_at TIMESTAMP             NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP             NOT NULL DEFAULT NOW()
);