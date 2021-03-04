/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

type SSHConnection struct {
	Username string
	Host     string
	Port     string
}

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect through SSH to the desired directory, generate a ZIP file in the server, copy that file to the local directory under './backup', delete the server file, and download the zip through SFTP",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		connectionParams := getConnectionParameters()
		sshConnection := establishSSHConection(connectionParams)
		session, filename := sshTerminalSession(sshConnection, connectionParams)

		defer downloadFileThroughSFTP(sshConnection, session, filename)
		defer fmt.Println("Done!")
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}

func getConnectionParameters() SSHConnection {
	fmt.Println("Please write your SSH username")
	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')
	fmt.Println("Now introduce the client host")
	host, _ := reader.ReadString('\n')
	fmt.Println("Now introduce the desired port")
	port, _ := reader.ReadString('\n')

	connection := SSHConnection{}
	connection.Username = strings.TrimSpace(username)
	connection.Host = strings.TrimSpace(host)
	connection.Port = strings.TrimSpace(port)

	return connection
}

func getKeyFile() (key ssh.Signer, err error) {
	usr, _ := user.Current()
	file := usr.HomeDir + "/.ssh/id_rsa"

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		return
	}

	return key, err
}

func establishSSHConection(sshParams SSHConnection) *ssh.Client {
	// Now in the main function DO:
	key, err := getKeyFile()

	sshConfig := &ssh.ClientConfig{
		User: sshParams.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		//probably not good solution, find secure one
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", sshParams.Host+":"+sshParams.Port, sshConfig)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}

	fmt.Println("Successfully connected to ssh server.")
	return client
}

func sshTerminalSession(conn *ssh.Client, sshParams SSHConnection) (*ssh.Session, string) {
	session, err := conn.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}

	now := time.Now().Format("2006_01_02_150405_")
	filename := now + sshParams.Host + ".zip"
	fmt.Println("File will be called: " + filename)
	fmt.Println("Zipping the site, please be patient...")
	err = session.Run("zip -r " + filename + " www/")

	if err != nil {
		log.Fatal(err)
	}

	return session, filename
}

func downloadFileThroughSFTP(conn *ssh.Client, session *ssh.Session, filename string) {
	defer session.Close()

	fmt.Println("Copying file from SFTP to your ./backups directory ...")
	connectViaSFTP(conn, filename)
	fmt.Println("Also done!")

	fmt.Println("Deleting zip file from the server...")
	session.Run("rm " + filename)
	fmt.Println("Exiting, bye! :D")
	session.Run("exit")
}

func connectViaSFTP(sshConnection *ssh.Client, filename string) {
	client, err := sftp.NewClient(sshConnection)
	if err != nil {
		log.Fatal(err)
	}

	w := client.Walk(".")
	path := "www"

	for w.Step() {

		if w.Path() == path {
			break
		}

		if w.Err() != nil {
			continue
		}
	}

	// Check that file exists
	_, err = client.Lstat(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Open the source file
	srcDir, err := client.Open("./" + filename)
	if err != nil {
		log.Fatal(err)
	}

	defer srcDir.Close()

	if _, err := os.Stat("./backups"); os.IsNotExist(err) {
		os.Mkdir("./backups", os.ModeDir|0755)
	}

	// Create the destination file
	dstDir, err := os.Create("./backups/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer dstDir.Close()

	srcDir.WriteTo(dstDir)
}
