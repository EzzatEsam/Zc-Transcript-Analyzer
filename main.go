package main



func main() {
	gh := NewGuiHandler()
	gh.Start()
}

// import (
// 	"fmt"
// 	"fyne.io/fyne/v2"
//     "fyne.io/fyne/v2/app"
//     "fyne.io/fyne/v2/container"
//     "fyne.io/fyne/v2/widget"

// )

// func main() {
	
// 	zc := ZcDataExtractor{}
// 	zc.GenerateDataFromPdf( Reader{ fileName: "UnofficialTranscript.pdf" } )
// 	report := zc.GetGradesReport()

// 	analyzer := ReportAnalyzer{ Report: report }
// 	gpa := analyzer.GetGpa( &ZcCourseGradeMapper{} )
// 	println(gpa)
// 	app := app.New()
//     win := app.NewWindow("Student Report")

//     var tabs []*container.TabItem

//     for term, grades := range report.Grades {
// 		var items []fyne.CanvasObject
	
// 		for courseCode, grade := range grades {
// 			courseLabel := widget.NewLabel(fmt.Sprintf("%-10s (%s)  Credits: %d", courseCode, grade.Name, grade.Credits))
// 			gradeEntry := widget.NewEntry()
// 			gradeEntry.Text = grade.Grade
// 			thatCourse := courseCode
// 			thatSemester := term
// 			updateButton := widget.NewButton("Update", func() {
// 				newGradeLetter := gradeEntry.Text
				
// 				fmt.Println(thatSemester, thatCourse, newGradeLetter)
// 				newCourseGrade := report.ModifySemesterCourseGrade( thatSemester,thatCourse, newGradeLetter)

				
// 				courseLabel.SetText(fmt.Sprintf("%-10s (%s) , Credits: %d", thatCourse, newCourseGrade.Name, newCourseGrade.Credits))
// 				//gpa := analyzer.GetGpa( &ZcCourseGradeMapper{} )
// 				//gpaLabel.SetText(fmt.Sprintf("GPA: %.2f", gpa))
// 			})
// 			fmt.Println(courseLabel )
// 			fmt.Println(gradeEntry )
// 			fmt.Println(updateButton )
// 			item := container.NewHBox(courseLabel, gradeEntry, updateButton)
// 			fmt.Println("a7a")
// 			items = append(items, item)
// 		}
// 		addNewCourseButton := widget.NewButton("Add Course", func() {
// 			// Add new course code here
// 		})
// 		items = append(items, addNewCourseButton)
// 		tab := container.NewTabItem(term, container.NewVBox(items...))
// 		tabs = append(tabs, tab)
// 	}

//     tabsContainer := container.NewAppTabs(tabs...)
//     win.SetContent(tabsContainer)
//     win.Resize(fyne.NewSize(800, 600))
//     win.ShowAndRun()
	

// }





