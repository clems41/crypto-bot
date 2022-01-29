package mailService

type SendRequest struct {
	From    string
	Subject string
	Body    string
}

type SendResponse struct {
}
