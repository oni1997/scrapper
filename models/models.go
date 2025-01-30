package models

type Product struct {
    ID        string  `json:"id"`
    Name      string  `json:"name"`
    Price     string  `json:"price"`
    ImageURL  string  `json:"image_url"`
    Promotion *string `json:"promotion,omitempty"`
}

type ProductResponse struct {
    Success   bool      `json:"success"`
    Message   string    `json:"message,omitempty"`
    Products  []Product `json:"products,omitempty"`
    Total     int       `json:"total"`
    SearchTerm string   `json:"search_term"`
}

type SearchRequest struct {
    SearchTerm string `json:"searchTerm"`
}