package response

type Response struct {
	Status     int                    `json:"status"`
	Message    string                 `json:"message"`
	Data       map[string]interface{} `json:"body"`
	TotalDocs  float32                `json:"totalDocs"`
	Limit      float32                `json:"limit"`
	TotalPages int                    `json:"totalpages"`
}
