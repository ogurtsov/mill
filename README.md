# mill
I built this tool as a hello-world to golang but ended up actually using it instead of the typical Fabric based approach


# Usage example:
0) Add the executable to sys path: `cp build/mill /usr/bin/mill`
1) Generate default json: `mill makedefaults`
2) Edit `~/.mill.json`
3) Use the tool `mill servername commandname`

# Note:
Optionally you may set your telegram api keys to get the notifications about the deployments. Leave these fields empty if you don't need this feature.
