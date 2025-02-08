package main

import (
	"bytes"
	"html/template"
	"log"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type EmailContact struct {
	Email string
	Name  string
}

type Mail struct {
	Domain     string
	Host       string
	Port       int
	Username   string
	Password   string
	From       EmailContact
	Encryption string
}

type Message struct {
	From        EmailContact
	To          EmailContact
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

func (m *Mail) SendSMTPMessage(msg Message) error {
	// default name and email
	if msg.From.Name == "" {
		msg.From.Name = m.From.Name
	}
	if msg.From.Email == "" {
		msg.From.Email = m.From.Email
	}

	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	htmlMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		log.Panic(err)
		return err
	}

	textMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		log.Panic(err)
		return err
	}

	log.Println(msg)
	log.Println(htmlMessage)
	log.Println(textMessage)

	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From.Email).
		AddTo(msg.To.Email).
		SetSubject(msg.Subject)
	email.SetBody(mail.TextPlain, textMessage)
	email.AddAlternative(mail.TextHTML, htmlMessage)

	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	err = email.Send(smtpClient)
	if err != nil {
		log.Panic(err)
		return err
	}

	return nil
}

func (m *Mail) getEncryption(s string) mail.Encryption {
	switch s {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	templateToRender := "./templates/mail.html.gohtml"

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		log.Panic("template.new.parsefiles", err)
		return "", err
	}

	var tpl bytes.Buffer

	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()

	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		log.Panic(err)
		return "", err
	}

	return formattedMessage, nil
}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	templateToRender := "./templates/mail.plain.gohtml"

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		log.Panic(err)
		return "", err
	}

	var tpl bytes.Buffer

	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}

func (m *Mail) inlineCSS(msg string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(msg, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}
