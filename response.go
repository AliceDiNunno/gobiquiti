package gobiquiti

type MetaResponse struct {
	Rc    string `json:"rc"`
	Count int    `json:"count"`
}

type Response struct {
	Meta MetaResponse  `json:"meta"`
	Data []interface{} `json:"data"`
}
