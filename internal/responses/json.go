package responses

type ErrResp struct {
	ErrCode int    `json:"ErrCode"`
	ErrMsg  string `json:"ErrMsg"`
}

func ERROR(errCode int, errMsg string) ErrResp {
	return ErrResp{
		ErrCode: errCode,
		ErrMsg:  errMsg,
	}
}
