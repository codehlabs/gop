package gop

// Configures driver behavior
type DriverConfig struct {
	UniqueUsername bool // unique database username defaults to true
	UniqueEmail    bool // unique database email defaults to true
	UniquePhone    bool // unique database phone defaults to true
}

var driver_config = &DriverConfig{true, true, true}

type ActionDriver interface {
	Login(username, email, phone string, password string) (id string, err error)
}

// Sets driver configuration
func SetDriverConfig(c DriverConfig) {
	driver_config = &c
}
