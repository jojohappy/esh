package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	etcd string
	home string
	cwd string
	options map[string]string
}

func (config Config) parsePrompt(prompt string) string {
	s := -1
    var env []string
    prompt = strings.Trim(prompt, "\"")
    prompt = strings.Trim(prompt, " ")
    for i, v := range prompt {
        c := string(v)
        if s != -1 && (c == "$" || c == " " || c == "\n") {
            env = append(env, prompt[s:i])
            s = -1
        } else if c == "$" {
            s = i
        }
    }
    for _, v := range env {
        key := strings.TrimPrefix(v, "$")
        prompt = strings.Replace(prompt, v, config.options[key], 1)
    }
	return prompt
}

func eshLoop(config Config) {
	options := config.options
	prompt := config.parsePrompt(options["prompt"])
	for {
		fmt.Print(prompt + " ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		strings.TrimSpace(text)
	}
}

func initConfig(etcd string) *Config {
	home := os.Getenv("HOME")
	path := home + "/.eshrc"
	cwd := "keys"
	options := map[string]string{
		"HOME": home,
		"ETCD_HOST": etcd,
		"CWD": cwd,
	}

	inFile, err := os.Open(path)
	defer inFile.Close()
	if err != nil {
		newFile, error := os.Create(path)
		defer newFile.Close()
		if error != nil {
			fmt.Print(error)
		}
		newFile.WriteString("prompt=\"[@$ETCD_HOST $CWD]$\"\n")
		newFile.Sync()
		options["prompt"] = "\"[@$ETCD_HOST $CWD]$\""
	} else {
		scanner := bufio.NewScanner(inFile)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			if line := strings.TrimSpace(scanner.Text()); line != "" && strings.HasPrefix(line, "#") == false {
				parameters := strings.Split(line, "=")
				options[parameters[0]] = parameters[1]
			}
		}
	}
	config := Config{etcd, home, cwd, options}
	return &config
}

func main() {
	var etcd string
	flag.StringVar(&etcd, "etcd", "localhost:4001", "Etcd Host. default: localhost:4001")
	flag.Parse()

	config := initConfig(etcd)
	eshLoop(*config)
}
