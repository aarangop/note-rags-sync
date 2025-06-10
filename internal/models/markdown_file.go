package models

type MarkdownFile struct {
	Path     string
	Content  []byte
	Checksum string
}

func (m *MarkdownFile) GetContent() []byte {
	return m.Content
}

func (m *MarkdownFile) GetPath() string {
	return m.Path
}

func (m *MarkdownFile) GetChecksum() string {
	return m.Checksum
}
