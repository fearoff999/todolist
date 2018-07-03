package todolist

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/NonerKao/color-aware-tabwriter"
	"github.com/fatih/color"
)

type ScreenPrinter struct {
	Writer *tabwriter.Writer
}

func NewScreenPrinter() *ScreenPrinter {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	formatter := &ScreenPrinter{Writer: w}
	return formatter
}

func (f *ScreenPrinter) Print(groupedTodos *GroupedTodos, printNotes bool) {
	blue := color.New(color.FgHiBlue).SprintFunc()

	var keys []string
	for key := range groupedTodos.Groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		fmt.Fprintf(f.Writer, "\n %s\n", blue(key))
		for _, todo := range groupedTodos.Groups[key] {
			f.printTodo(todo)
			if printNotes {
				for nid, note := range todo.Notes {
					fmt.Fprintf(f.Writer, "   %s\t%s\t\n",
						blue(strconv.Itoa(nid)), note)
				}
			}
		}
	}
	f.Writer.Flush()
}

func (f *ScreenPrinter) printTodo(todo *Todo) {
	white := color.New(color.FgWhite)
	if todo.IsPriority {
		white.Add(color.Bold, color.Italic)
	}
	fmt.Fprintf(f.Writer, " %s\t%s\t%s\t %s\n",
		white.SprintFunc()(strconv.Itoa(todo.Id)),
		f.formatCompleted(todo.Completed),
		f.formatDue(todo.Due, todo.IsPriority),
		f.formatSubject(todo.Subject, todo.IsPriority))
}

func (f *ScreenPrinter) formatDue(due string, isPriority bool) string {
	yellow := color.New(color.FgYellow)
	hiYellow := color.New(color.FgHiYellow)
	red := color.New(color.FgRed)
	hiRed := color.New(color.FgHiRed)
	white := color.New(color.FgWhite)
	green := color.New(color.FgGreen)
	hiGreen := color.New(color.FgHiGreen)

	if isPriority {
		yellow.Add(color.Bold, color.Italic)
		red.Add(color.Bold, color.Italic)
		hiYellow.Add(color.Bold, color.Italic)
		hiRed.Add(color.Bold, color.Italic)
		green.Add(color.Bold, color.Italic)
		hiGreen.Add(color.Bold, color.Italic)
	}

	if due == "" {
		return white.SprintFunc()(" ")
	}
	dueTime, err := time.Parse("2006-01-02", due)

	if err != nil {
		fmt.Println(err)
		fmt.Println("This may due to the corruption of .todos.json file.")
		os.Exit(-1)
	}

	if isToday(dueTime) {
		return hiRed.SprintFunc()("today")
	} else if isTomorrow(dueTime) {
		return yellow.SprintFunc()("tomorrow")
	} else if isNextWeek(dueTime) {
		return green.SprintFunc()("nextweek")
	} else if isThisWeek(dueTime) {
		return hiYellow.SprintFunc()("thisweek")
	} else if isPastDue(dueTime) {
		return red.SprintFunc()(dueTime.Format("Mon Jan 2"))
	} else {
		return hiGreen.SprintFunc()(dueTime.Format("Mon Jan 2"))
	}
}

func (f *ScreenPrinter) formatSubject(subject string, isPriority bool) string {

	hiBlue := color.New(color.FgHiBlue)
	hiYellow := color.New(color.FgHiYellow)
	hiWhite := color.New(color.FgHiWhite)
	white := color.New(color.FgWhite)

	if isPriority {
		hiBlue.Add(color.Bold, color.Italic)
		hiYellow.Add(color.Bold, color.Italic)
		hiWhite.Add(color.Bold, color.Italic)
		white.Add(color.Bold, color.Italic)
	}

	splitted := strings.Split(subject, " ")
	projectRegex, _ := regexp.Compile(`\+[\p{L}\d_]+`)
	contextRegex, _ := regexp.Compile(`\@[\p{L}\d_]+`)

	coloredWords := []string{}

	for _, word := range splitted {
		if projectRegex.MatchString(word) {
			coloredWords = append(coloredWords, hiBlue.SprintFunc()(word))
		} else if contextRegex.MatchString(word) {
			coloredWords = append(coloredWords, hiWhite.SprintFunc()(word))
		} else {
			coloredWords = append(coloredWords, white.SprintFunc()(word))
		}
	}
	return strings.Join(coloredWords, " ")

}

func (f *ScreenPrinter) formatCompleted(completed bool) string {
	green := color.New(color.FgHiGreen)
	red := color.New(color.FgHiRed)

	if completed {
		return "[" + green.SprintFunc()("v") + "]"
	} else {
		return "[" + red.SprintFunc()("x") + "]"
	}
}
