package error

type Response struct {
	Message string `json:"error"`
}

type Code struct {
	Message string
}
