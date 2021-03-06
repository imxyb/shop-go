package context

import (
	"github.com/dgrijalva/jwt-go"
	jwt2 "github.com/foxiswho/shop-go/consts/session/jwt"
	jwt3 "github.com/foxiswho/shop-go/module/jwt"
)

func (c *BaseContext) JwtTokenGetAdmin() map[string]interface{} {
	myMap := make(map[string]interface{})
	val := c.Get(jwt2.Jwt_Context_Key_admin)
	if val != nil {
		info := val.(*jwt.Token)
		if info != nil {
			return jwt3.GetJwtClaims(info)
		}
	}
	return myMap
}

func (c *BaseContext) JwtTokenGetUser() map[string]interface{} {
	myMap := make(map[string]interface{})
	val := c.Get(jwt2.Jwt_Context_Key_user)
	if val != nil {
		info := val.(*jwt.Token)
		if info != nil {
			return jwt3.GetJwtClaims(info)
		}
	}
	return myMap
}
