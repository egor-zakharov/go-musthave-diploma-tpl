package models

type User struct {
	UserID   string `db:"ID"`
	Login    string `db:"LOGIN"`
	Password string `db:"PASSWORD"`
}

func (u *User) IsValidLogin() bool {
	return u.Login != ""
}

func (u *User) IsValidPass() bool {
	return u.Password != ""
}
