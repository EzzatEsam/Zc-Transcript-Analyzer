package main
import "fmt"

type ReportAnalyzer struct {
    Report *StudentReport
}


func  (ra *ReportAnalyzer)GetGpa(mapper AbstractCourseGradeMapper ) float64 {
    gradesDict := make(map[string]*CourseGrade)

    for _, termGrades := range ra.Report.Grades {
        for course, grade := range termGrades {
            gradesDict[course] = grade
        }
    }

    totalWeights := 0.0
    totalGrades := 0.0

    for _, course := range gradesDict {
        grade , weight := mapper.GetCorrespondingGradeGpa(course)
        totalGrades += grade * weight
        totalWeights += weight
    }

    gpa := totalGrades / totalWeights

    fmt.Printf("Gpa : %.3f \n", gpa)
	return gpa
}

func (ra *ReportAnalyzer) GetAcquiredCredits() int {
    gradesDict := make(map[string]*CourseGrade)

    for _, termGrades := range ra.Report.Grades {
        for course, grade := range termGrades {
            gradesDict[course] = grade
        }
    }

    nonFulfillingGrades := []string{"F" ,"W" , "WP" , "WF"}
    totalCreds := 0
    for _, course := range gradesDict {
        if !isIn(course.Grade, nonFulfillingGrades) {
            totalCreds += course.Credits
        }
    }
    return totalCreds
}
