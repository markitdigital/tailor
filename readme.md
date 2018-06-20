# Tailor
Tailor acts as a drop-in replacement for the official Microsoft [IIS Service Monitor](https://github.com/Microsoft/IIS.ServiceMonitor) application with the additional ability to forward the contents of multiple log files (such as the ones created by IIS) to `stdout`, so they can be seen with `docker logs` and more easily collected with log aggregation systems.

It will monitor a given Windows service and exit when that service's state changes from `SERVICE_RUNNING`. In the background, it will collect the contents of any designated log files and write them to `stdout` in a structured JSON format. You might say it "stitches" together logfiles, hence the name :)

## Usage
For simple entrypoint usage, run like you'd run ServiceMonitor:

```bash
$ Tailor.exe w3svc
```

To automatically forward the contents of a log file to `stdout`, provide the path to the logfile relative to the current working directory:

```bash
$ Tailor.exe w3svc /app/logs/log.txt
```

For convenience, you can provide a glob-style path to tail multiple log files like so:

```bash
$ ServiceMonitor.exe w3svc c:\inetpub\logs\LogFiles\*\*.log
```

Output uses a structured JSON format including a timestamp and source:

```
...
{"level":"debug","msg":"2018-06-01 20:33:56 ::1 GET /ExampleWebApplication - 80 - ::1 Mozilla/5.0+(Windows+NT+10.0;+Win64;+x64;+rv:60.0)+Gecko/20100101+Firefox/60.0 - 500 0 0 1\r","source":"c:\\inetpub\\logs\\LogFiles\\W3SVC1\\u_ex180601.log","time":"2018-06-01T14:34:21-06:00"}
```
