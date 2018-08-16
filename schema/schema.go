package schema

type JsonSchema struct {
	Type     string `schema:"type"`
	Features []struct {
		Type       string `schema:"type"`
		Properties struct {
			MAPBLKLOT string `schema:"MAPBLKLOT"`
			BLKLOT    string `schema:"BLKLOT"`
			BLOCKNUM  string `schema:"BLOCK_NUM"`
			LOTNUM    string `schema:"LOT_NUM"`
			FROMST    string `schema:"FROM_ST"`
			TOST      string `schema:"TO_ST"`
			STREET    string `schema:"STREET"`
			STTYPE    string `schema:"ST_TYPE"`
			ODDEVEN   string `schema:"ODD_EVEN"`
		} `schema:"properties"`
		Geometry struct {
			Type        string        `schema:"type"`
			Coordinates [][][]float64 `schema:"coordinates"`
		} `schema:"geometry"`
	} `schema:"features"`
}
