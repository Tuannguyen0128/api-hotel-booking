package pheader

type err struct {
	ErrorCode string `json:"errorCode"`
	ErrorDesc string `json:"errorDesc"`
}

type Err interface {
	CodeName() string
	Desc(isDebug bool) string
}

func getErr(e Err) *err {
	if e != nil {
		return &err{
			ErrorCode: e.CodeName(),
			ErrorDesc: e.Desc(isDebug),
		}
	}

	return nil
}
