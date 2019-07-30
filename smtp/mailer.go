package smtp

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
)

// Mailer encapsulates the net/smtp mailer client
type Mailer struct {
	auth    smtp.Auth
	client  *smtp.Client
	address string
	host    string
	port    string
}

// NewMailer returns a new mailer instance set to the given address
func NewMailer(address string, opts ...func(*Mailer)) *Mailer {
	h, p, err := net.SplitHostPort(address)

	if err != nil {
		log.Fatal("HostPortStringSplitError: ", err)
	}

	m := &Mailer{
		address: address,
		host:    h,
		port:    p,
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

// Open opens a connection to the mail server
func (m *Mailer) Open() error {
	c, err := smtp.Dial(m.address)

	if err != nil {
		return err
	}

	m.client = c

	if m.auth != nil {
		err := m.client.Auth(m.auth)

		if err != nil {
			return err
		}
	}

	return nil
}

// OpenTLS opens the server using a tls connection
func (m *Mailer) OpenTLS() error {
	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         m.host,
	}

	conn, err := tls.Dial("tcp", m.address, tlsconfig)

	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, m.host)

	if err != nil {
		return err
	}

	m.client = c

	if m.auth != nil {
		if err := m.client.Auth(m.auth); err != nil {
			return err
		}
	}

	return nil
}

// Close closes an open smtp connection
func (m *Mailer) Close() error {
	return m.client.Quit()
}

// Send sends an email body to an address from an address
func (m *Mailer) Send(from, to, subject, body string) error {
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

	msg := makeHeaders(map[string]string{
		"From":    from,
		"To":      to,
		"Subject": subject,
	})

	msg += "\r\n" + body

	if _, err := wc.Write([]byte(msg)); err != nil {
		return err
	}

	return nil
}

// makeHeaders creates the email headers in the proper format
func makeHeaders(headers map[string]string) string {
	var res string
	for k, v := range headers {
		res += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	return res
}
