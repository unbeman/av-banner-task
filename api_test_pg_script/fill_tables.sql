insert into tag (id)
select nextval('tag_id_seq')
from generate_series(1, 100)
where not exists(select * from tag);

insert into feature (id)
select nextval('feature_id_seq')
from generate_series(1, 10)
where not exists(select * from feature);