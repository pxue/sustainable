package bettercotton

// Better Cotton initiative members
type Member struct {
	Name     string `json:"name"`
	Since    string `json:"since"`
	Category string `json:"category"`
	Country  string `json:"country"`
	Website  string `json:"website"`
}
