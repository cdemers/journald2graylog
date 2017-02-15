package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/cdemers/journald2graylog/blacklist"
	"github.com/cdemers/journald2graylog/gelf"
	"github.com/cdemers/journald2graylog/journald"
	rkgelf "github.com/robertkowalski/graylog-golang"
)

var (
	blacklistFlag = flag.String("J2G_BLACKLIST", os.Getenv("J2G_BLACKLIST"), "Blacklist Regex with ; separator ( e.g. : \"foo.*;bar.*\" )")
	hostname      = flag.String("J2G_HOSTNAME", os.Getenv("J2G_HOSTNAME"), "Hostname or IP of your Graylog server, it has no default and MUST be specified.")
	portStr       = flag.String("J2G_PORT", os.Getenv("J2G_PORT"), "Port of the UDP GELF input of the Graylog server, it will default to 12201")
	packetSizeStr = flag.String("J2G_PACKET_SIZE", os.Getenv("J2G_PACKET_SIZE"), "Maximum size of the TCP/IP packets you can use between the source (journald2graylg) and the destination (your Graylog server). Defaults to 1420")
)

func parseGraylogConfig() (hostname string, port int, packetSize int, err error) {
	if hostname == "" {
		err = fmt.Errorf("The Graylog server hostname is required but was not specified. The server hostname is expected to be specified via the J2G_HOSTNAME environment variable.")
		return "", 0, 0, err
	}

	if portStr == "" {
		port = 12201
	} else {
		port, err = strconv.Atoi(portStr)
		if err != nil {
			err = fmt.Errorf("Unable to parse the port number as an natural number.")
			return "", 0, 0, err
		}
	}

	if packetSizeStr == "" {
		packetSize = 1420
	} else {
		packetSize, err = strconv.Atoi(packetSizeStr)
		if err != nil {
			err = fmt.Errorf("Unable to parse the packet size as an natural number.")
			return "", 0, 0, err
		}
	}

	return hostname, port, packetSize, nil
}

func parseCommandLineFlags() (verboseFlag *bool, disableRawLogLine *bool) {
	verboseFlag = flag.Bool("verbose", false, "Wether journald2graylog will be verbose or not.")
	disableRawLogLine = flag.Bool("disable-rawlogline", false, "Wether journald2graylog will send the raw log line or not.")
	flag.Parse()
	return verboseFlag, disableRawLogLine
}

func main() {
	verbose, disableRawLogLine := parseCommandLineFlags()

	graylogHostname, graylogPort, graylogPacketSize, err := parseGraylogConfig()
	if err != nil {
		panic(err)
	}

	// Determine what will be the default value of the "hostname" field in the
	// GELF payload.
	defaultHostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	if defaultHostname == "localhost" || defaultHostname == "" {
		defaultHostname = "Unknown Host"
	}

	// Build the object that will allow us to transmit messages to the Graylog
	// server.
	graylog := rkgelf.New(rkgelf.Config{
		GraylogHostname: graylogHostname,
		GraylogPort:     graylogPort,
		Connection:      "wan",
		MaxChunkSizeLan: graylogPacketSize,
	})

	b := blacklist.PrepareBlacklist(blacklistFlag)

	if *verbose {
		log.Printf("Graylog host:\"%s\" port:\"%d\" packet size:\"%d\" blacklist:\"%v\"",
			graylogHostname, graylogPort, graylogPacketSize, b)
	}

	// Build the go reader of stdin from where the log stream will be comming from.
	reader := bufio.NewReader(os.Stdin)

	// Loop and process entries from stdin until EOF.
	for {
		line, overflow, err := reader.ReadLine()
		if err == io.EOF {
			os.Exit(0)
		}
		if err != nil {
			panic(err)
		}
		if overflow {
			log.Println("Got a log line that was bigger than the allocated buffer, it will be skiped.")
			continue
		}
		if b.IsBlacklisted(line) {
			continue
		}

		gelfPayload := prepareGelfPayload(disableRawLogLine, line, defaultHostname)
		if gelfPayload == "" {
			continue
		}

		if *verbose {
			log.Println(gelfPayload)
		}

		graylog.Log(gelfPayload)
	}

}

func prepareGelfPayload(disableRawLogLine *bool, line []byte, defaultHostname string) string {
	var logEntry journald.JournaldJSONLogEntry
	var gelfLogEntry gelf.GELFLogEntry

	err := json.Unmarshal(line, &logEntry)
	if err != nil {
		log.Printf("The following log line was not a correctly JSON encoded, it will be skiped: \"%s\"\n", line)
		return ""
	}

	if !*disableRawLogLine {
		gelfLogEntry.RawLogLine = string(line)
	}
	gelfLogEntry.Version = "1.1"
	if logEntry.Hostname == "" || logEntry.Hostname == "localhost" {
		gelfLogEntry.Host = defaultHostname
	} else {
		gelfLogEntry.Host = logEntry.Hostname
	}
	gelfLogEntry.Level, err = strconv.Atoi(logEntry.Priority)
	if err != nil {
		panic(err)
	}
	gelfLogEntry.ShortMessage = logEntry.Message
	var jts = logEntry.RealtimeTimestamp
	gelfLogEntry.Timestamp, err = strconv.ParseFloat(fmt.Sprintf("%s.%s", jts[:10], jts[10:]), 64)
	if (logEntry.SyslogFacility != "") && (logEntry.SyslogIdentifier != "") {
		gelfLogEntry.Facility = fmt.Sprintf("%s (%s)", logEntry.SyslogFacility, logEntry.SyslogIdentifier)
	} else if logEntry.SyslogFacility != "" {
		gelfLogEntry.Facility = logEntry.SyslogFacility
	} else if logEntry.SyslogIdentifier != "" {
		gelfLogEntry.Facility = logEntry.SyslogIdentifier
	} else {
		gelfLogEntry.Facility = "Undefined"
	}
	gelfLogEntry.BootID = logEntry.BootID
	gelfLogEntry.MachineID = logEntry.MachineID
	gelfLogEntry.PID = logEntry.PID
	gelfLogEntry.UID = logEntry.UID
	gelfLogEntry.GID = logEntry.GID
	gelfLogEntry.Executable = logEntry.Executable
	gelfLogEntry.CommandLine = logEntry.CommandLine
	var lineNumber int
	if logEntry.CodeLine != "" {
		lineNumber, err = strconv.Atoi(logEntry.CodeLine)
		if err != nil {
			return ""
		}
		// GELF: Line
		gelfLogEntry.Line = &lineNumber
		// GELF: File
		gelfLogEntry.File = logEntry.CodeFile
		// GELF: Function
		gelfLogEntry.Function = logEntry.CodeFunction
	}
	gelfLogEntry.Transport = logEntry.Transport
	gelfPayloadBytes, err := json.Marshal(gelfLogEntry)
	if err != nil {
		panic(err)
	}
	gelfPayload := string(gelfPayloadBytes)
	return gelfPayload
}
