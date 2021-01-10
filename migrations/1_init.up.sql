create table goods
(
    id          bigserial primary key,
    merchant_id bigint,
    offer_id    bigint,
    name        varchar(80),
    price       integer,
    quantity    integer,
    available   bool
);

create unique index UQE_merchant_id_offer_id
    on goods (merchant_id, offer_id);

create table tasks
(
    id           bigserial primary key not null,
    created_at   timestamp             not null,
    status       varchar(30)           not null,
    merchant_id  bigint                not null,
    url          varchar(1024)         not null,
    total_rows   integer               not null,
    updated_rows integer               not null,
    saved_rows   integer               not null,
    deleted_rows integer               not null,
    bad_rows     integer               not null
);

create unique index UQE_merchant_id_url
    on tasks (merchant_id, url);
