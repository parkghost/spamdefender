package mailfile

import (
	"bytes"
	"github.com/parkghost/pkg/net/mail"
	"os/exec"
	"path"
)

var (
	CmdPostcatPath = "/usr/sbin/postcat"
	CmdParameter   = []string{"-hb"} //print header and body only
)

type PostfixMail struct {
	filePath string
	subject  string
	content  string
	from     *mail.Address
	to       []*mail.Address
}

func (m *PostfixMail) Name() string {
	_, name := path.Split(m.filePath)
	return name
}

func (m *PostfixMail) Path() string {
	return m.filePath
}

func (m *PostfixMail) Subject() string {
	return m.subject
}

func (m *PostfixMail) Content() string {
	return m.content
}

func (m *PostfixMail) From() *mail.Address {
	return m.from
}

func (m *PostfixMail) To() []*mail.Address {
	return m.to
}

func (m *PostfixMail) Parse() (err error) {
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

	m.subject, err = parseSubject(message)
	if err != nil {
		return
	}

	m.from, err = parseFromAddress(message)
	if err != nil {
		return
	}

	m.to, err = parseToAddress(message)
	if err != nil {
		return
	}

	m.content, err = parseBoby(message)
	if err != nil {
		return
	}

	return
}

func (m *PostfixMail) String() string {
	return m.filePath
}

func NewPostfixMail(filePath string) Mail {
	return &PostfixMail{filePath: filePath}
}
