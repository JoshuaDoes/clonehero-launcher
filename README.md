# clonehero-launcher

### A custom launcher for Clone Hero that installs/updates the game for you before launching the game if needed

----

## Getting and using the launcher

1. Open the [releases page](https://github.com/JoshuaDoes/clonehero-launcher/releases) and download the latest appropriate version of Clone Hero Launcher to whichever folder on your computer you'd like to install Clone Hero to or update Clone Hero in.
2. From now on, run Clone Hero Launcher instead of regular Clone Hero.
3. Enjoy!

## What does it do?

When you run Clone Hero Launcher, it fetches the latest update data before then checking if the latest version number matches that of a specific file in the Clone Hero game files. If you're on the latest version already, it immediately launches the game. If you're not on the latest version or the game files don't exist, it will automatically download the latest version and install it to the current directory before launching the latest version.

## Why Clone Hero Launcher?

1. It keeps your Clone Hero updated automatically!
2. It's very helpful for quick one-click installs of Clone Hero on new or refreshed computers!
3. For both of the above, no need to fidgit with WinRAR/7-Zip or other RAR extraction utilities - Clone Hero Launcher does everything for you!

----

## Building it yourself

In order to build Clone Hero Launcher locally, you must have already installed a
working Golang environment on your development system and installed the package
dependencies that Clone Hero Launcher relies on to function properly. Clone Hero
Launcher is currently built using Golang `1.10.2`.

### Dependencies

| Package Name |
| ------------ |
| [googl](https://github.com/JoshuaDoes/googl) |
| [archiver](https://github.com/mholt/archiver) |
| [go-mega](https://github.com/xybydy/go-mega) |

### Building

Simply run `go build` in this repo's directory once all dependencies are satisfied.

### Acquiring necessary API keys

Clone Hero Launcher only requires one API key to function, and although my own is
supplied with this source, it is highly recommended to use your own in your own
builds of Clone Hero Launcher.

| Services | Requirements |
| -------- | ------------ |
| Googl | Googl API Key |

### Running Clone Hero Launcher

Finally, to run Clone Hero Launcher, simply type `./clonehero-launcher` in your
terminal/shell or `.\clonehero-launcher.exe` in your command prompt. If everything
goes well, you'll see a full installation of Clone Hero by time Clone Hero Launcher
finishes its job.

When pushing to your repo or submitting pull requests to this repo, it is highly
advised that you clean up the working directory to only contain `LICENSE`, `main.go`,
`README.md`, and the `.git` folder. A proper `.gitignore` will be written soon to
mitigate this requirement.

----

## Support
For help and support with Clone Hero Launcher, create an issue on the issues page. If you do not have a GitHub account, send me a message on Discord (@JoshuaDoes#1685).

## License
The source code for Clone Hero Launcher is released under the MIT License. See LICENSE for more details.

## Donations
All donations are highly appreciated. They help motivate me to continue working on side projects like these, especially when it comes to something you may really want added!

[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://paypal.me/JoshuaDoes)