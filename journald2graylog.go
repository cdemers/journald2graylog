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

	"github.com/cdemers/journald2graylog/gelf"
	"github.com/cdemers/journald2graylog/journald"
	rkgelf "github.com/robertkowalski/graylog-golang"
)

func parseGraylogConfig() (hostname string, port int, packetSize int, err error) {

	hostname = os.Getenv("J2G_HOSTNAME")
	portStr := os.Getenv("J2G_PORT")
	packetSizeStr := os.Getenv("J2G_PACKET_SIZE")

	if hostname == "" {
		err = fmt.Errorf("The Graylog server hostname is required but was not specified.")
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

func parseCommandLineFlags() (verboseFlag *bool) {
	verboseFlag = flag.Bool("verbose", false, "Wether journald2graylog will be verbose or not.")
	flag.Parse()
	return verboseFlag
}

func main() {
	verbose := *parseCommandLineFlags()

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

	if verbose {
		log.Printf("Graylog host:\"%s\" port:\"%d\" packet size:\"%d\".",
			graylogHostname, graylogPort, graylogPacketSize)
	}

	// Build the go reader of stdin from where the log stream will be comming
	// from.
	reader := bufio.NewReader(os.Stdin)

	// Loop and process entries from stdin until EOF.
	for {
		var logEntry journald.JournaldJSONLogEntry
		var gelfLogEntry gelf.GELFLogEntry

		line, overflow, err := reader.ReadLine()
		if err == io.EOF {
			os.Exit(0)
		}
		if err != nil {
			panic(err)
		}
		if overflow {
			log.Println("Got a log line that was bigger than our buffer, it will be skiped.")
			continue
		}

		err = json.Unmarshal(line, &logEntry)
		if err != nil {
			log.Printf("The following log line was not a correctly JSON encoded, it will be skiped: \"%s\"\n", line)
			continue
		}

		// Populating the new GELF structure with all the data we received from
		// the journald's JSON formatted data from stdin.
		gelfLogEntry.RawLogLine = string(line)

		// GELF: Version, mendatory.
		gelfLogEntry.Version = "1.1"

		// GELF: Hostname
		if logEntry.Hostname == "" || logEntry.Hostname == "localhost" {
			gelfLogEntry.Host = defaultHostname
		} else {
			gelfLogEntry.Host = logEntry.Hostname
		}

		// GELF: Log Priority/Level
		gelfLogEntry.Level, err = strconv.Atoi(logEntry.Priority)
		if err != nil {
			panic(err)
		}

		// GELF: Message (ShortMessage)
		gelfLogEntry.ShortMessage = logEntry.Message

		// GELF: Timestamp
		var jts = logEntry.RealtimeTimestamp
		gelfLogEntry.Timestamp, err = strconv.Atoi(fmt.Sprintf("%s", jts[:10]))

		// GELF: Facility
		gelfLogEntry.Facility = fmt.Sprintf("%s (%s)", logEntry.SyslogFacility, logEntry.SyslogIdentifier)

		// GELF: BootId
		gelfLogEntry.BootID = logEntry.BootID

		// GELF: MachineId
		gelfLogEntry.MachineID = logEntry.MachineID

		// GELF: PID
		gelfLogEntry.PID = logEntry.PID
		// GELF: UID
		gelfLogEntry.UID = logEntry.UID
		// GELF: GID
		gelfLogEntry.GID = logEntry.GID

		// GELF: Executable
		gelfLogEntry.Executable = logEntry.Executable
		// GELF: Command Line
		gelfLogEntry.CommandLine = logEntry.CommandLine

		var lineNumber int
		if logEntry.CodeLine != "" {
			lineNumber, err = strconv.Atoi(logEntry.CodeLine)
			if err != nil {
				break
			}
			// GELF: Line
			gelfLogEntry.Line = &lineNumber
			// GELF: File
			gelfLogEntry.File = logEntry.CodeFile
			// GELF: Function
			gelfLogEntry.Function = logEntry.CodeFunction
		}

		// GELF: Transport
		gelfLogEntry.Transport = logEntry.Transport

		// Prepare and send the GELF payload to the Graylog server.
		gelfPayloadBytes, err := json.Marshal(gelfLogEntry)
		if err != nil {
			panic(err)
		}

		gelfPayload := string(gelfPayloadBytes)

		if verbose {
			log.Println(gelfPayload)
		}

		graylog.Log(gelfPayload)
	}

}
