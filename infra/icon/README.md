## Install T-Bears Cli
Install python3.7:
```
brew install python@3.7
brew list | grep python
brew ls python@3.7
```
Add a soft link to /usr/local/bin/:
```
ln -s /usr/local/Cellar/python@3.7/3.7.9/bin/python3.7 /usr/local/bin/python3.7
python3.7 -V
```
Create a working directory
```
mkdir work
cd work
```
Create a Python virtual environment:
```
python3.7 -m venv venv37
```
Enter the virtual environment:
```
source venv37/bin/activate
```
Install T-Bears:
```
pip install tbears
```
## General Usage
Send a test JSON request:
```
curl -d '{"jsonrpc": "2.0", "method": "icx_getBlock", "id": 1234}' -H "Content-Type: application/json" -X POST http://127.0.0.1:9000/api/v3
```
## References
- [ICON JSON-RPC API v3](https://github.com/icon-project/icon-rpc-server/blob/master/docs/icon-json-rpc-v3.md)
