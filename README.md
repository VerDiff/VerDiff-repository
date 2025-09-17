## VerDiff

This repository contains data and code for our study on uncovering correctness issues in blockchain state tries.

### Environment Setup
- A 64-bit Linux operating system (e.g., Ubuntu 20.04) with at least 1 TB of free disk space for storing our experimental data, and preferably more than 8 TB for storing blockchain node data.
- Python 3.8+,
  Package dependencies for Python:
  - [icecream](https://github.com/gruns/icecream)
  - [matplotlib](https://matplotlib.org/)
- Basic environment for building [Go-ethereum](https://github.com/ethereum/go-ethereum)
  You need to modify your go runtime following the instructions like [here](https://github.com/dvyukov/go-fuzz/issues/354), otherwise, you may encounter unknown errors.
- Basic environment for building [Nethermind](https://github.com/NethermindEth/nethermind)
- Basic environment for building [EthereumJS](https://github.com/ethereumjs/ethereumjs-monorepo)

### Basic fuzzing dependencies
- Install go-fuzz and go-fuzz-build binaries:
  ```bash
  go install -v github.com/dvyukov/go-fuzz/go-fuzz@latest
  go install -v github.com/dvyukov/go-fuzz/go-fuzz-build@latest
  ```
- Add the go-fuzz and go-fuzz-build binaries to your PATH.

### Data preparation
We have provided intermediate data results in the `data` folder. If you want to reproduce the result from scratch, we also provide code for this purpose. Please follow the steps below.

- Step 1: Build the data collector
  
  We provide code in the `collector` folder for building the data collector tool. In this tool, we instrument and log EVM instructions like `SELFDESTRUCT` and `SSTORE` to collect basic data. For example, we log the Blockhight, Contract address, Slot, and Value for `SSTORE` instructions, e.g., `common.GlobalSlotWatcher.ChWaitTimes <- *verdiffCollectorslot.NewSLogData(log_type, log_blockheight, log_contractaddress, log_slot, log_slotvalue)`. In this tool, we use  `encoding/gob` to serialize the collected data and store them in `compressed_erigon_slotdata` folder. 
  ```go
  func opSstore(pc *uint64, interpreter *EVMInterpreter, scope *ScopeContext) ([]byte, error) {
    if interpreter.readOnly {
      return nil, ErrWriteProtection
    }
    loc := scope.Stack.Pop()
    val := scope.Stack.Pop()
    interpreter.hasherBuf = loc.Bytes32()
    interpreter.evm.IntraBlockState().SetState(scope.Contract.Address(), &interpreter.hasherBuf, val)

    // * common.Globalmevwatcher.ChWaitTimes <- *verdiffCollectormev.NewSWaitTimeData(tx.Hash().String(), time.Now())
    // note verdiffCollector - to log: Type, Blockhight, Contract address, Slot, Value
    // * hexLocStr := fmt.Sprintf("0x%x", loc)
    // * log.Println(hexLocStr)
    // * strconv.FormatUint(num, 10)
    log_type := "sstore"
    log_blockheight := strconv.FormatUint(interpreter.evm.Context.BlockNumber, 10)
    log_contractaddress := scope.Contract.Address().String()
    log_slot := loc.Hex()
    log_slotvalue := val.Hex()
    common.GlobalSlotWatcher.ChWaitTimes <- *verdiffCollectorslot.NewSLogData(log_type, log_blockheight, log_contractaddress, log_slot, log_slotvalue)

    return nil, nil
  }
  ```
  You can modify the path for storing the collected data in `./collector/verdiffCollector-slot/watcher.go`.
  ```go
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
  ```
- Step 2: Data collection
  
  The tool is built on top of [erigon](https://github.com/erigontech/erigon), and collects data during the normal node synchronization process. You may need to upgrade the main branch of erigon to the latest stable version to ensure the capability for blockchain synchronization.
  - This process may take a long time (e.g., several weeks) depending on your hardware and network conditions. 
  - You need to prepare sufficient disk space (e.g., more than 8 TB) to store the blockchain data and the collected data. The collected data may take more than 1 TB of disk space.
- Step 3: Data extraction
  
  We provide code in the `extractor_results.go` and `extractor_contracts.go` for extracting the collected data. The extracted data are used for further analysis and verification. You can check your results with our intermediate data in the `data` folder.


### Victim contracts with global dynamic array

The victim contracts threaten by the first class of issues are listed in `GlobalDynamicArrayContracts.md`. These contracts are derived from the popular contracts dataset in `dataset/dataset_2_popular_contracts.json`, which originates from the [RNVulDet dataset](https://github.com/Messi-Q/RNVulDet). We gratefully acknowledge the authors of RNVulDet for providing this resource.


### Reproduction for results of interesting state accesses

Enter into the root path, and run the scripts in `diverse_result.py`.
  ```bash
  python3 diverse_result.py
  ```

### Reproduction for results of different initial seed sizes

Enter into the root path, and run the scripts in `time_cost_of_init_corpus_sizes.py`.
  ```bash
  python3 time_cost_of_init_corpus_sizes.py
  ```

### Reproduction for results of trends of selfdestruct contracts

Enter into the root path, and run the scripts in `selfdestruct_contracts.py`.
  ```bash
  python3 selfdestruct_contracts.py
  ```

### Notice
This repository was prepared in a hurry and may raise unexpected errors during use. Please feel free to contact me if you have any questions.
