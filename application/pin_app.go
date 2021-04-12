package application

import (
	"fmt"
	"io"
	"pinterest/domain/entity"
	"pinterest/domain/repository"
)

type PinApp struct {
	p repository.PinRepository
}

func NewPinApp(p repository.PinRepository) *PinApp {
	return &PinApp{p}
}

type PinAppInterface interface {
	CreatePin(int, *entity.Pin) (int, error)
	AddPin(int, int) error             // Saving user's pin
	GetPin(int) (*entity.Pin, error)   // Get pin by pinID
	GetPins(int) ([]entity.Pin, error) // Get pins by boardID
	GetLastUserPinID(int) (int, error)
	SavePicture(pin *entity.Pin) error
	DeletePin(int, int, S3AppInterface) error           // Removes pin by ID
	UploadPicture(int, io.Reader, S3AppInterface) error // Upload pin
}

func (pn *PinApp) CreatePin(userID int, pin *entity.Pin) (int, error) {
	boardApp := BoardApp{}
	initBoard, err := boardApp.GetInitUserBoard(userID)
	if err != nil {
		return -1, err
	}
	pinID, err := pn.p.CreatePin(pin)
	if err != nil {
		return -1, err
	}

	err = pn.p.AddPin(initBoard.BoardID, pinID)
	if err != nil {
		pn.p.DeletePin(pinID, userID)
		return -1, err
	}

	return pinID, nil
}

// AddPin adds user's pin to database
// It returns pin's assigned ID and nil on success, any number and error on failure
func (pn *PinApp) AddPin(boardID int, pinID int) error {
	return pn.p.AddPin(boardID, pinID)
}

// GetPin returns pin with passed pinID
// It returns that pin and nil on success, nil and error on failure
func (pn *PinApp) GetPin(pinID int) (*entity.Pin, error) {
	return pn.p.GetPin(pinID)
}

// GetPins returns all the pins with passed boardID
// It returns slice of pins and nil on success, nil and error on failure
func (pn *PinApp) GetPins(boardID int) ([]entity.Pin, error) {
	return pn.p.GetPins(boardID)
}

// DeletePin deletes pin with passed pinID
// It returns nil on success and error on failure
func (pn *PinApp) DeletePin(pinID int, userID int, s3App S3AppInterface) error {
	pin, err := pn.p.GetPin(pinID)
	if err != nil {
		return err
	}

	err = pn.p.DeletePin(pinID, userID)
	if err != nil {
		return err
	}

	return s3App.DeleteFile(pin.ImageLink)
}

// SavePicture saves path to image of current pin in database
// It returns nil on success and error on failure
func (pn *PinApp) SavePicture(pin *entity.Pin) error {
	return pn.p.SavePicture(pin)
}

// GetLastUserPinID returns path to image of current pin in database
// It returns nil on success and error on failure
func (pn *PinApp) GetLastUserPinID(userID int) (int, error) {
	return pn.p.GetLastUserPinID(userID)
}

func (pn *PinApp) UploadPicture(userID int, file io.Reader, s3App S3AppInterface) error {
	pinID, err := pn.GetLastUserPinID(userID)
	if err != nil {
		return fmt.Errorf("No pin found to place picture")
	}

	pin, err := pn.GetPin(pinID)
	if err != nil {
		return fmt.Errorf("No pin found to place picture")
	}

	filenamePrefix, err := GenerateRandomString(40) // generating random image
	if err != nil {
		return fmt.Errorf("Could not generate filename")
	}

	picturePath := "pins/" + filenamePrefix + ".jpg"
	err = s3App.UploadFile(file, picturePath)
	if err != nil {
		return fmt.Errorf("File upload failed")
	}

	pin.ImageLink = picturePath

	err = pn.SavePicture(pin)
	if err != nil {
		s3App.DeleteFile(picturePath)
		return fmt.Errorf("Pin saving failed")
	}

	return nil
}
