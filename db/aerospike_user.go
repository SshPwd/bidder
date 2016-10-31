package db

import (
	"fmt"
	aerospike "github.com/aerospike/aerospike-client-go"
	"golang_bidder/djb2"
)

type SyncGuids struct {
	Uid  string
	IpUa string
	Ip   string
}

const (
	userMapExpire = 10 * 24 * 60 * 60
)

func SaveUserCombination(seatId int, userId, userAgent, guid []byte, ipStr string) {

	return // FIXME

	userAgentHash := djb2.Sum32Bytes(userAgent)

	currentGuid := getUser(seatId, userId, userAgentHash, ipStr)
	saveUser(currentGuid, seatId, userId, userAgentHash, guid, ipStr)
}

func saveUser(currentGuid SyncGuids, seatId int, userId []byte, userAgentHash uint32, guid []byte, ipStr string) {

	return // FIXME

	policy := aerospike.NewWritePolicy(1, userMapExpire)
	guidStr := string(guid)
	guidBin := aerospike.NewBin("guid", guidStr)

	var err error
	var key *aerospike.Key
	var keyStr string

	if currentGuid.Uid == "" || currentGuid.Uid != guidStr {

		keyStr = fmt.Sprintf("uid_%d_%s", seatId, userId)
		key, err = aerospike.NewKey("bidder", "user_map", keyStr)
		if err == nil {
			client.PutBins(policy, key, guidBin)
		}
	}

	if currentGuid.IpUa == "" || currentGuid.IpUa != guidStr {

		keyStr = fmt.Sprintf("ipua_%s_%x", ipStr, userAgentHash)
		key, err = aerospike.NewKey("bidder", "user_map", keyStr)
		if err == nil {
			client.PutBins(policy, key, guidBin)
		}
	}

	if currentGuid.Ip == "" || currentGuid.Ip != guidStr {

		keyStr = fmt.Sprintf("ip_%s", ipStr)
		key, err = aerospike.NewKey("bidder", "user_map", keyStr)
		if err == nil {
			client.PutBins(policy, key, guidBin)
		}
	}
}

func readFieldString(bins aerospike.BinMap, fieldName string) (string, bool) {

	if bins != nil {

		if fieldData, ok := bins[fieldName]; ok {

			value, ok := fieldData.(string)
			return value, ok
		}
	}
	return "", false
}

func getUser(seatId int, userId []byte, userAgentHash uint32, ipStr string) (currentGuid SyncGuids) {

	return // FIXME

	var err error
	var key *aerospike.Key
	keys := make([]*aerospike.Key, 0, 3)

	keyStr := fmt.Sprintf("uid_%d_%s", seatId, userId)
	key, err = aerospike.NewKey("bidder", "user_map", keyStr)
	if err == nil {
		keys = append(keys, key)
	}

	keyStr = fmt.Sprintf("ipua_%s_%x", ipStr, userAgentHash)
	key, err = aerospike.NewKey("bidder", "user_map", keyStr)
	if err == nil {
		keys = append(keys, key)
	}

	keyStr = fmt.Sprintf("ip_%s", ipStr)
	key, err = aerospike.NewKey("bidder", "user_map", keyStr)
	if err == nil {
		keys = append(keys, key)
	}

	if results, err := client.BatchGet(nil, keys); err == nil {

		for n := range results {

			fmt.Println(CM+"<<<", keys[n].Value(), CN)

			if results[n] != nil {
				fmt.Println(CG+">>>", results[n].Bins, CN)
			} else {
				fmt.Println(CG+">>>", results[n], CN)
			}
		}

		if results[0] != nil {
			if guid, ok := readFieldString(results[0].Bins, "guid"); ok {
				currentGuid.Uid = guid
			}
		}

		if results[1] != nil {
			if guid, ok := readFieldString(results[1].Bins, "guid"); ok {
				currentGuid.IpUa = guid
			}
		}

		if results[2] != nil {
			if guid, ok := readFieldString(results[2].Bins, "guid"); ok {
				currentGuid.Ip = guid
			}
		}
	}

	// fmt.Println("currentGuid:", currentGuid)
	return
}
