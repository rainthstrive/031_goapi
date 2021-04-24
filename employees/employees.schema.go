package employees

type (
	TUserRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	TUserToken struct {
		Token string `json:"token"`
	}

	TUserResponse struct {
		EmployeeID int    `json:"employeeId"`
		FirstName  string `json:"firstName"`
		Title      string `json:"title"`
	}
)