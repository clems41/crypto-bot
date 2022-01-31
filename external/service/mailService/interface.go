package mailService

type Service interface {
	Send(request SendRequest) (err error)
}
