package main

import (
	"context"
	"encoding/json"
	"errors"
	"go.etcd.io/etcd/clientv3"
	"hash/crc32"
	"strings"
	"time"
)

const (
	MinId = 1
	MaxId = 1 << 10
)

type MachineNode struct {
	Id               int    `json:"id"`
	Ip               string `json:"ip"`
	CreatedTimestamp int64  `json:"created_timestamp"`
	UpdatedTimestamp int64  `json:"updated_timestamp"`
}

func NewMachineNode() *MachineNode {
	ts := time.Now().Unix()
	return &MachineNode{
		CreatedTimestamp: ts,
		UpdatedTimestamp: ts,
	}
}

func (node *MachineNode) withIpId(ip string, id int) *MachineNode {
	node.Ip = ip
	node.Id = id
	return node
}

func (node *MachineNode) encode() ([]byte, error) {
	return json.Marshal(node)
}

type MachineData struct {
	Nodes     map[uint32]*MachineNode `json:"nodes"`
	ExtraData *MachineExtraData       `json:"extra"`
}

func NewMachineData() *MachineData {
	return &MachineData{
		Nodes:     make(map[uint32]*MachineNode),
		ExtraData: NewMachineExtraData(),
	}
}

func (data *MachineData) getNode(ip string) (*MachineNode, bool) {
	key := crc32.ChecksumIEEE([]byte(ip))
	node, ok := data.Nodes[key]
	return node, ok
}

func (data *MachineData) setNode(ip string) {
	key := crc32.ChecksumIEEE([]byte(ip))
	data.ExtraData.LastId++
	node := NewMachineNode().withIpId(ip, data.ExtraData.LastId)
	data.Nodes[key] = node
}

func (data *MachineData) setNodeById(ip string, id int) {
	key := crc32.ChecksumIEEE([]byte(ip))
	node := NewMachineNode().withIpId(ip, id)
	data.Nodes[key] = node
}

func (data *MachineData) JsonMarshal() ([]byte, error) {
	return json.Marshal(data)
}

type MachineExtraData struct {
	FreeSlots [MaxId]int `json:"free_slots"`
	LastId    int        `json:"last_id"`
}

func NewMachineExtraData() *MachineExtraData {
	return &MachineExtraData{
		LastId: 1,
	}
}

func (extra *MachineExtraData) findFreeSlot() int {
	for i := 1; i < MaxId; i++ {
		if extra.FreeSlots[i] == 0 {
			return i
		}
	}

	return 0
}

func (extra *MachineExtraData) updateFreeSlot(index int) {
	extra.FreeSlots[index] = 1
}

func (extra *MachineExtraData) updateLastId() {
	extra.LastId = extra.LastId + 1
}

func (extra *MachineExtraData) encode() ([]byte, error) {
	return json.Marshal(extra)
}

type EtcdService struct {
	conf      *Config
	prefixKey string
	extraKey  string
	cli       *clientv3.Client
}

