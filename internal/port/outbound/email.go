package outbound_port

//go:generate mockgen -source=email.go -destination=./../../../tests/mocks/port/mock_email.go
type EmailPort interface {
	SendEmail(to, subject, body string) error
}