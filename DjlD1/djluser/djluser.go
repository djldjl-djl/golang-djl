package djluser

import (
	"djl.com/DjlD1/jwt"
	"djl.com/DjlD1/sql"
	"errors"
)

type User1 struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User1) Login() (string, error) {
	var uu []sql.User
	d := sql.Lianjie()
	result := d.Db.Where("Username = ? AND Password = ?", u.Username, u.Password).Find(&uu)
	if result.Error != nil {
		return "", errors.New("查询失败")
	}
	if len(uu) == 0 {
		return "", errors.New("username or password error")
	}
	token, err := jwt.GenerateJWT(u.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}
