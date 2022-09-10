package commands_edit

type descriptionFlag struct{ *string }

func (flag *descriptionFlag) String() string {
	if flag.string == nil || *flag.string == "" {
		return ""
	}

	return *flag.string
}

func (flag *descriptionFlag) GetValue() *string {
	if flag.string == nil || *flag.string == "" {
		return nil
	}

	return flag.string
}

func (flag *descriptionFlag) Set(value string) error {
	*flag = descriptionFlag{string: &value}

	return nil
}

func (flag *descriptionFlag) Type() string {
	return "description"
}

func NewDescriptionFlag() *descriptionFlag {
	return &descriptionFlag{string: nil}
}
