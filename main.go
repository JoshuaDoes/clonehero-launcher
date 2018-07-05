package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/gen2brain/go-unarr"
	"github.com/xybydy/go-mega"
)

const (
	UpdateURL = "https://clonehero.cameronct.com/ingame/update.php"
)

type Update struct {
	Version  string `json:"version"`  //The latest available version
	Download string `json:"download"` //The Mega.nz folder containing builds of the latest available version
	Required bool   `json:"required"` //Whether or not the update is required
}

func main() {
	fmt.Println("CHUpdater © JoshuaDoes: 2018.")
	fmt.Println("Detected operating system: " + runtime.GOOS + "/" + runtime.GOARCH)
	fmt.Println("")

	fmt.Println("> Initializing Mega...")
	m := mega.New()

	fmt.Println("> Fetching update data...")
	update := &Update{}
	updateResult, err := http.Get(UpdateURL)
	if err != nil {
		defer runCloneHero()
		panic(err)
	}
	err = unmarshal(updateResult, update)
	if err != nil {
		defer runCloneHero()
		panic(err)
	}

	latestVersion, err := strconv.ParseFloat(update.Version, 32)
	if err != nil {
		defer runCloneHero()
		panic(err)
	}
	installUpdated := false

	fmt.Println("> Looking for Clone Hero data...")
	if _, err := os.Stat("Clone Hero_Data/data.unity3d"); err == nil {
		fmt.Println("> Reading Clone Hero data...")
		data, err := ioutil.ReadFile("Clone Hero_Data/data.unity3d")
		if err != nil {
			runCloneHero()
			panic(err)
		}
		fmt.Println("> Checking if Clone Hero is latest version...")
		if bytes.Contains(data, float32ToBytes(float32(latestVersion))) {
			installUpdated = true
		}
	}
	if installUpdated {
		fmt.Println("> Clone Hero is already up-to-date")
		runCloneHero()
		os.Exit(0)
	} else {
		fmt.Println("> An update is available!")
		fmt.Println("> Updating...")
	}

	fmt.Println("> Fetching Mega URL...")
	//Make a function that tells http to not run down the redirect line
	checkRedirect := func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	//Create an HTTP client that uses our custom checkRedirect function
	client := &http.Client{CheckRedirect: checkRedirect}
	//Create the request to the Bitly URL
	request, err := client.Head(update.Download)
	if err != nil {
		defer runCloneHero()
		panic(err)
	}
	//Finally, get the Mega URL from the redirect
	megaURL := request.Header.Get("Location")

	fmt.Println("> Setting MegaFS to Clone Hero folder...")
	_, _ = m.ReturnPublicNode(megaURL)

	fmt.Println("> Fetching MegaFS...")
	megaFS := m.FS

	fmt.Println("> Fetching MegaFS files...")
	megaFSNodes := megaFS.GetAllNodes()

	fmt.Println("> Looking for Clone Hero " + runtime.GOOS + "/" + runtime.GOARCH + "...")
	downloadFound := false
	downloadType := ""
	for _, v := range megaFSNodes {
		nodeName := v.GetName()

		switch nodeName {
		//RAR
		case "Windows (64).rar":
			if runtime.GOOS == "windows" && runtime.GOARCH == "amd64" {
				downloadFound = true
				downloadType = "rar"
				fmt.Println("> Downloading Clone Hero for " + runtime.GOOS + "/" + runtime.GOARCH + "...")
				m.DownloadFile(v, "clonehero.rar", nil, true)
				break
			}
		case "Windows (32).rar":
			if runtime.GOOS == "windows" && runtime.GOARCH == "386" {
				downloadFound = true
				downloadType = "rar"
				fmt.Println("> Downloading Clone Hero for " + runtime.GOOS + "/" + runtime.GOARCH + "...")
				m.DownloadFile(v, "clonehero.rar", nil, true)
				break
			}
		case "Linux.rar":
			if runtime.GOOS == "linux" {
				downloadFound = true
				downloadType = "rar"
				fmt.Println("> Downloading Clone Hero for " + runtime.GOOS + "/" + runtime.GOARCH + "...")
				m.DownloadFile(v, "clonehero.rar", nil, true)
				break
			}

		//7z
		case "Windows (64).7z":
			if runtime.GOOS == "windows" && runtime.GOARCH == "amd64" {
				downloadFound = true
				downloadType = "7z"
				fmt.Println("> Downloading Clone Hero for " + runtime.GOOS + "/" + runtime.GOARCH + "...")
				m.DownloadFile(v, "clonehero.7z", nil, true)
				break
			}
		case "Windows (32).7z":
			if runtime.GOOS == "windows" && runtime.GOARCH == "386" {
				downloadFound = true
				downloadType = "7z"
				fmt.Println("> Downloading Clone Hero for " + runtime.GOOS + "/" + runtime.GOARCH + "...")
				m.DownloadFile(v, "clonehero.7z", nil, true)
				break
			}
		case "Linux.7z":
			if runtime.GOOS == "linux" {
				downloadFound = true
				downloadType = "7z"
				fmt.Println("> Downloading Clone Hero for " + runtime.GOOS + "/" + runtime.GOARCH + "...")
				m.DownloadFile(v, "clonehero.7z", nil, true)
				break
			}
		}
	}
	if !downloadFound {
		defer runCloneHero()
		panic(errors.New("Error finding download for " + runtime.GOOS + "/" + runtime.GOARCH))
	}

	fmt.Println("> Removing previous Clone Hero game files...")
	removeCloneHero()

	fmt.Println("> Loading Clone Hero archive into memory...")
	archive, err := unarr.NewArchive("clonehero." + downloadType)
	if err != nil {
		os.Remove("clonehero." + downloadType)
		panic(err)
	}
	defer archive.Close()

	fmt.Println("> Extracting Clone Hero...")
	err = archive.Extract("")
	if err != nil {
		os.Remove("clonehero." + downloadType)
		removeCloneHero()
		panic(err)
	}

	fmt.Println("> Removing Clone Hero archive...")
	os.Remove("clonehero." + downloadType)

	if runtime.GOOS == "linux" {
		fmt.Println("> Moving game files to current working directory...")
		_ = os.Rename("Linux/Clone Hero_Data", "Clone Hero_Data")
		_ = os.Rename("Linux/Clone Hero.x86_64", "Clone Hero.x86_64")
		_ = os.Rename("Linux/README.txt", "README.txt")
		os.RemoveAll("Linux")
	} else if runtime.GOOS == "windows" {
		if runtime.GOARCH == "386" {
			fmt.Println("> Moving game files to current working directory...")
			os.Rename("Windows (32)/Clone Hero_Data", "Clone Hero_Data")
			os.Rename("Windows (32)/MonoBleedingEdge", "MonoBleedingEdge")
			os.Rename("Windows (32)/Songs", "Songs")
			os.Rename("Windows (32)/Clone Hero.exe", "Clone Hero.exe")
			os.Rename("Windows (32)/README.txt", "README.txt")
			os.Rename("Windows (32)/UnityCrashHandler32.exe", "UnityCrashHandler32.exe")
			os.Rename("Windows (32)/UnityPlayer.dll", "UnityPlayer.dll")
			os.RemoveAll("Windows (32)")
		} else if runtime.GOARCH == "amd64" {
			fmt.Println("> Moving game files to current working directory...")
			os.Rename("Windows (64)/Clone Hero_Data", "Clone Hero_Data")
			os.Rename("Windows (64)/MonoBleedingEdge", "MonoBleedingEdge")
			os.Rename("Windows (64)/Songs", "Songs")
			os.Rename("Windows (64)/Clone Hero.exe", "Clone Hero.exe")
			os.Rename("Windows (64)/README.txt", "README.txt")
			os.Rename("Windows (64)/UnityCrashHandler32.exe", "UnityCrashHandler32.exe")
			os.Rename("Windows (64)/UnityPlayer.dll", "UnityPlayer.dll")
			os.RemoveAll("Windows (64)")
		}
	}

	runCloneHero()
}

func runCloneHero() {
	fmt.Println("> Running Clone Hero...")
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("./Clone Hero.exe")
		err := cmd.Start()
		if err != nil {
			panic(err)
		}
	case "linux":
		cmd := exec.Command("./Clone Hero.x86_64")
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}
}

func removeCloneHero() {
	switch runtime.GOOS {
	case "windows":
		os.RemoveAll("Clone Hero_Data")
		os.Remove("Clone Hero.exe")
		os.Remove("UnityPlayer.dll")
	case "linux":
		os.RemoveAll("Clone Hero_Data")
		os.Remove("Clone Hero.x86_64")
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
