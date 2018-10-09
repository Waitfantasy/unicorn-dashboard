package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-contrib/multitemplate"
	"html/template"
	"io/ioutil"
	"log"
)

func initFlag() {
	flag.StringVar(&ConfigFilePath, "config", "", "config file path")
	flag.Parse()
}

func initConfig() *Config {
	if ConfigFilePath == "" {
		log.Fatal("config file path is empty!")
	}

	data, err := ioutil.ReadFile(ConfigFilePath)
	if err != nil {
		log.Fatalf("read config file error: %v", err)
	}

	config := new(Config)
	if err = json.Unmarshal(data, config); err != nil {
		log.Fatalf("the config file can not json.Unmarshal: %v", err)
	}
	fmt.Println(config)
	return config
}

func createHTMLRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index",
		"./templates/base.html",
		"./templates/layout/nav.html",
		"./templates/layout/sidebar.html",
		"./templates/layout/content-header.html",
		"./templates/layout/content.html",
		"./templates/layout/footer.html", )

	r.AddFromFilesFuncs("machine/index",
		template.FuncMap{
			"formatMachineIP": formatMachineIP,
			"formatDate":      formatDate,
		},
		"./templates/base.html",
		"./templates/layout/nav.html",
		"./templates/layout/sidebar.html",
		"./templates/layout/footer.html",
		"./templates/machine/index.html", )
	return r
}
