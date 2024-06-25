-- migrate:up
create index idx_due_date
    on loan (due_date);

create index idx_status
    on loan (status);

create index idx_user_id
    on loan (user_id);

-- migrate:down
drop index idx_due_date on loan;
drop index idx_status on loan;
drop index idx_user_id on loan;