package utils

import (
//"fmt"
"errors"
"encoding/json"
"io/ioutil"
)


type BrokerConfig struct {
	MessageBroker string `json:"brokertype"`  //We will keep brokertype modular so that it's easy in future to add other brokers like RMQ
	BrokerHost string  `json:"hostname"`
	DefaultQueue string `json:"queue"`
}


type Input struct {
	CheckName string `json:"check_name"`
	CheckCommand string `json:"check_command"`
}

type Graphite struct {
	Graphiteendpoints string `json:"endpoint"`
	Graphiteport int `json:"port"`
}

type Output struct {
	GraphiteOutput Graphite `json:"graphite,omitempty"`

}

type Config struct {
	Broker BrokerConfig `json:"broker"`
	Inputs []Input `json:"inputs,omitempty"`
	Outputs Output `json:"outputs"`
}

//LoadServerConfig will parse though the config.json file and get the redis broker endpoints and return the config.Config that can be used to instantiate a machinery Server. 
func LoadServerConfig(configfile string) (Config, error) {
	var conf Config

	//First read the json from the config file
	data, err := ioutil.ReadFile(configfile)
	if err != nil {
		
		return conf, errors.New("unable to read config file")
	}

	//now unmarshall the json to our Config struct

	err = json.Unmarshal(data, &conf)
	if err != nil {
		return conf, errors.New("Incorrect Json format")
	}

	

	return conf, nil  //return conf (pass as value). Configs are not very big so it's okay to pass as value, IMO


}
