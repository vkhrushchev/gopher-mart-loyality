create table users
(
	id serial constraint users_pk primary key,
	login varchar(255) not null,
	password_hash varchar(255) not null,
	salt varchar(255) not null
);

create unique index users_login_uindex on users (login);
create unique index users_password_hash_uindex on users (password_hash);

