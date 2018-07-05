# chupdater

### A custom launcher for Clone Hero that installs/updates the game for you pre-launch

----

## Getting and using CHUpdater

1. Open the [releases page](https://github.com/JoshuaDoes/chupdater/releases) and download the latest appropriate version of CHUpdater to whichever folder on your computer you'd like to install Clone Hero to or update Clone Hero in.
2. From now on, run CHUpdater instead of regular Clone Hero.
3. Enjoy!
X. Optional: Create a shortcut to CHUpdater and give it the Clone Hero icon.

## What does it do?

When you run CHUpdater, it fetches the latest update data and checks if the latest version number matches that of a specific file in the Clone Hero game files. If you're on the latest version already, it immediately launches the game. If you're not yet on the latest version or the game files don't exist, it will automatically download the latest version and install it to the current directory before launching it.

## Why CHUpdater?

1. It keeps your Clone Hero install updated automatically!
2. It's very helpful for quick one-click installs of Clone Hero on new or refreshed computers!
3. For both of the above, no need to fidgit with WinRAR/7-Zip or other RAR extraction utilities - CHUpdater does everything for you!

----

## Building it yourself

In order to build CHUpdater locally, you must have already installed a
working Golang environment on your development system and installed the package
dependencies that CHUpdater relies on to function properly.

CHUpdater is currently built using Golang `1.10.2`.

### Dependencies

| Package Name |
| ------------ |
| [go-unarr](https://github.com/gen2brain/go-unarr) |
| [go-mega](https://github.com/xybydy/go-mega) |

### Building

Simply run `go build` in this repo's directory once all dependencies are satisfied.

### Running CHUpdater

Finally, to run CHUpdater, simply type `./chupdater` in your
terminal/shell or `.\chupdater.exe` in your command prompt. If everything
goes well, you'll see a full installation of Clone Hero by time CHUpdater
finishes its job.

### Contributing notes

When pushing to your repo or submitting pull requests to this repo, it is highly
advised that you clean up the working directory to only contain `LICENSE`, `main.go`,
`make.bat`, `make.sh`, `README.md`, and the `.git` folder. A proper `.gitignore` will
be written soon to mitigate this requirement.

----

## Support
For help and support with CHUpdater, create an issue on the issues page. If you do not have a GitHub account, send me a message on Discord (@JoshuaDoes#1685).

## License
The source code for CHUpdater is released under the MIT License. See LICENSE for more details.

## Donations
All donations are highly appreciated. They help motivate me to continue working on side projects like these, especially when it comes to something you may really want added!

[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://paypal.me/JoshuaDoes)