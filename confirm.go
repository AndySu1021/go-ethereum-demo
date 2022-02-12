package main

import (
	"context"
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

func main()  {
	var err error

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

	// Fetch data
	var blockLogs []*models.BlockLog
	now := time.Now().Unix()
	err = databases.MySqlClient.Where("available_at < ?", now).Find(&blockLogs).Error
	utils.CheckErr(err)

	fmt.Printf("Job count: %d \n", len(blockLogs))

	var numberList []uint64
	for _, blockLog := range blockLogs {
		blockNumber := (*blockLog).BlockNumber
		scanAndUpdate(blockNumber, client)
		numberList = append(numberList, blockNumber)
	}

	// Clear completed job
	err = databases.MySqlClient.Where("block_number IN ?", numberList).Delete(models.BlockLog{}).Error
	utils.CheckErr(err)
}

func scanAndUpdate(blockNumber uint64, client *ethclient.Client) {
	wg := new(sync.WaitGroup)
	var i uint64
	wg.Add(constants.NeedConfirmNumber + 1)
	for i = blockNumber - constants.NeedConfirmNumber + 1; i <= blockNumber; i++ {
		go fetchBlockAndUpdate(i, client, wg)
	}
	go fetchLogAndUpdate(blockNumber - constants.NeedConfirmNumber + 1, blockNumber, client, wg)
	wg.Wait()
}

func fetchBlockAndUpdate(blockNumber uint64, client *ethclient.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	utils.CheckErr(err)

	modelBlock := models.Block{
		Hash: block.Hash().Hex(),
		ParentHash: block.ParentHash().Hex(),
		IsPending: 2,
		Time: block.Time(),
	}

	//databases.MySqlClient.Exec("UPDATE block SET is_pending = 0 WHERE number = ?", blockNumber)
	databases.MySqlClient.Model(&models.Block{}).Where("number = ?", blockNumber).Updates(modelBlock)

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

	databases.MySqlClient.Where("block_number = ?", blockNumber).Delete(models.Transaction{})
	if len(transactions) > 0 {
		databases.MySqlClient.Create(&transactions)
	}
}

func fetchLogAndUpdate(start uint64, end uint64, client *ethclient.Client, wg *sync.WaitGroup) {
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

	var numberList []uint64
	for i := start; i <= end; i++ {
		numberList = append(numberList, i)
	}

	databases.MySqlClient.Where("block_number IN ?", numberList).Delete(models.TransactionLog{})
	databases.MySqlClient.Create(&logs)
}
