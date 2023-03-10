package lessons

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"unicode"
)

func GenerateSelfStory(name string, age int, money float64) string {
	return fmt.Sprintf("Hello! My name is %s. I'm %d y.o. And I also have $%.2f in my wallet right now.", name, age, money)
}

func SumWorker(numsCh chan []int, sumCh chan int) {
	for ch := range numsCh {
		sum := 0
		for _, v := range ch {
			sum += v
		}
		sumCh <- sum
	}
}

func MaxSum(nums1, nums2 []int) []int {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		sum := 0
		for _, v := range nums1 {
			sum += v
		}
		ch1 <- sum
		close(ch1)
	}()

	go func() {
		sum := 0
		for _, v := range nums2 {
			sum += v
		}
		ch2 <- sum
		close(ch2)
	}()

	sum1, sum2 := <-ch1, <-ch2

	if sum1 >= sum2 {
		return nums1
	}
	return nums2
}

type MergeDictsJob struct {
	Dicts      []map[string]string
	Merged     map[string]string
	IsFinished bool
}

var (
	ErrNotEnoughDicts = errors.New("at least 2 dictionaries are required")
	ErrNilDict        = errors.New("nil dictionary")
)

func ExecuteMergeDictsJob(job *MergeDictsJob) (*MergeDictsJob, error) {
	job.IsFinished = true
	job.Merged = map[string]string{}
	if len(job.Dicts) < 2 {
		return job, ErrNotEnoughDicts
	}
	for _, v := range job.Dicts {
		if v == nil {
			return job, ErrNilDict
		}
		for key, val := range v {
			job.Merged[key] = val
		}
	}
	return job, nil
}

func GetErrorMsg(err error) string {
	if errors.Is(err, ErrBadConnection) {
		return ErrBadConnection.Error()
	}
	if errors.Is(err, ErrBadRequest) {
		return ErrBadRequest.Error()
	}
	if errors.As(err, &NonCriticalError{}) {
		return ""
	}
	return UnknownErrorMsg
}

type NonCriticalError struct{}

func (e NonCriticalError) Error() string {
	return "validation error"
}

const UnknownErrorMsg = "unknown error"

var (
	ErrBadConnection = errors.New("bad connection")
	ErrBadRequest    = errors.New("bad request")
)

type CreateUserRequest struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

var (
	ErrEmailRequired                = errors.New("email is required")
	ErrPasswordRequired             = errors.New("password is required")
	ErrPasswordConfirmationRequired = errors.New("password confirmation is required")
	ErrPasswordDoesNotMatch         = errors.New("password does not match with the confirmation")
)

func DecodeAndValidateRequest(requestBody []byte) (CreateUserRequest, error) {
	user := CreateUserRequest{}
	err := json.Unmarshal(requestBody, &user)
	if err != nil {
		return CreateUserRequest{}, err
	}
	if err = user.isValid(); err != nil {
		return CreateUserRequest{}, err
	}
	return user, nil
}

func (u *CreateUserRequest) isValid() error {
	if u.Email == "" {
		return ErrEmailRequired
	}
	if u.Password == "" {
		return ErrPasswordRequired
	}
	if u.PasswordConfirmation == "" {
		return ErrPasswordConfirmationRequired
	}
	if u.Password != u.PasswordConfirmation {
		return ErrPasswordDoesNotMatch
	}
	return nil
}

type Person struct {
	Age uint8
}

type PersonList []Person

func (pl PersonList) GetAgePopularity() map[uint8]int {
	pop := map[uint8]int{}
	for _, v := range pl {
		if _, ok := pop[v.Age]; ok {
			pop[v.Age]++
			continue
		}
		pop[v.Age] = 1
	}
	return pop
}

type ListNode struct {
	Next *ListNode
	Val  int
}

func (head *ListNode) Reverse() *ListNode {
	pCurr := head
	var pTop *ListNode = nil
	for {
		if pCurr == nil {
			break
		}
		pTemp := pCurr.Next
		pCurr.Next = pTop
		pTop = pCurr
		pCurr = pTemp
	}

	return pTop
}

type Parent struct {
	Name     string
	Children []Child
}

type Child struct {
	Name string
	Age  int
}

func CopyParent(p *Parent) Parent {
	if p == nil {
		return Parent{}
	}
	var child []Child
	for _, v := range p.Children {
		child = append(child, v)
	}
	return Parent{
		Name:     p.Name,
		Children: child,
	}
}

func MergeNumberLists(numberLists ...[]int) []int {
	result := []int{}
	for _, v := range numberLists {
		if len(v) == 0 {
			continue
		}
		result = append(result, v...)
	}
	return result
}

func LatinLetters(s string) string {
	var newStr []rune
	for _, v := range s {
		if unicode.Is(unicode.Latin, v) {
			newStr = append(newStr, v)
		}
	}
	return string(newStr)
}

func IsASCII(s string) bool {
	for _, v := range []byte(s) {
		if v > 127 {
			return false
		}
	}
	return true
}

func ShiftASCII(s string, step int) string {
	var newStr []byte
	for _, v := range []byte(s) {
		newStr = append(newStr, v+byte(step))
	}
	return string(newStr)
}

func MostPopularWord(words []string) string {
	m := map[string]int{}
	for _, word := range words {
		if _, ok := m[word]; ok {
			m[word]++
			continue
		}
		m[word] = 1
	}
	max, index := 0, -1
	for i, word := range words {
		if max < m[word] {
			max = m[word]
			index = i
		}
	}
	return words[index]
}

func UniqueSortedUserIDs(userIDs []int64) []int64 {
	sort.Slice(userIDs, func(i, j int) bool {
		return userIDs[i] < userIDs[j]
	})
	var res []int64
	for i := 0; i < len(userIDs)-1; i++ {
		if userIDs[i] == userIDs[i+1] {
			continue
		}
		res = append(res, userIDs[i])
	}
	return res
}

func modifySlice(nums []int) {
	nums[0] = 2
	nums[1] = 1
	nums[2] = 0
	nums = append(nums, 4)
}

func isPalindrome(x int) bool {
	var sl []int
	for x > 0 {
		sl = append(sl, x%10)
		x /= 10
	}
	l := len(sl)
	for i := 0; i < l; i++ {
		if sl[i] != sl[l-i] {
			return false
		}
	}
	return true
}

func Map(strs []string, mapFunc func(s string) string) []string {
	res := make([]string, len(strs))
	for i, v := range strs {
		res[i] = mapFunc(v)
	}
	return res
}

const (
	OK = iota
	CANCELLED
	UNKNOWN
)

func ErrorMessageToCode(msg string) int {
	switch msg {
	case "OK":
		return OK
	case "CANCELLED":
		return CANCELLED
	default:
		return UNKNOWN
	}
}

const (
	Foo = 42
	Bar = 12
)

func converter(n int) string {
	dict := map[int]string{
		1: "I",
		2: "II",
		3: "III",
		4: "IV",
	}
	return dict[n]
}

type UserCreateRequest struct {
	FirstName string // не может быть пустым; не может содержать пробелы
	Age       int    // не может быть 0 или отрицательным; не может быть больше 150
}

func Validate(req UserCreateRequest) string {
	if strings.Contains(req.FirstName, " ") || len(req.FirstName) < 1 || req.Age < 1 || req.Age > 150 {
		return "invalid request"
	}
	return ""
}
