package tradingAPI

var _ Api = (*api)(nil)

type Api interface {
}

type api struct {
}

func New() Api {
	return &api{}
}
