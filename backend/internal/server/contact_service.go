package server

import (
	"errors"
	"fmt"
	"net/mail"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"devfolio/backend/internal/domain"
	"devfolio/backend/internal/store"
)

var (
	ErrContactInvalid   = errors.New("contact submission is invalid")
	ErrContactRateLimit = errors.New("contact submission rate limited")
	ErrContactStore     = errors.New("contact submission could not be saved")
	ErrContactNotify    = errors.New("contact notification could not be delivered")
)

type ContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

type ContactNotifier interface {
	Notify(domain.ContactSubmission) error
}

type noopContactNotifier struct{}

func (noopContactNotifier) Notify(domain.ContactSubmission) error { return nil }

type SMTPContactNotifier struct {
	host     string
	port     int
	username string
	password string
	from     string
	to       string
}

func NewSMTPContactNotifierFromEnv() ContactNotifier {
	host := strings.TrimSpace(os.Getenv("DEVFOLIO_SMTP_HOST"))
	to := strings.TrimSpace(os.Getenv("DEVFOLIO_CONTACT_TO"))
	from := strings.TrimSpace(os.Getenv("DEVFOLIO_SMTP_FROM"))
	if host == "" || to == "" || from == "" {
		return noopContactNotifier{}
	}
	port, _ := strconv.Atoi(strings.TrimSpace(os.Getenv("DEVFOLIO_SMTP_PORT")))
	if port == 0 {
		port = 587
	}
	return SMTPContactNotifier{
		host:     host,
		port:     port,
		username: strings.TrimSpace(os.Getenv("DEVFOLIO_SMTP_USER")),
		password: strings.TrimSpace(os.Getenv("DEVFOLIO_SMTP_PASS")),
		from:     from,
		to:       to,
	}
}

func (n SMTPContactNotifier) Notify(submission domain.ContactSubmission) error {
	addr := fmt.Sprintf("%s:%d", n.host, n.port)
	var message strings.Builder
	message.WriteString(fmt.Sprintf("To: %s\r\n", n.to))
	message.WriteString(fmt.Sprintf("From: %s\r\n", n.from))
	message.WriteString(fmt.Sprintf("Subject: DevFolio contact from %s\r\n", submission.Name))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
	message.WriteString(fmt.Sprintf("Name: %s\nEmail: %s\nTime: %s\n\n%s\n", submission.Name, submission.Email, submission.CreatedAt.Format(time.RFC3339), submission.Message))

	var auth smtp.Auth
	if n.username != "" && n.password != "" {
		auth = smtp.PlainAuth("", n.username, n.password, n.host)
	}
	return smtp.SendMail(addr, auth, n.from, []string{n.to}, []byte(message.String()))
}

type ContactService struct {
	repo     store.Repository
	limiter  *RateLimiter
	notifier ContactNotifier
}

func NewContactService(repo store.Repository, limiter *RateLimiter, notifier ContactNotifier) *ContactService {
	if limiter == nil {
		limiter = NewRateLimiter(1*time.Hour, 5)
	}
	if notifier == nil {
		notifier = noopContactNotifier{}
	}
	return &ContactService{repo: repo, limiter: limiter, notifier: notifier}
}

func (s *ContactService) Submit(name, email, message string) (domain.ContactSubmission, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	message = strings.TrimSpace(message)
	if name == "" || email == "" || message == "" {
		return domain.ContactSubmission{}, ErrContactInvalid
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return domain.ContactSubmission{}, ErrContactInvalid
	}
	if !s.limiter.Allow(strings.ToLower(email)) {
		return domain.ContactSubmission{}, ErrContactRateLimit
	}

	submission := domain.ContactSubmission{
		Name:      name,
		Email:     email,
		Message:   message,
		CreatedAt: time.Now(),
	}
	if err := s.repo.SaveContactSubmission(submission); err != nil {
		return domain.ContactSubmission{}, fmt.Errorf("%w: %v", ErrContactStore, err)
	}
	if err := s.notifier.Notify(submission); err != nil {
		return submission, fmt.Errorf("%w: %v", ErrContactNotify, err)
	}
	return submission, nil
}

func (s *ContactService) List() []domain.ContactSubmission {
	return s.repo.ListContactSubmissions()
}