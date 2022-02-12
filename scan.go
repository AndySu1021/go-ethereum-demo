package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"go-ethereum-demo/config"
	"go-ethereum-demo/constants"
	"go-ethereum-demo/databases"
	"go-ethereum-demo/models"
	"go-ethereum-demo/utils"
	"math/big"
	"sync"
	"time"
)

const batchSize = 100

type args struct {
	startBlockNumber uint64
}

func getArgs() (args, error) {
	startBlockNumber := flag.Uint64("n", 16683000, "Start block number")
	flag.Parse()

	return args{*startBlockNumber}, nil
}

var EndNumber uint64

func main() {
	var err error

	// Showing useful information when the user enters the --help option
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] \nOptions:\n", "./scan")
		flag.PrintDefaults()
	}

	// Get the arguments that was entered by the user
	userArgs, err := getArgs()
	utils.CheckErr(err)

	// Init viper config
	err = config.InitConfig()
	utils.CheckErr(err)

	// Init database
	err = databases.InitMySql()
	utils.CheckErr(err)
	defer databases.Close()

	// Get rpc client
	rpcHost := viper.GetString("rpc.host")
	client, err := utils.GetRpcClient(rpcHost)
	utils.CheckErr(err)

	// Get current block number
	endNumber, err := client.BlockNumber(context.Background())
	EndNumber = endNumber
	utils.CheckErr(err)

	// Get start block number for scanning
	startNumber := getStartNumber(userArgs.startBlockNumber)
	fmt.Printf("Start: %d, End: %d", startNumber, endNumber)

	if startNumber > endNumber {
		return
	}

	remainder := (endNumber - startNumber + 1) % batchSize
	for num := startNumber; num <= endNumber - remainder; num += batchSize {
		batchScan(num, client, 0)
	}

	batchScan(endNumber - remainder + 1, client, endNumber)

	databases.MySqlClient.Create(&models.BlockLog{
		BlockNumber: endNumber,
		AvailableAt: uint64(time.Now().Add(300 * time.Second).Unix()),
	})
}

func getStartNumber(startBlockNumber uint64) uint64 {
	var startNumber uint64
	var block models.Block
	err := databases.MySqlClient.Order("number desc").Limit(1).Find(&block).Error
	utils.CheckErr(err)
	if block.Number != 0 {
		startNumber = block.Number + 1
	} else {
		startNumber = startBlockNumber
	}

	return startNumber
}

func batchScan(start uint64, client *ethclient.Client, end uint64) {
	wg := new(sync.WaitGroup)
	var i uint64
	if end == 0 {
		end = start + batchSize - 1
	}
	wg.Add(int(end - start + 2))
	for i = start; i <= end; i++ {
		go fetchBlock(i, client, wg)
	}
	go fetchLog(start, end, client, wg)
	wg.Wait()
}

func fetchBlock(blockNumber uint64, client *ethclient.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	utils.CheckErr(err)

	blockHash := block.Hash()
	modelBlock := models.Block{
		Number: blockNumber,
		Hash: blockHash.Hex(),
		ParentHash: block.ParentHash().Hex(),
		IsPending: 2,
		Time: block.Time(),
	}

	if blockNumber > EndNumber - constants.NeedConfirmNumber {
		modelBlock.IsPending = 1
	}

	databases.MySqlClient.Create(&modelBlock)

	var transactions []models.Transaction
	for _, transaction := range block.Transactions() {
		tmpTx := models.Transaction{
			BlockNumber: blockNumber,
			Hash: transaction.Hash().Hex(),
			From: "",
			To: "",
			Nonce: transaction.Nonce(),
			Data: "",
			Value: transaction.Value().String(),
		}

		if transaction.To() != nil {
			tmpTx.To = transaction.To().Hex()
		}

		baseFee := big.NewInt(100)
		chainID, err := client.ChainID(context.Background())
		utils.CheckErr(err)

		if msg, err := transaction.AsMessage(types.NewEIP155Signer(chainID), baseFee); err == nil {
			tmpTx.From = msg.From().Hex()
		}

		transactions = append(transactions, tmpTx)
	}

	if len(transactions) > 0 {
		databases.MySqlClient.Create(&transactions)
	}
}

func fetchLog(start uint64, end uint64, client *ethclient.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(start)),
		ToBlock: big.NewInt(int64(end)),
	}

	filteredLogs, err := client.FilterLogs(context.Background(), query)
	utils.CheckErr(err)

	var logs []models.TransactionLog
	for _, logData := range filteredLogs {
		tmpLog := models.TransactionLog{
			BlockNumber: logData.BlockNumber,
			TxHash: logData.TxHash.Hex(),
			Index:  logData.Index,
			Data:   "",
		}

		logs = append(logs, tmpLog)
	}

	databases.MySqlClient.Create(&logs)
}
