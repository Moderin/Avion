package color

const (
	green = "\033[1;32m"
	red   = "\033[1;31m"
	orange= "\033[38;5;208m"
	yellow= "\033[38;5;11m"
	reset = "\033[00m"
)

// Green colors console output
func Green(input string) string {
	return green + input + reset
}

// Red colors console output
func Red(input string) string {
	return red + input + reset
}

func Orange(input string) string {
	return orange + input + reset 
}

func Yellow(input string) string {
	return yellow + input + reset 
}