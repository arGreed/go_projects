package feedBack

type UsrFeedBack struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Pros        string `json:"Pros"`
	Cons        string `json:"Cons"`
	Experience  int    `json:"Experience"`
}
