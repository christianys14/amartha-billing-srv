-- migrate:up
alter table loan
    add version int(2) not null comment 'versioning';

alter table loan
    add updated_at timestamp not null on update current_timestamp comment 'updated_at of the transaction';

alter table loan
    add created_at timestamp not null comment 'created_at of the transaction';

create index idx_created_at
    on loan (created_at);

-- migrate:down
alter table loan drop column version;
alter table loan drop column updated_at;
alter table loan drop column created_at;
alter table loan drop index idx_created_at;