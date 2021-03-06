package cmd

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stephane-martin/bouncer/log"
)

var ConfigDir string
var LogLevel string
var RequestLogLevel string
var LogFilename string
var RequestLogFilename string
var PidFilename string
var Syslog bool
var LogInJSON bool
var ConsulAddr string
var ConsulPrefix string
var ConsulToken string
var ConsulDatacenter string
var ConsulLdapDatacenter string
var ConsulLdapServiceName string
var ConsulLdapTag string

var RootCmd = &cobra.Command{
	Use:   "bouncer",
	Short: "Authentication with LDAP backend for nginx",
	Long: `bouncer implements authentication with an LDAP backend for Nginx.
	It uses 'ngx_http_auth_request_module' module to perform the authentication.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(init_logging)
	RootCmd.PersistentFlags().StringVar(&ConfigDir, "config", "", "the configuration directory to search")
	RootCmd.PersistentFlags().StringVar(&LogLevel, "loglevel", "info", "set logging level")
	RootCmd.PersistentFlags().StringVar(&RequestLogLevel, "req-loglevel", "info", "set request logging level")
	RootCmd.PersistentFlags().StringVar(&LogFilename, "logfile", "", "if specified, write logs to that file instead of stdout/stderr")
	RootCmd.PersistentFlags().StringVar(&RequestLogFilename, "req-logfile", "", "if specified, write request logs to that file instead of stdout/stderr")
	RootCmd.PersistentFlags().StringVar(&PidFilename, "pidfile", "", "if specified, write PID there")
	RootCmd.PersistentFlags().BoolVar(&Syslog, "syslog", false, "if specified, send all logs to the local syslog instead of stdout/stderr")
	RootCmd.PersistentFlags().BoolVar(&LogInJSON, "json", false, "if specified, write logs in JSON format")

	RootCmd.PersistentFlags().StringVar(&ConsulAddr, "consul", "", "Consul scheme, host and port (eg http://A.B.C.D:PORT)")
	RootCmd.PersistentFlags().StringVar(&ConsulPrefix, "prefix", "bouncer", "bouncer prefix in Consul KV")
	RootCmd.PersistentFlags().StringVar(&ConsulToken, "token", "", "Consul token")
	RootCmd.PersistentFlags().StringVar(&ConsulDatacenter, "datacenter", "", "Consul datacenter that stores bouncer configuration")
	RootCmd.PersistentFlags().StringVar(&ConsulLdapDatacenter, "ldap-datacenter", "", "Consul datacenter that hosts the LDAP servers")
	RootCmd.PersistentFlags().StringVar(&ConsulLdapServiceName, "ldap-service-name", "", "Consul LDAP service name to discover")
	RootCmd.PersistentFlags().StringVar(&ConsulLdapTag, "ldap-tag", "", "Consul LDAP services filter tag")
}

func init_logging() {
	if LogInJSON {
		log.Log.Formatter = &logrus.JSONFormatter{}
		log.RequestLog.Formatter = &logrus.JSONFormatter{}
	}

	lvl, err := logrus.ParseLevel(LogLevel)
	if err != nil {
		log.Log.WithError(err).WithField("loglevel", LogLevel).Warn("Unknown log level. Selecting INFO instead.")
		log.Log.Level = logrus.InfoLevel
	} else {
		log.Log.Level = lvl
	}

	lvl, err = logrus.ParseLevel(RequestLogLevel)
	if err != nil {
		log.Log.WithError(err).WithField("req-loglevel", RequestLogLevel).Warn("Unknown log level. Selecting INFO instead.")
		log.RequestLog.Level = logrus.InfoLevel
	} else {
		log.RequestLog.Level = lvl
	}
}
