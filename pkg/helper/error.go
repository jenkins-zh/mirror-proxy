package helper

// CheckErr print a friendly error message
func CheckErr(printer Printer, err error) {
	switch {
	case err == nil:
		return
	default:
		printer.Println(err)
	}
}

// Printer
type Printer interface {
	Println(a ...interface{})
}
