package models

type PumpFunIdl struct {
	Version      string `json:"version"`
	Name         string `json:"name"`
	Instructions []struct {
		Name     string   `json:"name"`
		Docs     []string `json:"docs"`
		Accounts []struct {
			Name     string `json:"name"`
			IsMut    bool   `json:"isMut"`
			IsSigner bool   `json:"isSigner"`
		} `json:"accounts"`
		Args []any `json:"args"`
	} `json:"instructions"`

	Accounts []struct {
		Name string `json:"name"`
		Type struct {
			Kind   string `json:"kind"`
			Fields []struct {
				Name string `json:"name"`
				Type string `json:"type"`
			} `json:"fields"`
		} `json:"type"`
	} `json:"accounts"`

	Events []struct {
		Name   string `json:"name"`
		Fields []struct {
			Name  string `json:"name"`
			Type  string `json:"type"`
			Index bool   `json:"index"`
		} `json:"fields"`
	} `json:"events"`

	Errors []struct {
		Code int    `json:"code"`
		Name string `json:"name"`
		Msg  string `json:"msg"`
	} `json:"errors"`

	Metadata struct {
		Address string `json:"address"`
	} `json:"metadata"`
}
