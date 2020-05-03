package types

const (
	STATUS_WAITING = 0
	STATUS_SUCCESS = 1
	STATUS_FAILED  = -1
)

type RspVerify struct {
	VerifyId string `json:"verify_id"`
	Imag     string `json:"imag"`
}
