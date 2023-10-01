package app

type UserRegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
