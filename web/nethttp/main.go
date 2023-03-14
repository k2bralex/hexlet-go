package main

import (
	"fmt"
	"net/http"
	"strconv"
)

var courses = map[int64]string{
	1: "Introduction to programming",
	2: "Introduction to algorithms",
	3: "Data structures",
}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/courses/description", CourseDescHandler)

	http.ListenAndServe(":8080", nil)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Go to /courses/description"))
}

func CourseDescHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("course_id"), 10, 0)
	if err != nil {
		fmt.Println(err.Error())
	}
	w.Write([]byte(courses[id]))
}
