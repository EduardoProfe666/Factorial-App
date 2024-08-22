package colors

const (
	Reset       = "\033[0m"
	RedColor    = "\033[31m"
	GreenColor  = "\033[32m"
	YellowColor = "\033[33m"
	BlueColor   = "\033[34m"
	PurpleColor = "\033[35m"
	CyanColor   = "\033[36m"
	GrayColor   = "\033[37m"
	WhiteColor  = "\033[97m"
)

func Colorize(color, text string) string {
	return color + text + Reset
}

func Red(text string) string {
	return Colorize(RedColor, text)
}

func Green(text string) string {
	return Colorize(GreenColor, text)
}

func Yellow(text string) string {
	return Colorize(YellowColor, text)
}

func Blue(text string) string {
	return Colorize(BlueColor, text)
}

func Purple(text string) string {
	return Colorize(PurpleColor, text)
}

func Cyan(text string) string {
	return Colorize(CyanColor, text)
}

func Gray(text string) string {
	return Colorize(GrayColor, text)
}

func White(text string) string {
	return Colorize(WhiteColor, text)
}
