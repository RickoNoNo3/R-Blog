package bloghelper

func MakeTitle(title string) (titleRes string) {
	if title == "" {
		return "R崽的博客"
	} else {
		return title + " - R崽的博客"
	}
}
