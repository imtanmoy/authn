CREATE TABLE organizations
(
    id         BIGSERIAL PRIMARY KEY NOT NULL,
    name       VARCHAR(100)          NOT NULL,
    created_at TIMESTAMP             NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP             NOT NULL DEFAULT NOW()
);

CREATE TABLE users
(
    id              BIGSERIAL PRIMARY KEY NOT NULL,
    name            VARCHAR(100)          NOT NULL,
    designation     VARCHAR(70),
    email           VARCHAR(255)          NOT NULL,
    password        VARCHAR(255)          NOT NULL,
    enabled         BOOLEAN               NOT NULL DEFAULT TRUE,
    organization_id BIGINT                NOT NULL,
    created_by      BIGINT,
    updated_by      BIGINT,
    deleted_by      BIGINT,
    joined_at       TIMESTAMP             NULL,
    created_at      TIMESTAMP             NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP             NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMP             NULL

);

ALTER TABLE users
    ADD CONSTRAINT fk_users_organizations
        FOREIGN KEY (organization_id)
            REFERENCES organizations (id);
ALTER TABLE users
    ADD CONSTRAINT uk_users_email
        UNIQUE (email, deleted_at);
ALTER TABLE users
    ADD CONSTRAINT fk_users_created_by
        FOREIGN KEY (created_by)
            REFERENCES users (id);
ALTER TABLE users
    ADD CONSTRAINT fk_users_updated_by
        FOREIGN KEY (created_by)
            REFERENCES users (id);
ALTER TABLE users
    ADD CONSTRAINT fk_users_deleted_by
        FOREIGN KEY (created_by)
            REFERENCES users (id);