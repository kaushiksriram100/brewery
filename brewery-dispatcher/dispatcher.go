package main 


import (
"fmt"
machinery "github.com/RichardKnop/machinery/v1"
machineryconfig "github.com/RichardKnop/machinery/v1/config"
breweryutils "github.com/kaushiksriram100/brewery/brewery-utils"
tasks "github.com/RichardKnop/machinery/v1/tasks"
"flag"
)

func main() {

	configfile := flag.String("c", "/var/tmp/brewery-conf.json", "- json format config file with broker end points")
	flag.Parse()
	brewerycnf, err := breweryutils.LoadServerConfig(*configfile)  //cnf is of type *config.Config
	if err != nil {
		fmt.Printf("Some issues loading the config- %s", err)
		return
	}

	//Now broker endpoints are unmarshalled and put in the Config Struct. Let's add this to the brokerconfig struct we created and return
	var brokerconfig = &machineryconfig.Config{} //this config as mandated in machinery code. Don't confuse with the above config that is meant to parse the json config file

	(*brokerconfig).Broker = brewerycnf.Broker.MessageBroker+"://"+brewerycnf.Broker.BrokerHost
	(*brokerconfig).ResultBackend = brewerycnf.Broker.MessageBroker+"://"+brewerycnf.Broker.BrokerHost
	(*brokerconfig).DefaultQueue = brewerycnf.Broker.DefaultQueue

	//create a new server object. Even for a client, machinery requires a server to be instantiated before launching the worker. We will use our custom utils to do all that. 

	server, err := machinery.NewServer(brokerconfig)
	if err != nil {
		fmt.Printf("Can't Create Server.. exiting")
		return
	}

//start from here - you have all data in brewerycnf now. Need to send them to a new task (eg- metriccollector). Create a new task in brewery-tasks.go
	sayTask := &tasks.Signature{
		Name: "Say",
		Args: []tasks.Arg{tasks.Arg{Type: "string",Value: "sriram"}},
	}



	asyncResult, err := server.SendTask(sayTask)

	if err != nil {
		fmt.Println(err)
		fmt.Printf("didn't send task, exiting\n")
		return
	}

	taskState := asyncResult.GetState()
	fmt.Println(taskState.State)

}