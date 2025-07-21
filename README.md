<h1 align="center">mofetch</h1>
<p align="center">
<img src="./mofetch.jpg" width="80%" />
</p>

Neofetch-based app for movies!

Heavily based off of [Mufetch](https://github.com/ashish0kumar/mufetch) by [Ashish Kumar](https://github.com/ashish0kumar), it is indicated in the files when some or all of the code for that file was taken from his project. Go check out his project, he deserves tons of credit.

## Features
- Search for any movie with posters and detailed information
- Customizable poster dimensions
- Customizable amount of information displayed
- Cross-platform support, works on any computer

## Installation

### Releases Tab
Simply go to the releases tab on this page, and select the appropriate file for your system.

The ZIP/tarball files are just source code.

If you are unsure of what version to download, run 
```bash
# On Windows Command Prompt
echo %PROCESSOR_ARCHITECTURE%

#On the Mac/Linux terminal
arch=$(uname -m)
case "$arch" in
    x86_64) echo "amd64" ;;
    aarch64 | arm64) echo "arm64" ;;
    i386 | i686) echo "386" ;;
    *) echo "$arch" ;;
esac
```

### Build from source
```bash
git clone https://github.com/ashish0kumar/mufetch.git
cd mufetch/
go build
sudo mv mufetch /usr/local/bin/
mufetch --help
```

## First-Time Authentication
1. Go to [the OMDB API website](https://www.omdbapi.com/apikey.aspx)
2. Select 'FREE! (1,000 daily limit)'
3. Enter in your email and name, put anything in for 'Use'
4. Click Submit
5. Check your email for the API Key
6. Run ```mofetch auth``` and paste in your API key.
7. Hit enter, and you're good to go!

i will update this later
