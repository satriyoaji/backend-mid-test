package model

type JwtClaims struct {
	Aud  []string `json:"aud"`
	User struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Email      string `json:"email"`
		Superadmin bool   `json:"superadmin"`
		Roles      []struct {
			ID          int      `json:"id"`
			Name        string   `json:"name"`
			Code        string   `json:"code"`
			Permissions []string `json:"permissions"`
			App         struct {
				ID      int    `json:"id"`
				Name    string `json:"name"`
				Code    string `json:"code"`
				BaseURL string `json:"base_url"`
			} `json:"app"`
		} `json:"roles"`
	} `json:"user"`
	Exp int64 `json:"exp"`
	Iat int64 `json:"iat"`
}
