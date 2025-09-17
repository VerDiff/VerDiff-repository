package main

import (
	"encoding/gob"
	"fmt"
	"os"
)

type SWaitTimeData struct {
	// note verdiffCollector - to log: Type, Blockhight, Contract address, Slot, Value
	LogType         string
	BlockHeight     string
	ContractAddress string
	Slot            string
	SlotValue       string
}

func GobFinalize(retDir string, filename string, compresseddata []SWaitTimeData) {
	dataFile, err := os.Create(retDir + filename)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//os.Exit(0)
	// serialize the data
	dataEncoder := gob.NewEncoder(dataFile)
	dataEncoder.Encode(compresseddata)
	dataFile.Close()
}

func compressSstoreSlotValueHex(input string) string {
	if input == "0x0" {
		return "0x0"
	}
	return "0x1"
}

func compressOneFile(source_path string, target_path string, filename string) {
	var data = make([]SWaitTimeData, 1e5+10)      //1e3+10) // ! verdiffCollector - 1e5
	var compresseddata = make([]SWaitTimeData, 0) //1e3+10) // ! verdiffCollector - 1e5

	// open data file
	// dataFile, err := os.Open("/media/data/verdiffCollector/slotdata/erigon_slotdata/erigon-slots/1720160782851601")
	dataFile, err := os.Open(source_path + filename)
	retDir := target_path

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(&data)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dataFile.Close()
	//count := 0
	for _, elem := range data {
		//count += 1
		//if count%10000 == 0 {
		//	fmt.Println("Line Done", count)
		//}
		// fmt.Println(elem)
		/*fmt.Println(elem.LogType)         // string
		fmt.Println(elem.BlockHeight)     // uint string
		fmt.Println(elem.ContractAddress) // hex string
		fmt.Println(elem.Slot)            // hex string
		fmt.Println(elem.SlotValue)       // hex string*/

		switch elem.LogType {
		case "sstore":
			elem.SlotValue = compressSstoreSlotValueHex(elem.SlotValue)
		case "sload":
			continue // ! continue
		}
		// fmt.Println("compressed", elem)
		compressed_elem := SWaitTimeData{
			LogType:         elem.LogType,
			BlockHeight:     elem.BlockHeight,
			ContractAddress: elem.ContractAddress,
			Slot:            elem.Slot,
			SlotValue:       elem.SlotValue,
		}
		compresseddata = append(compresseddata, compressed_elem)

		//fmt.Println(compresseddata)
		//os.Exit(0)
	}
	GobFinalize(retDir, filename+".compress", compresseddata)
}

func compressInLoop() {
	var source_path = "/media/data/verdiffCollector/slotdata/backup/slots/"
	var target_path = "/media/data/verdiffCollector/slotdata/backup/compressed_slots/"

	uncompressed_files, err := os.ReadDir(source_path)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}
	compressed_files, err := os.ReadDir(target_path)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	compressed_fileDict := make(map[string]bool)

	for _, file := range compressed_files {
		compressed_fileDict[file.Name()[:len(file.Name())-9]] = true
	}
	fmt.Println("Files to be compressed:", len(uncompressed_files))
	fmt.Println("Files have be compressed:", len(compressed_fileDict))

	var count = 0
	for _, one_uncompressed_file := range uncompressed_files {
		one_uncompressed_filename := one_uncompressed_file.Name()
		count += 1
		//if count%2 == 0 {
		//	fmt.Println("File Done", count)
		//}

		if _, exists := compressed_fileDict[one_uncompressed_filename]; exists {
			continue
		}

		fmt.Println("File Done", one_uncompressed_filename)
		compressOneFile(source_path, target_path, one_uncompressed_filename)
	}
}

func main() {
	//compressInLoop()

	var data = make([]SWaitTimeData, 1e5+10) //1e3+10) // ! verdiffCollector - 1e5

	// open data file
	// dataFile, err := os.Open("/media/data/verdiffCollector/slotdata/erigon_slotdata/erigon-slots/1720160782851601")
	dataFile, err := os.Open("/media/data/verdiffCollector/slotdata/compressed_erigon_slotdata/1720279985246053.compress")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(&data)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dataFile.Close()
	count := 0
	for _, elem := range data {
		count += 1
		if elem.LogType == "selfdestruct" {
			fmt.Println(elem)
			os.Exit(0)
		}
		fmt.Println(elem)
	}
}