func NewEtcdService(conf *Config) *EtcdService {
	return &EtcdService{
		conf:      conf,
		prefixKey: "/unicorn_machine_data",
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

func (e *EtcdService) GetMachineData() (*MachineData, error) {
	res, err := e.cli.Get(context.Background(), e.prefixKey)
	if err != nil {
		return nil, err
	}

	for _, kv := range res.Kvs {
		if string(kv.Key) == e.prefixKey {
			return e.extractMachineData(kv.Value)
		}
	}

	return e.PutMachineData(NewMachineData())
}

func (e *EtcdService) PutMachineData(data *MachineData) (*MachineData, error) {
	b, err := data.JsonMarshal();
	if err != nil {
		return nil, err
	}

	if _, err = e.cli.Put(context.Background(), e.prefixKey, string(b)); err != nil {
		return nil, err
	}
	return data, nil
}

func (e *EtcdService) PutMachineNode(ip string) (*MachineData, error) {
	data, err := e.GetMachineData()
	if err != nil {
		return nil, err
	}

	_, ok := data.getNode(ip);
	if ok {
		return nil, errors.New("the machine ip already exists in etcd")
	}

	if data.ExtraData.LastId >= MaxId {
		if index := data.ExtraData.findFreeSlot(); index == 0 {
			return nil, errors.New("no machine id available")
		} else {
			data.setNodeById(ip, index)
			data.ExtraData.updateFreeSlot(index)
			return e.PutMachineData(data)
		}
	} else {
		data.setNode(ip)
		data.ExtraData.updateFreeSlot(data.ExtraData.LastId)
		return e.PutMachineData(data)
	}
}

func (e *EtcdService) GetMachineNode(ip string) (*MachineNode, error) {
	data, err := e.GetMachineData()
	if err != nil {
		return nil, err
	}

	node, ok := data.getNode(ip)
	if !ok {
		return nil, errors.New("no machine node info found by ip")
	}

	return node, nil
}

func (e *EtcdService) extractMachineData(data []byte) (*MachineData, error) {
	machineData := &MachineData{}
	if err := json.Unmarshal(data, machineData); err != nil {
		return nil, err
	}
	return machineData, nil
}

//func (e *EtcdService) GetMachineList() ([]map[string]string, error) {
//	b, _ := json.Marshal(&MachineExtraData{
//	})
//	fmt.Println(string(b))
//	result := make([]map[string]string, 0)
//	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
//	res, err := e.cli.Get(ctx, e.prefixKey, clientv3.WithPrefix())
//
//	if err != nil {
//		return nil, err
//	}
//
//	for _, kv := range res.Kvs {
//		if data, err := e.extractMachine(kv.Value); err != nil {
//			continue
//		} else {
//			result = append(result, data)
//		}
//	}
//
//	return result, nil
//}
//
//func (e *EtcdService) GetMachineId(key string) (map[string]string, error) {
//	res, err := e.cli.Get(context.Background(), key)
//	if err != nil {
//		return nil, err
//	}
//
//	for _, kv := range res.Kvs {
//		if string(kv.Key) == key {
//			return e.extractMachine(kv.Value)
//		}
//	}
//
//	return nil, nil
//}
//
//func (e *EtcdService) PutMachineId(ip, id string) (map[string]string, error) {
//
//	key := e.MachineKey(ip, id)
//	machineData, err := e.GetMachineId(key)
//	if err != nil {
//		return nil, err
//	}
//
//	if machineData != nil {
//		return nil, errors.New("the machine ip already exists in etcd")
//	}
//
//	if machineData["id"] == id {
//		return nil, errors.New("the machine id has been used in etcd")
//	}
//
//	if !ValidId(id) {
//		return nil, errors.New("the machine id invalid")
//	}
//
//	newMachineData := make(map[string]string)
//	ts := time.Now().Unix()
//	tsStr := strconv.Itoa(int(ts))
//	newMachineData["id"] = id
//	newMachineData["ip"] = ip
//	newMachineData["created_timestamp"] = tsStr
//	newMachineData["updated_timestamp"] = tsStr
//	newMachineDataByte, err := e.makeMachineData(newMachineData)
//	if err != nil {
//		return nil, err
//	}
//
//	if _, err = e.cli.Put(context.Background(), key, string(newMachineDataByte)); err != nil {
//		return nil, err
//	}
//
//	return newMachineData, nil
//}
//
//func (e *EtcdService) GetMachine(key string) (map[string]string, error) {
//	res, err := e.cli.Get(context.Background(), key)
//	if err != nil {
//		return nil, err
//	}
//
//	for _, kv := range res.Kvs {
//		if string(kv.Key) == key {
//			return e.extractMachine(kv.Value)
//		}
//	}
//
//	return nil, nil
//}
//
//func (e *EtcdService) PutMachine(ip string) (error) {
//	key := e.MachineKey(ip)
//	data, err := e.GetMachine(key)
//	if err != nil {
//		return err
//	}
//
//	if data != nil {
//		return errors.New("the machine ip already exists in etcd")
//	}
//
//	// assign id
//	extraData, err := e.GetMachineExtraData()
//	if err != nil {
//		return err
//	}
//
//	if extraData.LastId >= 1024 {
//		if index := extraData.findFreeSlot(); index == 0 {
//			return errors.New("no machine id available")
//		} else {
//			if data, err := e.putMachineDataById(ip, index); err != nil {
//				return err
//			}
//			extraData.updateFreeSlot(index)
//			if err := e.PutMachineExtraData(extraData); err != nil {
//				return err
//			}
//		}
//	}
//}
//
//func (e *EtcdService) DelMachine(ip string) {
//}
//
//func (e *EtcdService) putMachineDataById(ip string, id int) (*MachineData, error) {
//	data := NewMachineData()
//	b, err := data.withIpId(ip, id).encode();
//	if err != nil {
//		return nil, err
//	}
//	key := e.MachineKey(ip)
//	if _, err = e.cli.Put(context.Background(), key, string(b)); err != nil {
//		return nil, err
//	}
//	return data, nil
//}
//
//func (e *EtcdService) GetMachineExtraData() (*MachineExtraData, error) {
//	res, err := e.cli.Get(context.Background(), e.extraKey)
//	if err != nil {
//		return nil, err
//	}
//
//	for _, kv := range res.Kvs {
//		if string(kv.Key) == e.extraKey {
//			return e.extractMachineExtraData(kv.Value)
//		}
//	}
//
//	// create new machine extra data in etcd
//	extraData := &MachineExtraData{
//		FreeSlots: [MaxId]int{0: MaxId,},
//	}
//	if bs, err := json.Marshal(extraData); err != nil {
//		return nil, err
//	} else {
//		e.cli.Put(context.Background(), e.extraKey, string(bs))
//		return extraData, nil
//	}
//}
//
//func (e *EtcdService) PutMachineExtraData(data *MachineExtraData) error {
//	b, err := data.encode()
//	if err != nil {
//		return err
//	}
//
//	if _, err = e.cli.Put(context.Background(), e.extraKey, string(b)); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (e *EtcdService) MachineKey(ip string) string {
//	key := e.prefixKey + strconv.Itoa(int(crc32.ChecksumIEEE([]byte(ip))))
//	fmt.Println(key)
//	return key
//}
//
//func (e *EtcdService) makeMachineData(data map[string]string) ([]byte, error) {
//	return json.Marshal(data)
//}
//
//func (e *EtcdService) extractMachineExtraData(data []byte) (*MachineExtraData, error) {
//	result := new(MachineExtraData)
//	if err := json.Unmarshal(data, result); err != nil {
//		return nil, err
//	}
//
//	return result, nil
//}
//
//func ValidId(id string) bool {
//	intId, err := strconv.Atoi(id)
//	if err != nil {
//		return false
//	}
//
//	if intId < MinId || intId > MaxId {
//		return false
//	}
//
//	return true
//}
