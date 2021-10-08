# mytools

A tool for converting log file into plain text or JSON.
Requires Go >= 1.16 to build this tool.

## Usage

```
$ ./mytools [flags] path/to/file.log
```

Available flags:
- **t**: output type, available values are text and json, default is text if this flag not provided
- **o**: file path for the output

### Examples

```
$ ./mytools /var/log/file.log
```
By default it will print the output as plain text

```
$ ./mytools -t text /var/log/file.log
```
It will print the output as plain text
```
$ ./mytools -t json /var/log/file.log
```
It will print the output as JSON
```
$ ./mytools -o /path/to/outputfile.txt /var/log/file.log
```
It will save the output as plain text to /path/to/outputfile.txt
```
$ ./mytools -t json -o /path/to/outputfile.json /var/log/file.log
```
It will save the output as JSON to /path/to/outputfile.json
