package main

import (
	"bytes"
	"encoding/binary"
	//"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/JoshuaDoes/googl"
	"github.com/mholt/archiver"
	"github.com/xybydy/go-mega"
)

const (
	UpdateURL   = "https://clonehero.cameronct.com/ingame/update.php"
	GooglAPIKey = "AIzaSyBx1uipz3amt54YGfsPXZhopIjQlU4kveo"
)

type Update struct {
	Version  string `json:"version"`  //The latest available version
	Download string `json:"download"` //The Mega.nz folder containing builds of the latest available version
	Required bool   `json:"required"` //Whether or not the update is required
}

func main() {
	fmt.Println("Clone Hero Installer/Updater © JoshuaDoes: 2018.")
	fmt.Println("Detected operating system: " + runtime.GOOS + "/" + runtime.GOARCH)
	fmt.Println("")

	fmt.Println("> Initializing Mega...")
	m := mega.New()
	fmt.Println("> Initializing Googl...")
	g := googl.NewClient(GooglAPIKey)

	fmt.Println("> Fetching update data...")
	update := &Update{}
	updateResult, err := http.Get(UpdateURL)
	if err != nil {
		clonehero()
		panic(err)
	}
	err = unmarshal(updateResult, update)
	if err != nil {
		clonehero()
		panic(err)
	}

	latestVersion, err := strconv.ParseFloat(update.Version, 32)
	if err != nil {
		clonehero()
		panic(err)
	}
	installUpdated := false

	fmt.Println("> Looking for Clone Hero data...")
	if _, err := os.Stat("Clone Hero_Data/data.unity3d"); err == nil {
		fmt.Println("> Reading Clone Hero data...")
		data, err := ioutil.ReadFile("Clone Hero_Data/data.unity3d")
		if err != nil {
			clonehero()
			panic(err)
		}
		fmt.Println("> Checking if Clone Hero is latest version...")
		if bytes.Contains(data, float32ToBytes(float32(latestVersion))) {
			installUpdated = true
		}
	}
	if installUpdated {
		fmt.Println("> Clone Hero is already up-to-date")
		clonehero()
		os.Exit(0)
	} else {
		fmt.Println("> Updating...")
	}

	fmt.Println("> Fetching Mega.nz URL...")
	gExpand, err := g.Expand(update.Download)
	if err != nil {
		clonehero()
		panic(err)
	}
	megaURL := gExpand.LongUrl

	fmt.Println("> Setting MegaFS to folder node...")
	_, _ = m.ReturnPublicNode(megaURL)

	fmt.Println("> Fetching MegaFS...")
	megaFS := m.FS

	fmt.Println("> Fetching MegaFS nodes...")
	megaFSNodes := megaFS.GetAllNodes()

	fmt.Println("> Looking for Clone Hero " + runtime.GOOS + "/" + runtime.GOARCH + "...")
	downloadFound := false
	for _, v := range megaFSNodes {
		nodeName := v.GetName()

		switch nodeName {
		case "Windows (64).rar":
			if runtime.GOOS == "windows" && runtime.GOARCH == "amd64" {
				downloadFound = true
				fmt.Println("> Downloading Clone Hero...")
				m.DownloadFile(v, "clonehero.rar", nil, true)
				break
			}
		case "Windows (32).rar":
			if runtime.GOOS == "windows" && runtime.GOARCH == "i386" {
				downloadFound = true
				fmt.Println("> Downloading Clone Hero...")
				m.DownloadFile(v, "clonehero.rar", nil, true)
				break
			}
		case "Linux.rar":
			if runtime.GOOS == "linux" {
				downloadFound = true
				fmt.Println("> Downloading Clone Hero...")
				m.DownloadFile(v, "clonehero.rar", nil, true)
				break
			}
		}
	}
	if !downloadFound {
		clonehero()
		panic(errors.New("Error finding download for " + runtime.GOOS + "/" + runtime.GOARCH))
	}

	fmt.Println("> Extracting Clone Hero...")
	err = archiver.Rar.Open("clonehero.rar", "")
	if err != nil {
		os.Remove("clonehero.rar")
		os.RemoveAll("clonehero")
		panic(err)
	}

	fmt.Println("> Removing Clone Hero archive...")
	os.Remove("clonehero.rar")

	fmt.Println("> Running Clone Hero...")
	clonehero()
}

func clonehero() {
	fmt.Println("> Running Clone Hero...")
	if runtime.GOOS == "windows" {
		cmd := exec.Command("Clone Hero.exe")
		err := cmd.Start()
		if err != nil {
			panic(err)
		}
	}
	if runtime.GOOS == "linux" {
		cmd := exec.Command("Clone Hero")
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}
}

func unmarshal(body *http.Response, target interface{}) error {
	defer body.Body.Close()
	return json.NewDecoder(body.Body).Decode(target)
}

func float32ToBytes(f float32) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, f)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}
