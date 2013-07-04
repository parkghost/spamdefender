package mailfile

import (
	"bytes"
	"io"
	"net/mail"
	"os/exec"
)

var (
	CmdPostcatPath = "/usr/sbin/postcat"
	CmdParameter   = []string{"-hb"} //print header and body only
)

type PostfixMail struct {
	parsed   bool
	filePath string
	subject  string
	content  string
}

func (m *PostfixMail) Subject() string {
	if !m.parsed {
		m.parse()
	}

	return m.subject
}

func (m *PostfixMail) Content() string {
	if !m.parsed {
		m.parse()
	}

	return m.content
}

func (m *PostfixMail) parse() (err error) {

	cmd := &exec.Cmd{}
	cmdBuf := &bytes.Buffer{}

	cmd = exec.Command(CmdPostcatPath, append(CmdParameter, m.filePath)...)
	cmd.Stderr = cmdBuf
	cmd.Stdout = cmdBuf

	err = cmd.Run()
	if err != nil {
		return
	}

	message, err := mail.ReadMessage(cmdBuf)
	if err != nil {
		return
	}

	rawSubjectStr := message.Header.Get("Subject")

	m.subject, err = DecodeRFC2047String(rawSubjectStr)
	if err != nil {
		return
	}

	bodyContent := &bytes.Buffer{}
	io.Copy(bodyContent, message.Body)
	m.content = bodyContent.String()

	m.parsed = true

	return
}

func NewPostfixMail(filePath string) Mail {
	return &PostfixMail{filePath: filePath}
}
