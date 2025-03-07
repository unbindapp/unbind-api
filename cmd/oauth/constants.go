package oauth

import (
	"time"
)

const ACCESS_TOKEN_EXP = 1 * time.Minute
const REFRESH_TOKEN_EXP = 24 * time.Hour * 14 // 2 weeks
var ALLOWED_SCOPES = []string{"openid", "profile", "email", "offline_access"}
