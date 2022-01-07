package hasura

type ActionBody struct {
	Type string      `json:"type"`
	Args interface{} `json:"args"`
}
