// replication-manager - Replication Manager Monitoring and CLI for MariaDB
// Authors: Guillaume Lefranc <guillaume.lefranc@mariadb.com>
//          Stephane Varoqui  <stephane.varoqui@mariadb.com>
// This source code is licensed under the GNU General Public License, version 3.
// Redistribution/Reuse of this code is permitted under the GNU v3 license, as
// an additional term, ALL code must carry the original Author(s) credit in comment form.
// See LICENSE in this directory for the integral text.

package config

type Config struct {
	User               string
	Hosts              string
	Socket             string
	RplUser            string
	Interactive        bool
	Verbose            bool
	PreScript          string `mapstructure:"pre-failover-script"`
	PostScript         string `mapstructure:"post-failover-script"`
	MaxDelay           int64
	GtidCheck          bool
	PrefMaster         string
	IgnoreSrv          string `mapstructure:"ignore-servers"`
	WaitKill           int64  `mapstructure:"wait-kill"`
	WaitTrx            int64  `mapstructure:"wait-trx"`
	ReadOnly           bool
	MaxFail            int `mapstructure:"failcount"`
	Autorejoin         bool
	LogFile            string
	Timeout            int   `mapstructure:"connect-timeout"`
	FailLimit          int   `mapstructure:"failover-limit"`
	FailTime           int64 `mapstructure:"failover-time-limit"`
	CheckType          string
	MasterConn         string `mapstructure:"master-connection"`
	MultiMaster        bool
	Spider             bool
	BindAddr           string `mapstructure:"http-bind-address"`
	HttpPort           string `mapstructure:"http-port"`
	HttpServ           bool   `mapstructure:"http-server"`
	HttpRoot           string `mapstructure:"http-root"`
	Daemon             bool
	MailFrom           string `mapstructure:"mail-from"`
	MailTo             string `mapstructure:"mail-to"`
	MailSMTPAddr       string `mapstructure:"mail-smtp-addr"`
	MasterConnectRetry int    `mapstructure:"master-connect-retry"`
	RplChecks          bool
	FailSync           bool   `mapstructure:"failover-at-sync"`
	FailEventScheduler bool   `mapstructure:"failover-event-scheduler"`
	Heartbeat          bool   `mapstructure:"heartbeat-table"`
	MxsOn              bool   `mapstructure:"maxscale"`
	MxsHost            string `mapstructure:"maxscale-host"`
	MxsPort            string `mapstructure:"maxscale-port"`
	MxsUser            string `mapstructure:"maxscale-user"`
	MxsPass            string `mapstructure:"maxscale-pass"`
	HaproxyOn          bool   `mapstructure:"haproxy"`
	HaproxyWritePort   int    `mapstructure:"haproxy-write-port"`
	HaproxyReadPort    int    `mapstructure:"haproxy-read-port"`
	HaproxyWriteBindIp string `mapstructure:"haproxy-ip-write-bind"`
	HaproxyReadBindIp  string `mapstructure:"haproxy-ip-read-bind"`
	HaproxyBinaryPath  string `mapstructure:"haproxy-binary-path"`
	KeyPath            string
	LogLevel           int `mapstructure:"log-level"`
	Test               bool
}
