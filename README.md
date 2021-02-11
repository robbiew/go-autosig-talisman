# autosig
Auto Signature Mod for Talisman BBS

This go-based program allows you to add an Auto Sig editor to you Talisman BBS -- as a door. 

It incorporates a modified vesion of the kilo editor (https://github.com/bediger4000/kilo-in-go). 

Note: it won't compile on OSX. In fact, I've only tested it on Linux 64 (Ubuntu 20.04). It's possible that it wil compile on Windows and Raspbery Pi...

To set up:

- Compile the program (go build .) -- or, use the pre-built binary in this repo for Linux 64
- Create /bbs/doors/autosig directory
- Copy autosig (binary), header.ans, start.sh to the above directory 
- Add this to a menu, and it edit path (e.g. this is in /bbs/menus/message.toml)

[[menuitem]]
command = "RUNDOOR"
data = "/home/robbiew/bbs/doors/autosig/start.sh"
hotkey = "Z"

Talisman will automatically pass the node number to the start.sh file, and the autosig program will use this to grab the drop file (e.g. bbs/temp/$1/door.sys). The drop file contains the user's id and name, which will be used to retrieve the "signature" in Talisman's user database (sqlite3). If the signature value doesn't exist yet, it will create it. If it exists, it will updated it.

For display purposes, the autosig program translates between color "pipe" codes (e.g. "|02") and actual ansi escape codes.

I have not tested for background colors yet.

TO DO:

- A way to add extended ansi characters, like blocks and lines...
- Re-edit the sig before exiting
- match the editor's style to Talisman's internal full-screen editor

