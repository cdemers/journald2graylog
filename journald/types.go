package journald


type JournaldJsonLogEntry struct {
  // User Journal Fields (from freedesktop docs)
  Message string `json:"MESSAGE"`
  MessageId string `json:"MESSAGE_ID"`
  Priority string `json:"PRIORITY"`
  CodeFile string `json:"CODE_FILE"`
  CodeLine string `json:"CODE_LINE"`
  CodeFunction string `json:"CODE_FUNCTION"`
  Errno string `json:"ERRNO"`
  SyslogFacility string `json:"SYSLOG_FACILITY"`
  SyslogIdentifier string `json:"SYSLOG_IDENTIFIER"`
  SyslogPid string `json:"SYSLOG_PID"`

  // Trusted Journal Fields (from freedesktop docs)
  PID string `json:"_PID"`
  UID string `json:"_UID"`
  GID string `json:"_GID"`

  Command string `json:"_COMM"`
  Executable string `json:"_EXE"`
  CommandLine string `json:"_CMDLINE"`

  AuditSession string `json:"_AUDIT_SESSION"`
  AuditLoginUid string `json:"_AUDIT_LOGINUID"`

  SystemdCGroup string `json:"_SYSTEMD_CGROUP"`
  SystemdSession string `json:"_SYSTEMD_SESSION"`
  SystemdUnit string `json:"_SYSTEMD_UNIT"`
  SystemdUserUnit string `json:"_SYSTEMD_USER_UNIT"`
  SystemdOwnerUid string `json:"_SYSTEMD_OWNER_UID"`
  SystemdSlice string `json:"_SYSTEMD_SLICE"`

  EffectiveCapabilities string `json:"_CAP_EFFECTIVE"`
  SELinuxContext string `json:"_SELINUX_CONTEXT"`
  SourceRealtimeTimestamp string `json:"_SOURCE_REALTIME_TIMESTAMP"`
  BootId string `json:"_BOOT_ID"`
  MachineId string `json:"_MACHINE_ID"`
  Hostname string `json:"_HOSTNAME"`
  Transport string `json:"_TRANSPORT"` // One of 'audit' 'driver' 'syslog' 'journal' 'stdout' 'kernel'

  // Address Fields (from freedesktop docs)
  Cursor string `json:"__CURSOR"`
  RealtimeTimestamp string `json:"__REALTIME_TIMESTAMP"`
  MonotonicTimestamp string `json:"__MONOTONIC_TIMESTAMP"`

  // Other
  Unit string `json:"UNIT"` // ???
}
