package main

import (
	"fmt"
	"fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"image/color"
	"strconv"
	"os"
)

type GuiHandler struct {
	report *StudentReport
	oldReport *StudentReport
	analyzer *ReportAnalyzer
	app fyne.App
	window fyne.Window
	tabsContainer *container.AppTabs
	tabs []*container.TabItem
	//tabsDict map[string]*container.TabItem
	gpaLabel *widget.Label
	termGpaLabel *widget.Label
	specialGpaLabel *widget.Label
	credsLabel *widget.Label
}

func NewGuiHandler() *GuiHandler {
	app := app.New()
    win := app.NewWindow("Student Report")
	tabs := []*container.TabItem{}

	gpaLabel := widget.NewLabelWithStyle("" , fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	termGpaLabel := widget.NewLabelWithStyle("" , fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	specialGpaLabel := widget.NewLabelWithStyle("" , fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	credsLabel := widget.NewLabelWithStyle("" , fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	return &GuiHandler{
		app: app,
		window: win,
		tabs: tabs,
		gpaLabel: gpaLabel,
		termGpaLabel: termGpaLabel,
		specialGpaLabel: specialGpaLabel,
		credsLabel: credsLabel,
	}
}

func (gh *GuiHandler) openFile()  {
	var filename string
	var err error

	// Create a channel to signal that the user has selected a file
	

	fileDialog := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, gh.window)
			return
		}
		if file == nil {
			// File not selected
			dialog.ShowCustom("Error", "Ok",widget.NewLabel("No file selected") , gh.window)
			return
		}

		// Get the file name
		filename = file.URI().Path()
		fmt.Println("Selected file name:", filename)
		zc := ZcDataExtractor{}
		zc.GenerateDataFromPdf(Reader{ fileName: filename })
		report := zc.GetGradesReport()
		gh.InitDisplay(report)


		
	}, gh.window)

	// Set the filter to only accept PDF files
	fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".pdf"}))

	currentDir, err := os.Getwd()
	if err != nil {
		dialog.ShowError(err, gh.window)

	}
	listableURI, err := storage.ListerForURI(storage.NewFileURI(currentDir))
	if err != nil {
		dialog.ShowError(err, gh.window)

	}
	fileDialog.SetLocation(listableURI)
	// Show the file dialog
	fileDialog.Show()

}


