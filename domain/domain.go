package domain

const NotFoundMessage = "Item not found"

type Pokemon struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Power int    `json:"power"`
}

type ApiResponse struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous []string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}
