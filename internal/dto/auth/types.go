package auth

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
}

type UserInfo struct {
	ID       uint64   `json:"id"`
	Username string   `json:"username"`
	Nickname string   `json:"nickname"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Avatar   string   `json:"avatar"`
	Status   int      `json:"status"`
	Roles    []string `json:"roles"`
}

type LoginResponse struct {
	TokenType   string   `json:"token_type"`
	AccessToken string   `json:"access_token"`
	ExpiresAt   int64    `json:"expires_at"`
	User        UserInfo `json:"user"`
}

type UserMenu struct {
	ID         uint64     `json:"id"`
	ParentID   uint64     `json:"parent_id"`
	Name       string     `json:"name"`
	Title      string     `json:"title"`
	Path       string     `json:"path"`
	Component  string     `json:"component"`
	Icon       string     `json:"icon"`
	MenuType   string     `json:"menu_type"`
	Permission string     `json:"permission"`
	Sort       int        `json:"sort"`
	Hidden     int        `json:"hidden"`
	Children   []UserMenu `json:"children"`
}
