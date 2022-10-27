package greet

type Greet struct {
	Name  string `json:"name"`
	Pass  string `json:"pass"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}
type Resp struct {
	Message string `json:"message"`
}
