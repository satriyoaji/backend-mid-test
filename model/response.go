package model

type ResponseBody struct {
	Status       string      `json:"status"`
	Code         string      `json:"code"`
	Data         interface{} `json:"data"`
	Pagination   *Pagination `json:"pagination"`
	ErrorMessage *string     `json:"error_message"`
	ErrorRemark  *string     `json:"remark,omitempty"`
}

type ResponseBodyMutation struct {
	Status string      `json:"status"`
	Code   string      `json:"code"`
	Data   interface{} `json:"mutasi"`
}
