package main

import (
	"regexp"
	"strconv"
	"strings"
)

type TextDataExtractor interface {
    GenerateDataFromPdf(reader Reader) error
    GetGradesReport() *StudentReport
}

type ZcDataExtractor struct {
    GradesReport *StudentReport
    StudentName string
}

func (zce *ZcDataExtractor) GenerateDataFromPdf(reader Reader) error {
    text, err := reader.readLines()
    if err != nil {
        return err
    }

    lines := strings.Split(text, "\n")

    zce.StudentName = lines[6]

    report := NewStudentReport(zce.StudentName, "", make(map[string]map[string]*CourseGrade))

    idx := 0

    for idx < len(lines) {
        if zce.isTerm(lines[idx]) {
            term := lines[idx]
            report.AddNewTerm(term)
            idx += 8
            for idx < len(lines) && zce.isCourse(lines[idx]) {
                courseCode := strings.ToUpper(lines[idx])
                courseName := lines[idx+1]
                grade := strings.ToUpper(lines[idx+3])
                creditsFloat, _ :=  strconv.ParseFloat(lines[idx+4], 64)
				credits := int(creditsFloat)
				courseGrade := CourseGrade{
					CourseCode: courseCode,
					Name: courseName,
					Grade: grade,
					Credits: credits,
				}
                report.Grades[term][courseCode] = &courseGrade
                idx += 6
            }
        } else {
            idx++
        }
    }

    report.PrintReport()

    zce.GradesReport = report

    return nil
}

func (zce *ZcDataExtractor) GetGradesReport() *StudentReport {
    return zce.GradesReport
}

func (zce *ZcDataExtractor) isTerm(text string) bool {
    pattern := regexp.MustCompile(`\b\d{4}\s(Spring|Summer|Fall)\b`)
    return pattern.MatchString(text)
}

func (zce *ZcDataExtractor) isCourse(text string) bool {
    pattern := regexp.MustCompile(`\b[a-zA-Z]{3,4}\s\d{3}\b`)
    return pattern.MatchString(text)
}