func (gh *GuiHandler) Start()  {
	gh.window.SetContent(container.NewAdaptiveGrid(
		1,
		widget.NewButton("Open transcript file", func() {
			gh.openFile()		
		}) ,
		widget.NewButton("Create transcript manually", func() {
			report := NewStudentReport("Custom", "", make(map[string]map[string]*CourseGrade))
			gh.InitDisplay(report)
		}),
	))
	gh.window.Resize(fyne.NewSize(1300, 630))
	gh.window.CenterOnScreen()
	gh.window.ShowAndRun()
}
func (gh *GuiHandler) InitDisplay(report *StudentReport) {
	gh.report = report
	gh.oldReport = report.Copy()
	gh.analyzer = &ReportAnalyzer{ Report: gh.report }
	
	topBar := container.NewAdaptiveGrid( 
		2,
		container.NewAdaptiveGrid(
			1, 
			widget.NewLabelWithStyle(gh.report.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		),
		container.NewAdaptiveGrid(
			2,
			widget.NewButton("Light" , func() {
				gh.app.Settings().SetTheme(theme.LightTheme())
				// ugly hack for now
				gh.tabsContainer.OnSelected = func(tab *container.TabItem) {
					gh.UpdateSemesterGpa()
					gh.app.Settings().SetTheme(theme.LightTheme())
				}
			}),
			widget.NewButton("Dark" , func() {
				gh.app.Settings().SetTheme(theme.DarkTheme())
				// ugly hack for now
				gh.tabsContainer.OnSelected = func(tab *container.TabItem) {
					gh.UpdateSemesterGpa()
					gh.app.Settings().SetTheme(theme.DarkTheme())
				}
			}),
		),
	)
	
	gh.DrawTabs()

	gh.tabsContainer = container.NewAppTabs(gh.tabs...)
	gh.termGpaLabel.SetText(fmt.Sprintf("Term GPA: %.3f", gh.analyzer.GetTermGpa(gh.tabs[0].Text, &ZcCourseGradeMapper{})))
	gh.tabsContainer.OnSelected = func(tab *container.TabItem) {
		gh.UpdateSemesterGpa()
	}
	bottomBar := container.NewAdaptiveGrid(5, gh.termGpaLabel ,gh.gpaLabel,gh.specialGpaLabel ,gh.credsLabel ,widget.NewButton("Reset", func() {
		gh.Reset()
	}))
	content := container.NewVBox(topBar ,gh.tabsContainer,canvas.NewLine(color.Gray{50}) ,bottomBar)
	gh.window.SetContent(content)
}

func (gh *GuiHandler) UpdateGpa() {
	gpa := gh.analyzer.GetGpa( &ZcCourseGradeMapper{} )
	specialGpa , err := gh.analyzer.GetSpecialGpa( &ZcCourseGradeMapper{})
	creds := gh.analyzer.GetAcquiredCredits()
	gh.gpaLabel.SetText(fmt.Sprintf("Cumulative GPA: %.3f", gpa))

	gh.credsLabel.SetText(fmt.Sprintf("Acquired Credits: %d", creds))
	if  err == nil {
		gh.specialGpaLabel.SetText(fmt.Sprintf("Special GPA: %.3f", specialGpa))
	}
	
}

func (gh *GuiHandler) DrawTabs() {
	gh.UpdateGpa()
	terms := gh.report.SemestersOrdered
	for _,term := range terms {
		grades := gh.report.Grades[term]
		//fmt.Println(term)
		tab := gh.NewTab(term, grades)
		gh.tabs = append(gh.tabs, tab)
	}

	newTermEntry := widget.NewEntry()
	appendTab := container.NewTabItem(
		"New Term",
		container.NewVBox(
			widget.NewLabelWithStyle("Term title" , fyne.TextAlignCenter, fyne.TextStyle{Bold: true}) ,
			newTermEntry,
			widget.NewButton("Add Term", func() {
				if gh.report.AddNewTerm(newTermEntry.Text) {
					gh.RedrawTabs()
					gh.tabsContainer.Select(gh.tabs[len(gh.tabs)-2])
					
				}
			}),
		),
	)
	gh.tabs = append(gh.tabs, appendTab)
}

func (gh *GuiHandler) RedrawTabs() {
	gh.tabs = [] *container.TabItem{}
	gh.DrawTabs()
	gh.UpdateSemesterGpa()
	gh.tabsContainer.Items = gh.tabs
}
func (gh *GuiHandler) NewTab(term string , grades map[string]*CourseGrade) *container.TabItem {
	items := []fyne.CanvasObject{}
	topBar := container.NewAdaptiveGrid(
		4, 
		widget.NewLabelWithStyle("Course Code", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Course name", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Credits", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Grade", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)
	
	items = append(items, topBar)
		for courseCode, grade := range grades {
			item := gh.NewCourseRow( term,courseCode, grade)
			items = append(items, item)
		}

		newCourseCodeEntry := widget.NewEntry()
		
		newCourseNameEntry := widget.NewEntry()
		
		newCourseGradeEntry := widget.NewSelect([]string{ "A", "A-", "B+", "B", "B-", "C+", "C", "C-", "F"}, func(s string) {
		})
		
		newCourseCreditsEntry := widget.NewEntry()
		
		addNewCourseButton := widget.NewButton("Add Course", func() {
			creds , err := strconv.Atoi(newCourseCreditsEntry.Text)
			if err != nil {
				dialog.ShowError(err, gh.window)
				return
			}
			newCourseGrade := &CourseGrade{
				CourseCode: newCourseCodeEntry.Text,
				Name: newCourseNameEntry.Text,
				Grade: newCourseGradeEntry.Selected,
				Credits:  creds,
			}

			gh.AddCourseToTerm(term, newCourseGrade)
			
		})
		addLabels := container.NewAdaptiveGrid(
			5 ,
			widget.NewLabel("Course code") ,
			widget.NewLabel("Course name"),
			widget.NewLabel("Credits"),
			widget.NewLabel("Grade"),
		)
		cont := container.NewAdaptiveGrid(
			5,
			newCourseCodeEntry,
			newCourseNameEntry,
			newCourseCreditsEntry,
			newCourseGradeEntry,
			addNewCourseButton)
			
		items = append(items, container.NewVBox(
			canvas.NewLine(color.Gray{50}) ,
			widget.NewLabelWithStyle("Add new" , fyne.TextAlignCenter, fyne.TextStyle{Bold: true}), 
			addLabels ,
			cont,
			))
			
		return container.NewTabItem(term, container.NewVBox(items...))
		
}

func (gh *GuiHandler) NewCourseRow(term string, courseCode string, grade *CourseGrade) *fyne.Container {
	thatCourse := courseCode
	thatSemester := term
	//courseLabel := widget.NewLabel(fmt.Sprintf("%-10s (%s)  Credits: %d", courseCode, grade.Name, grade.Credits))
	CourseCodeLabel := widget.NewLabelWithStyle(courseCode , fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	courseNameLabel := widget.NewLabel(grade.Name)
	courseCreditsLabel := widget.NewLabel(fmt.Sprintf("%d", grade.Credits))
	gradeEntry := widget.NewSelect([]string{ "A", "A-", "B+", "B", "B-", "C+", "C", "C-", "F"}, func(s string) {
		newGradeLetter := s
		gh.report.ModifySemesterCourseGrade( thatSemester,thatCourse, newGradeLetter)
		gh.UpdateGpa()
		gh.UpdateSemesterGpa()
	})
	gradeEntry.Selected = grade.Grade
	
	return  container.NewAdaptiveGrid(4,CourseCodeLabel,courseNameLabel ,courseCreditsLabel,gradeEntry)
}
	
// AddCourseToTerm adds a course grade to the report for the given term and
// redraws the tabs if successful.
//
// term string: the term to add the course to.
// courseGrade *CourseGrade: the course grade to add to the report.
func (gh *GuiHandler) AddCourseToTerm(term string , courseGrade *CourseGrade) {
	if gh.report.AddCourseGrade(term , courseGrade) {
		gh.RedrawTabs()
		tabsIdx := Find(gh.report.SemestersOrdered, term)
		gh.tabsContainer.Select(gh.tabs[tabsIdx])
	}
}

// Reset resets the GuiHandler to its initial state by restoring the
// original report and redrawing the tabs.
//
// No parameters.
// No return type.
func (gh *GuiHandler) Reset() {
	gh.report = gh.oldReport.Copy()
	gh.analyzer = &ReportAnalyzer{ Report: gh.report }
	gh.RedrawTabs()
}

func (gh *GuiHandler) UpdateSemesterGpa()  {
	tab := gh.tabsContainer.CurrentTab()
	term := tab.Text
	termGpa := gh.analyzer.GetTermGpa( term ,&ZcCourseGradeMapper{})
	gh.termGpaLabel.SetText(fmt.Sprintf("Term GPA: %.3f", termGpa))
}


