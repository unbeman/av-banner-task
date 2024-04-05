create table if not exists public.feature
(
    id bigserial
    constraint feature_pk
    primary key
);

alter table public.feature
    owner to postgres;

create table if not exists public.tag
(
    id bigserial
    constraint tag_pk
    primary key
);

alter table public.tag
    owner to postgres;

create table if not exists public.banner
(
    id         bigserial
    constraint banner_pk
    primary key,
    feature_id integer                             not null
    constraint banner_feature_id_fk
    references public.feature,
    content    json                                not null,
    is_active  boolean   default false             not null,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    updated_at timestamp default CURRENT_TIMESTAMP not null
);

alter table public.banner
    owner to postgres;

create table if not exists public.banner_tags
(
    banner_id integer not null
    constraint banner_tags_banner_id_fk
    references public.banner,
    tag_id    integer not null
    constraint banner_tags_tag_id_fk
    references public.tag,
    constraint banner_tags_pk
    primary key (banner_id, tag_id)
    );

alter table public.banner_tags
    owner to postgres;

