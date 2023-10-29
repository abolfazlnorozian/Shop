package entities

type State struct {
	Name   string `json:"name"`
	Cities []City `json:"Cities"`
}

type City struct {
	Name string `json:"name"`
}
