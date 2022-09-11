package es

type CommonEsFields struct {
	ESId string `json:"_id"`
	Hits hits   `json:"hits"`
}

type hits struct {
	Total    total    `json:"total"`
	MaxScore int64    `json:"max_score"`
	Hits     []hitsEs `json:"hits"`
}

type hitsEs struct {
	Index  string      `json:"_index"`
	Type   string      `json:"_type"`
	Id     string      `json:"_id"`
	Score  int64       `json:"_score"`
	Source interface{} `json:"_source"`
}

type total struct {
	Value    int64  `json:"value"`
	Relation string `json:"relation"`
}
