package dto

// StandardResponse is a unified response structure for API responses
// The format is {errno:123,errmsg:"xxx",data:interface{},request_id="xxx"}
type StandardResponse struct {
	Errno     int         `json:"errno"`
	Errmsg    string      `json:"errmsg"`
	Data      interface{} `json:"data"`
	RequestID string      `json:"request_id"`
}
