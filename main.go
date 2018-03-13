package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os/user"
	"strings"

	"github.com/beevik/etree"
	"github.com/go-ini/ini"
)

type WinSCPSession struct {
	Name       string
	HostName   string
	UserName   string
	PortNumber string
	FSProtocol string
	Password   string
}

func ReadWinSCPIni(filepath string) []WinSCPSession {
	cfg, err := ini.Load(filepath)
	if err != nil {
		panic(err)
	}

	sessions := make([]WinSCPSession, 0)

	for _, s := range cfg.Sections() {
		if !strings.HasPrefix(s.Name(), "Sessions\\") || !s.HasKey("HostName") {
			continue
		}

		session := WinSCPSession{}
		session.Name = strings.TrimPrefix(s.Name(), "Sessions\\")
		session.HostName = s.Key("HostName").Value()
		if s.HasKey("UserName") {
			session.UserName = s.Key("UserName").Value()
		}
		session.Password = Decrypt(session.HostName, session.UserName, s.Key("Password").Value())
		if s.HasKey("FSProtocol") {
			session.FSProtocol = s.Key("FSProtocol").Value()
		}
		if s.HasKey("PortNumber") {
			session.PortNumber = s.Key("PortNumber").Value()
		}

		sessions = append(sessions, session)
	}

	return sessions
}

func getOrCreateFolder(nameParts []string, parent *etree.Element) *etree.Element {
	var folder *etree.Element
	if parent != nil {
		match := parent.FindElements(fmt.Sprintf("//Folder/text() = '%s'", nameParts[0]))
		if len(match) > 0 {
			folder = match[0]
		} else {
			folder = parent.CreateElement("Folder")
			folder.SetText(nameParts[0])
		}
	} else {
		folder = etree.NewElement("Folder")
		folder.SetText(nameParts[0])
	}

	nameParts = nameParts[1:]
	if len(nameParts) > 0 {
		getOrCreateFolder(nameParts, folder)
	}
	allFolders := folder.FindElements("//Folder")
	if len(allFolders) > 0 {
		return allFolders[len(allFolders)-1]
	}
	return folder
}

func WriteFileZillaXML(sessions []WinSCPSession, filepath string) {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version=1.0 encoding="UTF-8"`)
	rootNode := doc.CreateElement("FileZilla3")
	serversNode := rootNode.CreateElement("Servers")
	rootFolder := serversNode.CreateElement("Folder")
	rootFolder.SetText("Imported")

	for _, session := range sessions {
		nameParts := strings.Split(session.Name, "/")
		currentFolder := rootFolder
		if len(nameParts) > 1 {
			currentFolder = getOrCreateFolder(nameParts[:len(nameParts)-1], rootFolder)
		}
		serverNode := currentFolder.CreateElement("Server")
		node := serverNode.CreateElement("Name")
		node.SetText(nameParts[len(nameParts)-1])
		node = serverNode.CreateElement("Host")
		node.SetText(session.HostName)
		node = serverNode.CreateElement("Protocol")
		if session.FSProtocol != "" || session.FSProtocol == "5" {
			node.SetText("0") // FTP
		} else {
			node.SetText("1") // SFTP, default in WinSCP
		}
		node = serverNode.CreateElement("User")
		node.SetText(session.UserName)
		node = serverNode.CreateElement("Pass")
		node.SetText(base64.StdEncoding.EncodeToString([]byte(session.Password)))
		node.CreateAttr("encoding", "base64")
		node = serverNode.CreateElement("Port")
		if session.PortNumber != "" {
			node.SetText(session.PortNumber)
		} else {
			node.SetText("22")
		}
		node = serverNode.CreateElement("LocalDir")
		node = serverNode.CreateElement("RemoteDir")
		node = serverNode.CreateElement("Type")
		node.SetText("0")
		node = serverNode.CreateElement("Logontype")
		node.SetText("1")
		node = serverNode.CreateElement("TimezoneOffset")
		node.SetText("0")
		node = serverNode.CreateElement("PasvMode")
		node.SetText("MODE_DEFAULT")
		node = serverNode.CreateElement("MaximumMultipleConnections")
		node.SetText("0")
		node = serverNode.CreateElement("EncodingType")
		node.SetText("Auto")
		node = serverNode.CreateElement("BypassProxy")
		node.SetText("0")
		node = serverNode.CreateElement("SyncBrowsing")
		node.SetText("0")
		node = serverNode.CreateElement("Comments")
	}

	doc.Indent(2)
	doc.WriteToFile(filepath)
}

func getDefaultIniPath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir + "\\AppData\\Roaming\\winSCP.ini"
}

func main() {
	iniPath := flag.String("in", getDefaultIniPath(), fmt.Sprintf("Path to your WinScp.ini (default: %s)", getDefaultIniPath()))
	outPath := flag.String("out", "sites.xml", fmt.Sprintf("Output path for the resulting XML file. (default: %s)", "sites.xml"))

	flag.PrintDefaults()
	flag.Parse()

	sessions := ReadWinSCPIni(*iniPath)
	WriteFileZillaXML(sessions, *outPath)
}
