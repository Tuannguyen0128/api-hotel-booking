package user

type userInfo struct {
	Id          string `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	CompanyId   string `json:"companyId"`
	Role        string `json:"role"`
	Status      string `json:"status"`
	CreateBy    string `json:"createBy"`
	CreateDt    string `json:"createDt"`
	LastEditDt  string `json:"lastEditDt"`
	CompanyName string `json:"companyName"`
}
