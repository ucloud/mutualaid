package userutil_test

import (
	"testing"

	biz "github.com/ucloud/mutualaid/backend/infra/userutil"
)

var jwt *biz.JWT

func init() {
	jwt = biz.NewJWT()
}

func TestUseCase_Auth(t *testing.T) {
	szToken, err := jwt.Auth(12193117562208258, "odHRQ4xaLwmT43Nt_IKdXLLf_xDc")
	if err != nil {
		t.Log(err)
	}
	t.Log("jwt  ", szToken)
}
