package data

import "regexp"

func CheckPatternTelephone(telelphone string) bool {
	// pattern := "^[1]\\d{9}|[9]\\d{8}|[2-8]\\d{8,9}$"
	r, _ := regexp.Compile("^(0?|\\+84|\\+840)\\d{9}$")
	return r.MatchString(telelphone)
}

func CheckPatternOTP(otpVoice string) bool {
	// pattern := "^[1]\\d{9}|[9]\\d{8}|[2-8]\\d{8,9}$"
	r, _ := regexp.Compile("^\\d+$")
	return r.MatchString(otpVoice)
}

func CheckPatternNumeric(Key string) bool {
	// pattern := "^[1]\\d{9}|[9]\\d{8}|[2-8]\\d{8,9}$"
	r, _ := regexp.Compile("^\\d+$")
	return r.MatchString(Key)
}
func CheckPatternNumericText(Key string) bool {
	r, _ := regexp.Compile("^#\\d+$")
	return r.MatchString(Key)
}
