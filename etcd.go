package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"hash/crc32"
	"strconv"
	"strings"
	"time"
)

const (
	MinId = 1
	MaxId = 1 << 10
)

type EtcdService struct {
	conf      *Config
	prefixKey string
	idListKey string
	cli       *clientv3.Client
}

func NewEtcdService(conf *Config) *EtcdService {
	return &EtcdService{
		conf:      conf,
		prefixKey: "/unicorn_machine/",
		idListKey: "/unicorn_machine_ids",
	}
}

func (e *EtcdService) connection() error {
	cfg := clientv3.Config{
		Endpoints: strings.Split(e.conf.Etcd.Cluster, ","),
	}

	if cli, err := clientv3.New(cfg); err != nil {
		return err
	} else {
		e.cli = cli
		return nil
	}
}

func (e *EtcdService) GetMachineList() ([]map[string]string, error) {
	result := make([]map[string]string, 0)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	res, err := e.cli.Get(ctx, e.prefixKey, clientv3.WithPrefix())

	if err != nil {
		return nil, err
	}

	for _, kv := range res.Kvs {
		if data, err := e.extractMachineData(kv.Value); err != nil {
			continue
		} else {
			result = append(result, data)
		}
	}

	return result, nil
}

func (e *EtcdService) GetMachineId(key string) (map[string]string, error) {
	res, err := e.cli.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}

	for _, kv := range res.Kvs {
		if string(kv.Key) == key {
			return e.extractMachineData(kv.Value)
		}
	}

	return nil, nil
}

func (e *EtcdService) PutMachineId(ip, id string) (map[string]string, error) {
	key := e.MachineKey(ip, id)
	machineData, err := e.GetMachineId(key)
	if err != nil {
		return nil, err
	}

	if machineData != nil {
		return nil, errors.New("the machine ip already exists in etcd")
	}

	if machineData["id"] == id {
		return nil, errors.New("the machine id has been used in etcd")
	}

	if !ValidId(id) {
		return nil, errors.New("the machine id invalid")
	}

	newMachineData := make(map[string]string)
	ts := time.Now().Unix()
	tsStr := strconv.Itoa(int(ts))
	newMachineData["id"] = id
	newMachineData["ip"] = ip
	newMachineData["created_timestamp"] = tsStr
	newMachineData["updated_timestamp"] = tsStr
	newMachineDataByte, err := e.makeMachineData(newMachineData)
	if err != nil {
		return nil, err
	}

	if _, err = e.cli.Put(context.Background(), key, string(newMachineDataByte)); err != nil {
		return nil, err
	}

	return newMachineData, nil
}

func (e *EtcdService) MachineKey(ip, id string) string {
	key := e.prefixKey + strconv.Itoa(int(crc32.ChecksumIEEE([]byte(ip+"_"+id))))
	fmt.Println(key)
	return key
}

func (e *EtcdService) makeMachineData(data map[string]string) ([]byte, error) {
	return json.Marshal(data)
}

func (e *EtcdService) extractMachineData(data []byte) (map[string]string, error) {
	result := make(map[string]string)
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func ValidId(id string) bool {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return false
	}

	if intId < MinId || intId > MaxId {
		return false
	}

	return true
}
