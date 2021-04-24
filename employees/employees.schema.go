package employees

type (
	TUserResponse struct {
		EmployeeID int    `json:"employeeId"`
		FirstName  string `json:"firstName"`
		Title      string `json:"title"`
	}
)