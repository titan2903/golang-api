package user

type RegisterUserInput struct { //! struct yang digunakan untuk mapping dari inputan user
	Name       string `json:"name" binding:"required"`
	Occupation string `json:"occupation" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type CheckEmailInput struct {
	Email string `json:"email" binding:"required"`
}

type FormCreateUserInput struct {
	Name       string `form:"name" binding:"required"`
	Email      string `form:"email" binding:"required"`
	Occupation string `form:"occupation" binding:"required"`
	Password   string `form:"password" binding:"required"`
	Error      error
}

type FormUpdateUserInput struct {
	ID         int
	Name       string `form:"name" binding:"required"`
	Email      string `form:"email" binding:"required"`
	Occupation string `form:"occupation" binding:"required"`
	Error      error
}