package response

type Response struct {
	Status     int                    `json:"status"`
	Message    string                 `json:"message"`
	Data       map[string]interface{} `json:"body"`
	TotalDocs  float64                `json:"totalDocs"`
	Limit      float64                `json:"limit"`
	TotalPages int                    `json:"totalPages"`
}
