package tradingStrategy

type Algo interface {
	ShouldAddOrder() (ok bool, err error)
	Config() (config *Config, err error)
}
