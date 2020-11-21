package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	su "github.com/qumine/qumine-server-java/internal/updater/server"
	"github.com/sirupsen/logrus"
)

var (
	version = "dev"
	commit  = "none"
	date    = "uknown"
)

var (
	helpFlag    bool
	versionFlag bool
	debugFlag   bool
	traceFlag   bool

	serverTypeFlag          string
	serverVersionFlag       string
	serverDownloadApiFlag   string
	serverDownloadForceFlag bool

	forceDownloadPluginsFlag bool
)

func init() {
	flag.BoolVar(&helpFlag, "help", false, "Show this page")
	flag.BoolVar(&versionFlag, "version", false, "Show the current version")
	flag.BoolVar(&debugFlag, "debug", false, "Enable debugging log level")
	flag.BoolVar(&traceFlag, "trace", false, "Enable trace log level")

	flag.StringVar(&serverTypeFlag, "server-type", "vanilla", "Which server type to use.")
	flag.StringVar(&serverVersionFlag, "server-version", "latest", "Which server version to use.")
	flag.StringVar(&serverDownloadApiFlag, "server-download-api", "", "Url to the server download api.")
	flag.BoolVar(&serverDownloadForceFlag, "server-download-force", false, "Force the download of the server jar")

	flag.BoolVar(&forceDownloadPluginsFlag, "force-download-plugins", false, "Force the download of the server plugins")
	flag.Parse()
}

func main() {
	if helpFlag {
		showHelp()
	}

	if versionFlag {
		showVersion()
	}

	if debugFlag {
		enableDebug()
	}
	if traceFlag {
		enableTrace()
	}

	serverUpdater := su.NewVanillaUpdater(serverVersionFlag, serverDownloadApiFlag, serverDownloadForceFlag)
	serverUpdater.Update()

	logrus.Info("writing eula.txt")
	ioutil.WriteFile("./eula.txt", []byte("eula=true"), os.ModeAppend)

	logrus.Info("handing over to java")
	cmd := exec.Command("java", "-jar", "server.jar")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Start()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go cmd.Wait()

	<-c
}

func showHelp() {
	flag.Usage()
	os.Exit(0)
}

func showVersion() {
	fmt.Printf("%v, commit %v, built at %v", version, commit, date)
	os.Exit(0)
}

func enableDebug() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debug("debugging enabled")
}

func enableTrace() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.Debug("tracing enabled")
}
