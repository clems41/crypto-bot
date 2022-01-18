package tradingPlatformMock

type Position struct {
	ID           string
	DatasetIndex int
	Amount       float64
	Result       float64
}

type Dataset1 struct {
	Values []ValueDataset1 `json:"values"`
}

type ValueDataset1 struct {
	Ct int    `json:"ct"`
	Op string `json:"op"`
	Hp string `json:"hp"`
	Lp string `json:"lp"`
	Cp string `json:"cp"`
	V  string `json:"v"`
	Qv string `json:"qv"`
}
