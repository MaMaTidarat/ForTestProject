package main

import (
	"io"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func processExcelFile(file io.Reader) (map[string]interface{}, error) {
	excelFile, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}
	defer excelFile.Close()

	rows, err := excelFile.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	finalResult := make(map[string]interface{})

	planValues := []int{}

	for _, row := range rows[1:] { // Skipping header row
		if len(row) < 2 {
			continue
		}

		dataCell := row[0]
		factorCell := row[1]

		if strings.Contains(dataCell, "-") {
			parts := strings.Split(dataCell, "-")
			from, _ := strconv.Atoi(parts[0])
			to, _ := strconv.Atoi(parts[1])

			fieldName := strings.ToLower(factorCell)
			entry := map[string]int{
				"from": from,
				"to":   to,
			}

			if _, exists := finalResult[fieldName]; !exists {
				finalResult[fieldName] = []map[string]int{}
			}
			finalResult[fieldName] = append(finalResult[fieldName].([]map[string]int), entry)
		} else if strings.ToLower(factorCell) == "plan" {
			plan, _ := strconv.Atoi(dataCell)
			planValues = append(planValues, plan)
		}
	}

	// Handle multiple plan values
	if len(planValues) > 0 {
		finalResult["plan"] = planValues
	}

	return finalResult, nil
}
