The [Coveralls][] project can be updated with [goveralls][]. But currently Travis does not support [goveralls][]. So coverage needs to be run manually

```
$ goveralls -package="./flowdock" your_repos_coveralls_token
```

[goveralls]: https://github.com/mattn/goveralls
[Coveralls]: https://coveralls.io/r/wm/go-flowdock
