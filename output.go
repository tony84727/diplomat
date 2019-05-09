package diplomat

type Output interface {
	WriteFile(filename string, data []byte) error
}
