// this code will refresh the aws secrets in cluster

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	credsDir := "/lhome/kcharan/.aws/credentials"
	file, _ := os.Open(credsDir)
	scanner := bufio.NewScanner(file)
	var aws_access_key_id, aws_secret_access_key, aws_session_token string
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		stringsSli := strings.Split(line, "=")
		if len(stringsSli) > 1 {
			if strings.Contains(stringsSli[0], "aws_access_key_id") {
				aws_access_key_id = stringsSli[1]
			} else if strings.Contains(stringsSli[0], "aws_secret_access_key") {
				aws_secret_access_key = stringsSli[1]
			} else {
				aws_session_token = stringsSli[1]
			}
		}
	}
	cwd, _ := os.Getwd()
	shFile, _ := os.Create(cwd + "/sercretsupdate.sh")
	shFile.Write([]byte("chmod +x " + cwd + "/sercretsupdate.sh\n"))
	shFile.Write([]byte("kubectl delete secrets aws-credentials\n"))
	shFile.Write([]byte("kubectl create secret generic aws-credentials " +
		"--from-literal=aws_access_key_id=\"" + aws_access_key_id + "\" " +
		"--from-literal=aws_secret_access_key=\"" + aws_secret_access_key + "\" " +
		"--from-literal=aws_session_token=\"" + aws_session_token + "\""))
	cmd := exec.Command("/bin/bash", cwd+"/sercretsupdate.sh")
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
}
