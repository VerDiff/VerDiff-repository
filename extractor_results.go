package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"strconv"
)

type LogData struct {
	// note verdiffCollector - to log: Type, Blockhight, Contract address, Slot, Value
	LogType         string
	BlockHeight     string
	ContractAddress string
	Slot            string
	SlotValue       string
}

// {sstore 4766792 0x29d38FdF26d64Fa799276e6615759D27dB1F1fcD 0xa5baec7d73105a3c7298203bb205bbc41b63fa384ae73a6016b890a7ca29ae2f 0x0}
// {sstore 4766790 0x06012c8cf97BEaD5deAe237070F9587f8E7A266d 0xc7b975ddb51e69aba9319314355259a765f41a15409d2610a7f6eaa7cd0ac7aa 0x1}
// {selfdestruct 51921 0x412FdA7643b37d436cB40628F6DbBB80a07267ed 0x0 0x0}

type SlotStruct struct {
	Slot             map[string]string
	SelfdestructFlag bool
}

var contractUnzeroSlots map[string]SlotStruct
var target_BlockHeight int
var suicidedContracts map[string]bool

func parseInLoop() {
	var compressed_file_path = "./compressed_erigon_slotdata/"
	compressed_files, err := os.ReadDir(compressed_file_path)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	fmt.Println("Files have be compressed:", len(compressed_files))

	count := 0
	for _, one_compressed_file := range compressed_files {
		count++
		if count == 50 {
			break
		}
		one_compressed_filename := one_compressed_file.Name()
		one_compressed_file_path := compressed_file_path + one_compressed_filename
		parseOneFile(one_compressed_file_path)
	}
}

func checkKeyExists(key string) (SlotStruct, bool) {
	value, exists := contractUnzeroSlots[key]
	return value, exists
}

func parseOneFile(filename string) {
	var data = make([]LogData, 1e5+10) //1e3+10) // ! logdata size - 1e5

	dataFile, err := os.Open(filename)

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
		one_LogType := elem.LogType
		one_BlockHeight := elem.BlockHeight
		one_ContractAddress := elem.ContractAddress
		//one_Slot := elem.Slot
		//one_SlotValue := elem.SlotValue

		blockHeight_Int, err := strconv.Atoi(one_BlockHeight)
		if err != nil {
			fmt.Printf("Error converting one_BlockHeight to int: %v\n", err)
		}
		if blockHeight_Int > target_BlockHeight {
			target_BlockHeight = ((blockHeight_Int / 50000) + 1) * 50000

			fmt.Printf("Current BlockHeight: %d\n", blockHeight_Int)

			//fmt.Printf("Total contracts: %d\n", len(contractUnzeroSlots))

			totalSlots := 0
			selfDestructSlots := 0
			totalSuicideContracts := 0
			for _, slotStruct := range contractUnzeroSlots {
				if slotStruct.SelfdestructFlag {
					totalSuicideContracts = totalSuicideContracts + 1
				}
				totalSlots += len(slotStruct.Slot)

				if slotStruct.SelfdestructFlag {
					selfDestructSlots += len(slotStruct.Slot)
				}
			}
			//fmt.Println("Total slots: ", totalSlots)
			//fmt.Println("Selfdestruct slots: ", selfDestructSlots)

			logFile := "./output.result"
			file, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println("Error opening log file:", err)
				return
			}
			defer file.Close()
			fmt.Fprintf(file, "%d,%d,%d,%d,%d\n", blockHeight_Int, len(contractUnzeroSlots), totalSlots, selfDestructSlots, totalSuicideContracts)

		}

		/* if _, exists := suicidedContracts[one_ContractAddress]; !exists {
			continue
		} */

		switch one_LogType {
		/* case "sstore":
		if one_UnzeroSlots, exists := checkKeyExists(one_ContractAddress); exists {
			if _, slotExists := one_UnzeroSlots.Slot[one_Slot]; slotExists {
				if one_SlotValue == "0x0" {
					delete(one_UnzeroSlots.Slot, one_Slot)
				}
			} else {
				if one_SlotValue == "0x1" {
					one_UnzeroSlots.Slot[one_Slot] = one_SlotValue
				}
			}
			contractUnzeroSlots[one_ContractAddress] = one_UnzeroSlots
		} else {
			if one_SlotValue == "0x1" {
				contractUnzeroSlots[one_ContractAddress] = SlotStruct{
					Slot:             map[string]string{one_Slot: one_SlotValue},
					SelfdestructFlag: false,
				}
			}
		} */
		case "selfdestruct":
			contractUnzeroSlots[one_ContractAddress] = SlotStruct{
				Slot:             map[string]string{"": ""},
				SelfdestructFlag: true,
			}

			/* if slotStruct, exists := contractUnzeroSlots[one_ContractAddress]; exists {
				slotStruct.SelfdestructFlag = true
				contractUnzeroSlots[one_ContractAddress] = slotStruct
			} */
		}
	}
}

func main() {
	logFile := "./output.result"
	file, err := os.Create(logFile)
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return
	}
	file.Close()
	file, err = os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer file.Close()
	fmt.Fprintf(file, "BlockHeight,NumContracts,NumTotalSlots,NumSelfdestructSlots,NumSuicideContracts\n")

	contractUnzeroSlots = make(map[string]SlotStruct)
	target_BlockHeight = 50000

	suicidedContracts = make(map[string]bool)

	parseInLoop()

	suicidedContracts_File := "./temp_suicidedContracts.result"
	file2, err := os.Create(suicidedContracts_File)
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return
	}
	file2.Close()
	file2, err = os.OpenFile(suicidedContracts_File, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer file2.Close()
	fmt.Fprintf(file2, "Suicded_contracts_addresses\n")

	for contract, slotStruct := range contractUnzeroSlots {
		if slotStruct.SelfdestructFlag {
			fmt.Fprintf(file2, "%s\n", contract)
		}
	}
}
