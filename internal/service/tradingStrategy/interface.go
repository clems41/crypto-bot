package tradingStrategy

type Algo interface {
	// ShouldAddOrder determines if order should be open. If returned order is nil, order should not be open
	ShouldAddOrder(form ShouldAddOrderForm) (view ShouldAddOrderView, err error)
}
