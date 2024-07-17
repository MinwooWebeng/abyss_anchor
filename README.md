## Abyss anchor
An example of abyss identity anchor server

# Features

* STUN-like public address telling
* public identity board
* no persistence; peer info available only if alive

# API

GET https://0.0.0.0/id?name=$name </br>
GET https://0.0.0.0/stun </br>
POST https://0.0.0.0/reg </br>
name:mallang </br>
sha3-256:****************************** </br>
 </br>
loopback:1605 </br>
192.168.0.2:1605 </br>
 </br>
Hash value is base58 encoded. Each line contains locally detected ip:port address.
