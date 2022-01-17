package tradingPlatformMock

type Position struct {
	ID           string
	DatasetIndex int
	Amount       float32
}

type Dataset struct {
	Values []ValueDataset `json:"values"`
}

type ValueDataset struct {
	Ct int    `json:"ct"`
	Op string `json:"op"`
	Hp string `json:"hp"`
	Lp string `json:"lp"`
	Cp string `json:"cp"`
	V  string `json:"v"`
	Qv string `json:"qv"`
}
