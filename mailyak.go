package main

import (
	"bytes"
	"fmt"
	"net/smtp"
	"regexp"
	"strings"
	"time"
  
  "crypto/tls"
  
  "github.com/golang/glog"
)

// TODO: in the future, when aliasing is supported or we're making a breaking
// API change anyway, change the MailYak struct name to Email.

// MailYak represents an email.
type MailYak struct {
  client         *smtp.Client
	html           BodyPart
	plain          BodyPart

	toAddrs        []string
	ccAddrs        []string
	bccAddrs       []string
	subject        string
	fromAddr       string
	fromName       string
	replyTo        string
	headers        map[string]string // arbitrary headers
	attachments    []attachment
	auth           smtp.Auth
	trimRegex      *regexp.Regexp
	host           string
	writeBccHeader bool
	date           string
}

// New returns an instance of MailYak using host as the SMTP server, and
// authenticating with auth where required.
//
// host must include the port number (i.e. "smtp.itsallbroken.com:25")
//
// 		mail := mailyak.New("smtp.itsallbroken.com:25", smtp.PlainAuth(
// 			"",
// 			"username",
// 			"password",
// 			"stmp.itsallbroken.com",
//		))
//
func NewMail(info *SMTPInfo) *MailYak {
  if glog.V(9) {
    glog.Infof("DBG: NEW MAIL: Init (%v)", info)
  }
  m := MailYak{
    client: nil,
		headers:        map[string]string{},
		host:           info.ConnectStr,
		auth:           info.Auth,
		trimRegex:      regexp.MustCompile("\r?\n"),
		writeBccHeader: false,
		date:           time.Now().Format(time.RFC1123Z),
	}
  if info.EnableTLS {
    // Here is the key, you need to call tls.Dial instead of smtp.Dial
    // for smtp servers running on 465 that require an ssl connection
    // from the very beginning (no starttls)
    conn, err := tls.Dial("tcp", info.ConnectStr, &info.TLS)
    if err != nil {
      glog.Errorf("ERR: MAIL: New Mail Dial TLS(%s): %v", info.ConnectStr, err)
      return nil
    }

    m.client, err = smtp.NewClient(conn, info.Address)
    if err != nil {
      glog.Errorf("ERR: MAIL: New Mail Dial Client(%s): %v", info.Address, err)
      return nil
    }
    // Auth
    if err = m.client.Auth(info.Auth); err != nil {
      glog.Errorf("ERR: MAIL: New Mail Auth(%s): %v", info.ConnectStr, err)
      return nil
    }
    if glog.V(9) {
      glog.Infof("DBG: MAIL: Init TLS (%s)", info.ConnectStr)
    }
  }
  
	return &m
}

// Send attempts to send the built email via the configured SMTP server.
//
// Attachments are read when Send() is called, and any connection/authentication
// errors will be returned by Send().
func (m *MailYak) Send() error {
	buf, err := m.buildMime()
	if err != nil {
		return err
	}
  if m.client != nil {
    if glog.V(9) {
      glog.Infof("DBG: MAIL: Sending TLS (%s)", m.host)
    }
    
    if err = m.client.Mail(m.fromAddr); err != nil {
      glog.Errorf("ERR: MAIL: m.fromAddr(%s): %v", m.fromAddr, err)
      return err
    }

    if err = m.client.Rcpt(m.toAddrs[0]); err != nil {
      glog.Errorf("ERR: MAIL: m.toAddrs(%s): %v", m.toAddrs, err)
      return err
    }

    // Data
    w, err := m.client.Data()
    if err != nil {
      glog.Errorf("ERR: MAIL: m.Data: %v", err)
      return err
    }

    _, err = w.Write(buf.Bytes())
    if err != nil {
      glog.Errorf("ERR: MAIL: m.Data.Write: %v", err)
      return err
    }

    err = w.Close()
    if err != nil {
      glog.Errorf("ERR: MAIL: Close: %v", err)
      return err
    }
    return nil
  }
	return smtp.SendMail(
		m.host,
		m.auth,
		m.fromAddr,
		append(append(m.toAddrs, m.ccAddrs...), m.bccAddrs...),
		buf.Bytes(),
	)
}

// MimeBuf returns the buffer containing all the RAW MIME data.
//
// MimeBuf is typically used with an API service such as Amazon SES that does
// not use an SMTP interface.
func (m *MailYak) MimeBuf() (*bytes.Buffer, error) {
	buf, err := m.buildMime()
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// String returns a redacted description of the email state, typically for
// logging or debugging purposes.
//
// Authentication information is not included in the returned string.
func (m *MailYak) String() string {
	var (
		att    []string
		custom string
	)
	for _, a := range m.attachments {
		att = append(att, "{filename: "+a.filename+"}")
	}

	if len(m.headers) > 0 {
		var hdrs []string
		for k, v := range m.headers {
			hdrs = append(hdrs, fmt.Sprintf("%s: %q", k, v))
		}
		custom = strings.Join(hdrs, ", ") + ", "
	}
	return fmt.Sprintf(
		"&MailYak{date: %q, from: %q, fromName: %q, html: %v bytes, plain: %v bytes, toAddrs: %v, "+
			"bccAddrs: %v, subject: %q, %vhost: %q, attachments (%v): %v, auth set: %v}",
		m.date,
		m.fromAddr,
		m.fromName,
		len(m.HTML().String()),
		len(m.Plain().String()),
		m.toAddrs,
		m.bccAddrs,
		m.subject,
		custom,
		m.host,
		len(att),
		att,
		m.auth != nil,
	)
}

// HTML returns a BodyPart for the HTML email body.
func (m *MailYak) HTML() *BodyPart {
	return &m.html
}

// Plain returns a BodyPart for the plain-text email body.
func (m *MailYak) Plain() *BodyPart {
	return &m.plain
}

func (m *MailYak) Close() {
  if m.client != nil {
    m.client.Quit()
  }
}
