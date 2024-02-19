package kbheader

type KErr struct {
	ErrorCode     string  `json:"errorCode"`
	ErrorDesc     string  `json:"errorDesc"`
	ErrorAppId    *string `json:"errorAppId"`
	ErrorAppAbbrv *string `json:"errorAppAbbrv"`
}

type Err interface {
	CodeName() string
	Desc(isDebug bool) string
}

func getKErr(err Err) KErr {
	return KErr{
		ErrorCode:     err.CodeName(),
		ErrorDesc:     err.Desc(isDebug),
		ErrorAppId:    &appId,
		ErrorAppAbbrv: &appAbbr,
	}
}
