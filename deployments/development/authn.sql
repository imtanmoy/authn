-- organizations start
CREATE TABLE organizations
(
    id         BIGSERIAL PRIMARY KEY NOT NULL,
    name       VARCHAR(100)          NOT NULL,
--     owner_id   BIGINT                NOT NULL,
    created_at TIMESTAMP             NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP             NOT NULL DEFAULT NOW()
);

-- ALTER TABLE organizations
--     ADD CONSTRAINT fk_organizations_owner
--         FOREIGN KEY (owner_id)
--             REFERENCES users (id);

-- organizations end

-- users start
CREATE TABLE users
(
    id         BIGSERIAL PRIMARY KEY NOT NULL,
    name       VARCHAR(100)          NOT NULL,
    email      VARCHAR(255)          NOT NULL,
    password   VARCHAR(255)          NOT NULL,
    created_at TIMESTAMP             NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP             NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP             NULL
);

ALTER TABLE users
    ADD CONSTRAINT uk_users_email
        UNIQUE (email, deleted_at);
-- users end


-- user_organization start
CREATE TABLE users_organizations
(
    id              BIGSERIAL NOT NULL,
    user_id         BIGINT    NOT NULL,
    organization_id BIGINT    NOT NULL,
    joined_at       TIMESTAMP,
    enabled         BOOLEAN   NOT NULL DEFAULT TRUE,
    created_by      BIGINT,
    updated_by      BIGINT,
    deleted_by      BIGINT
);

ALTER TABLE users_organizations
    ADD CONSTRAINT pk_users_organizations
        PRIMARY KEY (id);

ALTER TABLE users_organizations
    ADD CONSTRAINT fk_users_organizations_users
        FOREIGN KEY (user_id)
            REFERENCES users (id);

ALTER TABLE users_organizations
    ADD CONSTRAINT fk_users_organizations_organizations
        FOREIGN KEY (organization_id)
            REFERENCES organizations (id);

ALTER TABLE users_organizations
    ADD CONSTRAINT fk_users_organizations_created_by
        FOREIGN KEY (created_by)
            REFERENCES users (id);

ALTER TABLE users_organizations
    ADD CONSTRAINT fk_users_organizations_updated_by
        FOREIGN KEY (updated_by)
            REFERENCES users (id);

ALTER TABLE users_organizations
    ADD CONSTRAINT fk_users_organizations_deleted_by
        FOREIGN KEY (deleted_by)
            REFERENCES users (id);

-- user_organization end

-- invitations start
create type invitation_status as enum ('pending', 'successful','canceled');
CREATE TABLE invitations
(
    id              BIGSERIAL         NOT NULL,
    email           VARCHAR(255)      NOT NULL,
    token           VARCHAR(32)       NOT NULL,
    status          invitation_status NOT NULL DEFAULT 'pending',
    organization_id BIGINT            NOT NULL,
    user_id         BIGINT            NULL,
    invited_by      BIGINT            NOT NULL,
    accepted_at     TIMESTAMP         NULL,
    created_at      TIMESTAMP         NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP         NOT NULL DEFAULT NOW()
);

ALTER TABLE invitations
    ADD CONSTRAINT pk_invitations_id
        PRIMARY KEY (id);

ALTER TABLE invitations
    ADD CONSTRAINT fk_invitations_organizations
        FOREIGN KEY (organization_id)
            REFERENCES organizations (id);

ALTER TABLE invitations
    ADD CONSTRAINT fk_invitations_users
        FOREIGN KEY (user_id)
            REFERENCES users (id);

ALTER TABLE invitations
    ADD CONSTRAINT uk_invitations_token
        UNIQUE (token);

ALTER TABLE invitations
    ADD CONSTRAINT fk_invitations_invited_by_users
        FOREIGN KEY (invited_by)
            REFERENCES users (id);

ALTER TABLE invitations
    ADD CONSTRAINT uk_invitations_email_organization_user
        UNIQUE (email, organization_id, user_id);
-- invitations end