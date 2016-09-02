package gelf

import "fmt"

// GELFLogEntry is the structure that maps all the GELF fields that will be
// sent to the Graylog server.
type GELFLogEntry struct {
	// Standard GELF Fields
	Version      string  `json:"version"`
	Host         string  `json:"host"`
	ShortMessage string  `json:"short_message"`
	FullMessage  string  `json:"full_message"`
	Timestamp    float64 `json:"timestamp"`
	Level        int     `json:"level"`
	Facility     string  `json:"facility"`
	Line         *int    `json:"line"`
	File         string  `json:"file"`

	// Systemd Extended Fields
	BootID    string `json:"_BootID"`
	MachineID string `json:"_MachineID"`
	UID       string `json:"_UID"`
	GID       string `json:"_GID"`
	PID       string `json:"_PID"`

	Command     string `json:"_Command"`
	Executable  string `json:"_Executable"`
	CommandLine string `json:"_CommandLine"`

	Unit string `json:"_Unit"`

	Transport string `json:"_LogTransport"`

	Function string `json:"_function"`

	// Metadata
	RawLogLine string `json:"_RawLogLine"`
}

func (log *GELFLogEntry) String() (output string) {
	output = fmt.Sprintf("GELF:v%s Host:%s Timestamp:%f Level:%d Facility:%s Line:%d File:%s Message:\"%s\"",
		log.Version, log.Host, log.Timestamp, log.Level, log.Facility, log.Line, log.File, log.ShortMessage)
	return output
}
