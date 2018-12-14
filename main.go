package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/gen2brain/go-unarr"
)

const (
	UpdateURL   = "https://clonehero.net/ingame/update.php"
	DownloadUrl = "https://clonehero.net/download?all=true"
	ArchiveName = "clonehero.7z"
)

var (
	folderName string
	fileName   string
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
	establishFileName()

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

	fmt.Println("> Fetching Download URL...")

	resp, _ := http.Get(DownloadUrl)
	defer resp.Body.Close()

	fmt.Println("> Looking for Clone Hero " + runtime.GOOS + "/" + runtime.GOARCH + "...")
	//Finally, get the Mega URL from the redirect
	downloadUrl, err := getDownloadUrlFromHTML(resp.Body)
	out, err := os.Create(ArchiveName)
	if err != nil {
		runCloneHero()
		panic(err)
	}

	resp, err = http.Get(downloadUrl)
	if err != nil {
		runCloneHero()
		panic(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		runCloneHero()
		panic(err)
	}

	fmt.Println("> Removing previous Clone Hero game files...")
	_ = removeCloneHero()

	fmt.Println("> Loading Clone Hero archive into memory...")
	archive, err := unarr.NewArchive(ArchiveName)
	if err != nil {
		err1 := archive.Close()
		err2 := os.Remove(ArchiveName)
		errString := "Archive Open: " + err.Error() + "\n"
		if err1 != nil {
			errString += "Archive close Error: " + err1.Error() + "\n"
		}
		if err2 != nil {
			errString += "Archive delete Error: " + err2.Error() + "\n"
		}
		panic(errors.New(errString))
	}

	fmt.Println("> Extracting Clone Hero...")
	err = archive.Extract("")
	err1 := archive.Close()

	if err != nil {
		err2 := os.Remove(ArchiveName)
		_ = removeCloneHero()
		errString := "Extract Error: " + err.Error() + "\n"
		if err1 != nil {
			errString += "Archive close Error: " + err1.Error() + "\n"
		}
		if err2 != nil {
			errString += "Archive delete Error: " + err2.Error() + "\n"
		}
		panic(errors.New(errString))
	}

	fmt.Println("> Removing Clone Hero archive...")
	_ = os.Remove(ArchiveName)

	moveFiles()

	runCloneHero()
}

func moveFiles() {
	fmt.Println("> Moving game files to current working directory...")
	files, _ := ioutil.ReadDir(folderName)

	if files != nil {
		for _, file := range files {
			_ = os.Rename(folderName+"/"+file.Name(), file.Name())
		}
		_ = os.RemoveAll(folderName)
	}
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

func removeCloneHero() (err error) {
	switch runtime.GOOS {
	case "windows":
		err = os.RemoveAll("Clone Hero_Data")
		err = os.Remove("Clone Hero.exe")
		err = os.Remove("UnityPlayer.dll")
	case "linux":
		err = os.RemoveAll("Clone Hero_Data")
		err = os.Remove("Clone Hero.x86_64")
	}

	return
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

func establishFileName() {
	if runtime.GOOS == "windows" && runtime.GOARCH == "amd64" {
		folderName = "clonehero-win64"
	}
	if runtime.GOOS == "windows" && runtime.GOARCH == "386" {
		folderName = "clonehero-win32"
	}
	if runtime.GOOS == "linux" {
		folderName = "clonehero-linux"
	}

	fileName = folderName + ".7z"
}

func getDownloadUrlFromHTML(body io.ReadCloser) (string, error) {
	z := html.NewTokenizer(body)
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return "", errors.New("unable to find download link")
		case tt == html.StartTagToken:
			t := z.Token()

			if t.Data == "a" {
				for _, a := range t.Attr {
					if a.Key == "href" {
						urlParts, _ := url.Parse(a.Val)
						if urlParts.Host == "dl.clonehero.net" && urlParts.Path == "/"+fileName {
							return a.Val, nil
						}
					}
				}
			}
		}
	}
}
