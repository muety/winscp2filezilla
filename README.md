# winscp2filezilla

A tool that migrates your saved servers from [WinSCP](https://winscp.net/) to [FileZilla](https://filezilla-project.org/). To be precise, it will "only" migrate your folder structure and all entries including host name, user name, password, port and protocol (FTP / SFTP only).

[![Buy me a coffee](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://buymeacoff.ee/n1try)

## Usage
1. You need to tell WinSCP to save its config to a `.ini` file, instead of Windows registry (although this is way less secure!). Do so according to the screenshot below.
2. Either do `go get github.com/n1try/winscp2filezilla` or download the latest [release](https://github.com/n1try/winscp2filezilla/releases).
3. `winscp2filezilla --in C:\Users\You\AppData\Roaming\WinSCP.ini --out sites.xml` (arguments can be left out in order to use defaults)
4. Import `sites.xml` in FileZilla

![](https://anchr.io/i/4iRuH.png)

## Credits
Some parts of this code are borrowed from [anoopengineer/winscppasswd](https://github.com/anoopengineer/winscppasswd), so thank you for that!

## To Do
* Migrate `RemoteDir` and `LocalDir`
* Add tests

## License
MIT @ [Ferdinand Mütsch](https://ferdinand-muetsch.de)