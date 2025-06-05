package tempfile

import "os"

type TempFile struct {
	Pattern string
	Path    *os.File
}

func NewTempFile(pattern string) (*TempFile, error) {
	tempFile := &TempFile{Pattern: pattern}
	var err error
	tempFile.Path, err = os.CreateTemp("", pattern)
	if err != nil {
		return nil, err
	}

	return tempFile, nil
}

func (t *TempFile) Remove() error {
	if t.Path != nil {
		err := os.Remove(t.Path.Name())
		if err != nil {
			return err
		}
		t.Path = nil
	}

	return nil
}

func (t *TempFile) ReadFile() ([]byte, error) {
	if t.Path == nil {
		return nil, os.ErrInvalid
	}

	data, err := os.ReadFile(t.Path.Name())
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (t *TempFile) GetFileName() string {
	if t.Path == nil {
		return ""
	}

	return t.Path.Name()
}
