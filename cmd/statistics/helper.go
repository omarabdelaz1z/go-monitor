package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

var monthNames map[string]string = map[string]string{
	"01": "January",
	"02": "February",
	"03": "March",
	"04": "April",
	"05": "May",
	"06": "June",
	"07": "July",
	"08": "August",
	"09": "September",
	"10": "October",
	"11": "November",
	"12": "December",
}

func transformMonthName(months ...string) []string {
	var names []string

	for _, month := range months {
		names = append(names, fmt.Sprintf("%s - %s", month, monthNames[month]))
	}

	return names
}

func selectPrompt(label string, items ...string) (int, string, error) {
	ps := promptui.Select{
		Label: label,
		Items: items,
	}

	index, result, err := ps.Run()

	if err != nil {
		if err == promptui.ErrInterrupt {
			return -1, "", nil
		}

		return -1, "", fmt.Errorf("failed to run prompt: %w", err)
	}

	return index, result, nil
}
