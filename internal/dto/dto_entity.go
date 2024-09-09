package dto

type UserEntity struct {
	Id           int64  `db:"id"`
	Login        string `db:"login"`
	PasswordHash string `db:"password_hash"`
	Salt         string `db:"salt"`
}
