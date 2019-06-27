// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package coding_carefree

type LoginInVo struct {
	Usr string `json:"usr"`
	Pwd string `json:"pwd"`
}

type LoginOutVo struct {
	Code  int    `json:"code"`
	Token string `json:"token"`
}

type RegisterInVo struct {
	Usr  string `json:"usr"`
	Pwd  string `json:"pwd"`
	Code string `json:"code"`
}

type SendMailInVo struct {
	Type int    `json:"type"`
	Usr  string `json:"usr"`
}

type User struct {
	ID string `json:"id"`
}

type Void struct {
	Code int `json:"code"`
}