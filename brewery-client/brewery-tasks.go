package main 

import (
"os/exec"
"context"
"time"
"errors"
"fmt"
"os"
"strings"
)


func Say(name string) (string, error) {

	return "sriramkaushik", nil
}

//Function metric collector gets all the data that we pass to sendtask in machinery. 
func CommandExecutor(checkcommand string) (string, error) {

	var output string 

	chkcmd := strings.Fields(checkcommand)

	head := chkcmd[0]
	chkcmd = chkcmd[1:len(chkcmd)] //now chkcmd is a slice with each ags. 

	ctx, cancel := context.WithTimeout(context.Background(), 3100*time.Second)


	defer cancel()

	cmd := exec.CommandContext(ctx, head, chkcmd...)  //pass each element in chkcmd to COmmandContext. CommandContext receives args as ...string. So we can pass string...

	env := os.Environ()
	cmd.Env = env

	result, err := cmd.Output()

	if err != nil {
		return output, errors.New("some issues")
	}

	return string(result), nil
}