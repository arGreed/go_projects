package simpleWeb

type Note struct {
	Id          int64  `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}
