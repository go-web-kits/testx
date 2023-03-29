package let

import (
	"errors"

	auth "github.com/go-web-kits/rbac/rbac_auth"
	"github.com/go-web-kits/testx"
	"github.com/go-web-kits/utils"
)

var AuthRequestMethod func(ad, token string) (auth.Subject, error)

func AuthPass(permission ...string) *testx.MonkeyPatches {
	if len(permission) == 0 {
		permission = append(permission, "test")
	}

	p := testx.IsExpectedToCall(utils.GCACTokenDecode).AndReturn("current_admin", "token", nil)
	p.IsExpectedToCall(AuthRequestMethod).AndReturn(auth.Subject{Roles: []auth.Role{
		{Name: "role", Permissions: []auth.Permission{{Action: permission[0]}}}}}, nil)
	return p
}

func AuthFail() *testx.MonkeyPatches {
	p := testx.IsExpectedToCall(utils.GCACTokenDecode).AndReturn("current_admin", "token", nil)
	p.IsExpectedToCall(AuthRequestMethod).AndReturn(auth.Subject{Roles: []auth.Role{}}, errors.New(""))
	return p
}
