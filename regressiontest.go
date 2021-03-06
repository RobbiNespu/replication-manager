// replication-manager - Replication Manager Monitoring and CLI for MariaDB
// Authors: Guillaume Lefranc <guillaume.lefranc@mariadb.com>
//          Stephane Varoqui  <stephane.varoqui@mariadb.com>
// This source code is licensed under the GNU General Public License, version 3.
// Redistribution/Reuse of this code is permitted under the GNU v3 license, as
// an additional term, ALL code must carry the original Author(s) credit in comment form.
// See LICENSE in this directory for the integral text.

package main

import "github.com/tanji/replication-manager/dbhelper"

import "time"
import "sort"

//import "encoding/json"
//import "net/http"

const recover_time = 8

func testSwitchOverLongTransactionNoRplCheckNoSemiSync() bool {
	conf.RplChecks = false
	conf.MaxDelay = 8
	logprintf("TESTING : Starting Test %s", "testSwitchOverLongTransactionNoRplCheckNoSemiSync")
	for _, s := range servers {
		_, err := s.Conn.Exec("set global rpl_semi_sync_master_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
		_, err = s.Conn.Exec("set global rpl_semi_sync_slave_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
	}

	SaveMasterURL := master.URL
	masterTest, _ := newServerMonitor(master.URL)
	defer masterTest.Conn.Close()
	go masterTest.Conn.Exec("start transaction")
	time.Sleep(12 * time.Second)
	for i := 0; i < 1; i++ {

		logprintf("INFO :  Master is %s", master.URL)

		swChan <- true

		wait_failover_end()
		logprintf("INFO : New Master  %s ", master.URL)

	}

	time.Sleep(2 * time.Second)
	if master.URL != SaveMasterURL {
		logprintf("INFO : Saved Prefered master %s <>  from saved %s  ", SaveMasterURL, master.URL)
		return false
	}
	return true
}

func testSwitchOverLongQueryNoRplCheckNoSemiSync() bool {
	conf.RplChecks = false
	conf.MaxDelay = 8
	logprintf("TESTING : Starting Test %s", "testSwitchOverLongQueryNoRplCheckNoSemiSync")
	for _, s := range servers {
		_, err := s.Conn.Exec("set global rpl_semi_sync_master_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
		_, err = s.Conn.Exec("set global rpl_semi_sync_slave_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
	}

	SaveMasterURL := master.URL
	go dbhelper.InjectLongTrx(master.Conn, 20)
	for i := 0; i < 1; i++ {

		logprintf("INFO :  Master is %s", master.URL)

		swChan <- true

		wait_failover_end()
		logprintf("INFO : New Master  %s ", master.URL)

	}

	time.Sleep(20 * time.Second)
	if master.URL != SaveMasterURL {
		logprintf("INFO : Saved Prefered master %s <>  from saved %s  ", SaveMasterURL, master.URL)
		return false
	}
	return true
}

func testSwitchOverLongTransactionWithoutCommitNoRplCheckNoSemiSync() bool {
	conf.RplChecks = false
	conf.MaxDelay = 8
	logprintf("TESTING : Starting Test %s", "testSwitchOverLongTransactionNoRplCheckNoSemiSync")
	for _, s := range servers {
		_, err := s.Conn.Exec("set global rpl_semi_sync_master_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
		_, err = s.Conn.Exec("set global rpl_semi_sync_slave_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
	}

	SaveMasterURL := master.URL
	go dbhelper.InjectLongTrx(master.Conn, 20)
	for i := 0; i < 1; i++ {

		logprintf("INFO :  Master is %s", master.URL)

		swChan <- true

		wait_failover_end()
		logprintf("INFO : New Master  %s ", master.URL)

	}
	for _, s := range slaves {
		dbhelper.StartSlave(s.Conn)
	}
	time.Sleep(2 * time.Second)
	if master.URL != SaveMasterURL {
		logprintf("INFO : Saved Prefered master %s <>  from saved %s  ", SaveMasterURL, master.URL)
		return false
	}
	return true
}

func testSlaReplAllDelay() bool {
	return false
}

func testFailoverReplAllDelayInteractive() bool {
	return false
}

func testFailoverReplAllDelayAuto() bool {
	return false
}

func testSwitchoverReplAllDelay() bool {
	return false
}

func testSlaReplAllSlavesStopNoSemiSync() bool {
	logprintf("TESTING : Starting Test %s", "testSlaReplAllSlavesStopNoSemySync")
	conf.MaxDelay = 0
	for _, s := range servers {
		_, err := s.Conn.Exec("set global rpl_semi_sync_master_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
		_, err = s.Conn.Exec("set global rpl_semi_sync_slave_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
	}

	sme.ResetUpTime()
	time.Sleep(3 * time.Second)
	sla1 := sme.GetUptimeFailable()
	for _, s := range slaves {
		dbhelper.StopSlave(s.Conn)
	}
	time.Sleep(recover_time * time.Second)
	sla2 := sme.GetUptimeFailable()
	for _, s := range slaves {
		dbhelper.StartSlave(s.Conn)
	}
	for _, s := range servers {
		_, err := s.Conn.Exec("set global rpl_semi_sync_master_enabled='ON'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
		_, err = s.Conn.Exec("set global rpl_semi_sync_slave_enabled='ON'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
	}
	if sla2 == sla1 {
		return false
	} else {
		return true
	}
}

func testSlaReplOneSlavesStop() bool {
	for _, s := range slaves {
		dbhelper.StopSlave(s.Conn)
	}
	return false
}

func testSwitchOverReadOnlyNoRplCheck() bool {
	conf.RplChecks = false
	conf.MaxDelay = 0
	logprintf("TESTING : Starting Test %s", "testSwitchOverReadOnlyNoRplCheck")
	logprintf("INFO : Master is %s", master.URL)
	conf.ReadOnly = true
	for _, s := range slaves {
		_, err := s.Conn.Exec("set global read_only=1")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
	}
	swChan <- true
	wait_failover_end()
	logprintf("INFO : New Master is %s ", master.URL)
	for _, s := range slaves {
		logprintf("INFO : Server  %s is %s", s.URL, s.ReadOnly)
		if s.ReadOnly == "OFF" {
			return false
		}
	}
	return true
}

func testSwitchOverNoReadOnlyNoRplCheck() bool {
	conf.RplChecks = false
	conf.MaxDelay = 0
	logprintf("TESTING : Starting Test %s", "testSwitchOverNoReadOnlyNoRplCheck")
	logprintf("INFO : Master is %s", master.URL)
	conf.ReadOnly = false
	for _, s := range servers {
		_, err := s.Conn.Exec("set global read_only=0")
		if err != nil {
			logprintf("ERROR : %s", err.Error())
		}
	}
	SaveMasterURL := master.URL
	swChan <- true
	wait_failover_end()
	logprintf("INFO : New Master is %s ", master.URL)
	if SaveMasterURL == master.URL {
		logprintf("INFO : same server URL after switchover")
		return false
	}
	for _, s := range slaves {
		logprintf("INFO : Server  %s is %s", s.URL, s.ReadOnly)
		if s.ReadOnly != "OFF" {
			return false
		}
	}
	return true
}

func testSwitchOver2TimesReplicationOkNoSemiSyncNoRplCheck() bool {
	conf.RplChecks = false
	conf.MaxDelay = 0
	logprintf("TESTING : Starting Test %s", "testSwitchOver2TimesReplicationOk")
	for _, s := range servers {
		_, err := s.Conn.Exec("set global rpl_semi_sync_master_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
		_, err = s.Conn.Exec("set global rpl_semi_sync_slave_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
	}
	time.Sleep(2 * time.Second)

	for i := 0; i < 2; i++ {
		result, err := dbhelper.WriteConcurrent2(master.DSN, 10)
		if err != nil {
			logprintf("BENCH : %s %s", err.Error(), result)
		}
		logprintf("INFO : New Master  %s ", master.URL)
		SaveMasterURL := master.URL
		swChan <- true

		wait_failover_end()
		logprintf("INFO : New Master  %s ", master.URL)

		if SaveMasterURL == master.URL {
			logprintf("INFO : same server URL after switchover")
			return false
		}
	}

	for _, s := range slaves {
		if s.IOThread != "Yes" || s.SQLThread != "Yes" {
			logprintf("INFO : Slave  %s issue on replication  SQL Thread % IO %s ", s.URL, s.SQLThread, s.IOThread)
			return false
		}
		if s.MasterServerID != master.ServerID {
			logprintf("INFO :  Replication is  pointing to wrong master %s ", master.ServerID)
			return false
		}
	}
	return true
}

func testSwitchOver2TimesReplicationOkSemiSyncNoRplCheck() bool {
	conf.RplChecks = false
	conf.MaxDelay = 0
	logprintf("TESTING : Starting Test %s", "testSwitchOver2TimesReplicationOkSemisync")
	for _, s := range servers {
		_, err := s.Conn.Exec("set global rpl_semi_sync_master_enabled='ON'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
		_, err = s.Conn.Exec("set global rpl_semi_sync_slave_enabled='ON'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
	}
	time.Sleep(2 * time.Second)

	for i := 0; i < 2; i++ {
		result, err := dbhelper.WriteConcurrent2(master.DSN, 10)
		if err != nil {
			logprintf("BENCH : %s %s", err.Error(), result)
		}
		logprintf("INFO : New Master  %s ", master.URL)
		SaveMasterURL := master.URL
		swChan <- true

		wait_failover_end()
		logprintf("INFO : New Master  %s ", master.URL)

		if SaveMasterURL == master.URL {
			logprintf("INFO : same server URL after switchover")
			return false
		}
	}

	for _, s := range slaves {
		if s.IOThread != "Yes" || s.SQLThread != "Yes" {
			logprintf("INFO : Slave  %s issue on replication  SQL Thread % IO %s ", s.URL, s.SQLThread, s.IOThread)
			return false
		}
		if s.MasterServerID != master.ServerID {
			logprintf("INFO :  Replication is  pointing to wrong master %s ", master.ServerID)
			return false
		}
	}
	return true
}

func testSwitchOverBackPreferedMasterNoRplCheckSemiSync() bool {
	conf.RplChecks = false
	conf.MaxDelay = 0
	logprintf("TESTING : Starting Test %s", "testSwitchOverBackPreferedMasterNoRplCheckSemiSync")
	for _, s := range servers {
		_, err := s.Conn.Exec("set global rpl_semi_sync_master_enabled='ON'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
		_, err = s.Conn.Exec("set global rpl_semi_sync_slave_enabled='ON'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
	}
	conf.PrefMaster = master.URL
	logprintf("TESTING : Set conf.PrefMaster %s", "conf.PrefMaster")
	time.Sleep(2 * time.Second)
	SaveMasterURL := master.URL
	for i := 0; i < 2; i++ {

		logprintf("INFO : New Master  %s Failover counter %d", master.URL, i)

		swChan <- true

		wait_failover_end()
		logprintf("INFO : New Master  %s ", master.URL)

	}
	if master.URL != SaveMasterURL {
		logprintf("INFO : Saved Prefered master %s <>  from saved %s  ", SaveMasterURL, master.URL)
		return false
	}
	return true
}

func testSwitchOverAllSlavesStopRplCheckNoSemiSync() bool {
	conf.MaxDelay = 0
	conf.RplChecks = true
	logprintf("TESTING : Starting Test %s", "testSwitchOverAllSlavesStopRplCheckNoSemiSync")
	for _, s := range servers {
		_, err := s.Conn.Exec("set global rpl_semi_sync_master_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
		_, err = s.Conn.Exec("set global rpl_semi_sync_slave_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
	}
	for _, s := range slaves {
		dbhelper.StopSlave(s.Conn)
	}
	time.Sleep(5 * time.Second)

	SaveMasterURL := master.URL
	for i := 0; i < 1; i++ {

		logprintf("INFO :  Master is %s", master.URL)

		swChan <- true
		wait_failover_end()

		logprintf("INFO : New Master  %s ", master.URL)

	}
	for _, s := range slaves {
		dbhelper.StartSlave(s.Conn)
	}
	time.Sleep(2 * time.Second)
	if master.URL != SaveMasterURL {
		logprintf("INFO : Saved Prefered master %s <>  from saved %s  ", SaveMasterURL, master.URL)
		return false
	}
	return true
}

func testSwitchOverAllSlavesStopNoSemiSyncNoRplCheck() bool {
	conf.RplChecks = false
	conf.MaxDelay = 0
	logprintf("TESTING : Starting Test %s", "testSwitchOverAllSlavesStopNoRplCheck")
	for _, s := range servers {
		_, err := s.Conn.Exec("set global rpl_semi_sync_master_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
		_, err = s.Conn.Exec("set global rpl_semi_sync_slave_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
	}
	for _, s := range slaves {
		dbhelper.StopSlave(s.Conn)
	}
	time.Sleep(2 * time.Second)

	SaveMasterURL := master.URL
	for i := 0; i < 1; i++ {

		logprintf("INFO : Master is %s", master.URL)

		swChan <- true

		wait_failover_end()
		logprintf("INFO : New Master  %s ", master.URL)

	}
	for _, s := range slaves {
		dbhelper.StartSlave(s.Conn)
	}
	time.Sleep(2 * time.Second)
	if master.URL == SaveMasterURL {
		logprintf("INFO : Saved Prefered master %s <>  from saved %s  ", SaveMasterURL, master.URL)
		return false
	}
	return true
}

func testSwitchOverAllSlavesDelayRplCheckNoSemiSync() bool {
	conf.RplChecks = true
	conf.MaxDelay = 8
	logprintf("TESTING : Starting Test %s", "testSwitchOverAllSlavesDelayRplCheckNoSemySync")
	for _, s := range servers {
		_, err := s.Conn.Exec("set global rpl_semi_sync_master_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
		_, err = s.Conn.Exec("set global rpl_semi_sync_slave_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
	}
	for _, s := range slaves {
		dbhelper.StopSlave(s.Conn)
	}
	time.Sleep(15 * time.Second)

	SaveMasterURL := master.URL
	for i := 0; i < 1; i++ {

		logprintf("INFO :  Master is %s", master.URL)

		swChan <- true

		wait_failover_end()
		logprintf("INFO : New Master  %s ", master.URL)

	}
	for _, s := range slaves {
		dbhelper.StartSlave(s.Conn)
	}
	time.Sleep(2 * time.Second)
	if master.URL != SaveMasterURL {
		logprintf("INFO : Saved Prefered master %s <>  from saved %s  ", SaveMasterURL, master.URL)
		return false
	}
	return true
}

func testSwitchOverAllSlavesDelayNoRplChecksNoSemiSync() bool {
	conf.RplChecks = false
	conf.MaxDelay = 8
	logprintf("TESTING : Starting Test %s", "testSwitchOverAllSlavesDelayNoRplChecksNoSemiSync")
	for _, s := range servers {
		_, err := s.Conn.Exec("set global rpl_semi_sync_master_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
		_, err = s.Conn.Exec("set global rpl_semi_sync_slave_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
	}
	for _, s := range slaves {
		dbhelper.StopSlave(s.Conn)
	}
	time.Sleep(15 * time.Second)

	SaveMasterURL := master.URL
	for i := 0; i < 1; i++ {

		logprintf("INFO :  Master is %s", master.URL)

		swChan <- true

		wait_failover_end()
		logprintf("INFO : New Master  %s ", master.URL)

	}
	for _, s := range slaves {
		dbhelper.StartSlave(s.Conn)
	}
	time.Sleep(2 * time.Second)
	if master.URL == SaveMasterURL {
		logprintf("INFO : Saved Prefered master %s <>  from saved %s  ", SaveMasterURL, master.URL)
		return false
	}
	return true
}

func testSwitchOverAllSlavesDelayRplChecksNoSemiSync() bool {
	conf.RplChecks = true
	conf.MaxDelay = 8
	logprintf("TESTING : Starting Test %s", "testSwitchOverAllSlavesDelayNoRplChecksNoSemiSync")
	for _, s := range servers {
		_, err := s.Conn.Exec("set global rpl_semi_sync_master_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
		_, err = s.Conn.Exec("set global rpl_semi_sync_slave_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
	}
	for _, s := range slaves {
		dbhelper.StopSlave(s.Conn)
	}
	time.Sleep(10 * time.Second)

	SaveMasterURL := master.URL
	for i := 0; i < 1; i++ {

		logprintf("INFO :  Master is %s", master.URL)

		swChan <- true

		time.Sleep(recover_time * time.Second)
		logprintf("INFO : New Master  %s ", master.URL)

	}
	for _, s := range slaves {
		dbhelper.StartSlave(s.Conn)
	}
	time.Sleep(2 * time.Second)
	if master.URL != SaveMasterURL {
		logprintf("INFO : Saved Prefered master %s <>  from saved %s  ", SaveMasterURL, master.URL)
		return false
	}
	return true
}

func testFailOverAllSlavesDelayNoRplChecksNoSemiSync() bool {

	bootstrap()
	time.Sleep(5 * time.Second)

	logprintf("TESTING : Starting Test %s", "testFailOverAllSlavesDelayNoRplChecksNoSemiSync")
	for _, s := range servers {
		_, err := s.Conn.Exec("set global rpl_semi_sync_master_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
		_, err = s.Conn.Exec("set global rpl_semi_sync_slave_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
	}
	SaveMasterURL := master.URL
	logprintf("BENCH: Stopping replication")
	for _, s := range slaves {
		dbhelper.StopSlave(s.Conn)
	}
	result, err := dbhelper.WriteConcurrent2(master.DSN, 10)
	if err != nil {
		logprintf("BENCH : %s %s", err.Error(), result)
	}
	logprintf("BENCH : Write Concurrent Insert")

	dbhelper.InjectLongTrx(master.Conn, 10)
	logprintf("BENCH : Inject Long Trx")
	time.Sleep(10 * time.Second)
	logprintf("BENCH : Sarting replication")
	for _, s := range slaves {
		dbhelper.StartSlave(s.Conn)
	}

	logprintf("INFO :  Master is %s", master.URL)
	conf.Interactive = false
	master.FailCount = conf.MaxFail
	master.State = stateFailed
	conf.FailLimit = 5
	conf.FailTime = 0
	failoverCtr = 0
	conf.RplChecks = false
	conf.MaxDelay = 4
	checkfailed()

	wait_failover_end()
	logprintf("INFO : New Master  %s ", master.URL)

	time.Sleep(2 * time.Second)
	if master.URL == SaveMasterURL {
		logprintf("INFO : Old master %s ==  New master %s  ", SaveMasterURL, master.URL)

		return false
	}

	return true
}

func testFailOverAllSlavesDelayRplChecksNoSemiSync() bool {

	bootstrap()
	wait_failover_end()

	logprintf("TESTING : Starting Test %s", "testFailOverAllSlavesDelayRplChecksNoSemiSync")
	for _, s := range servers {
		_, err := s.Conn.Exec("set global rpl_semi_sync_master_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
		_, err = s.Conn.Exec("set global rpl_semi_sync_slave_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
	}
	SaveMasterURL := master.URL

	for _, s := range slaves {
		dbhelper.StopSlave(s.Conn)
	}
	result, err := dbhelper.WriteConcurrent2(master.DSN, 10)
	if err != nil {
		logprintf("BENCH : %s %s", err.Error(), result)
	}
	dbhelper.InjectLongTrx(master.Conn, 10)
	time.Sleep(10 * time.Second)
	for _, s := range slaves {
		dbhelper.StartSlave(s.Conn)
	}
	logprintf("INFO :  Master is %s", master.URL)

	master.State = stateFailed
	conf.Interactive = false
	master.FailCount = conf.MaxFail
	conf.FailLimit = 5
	conf.FailTime = 0
	failoverCtr = 0
	conf.RplChecks = true
	conf.MaxDelay = 4
	checkfailed()

	wait_failover_end()
	logprintf("INFO : New Master  %s ", master.URL)

	time.Sleep(2 * time.Second)
	if master.URL != SaveMasterURL {
		logprintf("INFO : Old master %s ==  New master %s  ", SaveMasterURL, master.URL)

		return false
	}
	bootstrap()
	wait_failover_end()
	return true
}

func testFailOverNoRplChecksNoSemiSync() bool {
	conf.MaxDelay = 0
	bootstrap()
	wait_failover_end()

	logprintf("TESTING : Starting Test %s", "testFailOverNoRplChecksNoSemiSync")
	for _, s := range servers {
		_, err := s.Conn.Exec("set global rpl_semi_sync_master_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
		_, err = s.Conn.Exec("set global rpl_semi_sync_slave_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
	}
	SaveMasterURL := master.URL

	logprintf("INFO :  Master is %s", master.URL)
	master.State = stateFailed
	conf.Interactive = false
	master.FailCount = conf.MaxFail
	conf.FailLimit = 5
	conf.FailTime = 0
	failoverCtr = 0
	conf.RplChecks = false
	conf.MaxDelay = 4
	checkfailed()

	wait_failover_end()
	logprintf("INFO : New Master  %s ", master.URL)
	if master.URL == SaveMasterURL {
		logprintf("INFO : Old master %s ==  Next master %s  ", SaveMasterURL, master.URL)

		return false
	}

	return true
}

func testNumberFailOverLimitReach() bool {
	conf.MaxDelay = 0
	bootstrap()
	wait_failover_end()

	logprintf("TESTING : Starting Test %s", "testNumberFailOverLimitReach")
	for _, s := range servers {
		_, err := s.Conn.Exec("set global rpl_semi_sync_master_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
		_, err = s.Conn.Exec("set global rpl_semi_sync_slave_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
	}
	SaveMasterURL := master.URL

	logprintf("INFO :  Master is %s", master.URL)
	master.State = stateFailed
	conf.Interactive = false
	master.FailCount = conf.MaxFail
	conf.FailLimit = 3
	conf.FailTime = 0
	failoverCtr = 3
	conf.RplChecks = false
	conf.MaxDelay = 20
	checkfailed()

	wait_failover_end()
	logprintf("INFO : New Master  %s ", master.URL)
	if master.URL != SaveMasterURL {
		logprintf("INFO : Old master %s ==  Next master %s  ", SaveMasterURL, master.URL)

		return false
	}

	return true
}

func testFailOverTimeNotReach() bool {
	conf.MaxDelay = 0
	bootstrap()
	wait_failover_end()

	logprintf("TESTING : Starting Test %s", "testFailOverTimeNotReach")
	for _, s := range servers {
		_, err := s.Conn.Exec("set global rpl_semi_sync_master_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
		_, err = s.Conn.Exec("set global rpl_semi_sync_slave_enabled='OFF'")
		if err != nil {
			logprintf("TESTING : %s", err)
		}
	}
	SaveMasterURL := master.URL

	logprintf("INFO :  Master is %s", master.URL)
	master.State = stateFailed
	conf.Interactive = false
	master.FailCount = conf.MaxFail
	failoverTs = time.Now().Unix()
	conf.FailLimit = 3
	conf.FailTime = 20
	failoverCtr = 1
	conf.RplChecks = false
	conf.MaxDelay = 20
	checkfailed()

	wait_failover_end()
	logprintf("INFO : New Master  %s ", master.URL)
	if master.URL != SaveMasterURL {
		logprintf("INFO : Old master %s ==  Next master %s  ", SaveMasterURL, master.URL)

		return false
	}

	return true
}

func getTestResultLabel(res bool) string {
	if res == false {
		return "FAILED"
	} else {
		return "PASS"
	}
}

func runAllTests() bool {

	var allTests = map[string]string{}
	cleanall = true
	bootstrap()
	wait_failover_end()
	ret := true
	var res bool

	res = testSwitchOverLongTransactionNoRplCheckNoSemiSync()
	allTests["1 Switchover Concurrent Long Transaction <conf.ReadOnly=false> <conf.RplChecks=true>"] = getTestResultLabel(res)
	if res == false {
		ret = res
	}

	res = testSwitchOverLongQueryNoRplCheckNoSemiSync()
	allTests["1 Switchover Concurrent Long Query <conf.ReadOnly=false> <conf.RplChecks=true>"] = getTestResultLabel(res)
	if res == false {
		ret = res
	}

	res = testSwitchOverNoReadOnlyNoRplCheck()
	allTests["1 Switchover <conf.ReadOnly=false> <conf.RplChecks=false>"] = getTestResultLabel(res)
	if res == false {
		ret = res
	}

	res = testSwitchOverReadOnlyNoRplCheck()
	allTests["1 Switchover <conf.ReadOnly=true> <conf.RplChecks=false>"] = getTestResultLabel(res)
	if res == false {
		ret = res
	}

	res = testSwitchOver2TimesReplicationOkNoSemiSyncNoRplCheck()
	allTests["2 Switchover Replication Ok <2 threads benchmark> <semisync=false> <conf.RplChecks=false>"] = getTestResultLabel(res)
	if res == false {
		ret = res
	}

	res = testSwitchOver2TimesReplicationOkSemiSyncNoRplCheck()
	allTests["2 Switchover Replication Ok <2 threads benchmark> <semisync=true> <conf.RplChecks=false>"] = getTestResultLabel(res)
	if res == false {
		ret = res
	}

	res = testSwitchOverBackPreferedMasterNoRplCheckSemiSync()
	allTests["2 Switchover Back Prefered Master <semisync=true> <conf.RplChecks=false>"] = getTestResultLabel(res)
	if res == false {
		ret = res
	}

	res = testSwitchOverAllSlavesStopRplCheckNoSemiSync()
	allTests["Can't Switchover All Slaves Stop  <semisync=false> <conf.RplChecks=true>"] = getTestResultLabel(res)
	if res == false {
		ret = res
	}

	res = testSwitchOverAllSlavesStopNoSemiSyncNoRplCheck()
	allTests["Can Switchover All Slaves Stop <semisync=false> <conf.RplChecks=false>"] = getTestResultLabel(res)
	if res == false {
		ret = res
	}

	res = testSwitchOverAllSlavesDelayRplCheckNoSemiSync()
	allTests["Can't Switchover All Slaves Delay <semisync=false> <conf.RplChecks=true>"] = getTestResultLabel(res)
	if res == false {
		ret = res
	}

	res = testSwitchOverAllSlavesDelayNoRplChecksNoSemiSync()
	allTests["Can Switchover All Slaves Delay <semisync=false> <conf.RplChecks=false>"] = getTestResultLabel(res)
	if res == false {
		ret = res
	}

	res = testSlaReplAllSlavesStopNoSemiSync()
	allTests["SLA Decrease Can't Switchover All Slaves Stop <Semisync=false>"] = getTestResultLabel(res)
	if res == false {
		ret = res
	}

	res = testFailOverNoRplChecksNoSemiSync()
	allTests["1 Failover <conf.RplChecks=false> <Semisync=false> "] = getTestResultLabel(res)
	if res == false {
		ret = res
	}

	res = testFailOverAllSlavesDelayNoRplChecksNoSemiSync()
	allTests["1 Failover All Slave Delay <conf.RplChecks=false> <Semisync=false> "] = getTestResultLabel(res)
	if res == false {
		ret = res
	}

	res = testFailOverAllSlavesDelayRplChecksNoSemiSync()
	allTests["1 Failover All Slave Delay <conf.RplChecks=true> <Semisync=false> "] = getTestResultLabel(res)
	if res == false {
		ret = res
	}

	res = testNumberFailOverLimitReach()
	allTests["1 Failover Number of Failover Reach <conf.RplChecks=false> <Semisync=false> "] = getTestResultLabel(res)
	if res == false {
		ret = res
	}

	res = testFailOverTimeNotReach()
	allTests["1 Failover Before Time Limit <conf.RplChecks=false> <Semisync=false> "] = getTestResultLabel(res)
	if res == false {
		ret = res
	}

	keys := make([]string, 0, len(allTests))
	for key := range allTests {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, v := range keys {
		logprintf("TESTS : Result  %s -> %s", v, allTests[v])
	}

	cleanall = false
	return ret
}

func wait_failover_end() {
	for sme.IsInFailover() {
		time.Sleep(time.Second)
	}
	time.Sleep(recover_time * time.Second)
}
