package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func main() {
	logger.SetOutput(os.Stdout)
	file, err := os.OpenFile(".log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()
	logger.SetOutput(file)

	http.HandleFunc("/sum", Sum)

	logWithPort := logrus.WithFields(logrus.Fields{
		"port": 8080,
	})
	logWithPort.Info("Starting a web-server on port")
	logWithPort.Fatal(http.ListenAndServe(":8080", nil))
}

func Sum(w http.ResponseWriter, r *http.Request) {
	x, err := strconv.Atoi(r.URL.Query().Get("x"))
	if err != nil {
		logger.Println(err)
	}

	y, err := strconv.Atoi(r.URL.Query().Get("y"))
	if err != nil {
		logger.Println(err)
	}

	if x+y < 0 {
		logger.WithFields(logrus.Fields{
			"x": x,
			"y": y,
		}).Warning("Sum overflows int")
		w.Write([]byte("-1"))
	} else {
		w.Write([]byte(fmt.Sprintf("%d", x+y)))
	}
}