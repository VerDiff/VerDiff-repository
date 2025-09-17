package verdiffCollectorslot

import (
	"encoding/gob"
	"fmt"
	"os"
	"strconv"
	"time"
)

type SWatcher struct {
	ChWaitTimes    chan SWaitTimeData
	WaitTimeRecord []SWaitTimeData
	RetDir         string
}

func NewSWatcher() *SWatcher {
	w := &SWatcher{
		ChWaitTimes:    make(chan SWaitTimeData, 5*1e6),  // ! verdiffCollector - 1e6
		WaitTimeRecord: make([]SWaitTimeData, 0, 1e5+10), // 1e3+10),  // ! verdiffCollector - 1e5
		RetDir:         "./compressed_erigon_slotdata/",
	}
	err := os.MkdirAll(w.RetDir, os.ModePerm)
	if err != nil {
		fmt.Printf("dir create failed  -> %v\n", err)
		os.Exit(17)
	} else {
		fmt.Println("dir create success!")
	}
	go w.Output()
	return w
}

type SWaitTimeData struct {
	// note verdiffCollector - to log: Type, Blockhight, Contract address, Slot, Value
	LogType         string
	BlockHeight     string
	ContractAddress string
	Slot            string
	SlotValue       string
}

func NewSLogData(logtype string, blockheight string, contractaddress string, slot string, slotvalue string) *SWaitTimeData {
	waittime := &SWaitTimeData{
		LogType:         logtype,
		BlockHeight:     blockheight,
		ContractAddress: contractaddress,
		Slot:            slot,
		SlotValue:       slotvalue,
	}
	return waittime
}

func (watcher *SWatcher) Output() {
	counter := 0
	filenamecounter := 0
	for {
		WaitTimeinfo, ok := <-watcher.ChWaitTimes
		watcher.WaitTimeRecord = append(watcher.WaitTimeRecord, WaitTimeinfo)
		if !ok {
			// * strconv.FormatInt(num, 10)
			// * strconv.Itoa(time.Now().UnixMicro())
			watcher.GobFinalize(strconv.FormatInt(time.Now().UnixMicro(), 10))
			break
		}
		counter++
		if counter == 1e5 { //1e3 { // ! verdiffCollector - 1e5
			// todo name the file with the nano time
			watcher.GobFinalize(strconv.FormatInt(time.Now().UnixMicro(), 10))
			// * verdiffCollector - set how frequently the data is written to disk
			watcher.WaitTimeRecord = make([]SWaitTimeData, 0, 1e5+10) // 1e3+10) // ! verdiffCollector - 1e5
			counter = 0
			filenamecounter++
		}
	}

}

func (watcher *SWatcher) GobFinalize(filename string) {
	dataFile, err := os.Create(watcher.RetDir + filename)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//os.Exit(0)
	// serialize the data
	dataEncoder := gob.NewEncoder(dataFile)
	dataEncoder.Encode(watcher.WaitTimeRecord)
	dataFile.Close()
}
