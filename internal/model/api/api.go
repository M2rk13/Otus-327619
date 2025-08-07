package api

type Request struct {
	Id     string  `json:"id"`
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

func (r *Request) GetId() string {
	return r.Id
}

type Info struct {
	Timestamp int64   `json:"timestamp"`
	Quote     float64 `json:"quote"`
}

type Response struct {
	Id      string  `json:"id"`
	Success bool    `json:"success"`
	Terms   string  `json:"terms"`
	Privacy string  `json:"privacy"`
	Query   Request `json:"query"`
	Info    Info    `json:"info"`
	Result  float64 `json:"result"`
}

func (r *Response) GetId() string {
	return r.Id
}
