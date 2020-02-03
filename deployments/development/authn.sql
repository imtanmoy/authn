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
CREATE TABLE user_organization
(
    user_id         BIGINT NOT NULL,
    organization_id BIGINT NOT NULL
);

ALTER TABLE user_organization
    ADD CONSTRAINT fk_user_organization_users
        FOREIGN KEY (user_id)
            REFERENCES users (id);

ALTER TABLE user_organization
    ADD CONSTRAINT fk_user_organization_organizations
        FOREIGN KEY (organization_id)
            REFERENCES organizations (id);

-- user_organization end

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
    accepted_at     TIMESTAMP     NULL,
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