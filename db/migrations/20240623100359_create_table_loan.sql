-- migrate:up
create table loan
(
    id       bigint auto_increment,
    status   varchar(10)    not null COMMENT 'PENDING (not yet paid), PAID',
    user_id  varchar(50)    not null COMMENT 'user id of the customer',
    due_date date           not null COMMENT 'due date of customer should pay',
    amount   decimal(20, 2) not null COMMENT 'amount of customer should pay',
    constraint pk_id primary key (id)
);

-- migrate:down
drop table loan;
