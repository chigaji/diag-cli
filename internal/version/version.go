package version

var (
	Version   = "dev"
	Commit    = ""
	BuildDate = ""
)

func FullVersion() string {
	return Version + " (" + Commit + ") built on " + BuildDate
}
