DROP TABLE IF EXISTS public.tg_users;

CREATE TABLE IF NOT EXISTS public.tg_users
(
    id          UUID PRIMARY KEY DEFAULT get_rendom_uuid(),
    tg_uid      INT NOT NULL UNIQUE,
    login       VARCHAR(128) UNIQUE,
    first_name  VARCHAR(128),
    last_name   VARCHAR(128)
);