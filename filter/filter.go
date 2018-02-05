package filter

import (
	"encoding/json"
	// "fmt"
	//"regexp"
)

type Beat struct {
	Name     string
	Hostname string
	Version  string
}

type Log struct {
	Tiimestamp string
	Metadata   string
	Offset     int64
	Message    string
	Source     string
	Prospector string
	Fields     string
	Beat       Beat
}

func JsontoStr(b []byte) Log {
	var log Log
	json.Unmarshal(b, &log)
	return log
}

/*
func LogReg(pattern ,s ,logfilepattern ,logfile,string) (b bool) {
	b, _ := regexp.MatchString(pattern, s)
	f, _ := regexp.MatchString(logfilepattern, logfile)
	if b && f {
		return true
	}else{
		return false
	}

}
*/
