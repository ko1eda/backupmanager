package smtp

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"

	"github.com/pkg/errors"
)

// Mailer encapsulates the net/smtp mailer client
type Mailer struct {
	auth    smtp.Auth
	from    string
	address string
	host    string
	port    string
}

// NewMailer returns a new mailer instance set to the given valid server address
// ex mail.mysite.com:25
// from sets the default sender field for the mailer
func NewMailer(address, from string, opts ...func(*Mailer)) *Mailer {
	h, p, err := net.SplitHostPort(address)

	if err != nil {
		log.Fatal("HostPortStringSplitError: ", err)
	}

	m := &Mailer{
		from:    from,
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

// DialAndSend Dials a connection to the mail server, sends a message and then closes the connection
func (m *Mailer) DialAndSend(from, to, subject, body string) error {
	client, err := m.openTLS()

	if err != nil {
		return errors.Wrap(err, "OpenTlsError: Could not open smtp client connection")
	}

	defer client.close()

	if err := client.send(from, to, subject, body); err != nil {
		return errors.Wrap(err, "DialAndSendError: Could not send email message")
	}

	return nil
}

// OpenTLS opens the server using a tls connection
func (m *Mailer) openTLS() (*smtpClient, error) {
	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         m.host,
	}

	conn, err := tls.Dial("tcp", m.address, tlsconfig)

	if err != nil {
		return nil, err
	}

	c, err := smtp.NewClient(conn, m.host)

	if err != nil {
		return nil, err
	}

	if m.auth != nil {
		if err := c.Auth(m.auth); err != nil {
			return nil, err
		}
	}

	return &smtpClient{c, m}, nil
}

// smtpClient wraps go's smtp client
type smtpClient struct {
	*smtp.Client
	m *Mailer
}

// Send sends an email body to an address from an address
// This will use the default mailer from address if no address is specified
func (c *smtpClient) send(from, to, subject, body string) error {
	if err := c.Mail(from); err != nil {
		return err
	}

	if err := c.Rcpt(to); err != nil {
		return err
	}

	wc, err := c.Data()

	if err != nil {
		return err
	}

	defer wc.Close()

	headers := map[string]string{
		"From":    c.m.from,
		"To":      to,
		"Subject": subject,
	}

	if from != "" {
		headers["From"] = from
	}

	msg := makeHeaders(headers)

	msg += "\r\n" + body

	if _, err := wc.Write([]byte(msg)); err != nil {
		return err
	}

	return nil
}

// Close closes an open smtp connection
func (c *smtpClient) close() error {
	return c.Quit()
}

// makeHeaders creates the email headers in the proper format
func makeHeaders(headers map[string]string) string {
	var res string
	for k, v := range headers {
		res += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	return res
}
