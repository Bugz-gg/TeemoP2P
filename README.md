# Eirb2Share

## Git
```bash
    git clone ssh://git@thor.enseirb-matmeca.fr:2222/git/free-eirb2share-T4
```


Tracker:
`CTRL+c` : Exit the tracker.

Tracker's config.ini :
* `tracker-ip` sets the IP for the tracker if `tracker-ip-mode` is set as `0`
* `tracker-ip-mode` specifies which interface to use (`0` by default): 
* * `0`: value of `tracker-ip` (localhost by default)
* * `1`: loopback only
* * `2` (or any other int): all the interfaces (check `ifconfig`)
* `tracker-port` specifies which port to use (`9000` by default)

Example :
```config.ini
tracker-ip = 10.0.0.126
tracker-port = 9001
```
