package trading

var _ Service = (*service)(nil)

type Service interface {
	Start() (err error)
	Stop() (err error)
}

type service struct {
	cryptoAPI cryptoAPI
}

func NewService(cryptoAPI cryptoAPI) Service {
	return &service{
		cryptoAPI: cryptoAPI,
	}
}

type cryptoAPI interface {
}

func (svc *service) Start() (err error) {
	return
}

func (svc *service) Stop() (err error) {
	return
}
