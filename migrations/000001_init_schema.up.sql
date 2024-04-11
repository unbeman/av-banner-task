create table if not exists feature
(
    id bigserial
        constraint feature_pk
            primary key
);

create table if not exists tag
(
    id bigserial
        constraint tag_pk
            primary key
);

create table if not exists banner
(
    id         bigserial
        constraint banner_pk
            primary key,
    content    varchar                             not null,
    is_active  boolean   default false             not null,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    updated_at timestamp default CURRENT_TIMESTAMP not null,
    feature_id integer                             not null
        constraint banner_feature_id_fk
            references feature
);

create index if not exists banner_feature_id_index
    on banner (feature_id);

create table if not exists banner_feature_tags
(
    banner_id  integer not null
        constraint banner_feature_tags_banner_id_fk
            references banner,
    tag_id     integer not null
        constraint banner_feature_tags_feature_id_fk
            references feature,
    feature_id integer not null,
    constraint banner_feature_tags_pk
        primary key (feature_id, tag_id)
);

create or replace function trigger_time_update()
    returns trigger as $$
begin
    update banner set updated_at = now() where id = new.id;
    return new;
end;
$$ language plpgsql;

create or replace trigger banner_update
    after update on banner
    for each row
    WHEN (row(OLD.*) IS DISTINCT FROM row(NEW.*))
execute procedure trigger_time_update();