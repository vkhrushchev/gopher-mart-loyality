create table user_balance
(
	id bigserial constraint user_balance_pk primary key,
	user_login varchar(255) not null,
	total_sum float default 0.0 not null,
	total_withdrawal_sum float default 0.0 not null
);

create unique index user_balance_user_login_uindex on user_balance (user_login);

