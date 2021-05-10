package userhelper

import (
	"regexp"
	"strings"
)

//goland:noinspection SpellCheckingInspection
var mobiles = []string{"mobile explorer", "palm", "motorola", "nokia", "palm", "iphone", "ipad", "ipod touch", "sony ericsson", "sony ericsson", "blackberry", "o2 cocoon", "treo", "lg", "amoi", "xda", "mda", "vario", "htc", "samsung", "sharp", "siemens", "alcatel", "benq", "hp ipaq", "motorola", "playstation portable", "playstation 3", "playstation 4", "playstation 5", "playstation vita", "danger hiptop", "nec", "panasonic", "philips", "sagem", "sanyo", "spv", "zte", "sendo", "nintendo dsi", "nintendo ds", "nintendo 3ds", "nintendo wii", "open web", "openweb", "android", "symbian", "symbianos", "palm", "symbian s60", "windows ce", "obigo", "netfront browser", "openwave browser", "mobile explorer", "opera mini", "opera mobile", "firefox mobile", "digital paths", "avantgo", "xiino", "novarra transcoder", "vodafone", "ntt docomo", "o2", "mobile", "wireless", "j2me", "midp", "cldc", "up.link", "up.browser", "smartphone", "cellphone", "generic mobile"}

func CheckUA(ua string) (isMobile, isSpecial bool) {
	ua = strings.ToLower(ua)
	isMobile = false
	isSpecial = false

	if regexp.MustCompile(`(edge|trident)`).MatchString(ua) {
		isSpecial = true
	} else if regexp.MustCompile(`(firefox|chrome)`).MatchString(ua) {
		isSpecial = false
	} else if regexp.MustCompile(`safari`).MatchString(ua) {
		isSpecial = false
	} else {
		isSpecial = true
	}

	for _, device := range mobiles {
		if strings.Index(ua, device) > -1 {
			isMobile = true
		}
	}

	return
}
