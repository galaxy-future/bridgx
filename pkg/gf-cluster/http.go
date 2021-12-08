package gf_cluster

var (
	StatusSuccess = "success"
	StatusFailed  = "failed"
	ModuleName    = "bridgx/containers-cloud"
	Version       = ""
)

type ResponseBase struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Pager struct {
	PageNumber int `json:"page_number"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
}

func NewSuccessResponse() *ResponseBase {
	return &ResponseBase{
		Status:  StatusSuccess,
		Message: "",
	}
}

func NewFailedResponse(message string) *ResponseBase {
	return &ResponseBase{
		Status:  StatusFailed,
		Message: message,
	}
}

//PingResponse 用于测试服务是否可用
type PingResponse struct {
	*ResponseBase
	Module  string
	Version string
}

func NewPingResponse() *PingResponse {
	return &PingResponse{
		ResponseBase: NewSuccessResponse(),
		Module:       ModuleName,
		Version:      Version,
	}
}
