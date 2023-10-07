package mysql

type noopScanner struct{}

func (noopScanner) Scan(interface{}) error {
	return nil
}
