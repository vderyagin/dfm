package dotfile

type FailError string
type SkipError string

func (e FailError) Error() string {
	return string(e)
}

func (e SkipError) Error() string {
	return string(e)
}
