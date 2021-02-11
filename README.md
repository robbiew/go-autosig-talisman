# autosig
Auto Signature Mod for [Talisman BBS](https://talismanbbs.com/)

This simple, go-based program allows you to add an Auto Sig editor to your Talisman BBS menus -- as a door. 

It incorporates a modified vesion of the [kilo editor](https://github.com/bediger4000/kilo-in-go). 

Note: it **won't compile on OSX**. In fact, I've only tested it on Linux 64 (Ubuntu 20.04). It's *possible* that it will compile on Windows and Raspbery Pi...

To set up:

- Make sure [Talisman](https://talismanbbs.com/) is installed and working 😃
- Compile the program (`go build .`) -- or, use the pre-built binary in this repo for Linux 64 only
- Create /bbs/doors/autosig directory
- Copy autosig (binary), header.ans, start.sh to the above directory
- Set executable permissions (`chmod +x autosig start.sh`)
- Add AutoSig item to a Talisman menu, and its edit path (e.g. this located in /bbs/menus/message.toml if you want it on the Message Menu)

```
[[menuitem]]
command = "RUNDOOR"
data = "/home/robbiew/bbs/doors/autosig/start.sh"
hotkey = "Z"
```

Talisman will automatically pass the node number to the start.sh file, and the autosig program will use this to grab the drop file (e.g. bbs/temp/$1/door.sys). The drop file contains that node's logged-in user id and name, and the id will be used to retrieve the "signature" row in Talisman's user database (sqlite3). If the 'signature' value doesn't exist, it will create it on save. If it exists, it will updated it.

For display purposes, the autosig program translates between color "pipe" codes (e.g. "|02") that Talisman uses internally, and actual ansi escape codes for display in the terminal program.

I have not tested for background colors yet.

TO DO:
- [x] Publish on github
- [ ] Test on Windows, Pi
- [ ] Add way to add extended ansi characters, like blocks and lines...
- [ ] Allow for re-editing the sig before exiting
- [ ] Match the editor's style to Talisman's internal full-screen editor

