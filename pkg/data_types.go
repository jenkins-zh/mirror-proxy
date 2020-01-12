package pkg

// PluginDownloadData represents the plugins download data
type PluginDownloadData struct {
	Year    string
	Plugins map[string]PluginData
}

// PluginData represents a plugin data
type PluginData struct {
	Data map[string]int64
}
