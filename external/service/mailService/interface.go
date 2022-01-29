package mailService

type Service interface {
	Send(request SendRequest) (response SendResponse, err error)
}
