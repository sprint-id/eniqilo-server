package response

type (
	SuccessReponse struct {
		Message string `json:"message"`
		Data    any    `json:"data"`
	}

	SuccessPageReponse struct {
		Message string `json:"message"`
		Data    any    `json:"data"`
		Meta    Meta   `json:"meta"`
	}

	Meta struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
	}
)
