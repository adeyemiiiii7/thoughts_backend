package thoughts

type createThoughtRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type updateThoughtRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
