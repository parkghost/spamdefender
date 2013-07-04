package mailfile

type Mail interface {
	Subject() string
	Content() string
}
