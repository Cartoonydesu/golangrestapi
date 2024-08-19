package response

type Success struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
}

type Fail struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
