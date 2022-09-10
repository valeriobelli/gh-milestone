package commands_edit

type titleFlag struct{ *string }

func (flag *titleFlag) String() string {
	if flag.string == nil || *flag.string == "" {
		return ""
	}

	return *flag.string
}

func (flag *titleFlag) GetValue() *string {
	if flag.string == nil || *flag.string == "" {
		return nil
	}

	return flag.string
}

func (flag *titleFlag) Set(value string) error {
	*flag = titleFlag{string: &value}

	return nil
}

func (flag *titleFlag) Type() string {
	return "title"
}

func NewTitleFlag() *titleFlag {
	return &titleFlag{string: nil}
}
