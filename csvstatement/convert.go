package csvstatement

import "io"

func Convert(source io.Reader, target io.Writer, sourceFormat, targetFormat Format) error {
	p := NewParser(sourceFormat)
	stmt, err := p.Parse(source)
	if err != nil {
		return err
	}

	err = WriteStatement(target, stmt, targetFormat)
	if err != nil {
		return err
	}
	return nil
}
