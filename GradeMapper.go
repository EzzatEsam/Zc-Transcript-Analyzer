package main
import "strings"

type AbstractCourseGradeMapper interface {
    GetCorrespondingGradeGpa(grade *CourseGrade) (float64, float64)
}


type ZcCourseGradeMapper struct {}

func (m *ZcCourseGradeMapper) GetCorrespondingGradeGpa(grade *CourseGrade) (float64, float64) {
    weight := float64(grade.Credits) 

    gradeStr := strings.TrimSpace(grade.Grade)
    if gradeStr == "" {
        return 0.0, 0.0
    }

    var letter string
    if strings.HasPrefix(gradeStr, "[") && strings.HasSuffix(gradeStr, "]") {
        letter = strings.TrimPrefix(strings.TrimSuffix(gradeStr, "]"), "[")
    } else {
        letter = gradeStr
    }

    switch letter {
    case "A":
        return 4.0, weight
    case "A-":
        return 3.7, weight
    case "B+":
        return 3.3, weight
    case "B":
        return 3.0, weight
    case "B-":
        return 2.7, weight
    case "C+":
        return 2.3, weight
    case "C":
        return 2.0, weight
    case "C-":
        return 1.7, weight
    case "F":
        return 0.0, weight
    default:
        return 0.0, 0.0
    }

}
