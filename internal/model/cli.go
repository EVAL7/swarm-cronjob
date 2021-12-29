package model

import "github.com/alecthomas/kong"

// Cli holds command line args, flags and cmds
type Cli struct {
	Version      kong.VersionFlag
	LogLevel     string `kong:"name='log-level',env='LOG_LEVEL',default='info',help='Set log level.'"`
	LogJSON      bool   `kong:"name='log-json',env='LOG_JSON',default='false',help='Enable JSON logging output.'"`
	EventPort    string `kong:"name='event-port',env='EVENT_PORT',default='8080',help='Port for incoming event requests'"`
	EventTimeout string `kong:"name='event-timeout',env='EVENT_TIMEOUT',default='1h',help='Max time for job'"`
}
