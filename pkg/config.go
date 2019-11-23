package pkg

import "github.com/spf13/viper"

// GetProviders get all providers
func GetProviders() (providers []string) {
	providers = viper.GetStringSlice("providers")
	return
}

// GetJSONServers get all JSON servers
func GetJSONServers() map[string]string {
	return viper.GetStringMapString("jsonServers")
}
