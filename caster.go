package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type SMTPAuth struct {
	Sender string `json:"email"`
	Pass   string `json:"pass"`
	Server string `json:"smtp"`
	Port   string `json:"port"`
}

type Email struct {
	Spoof      string
	Recipients []string
	Subject    string
	Message    []byte
}

type SpamMail struct {
	Score       string
	TestName    string
	Description string
}

// used for homograph attack(s), cryillic is just the first thing I think of so educate me!!
var latinToCyrillic = map[rune]rune{
	'a': 'а',
	'e': 'е',
	'o': 'о',
	'p': 'р',
	'c': 'с',
	'y': 'у',
	'x': 'х',
}

func main() {
	target := flag.String("target", "", "specify target(s) email address [filename or seperated by commas]")
	spoof := flag.String("spoof", "", "specify address to spoof email from [keep spam in mind]")
	subject := flag.String("subject", "", "specify a subject to add to email")
	template := flag.String("template", "", "specify a template for the email")
	homograph := flag.Bool("homograph", false, "specify option to replace chars with cryillic")
	spamScore := flag.Bool("spamfilter", false, "enable to get a given templates spam score")
	help := flag.Bool("help", false, "show usage")
	flag.Parse()
	if *help {
		flag.Usage()
	}

	if *template == "" {
		flag.Usage()
		return
	}

	if *homograph {
		templateModified := replaceCharsInHTMLFile(*template)
		fmt.Println("[*] Template modified for homograph attack [cryillic]...")
		// save modified template to homograph_{...}
		err := os.WriteFile("homograph_"+*template, []byte(templateModified), 0644)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(" -> File Name: homograph_" + *template)

		return // nothing more nothing less
	}

	user := smtpLogin()
	auth := smtp.PlainAuth("", user.Sender, user.Pass, user.Server)

	var email Email // struct to hold message and shii

	email.Spoof = *spoof
	email.Subject = *subject

	if *spamScore { // spam score
		mailFilter(user, auth, email, *template)
		return
	}

	// send email(s)
	getRecipients(*target, &email)
	for _, recipient := range email.Recipients {
		err := buildMessage(&email, *template, recipient)
		if err != nil {
			log.Fatalf("[-] Error opening %s for targets\n -> %v", *template, err)
		}
		sendEmail(user, auth, email, recipient)
	}
}

// function to send the email(s)
func sendEmail(user SMTPAuth, auth smtp.Auth, email Email, recipient string) {
	fmt.Printf("\r[*] Sending email to %s\r", recipient)
	err := smtp.SendMail(fmt.Sprintf("%s:%s", user.Server, user.Port),
		auth, user.Sender,
		[]string{recipient},
		email.Message)
	if err != nil {
		log.Fatalf("[-] Error: Failed to send email to %s\n -> %v", recipient, err)
	} else {
		fmt.Printf("[+] Email sent successfully to %s!\n", recipient)
	}
}

// function to compose message with template and subject
func buildMessage(email *Email, templatePath string, recipient string) error {
	templateBytes, err := os.ReadFile(templatePath)
	if err != nil {
		return err
	}

	email.Message = []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n"+
			"%s",
		email.Spoof,
		recipient,
		email.Subject,
		string(templateBytes),
	))

	return nil
}

// function to get auth from config file
func smtpLogin() SMTPAuth {
	var user SMTPAuth
	smtpJson, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer smtpJson.Close()

	bytesJson, err := io.ReadAll(smtpJson)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(bytesJson, &user)

	return user
}

// function to handle the recipients infile or as args
func getRecipients(targets string, email *Email) {
	readFile, err := os.Open(targets)
	if err != nil { // error indicating it's not a file.. ugly asf I know..
		email.Recipients = strings.Split(targets, ",") // seperate by comma
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() { // append each line to file.. rough world out here..
		email.Recipients = append(email.Recipients, fileScanner.Text())
	}
}

// function to replace Latin chars with Cyrillic
func replaceWithCyrillic(input string) string {
	var result strings.Builder

	for _, char := range input {
		if cyrillicChar, found := latinToCyrillic[char]; found {
			result.WriteRune(cyrillicChar) // replace char with Cyrillic char
		} else {
			result.WriteRune(char) // write original char if no Cyrillic replacement found
		}
	}

	return result.String()
}

// function to replace chars with cryillic chars
func replaceCharsInHTMLFile(templatePath string) string {
	templateBytes, err := os.ReadFile(templatePath)
	if err != nil {
		log.Fatalf("[-] Failed to read template file: %v", err)
	}

	htmlContent := string(templateBytes)

	re := regexp.MustCompile(`<(h[1-6]|p)[^>]*>(.*?)</(h[1-6]|p)>`)

	modifiedContent := re.ReplaceAllStringFunc(htmlContent, func(match string) string {
		openingTagEnd := strings.Index(match, ">") + 1
		closingTagStart := strings.LastIndex(match, "<")
		textContent := match[openingTagEnd:closingTagStart]
		replacedText := replaceWithCyrillic(textContent)

		return match[:openingTagEnd] + replacedText + match[closingTagStart:]
	})

	return modifiedContent
}

// function to send email to mail-tester and output spam score
func mailFilter(user SMTPAuth, auth smtp.Auth, email Email, template string) error {

	tempMail := "test-sxzd09jk9@srv1.mail-tester.com" // mail-tester email.. but not sure if it'll last forever!
	tempServ := strings.Split(tempMail, "@")[0]
	url := fmt.Sprintf("https://www.mail-tester.com/%s", tempServ)

	err := buildMessage(&email, template, tempMail)
	if err != nil {
		log.Fatalf("[-] Error opening %s for targets\n -> %v", template, err)
	}
	sendEmail(user, auth, email, tempMail)

	time.Sleep(3 * time.Second) // wait for 3 seconds before getting score

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("[-] Failed to create request: %w", err)
	}

	client := &http.Client{}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.6613.120 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("[-] Failed to send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("[*] Status code: %v\n -> Try other request methods.. the antibot is cookin", resp.StatusCode)
	}

	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	mailTester := doc.Find("span.score").Text()
	spamAssassin := doc.Find("div.about").Text()

	re := regexp.MustCompile(`Score: ([\d\.\-]+)`)
	matches := re.FindStringSubmatch(spamAssassin)
	if len(matches) > 0 {
		spamAssassin = matches[1]
	}

	var spamResults []SpamMail
	doc.Find("tr.sa-test").Each(func(i int, t *goquery.Selection) {
		score := t.Find("td.sa-test-score").Text()
		testName := t.Find("td.sa-test-name samp").Text()
		description := t.Find("td.sa-test-description").Text()

		spamResults = append(spamResults, SpamMail{
			Score:       score,
			TestName:    testName,
			Description: description,
		})
	})

	fmt.Printf("[+] Mail-tester Score: %s\n    Spam Assassin Score: %s\n", mailTester, spamAssassin)

	fmt.Println("[*] Mail-Tester Results:")

	for _, key := range spamResults {
		fmt.Printf("Score: %s\tTest Name: %s\nDescription: %s\n\n", key.Score, key.TestName, key.Description)
	}

	return nil
}
