package main
import "fmt"
type CourseGrade struct {
    CourseCode string
    Name string
    Grade string
    Credits int
}



func (cg *CourseGrade) String() string {
    return fmt.Sprintf("%s %s %s %d", cg.CourseCode, cg.Name, cg.Grade, cg.Credits)
}

func (cg *CourseGrade) Copy() *CourseGrade {
    return &CourseGrade{
        CourseCode: cg.CourseCode,
        Name: cg.Name,
        Grade: cg.Grade,
        Credits: cg.Credits,
    }
}