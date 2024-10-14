package main

import (
	"encoding/hex"
	"fmt"
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/storage/cache"
)

func printLogs(logCache *cache.LogCache) {
	for _, l := range *logCache {
		fmt.Println("log address: ", l.Address)
		for _, t := range l.Topics {
			fmt.Println("topic:", t)
		}
		fmt.Println("data bytes:", hex.EncodeToString(l.Data))
		fmt.Println("data as string:", string(l.Data))
	}
}

func printSlots(acc *environment.Account) {
	fmt.Println("address", acc.Address)
	for slot, val := range acc.Slots {
		fmt.Println("slot", slot.Int256(), " = value", val)
	}
}
