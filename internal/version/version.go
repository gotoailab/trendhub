package version

// 这些变量会在构建时通过 -ldflags 注入
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
	GoVersion = "unknown"
)

// GetVersion 获取版本信息
func GetVersion() string {
	return Version
}

// GetBuildInfo 获取构建信息
func GetBuildInfo() map[string]string {
	return map[string]string{
		"version":   Version,
		"buildTime": BuildTime,
		"gitCommit": GitCommit,
		"goVersion": GoVersion,
	}
}

