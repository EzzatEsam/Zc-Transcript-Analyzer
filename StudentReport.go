package main

import (
	"fmt"
)

type StudentReport struct {
    Name string
    Major string
    SemestersOrdered []string
    Grades map[string]map[string]*CourseGrade
}

func NewStudentReport(name string, major string, grades map[string]map[string]*CourseGrade) *StudentReport {
    return &StudentReport{
        Name: name,
        Major: major,
        Grades: grades,
    }
}

// make a copy of the student report
func (sr * StudentReport) Copy() *StudentReport {
    newGrades := make(map[string]map[string]*CourseGrade)
    for term, courseGrades := range sr.Grades {
        newGrades[term] = make(map[string]*CourseGrade)
        for courseCode, courseGrade := range courseGrades {
            newGrades[term][courseCode] = courseGrade.Copy()
        }
    }
    return &StudentReport{
        Name: sr.Name,
        Major: sr.Major,
        SemestersOrdered: sr.SemestersOrdered,
        Grades: newGrades,
    }
}

func (sr *StudentReport) PrintReport() {
    fmt.Printf("Student Name: %s\n", sr.Name)
    fmt.Println(sr.Major)
    fmt.Printf("Total credits: %d\n", sr.GetTotalCredits())
    for term := range sr.Grades {
        fmt.Println(term)
        for _, courseGrade := range sr.Grades[term] {
            fmt.Println(courseGrade.String())
        }
        fmt.Println("-------------------------------------")
    }
}
func (sr *StudentReport) AddNewTerm(term string) bool {
    if sr.Grades[term] == nil{
        sr.SemestersOrdered = append(sr.SemestersOrdered, term)
        sr.Grades[term] = make(map[string]*CourseGrade)
        return true
    } 
    return false
}

func (sr *StudentReport) AddCourseGrade(semester string , courseGrade *CourseGrade) bool {
    if sr.Grades == nil {
        sr.Grades = make(map[string]map[string]*CourseGrade)
    }
    if sr.Grades[semester] == nil {
        sr.Grades[semester] = make(map[string]*CourseGrade)
    }
    if sr.Grades[semester][courseGrade.CourseCode] != nil {
        return false
    }
    sr.Grades[semester][courseGrade.CourseCode] = courseGrade
    return true
}

func (sr *StudentReport) ModifySemesterCourseGrade(semester string, courseCode string, newGrade string) *CourseGrade {
    sr.Grades[semester][courseCode].Grade = newGrade
    return sr.Grades[semester][courseCode]
}

func (sr *StudentReport) ModifyCourseGrade( courseCode string, newGrade string) {
    // loop through all semesters
	for _, termGrades := range sr.Grades {
		for _, courseGrade := range termGrades {
			if courseGrade.CourseCode == courseCode {
				courseGrade.Grade = newGrade
			}
		}
	}
}


func (sr *StudentReport) GetTotalCredits() int {
    total := 0
    for term := range sr.Grades {
        termTotal := 0
        for _, courseGrade := range sr.Grades[term] {
            termTotal += courseGrade.Credits
        }
        fmt.Printf("Term %s: %d\n", term, termTotal)
        total += termTotal
    }
    return total
}