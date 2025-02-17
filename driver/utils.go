package driver

func LoginViaEmail(email, password string, d ActionDriver) (string, error) {
	return d.Login("", email, "", password)
}

// Authenti into the system using username as the identifier
func LoginViaUsername(username, password string, d ActionDriver) (string, error) {
	return d.Login(username, "", "", password)
}

// Authenticate into the system using the phone number as the identifier
// NOTE: may remove this one
func LoginViaPhone(phone, password string, d ActionDriver) (string, error) {
	return d.Login("", "", phone, password)
}
