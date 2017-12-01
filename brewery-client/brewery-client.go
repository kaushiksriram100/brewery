package main 


import (
"fmt"
machinery "github.com/RichardKnop/machinery/v1"
breweryutils "github.com/kaushiksriram100/brewery/brewery-utils"
machineryconfig "github.com/RichardKnop/machinery/v1/config"
"flag"
"errors"
"os"
)

func main() {

	configfile := flag.String("c", "/var/tmp/brewery-conf.json", "- json format config file with broker end points")
	flag.Parse()
	brewerycnf, err := breweryutils.LoadServerConfig(*configfile)  //cnf is of type *config.Config
	if err != nil {
		fmt.Printf("Some issues loading the config- %s\n", err)
		return
	}

	var brokerconfig = &machineryconfig.Config{} //this config as mandated in machinery code.

	(*brokerconfig).Broker = brewerycnf.Broker.MessageBroker+"://"+brewerycnf.Broker.BrokerHost
	(*brokerconfig).ResultBackend = brewerycnf.Broker.MessageBroker+"://"+brewerycnf.Broker.BrokerHost
	(*brokerconfig).DefaultQueue = brewerycnf.Broker.DefaultQueue

	//create a new server object. Even for a client, machinery requires a server to be instantiated before launching the worker. We will use our custom utils to do all that. 

	server, err := machinery.NewServer(brokerconfig)
	if err != nil {
		fmt.Printf("Can't create server.. exiting\n")
		return
	}

	//Register all the tasks as defined in tasks.go. Not using breweryutils for this. Just lazy

	server.RegisterTask("Say", Say)
	server.RegisterTask("CommandExecutor", CommandExecutor)

	workername, err := os.Hostname()
	if err != nil {
		fmt.Printf("can't get hostname\n")
		return 
	}
	worker := server.NewWorker(workername,1000)

	err = worker.Launch()

	if err != nil {
		errors.New("could not launch worker\n")
		return
	}

}