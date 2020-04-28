package types

const (
	STATUS_WAITING = 0
	STATUS_SUCCESS = 1
	STATUS_FAILED  = -1
)

type RepPost struct {
	Name string `json:"my_name"`
}

type RspPost struct {
	Name string `json:"your_name"`
}

type ReqLogin struct {
	UserName string `json:"user_name"`
	PassWord string `json:"user_pass"`
}

type RspLogin struct {
}

type ReqRegist struct {
	UserName string `json:"user_name"`
	PassWord string `json:"user_pass"`
}
