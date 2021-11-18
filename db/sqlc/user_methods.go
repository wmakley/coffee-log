package sqlc

import "coffee-log/util"

func (user *User)BasicCredentials() util.BasicCredentials {
	return util.NewBasicCredentials(user.Username, user.Password)
}
