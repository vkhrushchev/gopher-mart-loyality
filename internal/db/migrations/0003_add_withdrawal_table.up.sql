create table balance_withdrawals
(
	id bigserial constraint balance_withdrawals_pk primary key,
	user_login varchar(255) not null,
	order_number varchar(16) not null,
	withdrawal float not null,
	processed_at timestamptz not null
);

create unique index balance_withdrawals_order_number_uindex on balance_withdrawals (order_number);

create unique index balance_withdrawals_user_login_uindex on balance_withdrawals (user_login);

