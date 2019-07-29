package smtp

import (
	"net/smtp"
)

// Mailer encapsulates the net/smtp mailer client
type Mailer struct {
	auth    smtp.Auth
	client  *smtp.Client
	host    string
	address string
}

// NewMailer returns a new mailer instance set to the given address
func NewMailer(host string, opts ...func(*Mailer)) *Mailer {
	m := &Mailer{
		host:    host,
		address: "25",
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// WithCredentials adds Plain authentication credentials to your mailer requests
func WithCredentials(username, password string) func(*Mailer) {
	return func(m *Mailer) {
		m.auth = smtp.PlainAuth("", username, password, m.host)
	}
}

// WithAddress adds a port to the mailserver request, default is 25
func WithAddress(address string) func(*Mailer) {
	return func(m *Mailer) {
		m.address = address
	}
}

// Open opens a connection to the mail server
func (m *Mailer) Open() error {
	c, err := smtp.Dial(m.host + ":" + m.address)

	if err != nil {
		return err
	}

	m.client = c

	if m.auth != nil {
		m.client.Auth(m.auth)
	}

	return nil
}

// Close closes an open smtp connection
func (m *Mailer) Close() error {
	return m.client.Close()
}

// Hello test
func (m *Mailer) Hello() error {
	return m.client.Hello("test")
}

// Send sends an email body to an address from an address
func (m *Mailer) Send(to, from, body string) error {
	if err := m.client.Mail(from); err != nil {
		return err
	}

	if err := m.client.Rcpt(to); err != nil {
		return err
	}

	wc, err := m.client.Data()

	if err != nil {
		return err
	}

	defer wc.Close()

	if _, err := wc.Write([]byte(body)); err != nil {
		return err
	}

	return nil
}
