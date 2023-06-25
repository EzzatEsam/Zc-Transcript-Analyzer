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