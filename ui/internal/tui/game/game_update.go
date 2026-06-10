package tui_game

import (
	// "fmt"
	"strconv"
	
	Config "cheat-codex/internal/config"
)

// Helper function to parse uint from EditInput value and commit to memory map
func (model GameModel) updateUint() error {
	var err error = nil
	if model.TableRows[model.Cursor].Type == "uint16" {
		_, err = strconv.ParseUint(
			model.EditInput.Value(), 10, 16,
		)
	} else {
		_, err = strconv.ParseUint(
			model.EditInput.Value(), 10, 8,
		)
	}
	if err != nil {
		model = model.updateLogMessage("ERROR", err.Error())
		return err
	}

	model.TableRows[model.Cursor].CurrentValue = model.EditInput.Value()
	if err := model.SelectedGame.Map.UpdateMapFromTableRows(
		model.TableRows[model.Cursor],
	); err != nil {
		return err
	}

	model.Editing = false
	model.EditInput.Blur()


	return nil
}

// Helper function to calculate new boolean value, and update the memory map
func (model GameModel) updateBool() {
	// Calculate new boolean value and update
	valueInt, _ := strconv.Atoi(model.TableRows[model.Cursor].CurrentValue)
	newValue := strconv.Itoa(valueInt ^ 1)
	model.TableRows[model.Cursor].CurrentValue = newValue

	// Update the entry with matching offset
	if err := model.SelectedGame.Map.UpdateMapFromTableRows(
		model.TableRows[model.Cursor],
	); err != nil {
		model = model.updateLogMessage("ERROR", err.Error())
		return
	}
}

// Helper function to update model's current LogMessage
func (model GameModel) updateLogMessage(severity, message string) GameModel {
	model.LogMessage = Config.LogStruct{
		Severity: severity,
		Message: message,
	}
	return model
}