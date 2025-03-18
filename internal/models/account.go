package models

type Account struct {
	Id           int
	Username     string
	Email        string
	Bio          string
	Profile_pic  string
	Is_admin     bool
	Is_logged_in bool
}

func (a *Account) Fill_default() {
	a.Id = 0
	a.Username = ""
	a.Email = ""
	a.Is_admin = false
	a.Is_logged_in = false
	a.Bio = ""
	a.Profile_pic = ""
}
