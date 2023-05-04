[![Development Status](https://img.shields.io/badge/DEVELOPMENT-PAUSED-orange?style=for-the-badge)](#development-status)
#### Badges are clickable for more information
&nbsp;

# dothelp
dothelp is a simple console utility that makes using the dotnet CLI easier, written in GO. As of right now, development is paused. In its current state, it is not very feature rich, coverting only the project creation and building.

This was mostly a learning exercise and it is not recommended that you actually use this program in your workflow as it might be unstable or produce unexpected behaviour.

## Installation
dothelp can be installed either using the binaries provided under the releases, or by building it from source.

### Building from source
Prequisites
- Golang
- dotnet CLI

```bash
# Clone the repo
git clone https://github.com/underthefoxtree/dothelp
cd dothelp/src

# Build it
go build -o ../ -ldflags '-s -w'
```

### Manual installation
Copy the built executable into `/usr/bin/`

### Using fpm and your package manager (Recommended-ish)
Prequisites
- [fpm](https://github.com/jordansissel/fpm)

First generate the packaged file.
```bash
# Move to the output directory
cd ..

# On Debian/Ubuntu based systems
fpm -t deb

# On Red Hat/Fedora based systems
fpm -t rpm

# On Arch based systems (untested)
fpm -t pacman
```
Then install the packaged file using your package manager.

## Development Status
Status | Description
---|---
Active | The project is actively being worked on and new features are being added
On Demand | Bugs and other Issues will be fixed, but no new features will be added
Paused | No development will take place at the moment, but this may change in the future
Ceased | The project will not be worked on AT ALL

The development status can change at any time in both ways (more/less work being done).

## License
This software is licensed under the [MIT License](LICENSE.md).
