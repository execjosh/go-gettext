package domain

type Domain map[string]string

func (d Domain) Gettext(msgid string) string {
	if msg, ok := d[msgid]; ok {
		return msg
	}
	return msgid
}
