package main

import (
	"context"
	"encoding/json"
	"errors"
	"go.etcd.io/etcd/clientv3"
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
	cli       *clientv3.Client
}

func NewEtcdService(conf *Config) *EtcdService {
	return &EtcdService{
		conf:      conf,
		prefixKey: "/unicorn_machine/",
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
	res, err := e.cli.Get(context.Background(), e.prefixKey, clientv3.WithPrefix())

	if err != nil {
		return nil, err
	}

	for _, kv := range res.Kvs {
		if data, err := e.extractMachineData(kv.Value); err != nil {
			continue
		} else {
			data["ip"] = string(kv.Key)
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

func (e *EtcdService) PutMachineId(ip string, id int) error {
	key := e.machineKey(ip)
	machineData, err := e.GetMachineId(ip)
	if err != nil {
		return err
	}

	if ValidStrId(machineData["id"]) {
		return errors.New("this machine already exists in etcd")
	}

	if !ValidId(id) {
		return errors.New("the machine id invalid")
	}

	value, err := e.makeMachineData(id)
	if err != nil {
		return err
	}

	if _, err = e.cli.Put(context.Background(), key, string(value)); err != nil {
		return err
	}

	return nil
}

func (e *EtcdService) machineKey(ip string) string {
	return e.prefixKey + ip
}

func (e *EtcdService) makeMachineData(id int) ([]byte, error) {
	data := make(map[string]string)
	ts := time.Now().Unix()
	tsStr := strconv.Itoa(int(ts))
	data["id"] = strconv.Itoa(id)
	data["created_timestamp"] = tsStr
	data["updated_timestamp"] = tsStr
	return json.Marshal(data)
}

func (e *EtcdService) extractMachineData(data []byte) (map[string]string, error) {
	result := make(map[string]string)
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func ValidId(id int) bool {
	if id < MinId || id > MaxId {
		return false
	}

	return true
}

func ValidStrId(id string) bool {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return false
	}

	if intId < MinId || intId > MaxId {
		return false
	}

	return true
}
