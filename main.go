package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/flosch/pongo2"
	"github.com/jung-kurt/gofpdf"
	gomail "gopkg.in/gomail.v2"
)

type target struct {
	name  string
	email string
}

func (t *target) to() string {
	return fmt.Sprintf("%s <%s>", t.name, t.email)
}

var (
	file  *string
	dummy *bool
	send  *bool
	conf  app
)

func init() {
	file = flag.String("file", "", "data file")
	dummy = flag.Bool("dummy", true, "send dummy to testing email")
	send = flag.Bool("send", false, "send to participant")
}

func main() {
	flag.Parse()
	if *file == "" {
		log.Println("Please specify file")
		os.Exit(1)
	}
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		log.Println("Configuration file err:", err.Error())
		os.Exit(1)
	}
	f, err := os.Open(*file)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	r := csv.NewReader(f)
	rows, _ := r.ReadAll()
	var targets []target
	for _, row := range rows {
		targets = append(targets, target{row[0], row[1]})
	}
	if *send {
		for _, t := range targets {
			if conf.dummy.enable {
				t = target{t.name, conf.dummy.target}
			}
			msg := "Sending email to " + t.to() + " "
			if err := sendmail(t); err != nil {
				msg += "failed: " + err.Error()
				log.Println(msg)
				continue
			}
			msg += "gotcha!"
			log.Println(msg)
		}
	}
}

func sendmail(t target) error {
	if err := Generate(t.name); err != nil {
		return err
	}
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "Support from Valutac <support@valutac.com>")
	if strings.ContainsAny(t.name, ", & ) .") {
		mailer.SetHeader("To", t.email)
	} else {
		mailer.SetHeader("To", t.to())
	}
	mailer.SetHeader("Subject", "Certificate Meetup X")

	buff, err := ioutil.ReadFile("email.html")
	if err != nil {
		return err
	}
	tpl, err := pongo2.FromString(string(buff))
	if err != nil {
		return err
	}
	params := pongo2.Context{
		"name": t.name,
	}
	out, err := tpl.Execute(params)
	if err != nil {
		return err
	}
	mailer.SetBody("text/html", out)
	mailer.Attach(fmt.Sprintf("pdf/%s.pdf", t.name))

	d := gomail.NewDialer(
		conf.email.host,
		conf.email.port,
		conf.email.username,
		conf.email.password,
	)

	if err := d.DialAndSend(mailer); err != nil {
		return err
	}

	return nil
}

func Generate(name string) error {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetHeaderFunc(func() {
		pdf.ImageOptions("template.png", 0, 0, 297, 0, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
		pdf.AddFont("ChopinScript", "", "ChopinScript.json")
		pdf.SetTextColor(13, 116, 185)
		if len(strings.Split(name, " ")) > 3 {
			pdf.SetFont("ChopinScript", "", 48)
			pdf.SetY(33)
		} else {
			pdf.SetFont("ChopinScript", "", 72)
			pdf.SetY(33)
		}
		html := pdf.HTMLBasicNew()
		html.Write(36, fmt.Sprintf("<center>%s</center>", name))
	})
	return pdf.OutputFileAndClose(fmt.Sprintf("pdf/%s.pdf", name))

}
