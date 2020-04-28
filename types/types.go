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

type ReqVerify struct {
	VerifyId   string `json:"verify_id"`
	VerifyCode string `json:"verify_code"`
}

type ReqLogin struct {
	ReqVerify
	UserName string `json:"user_name"`
	PassWord string `json:"user_pass"`
}

type ReqRegist struct {
	ReqVerify
	UserName string `json:"user_name"`
	PassWord string `json:"user_pass"`
	FsId     string `json:"father"`
}
