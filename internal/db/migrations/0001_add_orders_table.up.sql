create table orders
(
	id bigserial constraint orders_pk primary key,
	user_login varchar(255) not null
		constraint orders_users_login_fk references users (login),
	number varchar(16),
	accrual float,
	status varchar(50) not null,
	uploaded_at timestamptz
);

