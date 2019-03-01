package args

// Args has flags for cli
type Args struct {
	FilePath         []string
	DownloadPath     string
	ResumeCapability bool
	Resume           bool
	Verbose          bool
}

// ARGS is an instance of Args
var ARGS = &Args{make([]string, 0), "", true, false, true}
