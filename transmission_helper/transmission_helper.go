package transmission_helper

type TransmissionRequest struct {
	Method    string                 `json:"method"`
	Arguments map[string]interface{} `json:"arguments"`
}

type TransmissionResponse struct {
	Result    string                 `json:"result"`
	Arguments map[string]interface{} `json:"arguments"`
}
