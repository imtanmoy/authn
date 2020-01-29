-- organizations start
CREATE TABLE organizations
(
    id         BIGSERIAL PRIMARY KEY NOT NULL,
    name       VARCHAR(100)          NOT NULL,
    created_at TIMESTAMP             NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP             NOT NULL DEFAULT NOW()
);
-- organizations end

-- users start
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
-- users end

-- invitations start
create type invite_status as enum ('pending', 'successful','canceled');
CREATE TABLE invites
(
    id              BIGSERIAL     NOT NULL,
    email           VARCHAR(255)  NOT NULL,
    token           VARCHAR(32)   NOT NULL,
    status          invite_status NOT NULL DEFAULT 'pending',
    organization_id BIGINT        NOT NULL,
    user_id         BIGINT        NULL,
    invited_by      BIGINT        NOT NULL,
    created_at      TIMESTAMP     NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP     NOT NULL DEFAULT NOW()
);

ALTER TABLE invites
    ADD CONSTRAINT pk_invites_id
        PRIMARY KEY (id);

ALTER TABLE invites
    ADD CONSTRAINT fk_invites_organizations
        FOREIGN KEY (organization_id)
            REFERENCES organizations (id);

ALTER TABLE invites
    ADD CONSTRAINT fk_invites_users
        FOREIGN KEY (user_id)
            REFERENCES users (id);

ALTER TABLE invites
    ADD CONSTRAINT uk_invites_token
        UNIQUE (token);

ALTER TABLE invites
    ADD CONSTRAINT fk_invites_invited_by_users
        FOREIGN KEY (invited_by)
            REFERENCES users (id);

ALTER TABLE invites
    ADD CONSTRAINT uk_invites_email_organization_user
        UNIQUE (email, organization_id, user_id);
-- invitations end