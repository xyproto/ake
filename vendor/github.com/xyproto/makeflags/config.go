package makeflags

import (
	"flag"
	"os"
	"path/filepath"
)

const (
	// The default number of concurrent jobs
	defaultJobs = 99

	// The default maximum load average
	defaultMaxLoadAvg = 99.0

	// The default job synchronization type
	defaultSyncType = "none" // none | line | target | recurse
)

// Config contains the same information as make gets after parsing the given flags and arguments
type Config struct {
	AlwaysMake          bool
	Directory           string
	DebugInfo           bool
	DebugFlags          string
	EnvironmentOverride bool
	Evaluate            string
	Makefile            string
	IgnoreErrors        bool
	IncludeSearchPath   string
	Jobs                int
	Targets             []string
	KeepGoing           bool
	MaxLoadAvg          float64
	CheckSymlinkTime    bool
	DryRun              bool
	KeepThisFile        string
	SyncType            string
	PrintInternalDB     bool
	StatusOnly          bool
	NoBuiltinRules      bool
	NoBuiltinVars       bool
	Silent              bool
	TouchTargets        bool
	PrintTrace          bool
	PrintDirectory      bool
	AssumeNewFile       string
	WarnUndefined       bool
	VersionInfoAndExit  bool
}

// exists checks if the given path exists
func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// New will parse the given arguments and return a Config struct
func New() *Config {
	config := &Config{}

	_ = flag.Bool("b", false, "Ignored for compatibility")
	_ = flag.Bool("m", false, "Ignored for compatibility")

	alwaysMakeShort := flag.Bool("B", false, "Unconditionally make all targets")
	alwaysMakeLong := flag.Bool("always-make", false, "Unconditionally make all targets")

	dirShort := flag.String("C", "", "Change to DIRECTORY before doing anything")
	dirLong := flag.String("directory", "", "Change to DIRECTORY before doing anything")

	debugInfo := flag.Bool("d", false, "Print lots of debugging information")

	debugFlags := flag.String("debug", "", "Print various types of debugging information")

	envShort := flag.Bool("e", false, "Environment variables override makefiles")
	envLong := flag.Bool("environment-overrides", false, "Environment variables override makefiles")

	eval := flag.String("eval", "", "Evaluate STRING as a makefile statement")

	filename1 := flag.String("f", "", "Read FILE as a makefile")
	filename2 := flag.String("file", "", "Read FILE as a makefile")
	filename3 := flag.String("makefile", "", "Read FILE as a makefile")

	ignoreShort := flag.Bool("i", false, "Ignore errors from recipes")
	ignoreLong := flag.Bool("ignore-errors", false, "Ignore errors from recipes")

	includeShort := flag.String("I", "", "Search DIRECTORY for included makefiles")
	includeLong := flag.String("include-dir", "", "Search DIRECTORY for included makefiles")

	jobsShort := flag.Int("j", defaultJobs, "Allow N jobs at once; infinite jobs with no arg")
	jobsLong := flag.Int("jobs", defaultJobs, "Allow N jobs at once; infinite jobs with no arg")

	keepGoingShort := flag.Bool("k", false, "Keep going when some targets can't be made")
	keepGoingLong := flag.Bool("keep-going", false, "Keep going when some targets can't be made")

	maxLoadAvg1 := flag.Float64("l", defaultMaxLoadAvg, "Don't start multiple jobs unless load is below N")
	maxLoadAvg2 := flag.Float64("load-average", defaultMaxLoadAvg, "Don't start multiple jobs unless load is below N")
	maxLoadAvg3 := flag.Float64("max-load", defaultMaxLoadAvg, "Don't start multiple jobs unless load is below N")

	checkSymlinkTimesShort := flag.Bool("L", false, "Use the latest mtime between symlinks and target")
	checkSymlinkTimesLong := flag.Bool("check-symlink-times", false, "Use the latest mtime between symlinks and target")

	dryRun1 := flag.Bool("n", false, "Don't actually run any recipe; just print them")
	dryRun2 := flag.Bool("just-print", false, "Don't actually run any recipe; just print them")
	dryRun3 := flag.Bool("dry-run", false, "Don't actually run any recipe; just print them")
	dryRun4 := flag.Bool("recon", false, "Don't actually run any recipe; just print them")

	oldFile1 := flag.String("o", "", "Consider FILE to be very old and don't remake it")
	oldFile2 := flag.String("old-file", "", "Consider FILE to be very old and don't remake it")
	oldFile3 := flag.String("assume-old", "", "Consider FILE to be very old and don't remake it")

	syncTypeShort := flag.String("O", "", "Synchronize output of parallel jobs by TYPE")
	syncTypeLong := flag.String("output-sync", "", "Synchronize output of parallel jobs by TYPE")

	printInternalShort := flag.Bool("p", false, "Print make's internal database")
	printInternalLong := flag.Bool("print-data-base", false, "Print make's internal database")

	statusOnlyShort := flag.Bool("q", false, "Run no recipe; exit status says if up to date")
	statusOnlyLong := flag.Bool("question", false, "Run no recipe; exit status says if up to date")

	noRulesShort := flag.Bool("r", false, "Disable the built-in implicit rules")
	noRulesLong := flag.Bool("no-builtin-rules", false, "Disable the built-in implicit rules")

	noVarsShort := flag.Bool("R", false, "Disable the built-in variable settings")
	noVarsLong := flag.Bool("no-builtin-variables", false, "Disable the built-in variable settings")

	silent1 := flag.Bool("s", false, "Don't echo recipes")
	silent2 := flag.Bool("silent", false, "Don't echo recipes")
	silent3 := flag.Bool("quiet", false, "Don't echo recipes")

	noKeepGoing1 := flag.Bool("S", false, "Turns off -k")
	noKeepGoing2 := flag.Bool("no-keep-going", false, "Turns off -k")
	noKeepGoing3 := flag.Bool("stop", false, "Turns off -k")

	touchShort := flag.Bool("t", false, "Touch targets instead of remaking them")
	touchLong := flag.Bool("touch", false, "Touch targets instead of remaking them")

	printTrace := flag.Bool("trace", false, "Print tracing information")

	printDirectoryShort := flag.Bool("w", false, "Print the current directory")
	printDirectoryLong := flag.Bool("print-directory", false, "Print the current directory")

	noPrintDirectory := flag.Bool("no-print-directory", false, "Turn off -w, even if it was turned on implicitly")

	assumeNewFile1 := flag.String("W", "", "Consider FILE to be infinitely new")
	assumeNewFile2 := flag.String("what-if", "", "Consider FILE to be infinitely new")
	assumeNewFile3 := flag.String("new-file", "", "Consider FILE to be infinitely new")
	assumeNewFile4 := flag.String("assume-new", "", "Consider FILE to be infinitely new")

	warnUndefined := flag.Bool("warn-undefined-variables", false, "Warn when an undefined variable is referenced")

	versionShort := flag.Bool("v", false, "Print the version number of make and exit")
	versionLong := flag.Bool("version", false, "Print the version number of make and exit")

	flag.Parse()

	config.AlwaysMake = *alwaysMakeShort || *alwaysMakeLong

	if *dirShort != "" {
		config.Directory = *dirShort
	}
	if *dirLong != "" {
		config.Directory = *dirLong
	}

	config.DebugInfo = *debugInfo

	config.DebugFlags = *debugFlags

	config.EnvironmentOverride = *envShort || *envLong

	config.Evaluate = *eval

	if *filename1 != "" {
		if filename := filepath.Join(config.Directory, *filename1); exists(filename) {
			config.Makefile = filename
		}
	}
	if *filename2 != "" {
		if filename := filepath.Join(config.Directory, *filename2); exists(filename) {
			config.Makefile = filename
		}
	}
	if *filename3 != "" {
		if filename := filepath.Join(config.Directory, *filename3); exists(filename) {
			config.Makefile = filename
		}
	}
	if config.Makefile == "" {
		if filename := filepath.Join(config.Directory, "GNUmakefile"); exists(filename) {
			config.Makefile = filename
		} else if filename := filepath.Join(config.Directory, "makefile"); exists(filename) {
			config.Makefile = filename
		} else if filename := filepath.Join(config.Directory, "Makefile"); exists(filename) {
			config.Makefile = filename
		}
	}

	config.IgnoreErrors = *ignoreShort || *ignoreLong

	if *includeShort != "" {
		config.IncludeSearchPath = *includeShort
	}
	if *includeLong != "" {
		config.IncludeSearchPath = *includeLong
	}

	config.Jobs = defaultJobs
	if *jobsShort < defaultJobs {
		config.Jobs = *jobsShort
	}
	if *jobsLong < defaultJobs {
		config.Jobs = *jobsLong
	}

	config.KeepGoing = *keepGoingShort || *keepGoingLong

	config.MaxLoadAvg = defaultMaxLoadAvg
	if *maxLoadAvg1 < defaultMaxLoadAvg {
		config.MaxLoadAvg = *maxLoadAvg1
	}
	if *maxLoadAvg2 < defaultMaxLoadAvg {
		config.MaxLoadAvg = *maxLoadAvg2
	}
	if *maxLoadAvg3 < defaultMaxLoadAvg {
		config.MaxLoadAvg = *maxLoadAvg3
	}

	config.CheckSymlinkTime = *checkSymlinkTimesShort || *checkSymlinkTimesLong

	config.DryRun = *dryRun1 || *dryRun2 || *dryRun3 || *dryRun4

	if *oldFile1 != "" {
		config.KeepThisFile = *oldFile1
	}
	if *oldFile2 != "" {
		config.KeepThisFile = *oldFile2
	}
	if *oldFile3 != "" {
		config.KeepThisFile = *oldFile3
	}

	if *syncTypeShort != "" {
		config.SyncType = *syncTypeShort
	}
	if *syncTypeLong != "" {
		config.SyncType = *syncTypeLong
	}

	switch config.SyncType {
	case "":
		// Check if the sync flag is given (without a value)
		outputSyncFlagGiven := false
		flag.Visit(func(f *flag.Flag) {
			if f.Name == "O" || f.Name == "output-sync" {
				outputSyncFlagGiven = true
			}
		})
		if outputSyncFlagGiven {
			config.SyncType = "target"
		} else {
			config.SyncType = defaultSyncType
		}
	case "none", "line", "target", "resurse": // the allowed values
		break
	default:
		config.SyncType = defaultSyncType
	}

	config.PrintInternalDB = *printInternalShort || *printInternalLong

	config.StatusOnly = *statusOnlyShort || *statusOnlyLong

	config.NoBuiltinRules = *noRulesShort || *noRulesLong

	config.NoBuiltinVars = *noVarsShort || *noVarsLong

	config.Silent = *silent1 || *silent2 || *silent3

	noKeepGoing := *noKeepGoing1 || *noKeepGoing2 || *noKeepGoing3
	if noKeepGoing {
		config.KeepGoing = false
	}

	config.TouchTargets = *touchShort || *touchLong

	config.PrintTrace = *printTrace

	config.VersionInfoAndExit = *versionShort || *versionLong

	config.PrintDirectory = *printDirectoryShort || *printDirectoryLong
	if *noPrintDirectory {
		config.PrintDirectory = false
	}

	if *assumeNewFile1 != "" {
		config.AssumeNewFile = *assumeNewFile1
	}
	if *assumeNewFile2 != "" {
		config.AssumeNewFile = *assumeNewFile2
	}
	if *assumeNewFile3 != "" {
		config.AssumeNewFile = *assumeNewFile3
	}
	if *assumeNewFile4 != "" {
		config.AssumeNewFile = *assumeNewFile4
	}

	config.WarnUndefined = *warnUndefined

	config.Targets = flag.Args()

	return config
}
