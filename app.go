package main

var ConfigFilePath string

type Config struct {
	Etcd struct {
		Cluster string `json:"cluster"`
	} `json:"etcd"`
}
