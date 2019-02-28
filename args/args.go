package args

type Args struct {
	FilePath []string
	DownloadPath string
	ResumeCapability bool
	Resume bool
	Verbose bool
}

var ARGS = &Args { make([]string, 0) ,"", true , false, true}