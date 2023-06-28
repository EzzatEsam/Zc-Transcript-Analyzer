package main
import (
    "fmt"
    "strconv"
    "errors"
)

type ReportAnalyzer struct {
    Report *StudentReport
}


func  (ra *ReportAnalyzer)GetGpa(mapper AbstractCourseGradeMapper ) float64 {
    gradesDict := make(map[string]*CourseGrade)

    for _, sem := range ra.Report.SemestersOrdered {
        termGrades := ra.Report.Grades[sem]
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

func (ra * ReportAnalyzer) GetSpecialGpa(mapper AbstractCourseGradeMapper) (float64  , error){
    sems := ra.Report.SemestersOrdered
    firstYear , _ := GetDateFromCourseCode(sems[0])
    firstYear ++ // the start of second year
    idx := 1
    found := false
    for idx < len(sems) {
        sem := sems[idx]
        semYear , _ := GetDateFromCourseCode(sem)
        //fmt.Println(sem[4: len(sem)-1])
        if semYear > firstYear || (semYear == firstYear && sem[5:] == "Fall" ) {
            found = true
            break
        }
        idx++
            
    }
    if ! found {
        return -1 , errors.New("still in foundation")
    }

    gradesDict := make(map[string]*CourseGrade)

    for i := idx; i < len(sems); i++ {
        sem := sems[i]
        fmt.Println("special gpa semester" , sem)
        termGrades := ra.Report.Grades[sem]
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

    specialGpa := totalGrades / totalWeights

    fmt.Printf("Special Gpa : %.3f \n", specialGpa)
	return specialGpa , nil

}

func (ra *ReportAnalyzer) GetTermGpa( term string , mapper AbstractCourseGradeMapper) float64  {
    totalWeights := 0.0
    totalGrades := 0.0

    // just for testing
    // TODO remove later

    // gradesDict := make(map[string]*CourseGrade)
    // found := false
    // for _, sem := range ra.Report.SemestersOrdered {
    //     if found {
    //         break
    //     }
    //     if sem == term {
    //         found = true
    //     }
    //     termGrades := ra.Report.Grades[sem]
    //     for course, grade := range termGrades {
    //         gradesDict[course] = grade
    //     }
    // }

    

    // for _, course := range gradesDict {
    //     grade , weight := mapper.GetCorrespondingGradeGpa(course)
    //     totalGrades += grade * weight
    //     totalWeights += weight
    // }

    // test_gpa := totalGrades / totalWeights

    // fmt.Printf("test Gpa : %.3f \n", test_gpa)

    // totalWeights = 0.0
    // totalGrades = 0.0
    
    // end of testing
    // 
    for _, course := range ra.Report.Grades[term] {
        grade , weight := mapper.GetCorrespondingGradeGpa(course)
        totalGrades += grade * weight
        totalWeights += weight
    }

    

    gpa := totalGrades / totalWeights
    fmt.Print("Term :" , term)
    fmt.Printf("Gpa : %.3f \n", gpa)
	return gpa
}

func (ra *ReportAnalyzer) GetAcquiredCredits() int {
    gradesDict := make(map[string]*CourseGrade)

    for _, sem := range ra.Report.SemestersOrdered {
        termGrades := ra.Report.Grades[sem]
        for course, grade := range termGrades {
            gradesDict[course] = grade
        }
    }

    nonFulfillingGrades := []string{"F" ,"W" , "WP" , "WF" , "[F]" , "[W]" , "[WP]" ,"[WF]"}
    totalCreds := 0
    for _, course := range gradesDict {
        if !isIn(course.Grade, nonFulfillingGrades) {
            totalCreds += course.Credits
        }
    }
    return totalCreds
}

func GetDateFromCourseCode(courseCode string) (int , error) {
    return strconv.Atoi(courseCode[0:4])
}



