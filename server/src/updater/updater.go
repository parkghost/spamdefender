package updater

import "mailfile"

type Updater interface {
	Update(mailfile.Mail)
}
