# Accent

Send certificate to your meetup or event participant

## BUILD & RUN

```
$ go build -o accent
$ ./accent -dummy=false -file=source.csv -send=true
```

### Parameters


- `dummy` has two value, true/false. Whenever the value is false, it will send to support@valutac.com. You can change the target on the code and rebuild the binary. Use this dummy option to check the email which will be sent to participant.
- `file` is the data input, please check `dummy.csv` file to see the format file.
- `send` if the value is true, it will send the certificate.


### Template

There is two template in this application:

- Email template `email.html`
- Certificate template `template.png`

## LICENSE

<a href="LICENSE">
<img src="mit.png" width="75"></img>
</a>
