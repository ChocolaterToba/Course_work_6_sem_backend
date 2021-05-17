package application

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"log"
	"pinterest/domain/entity"
	grpcPins "pinterest/services/pins/proto"
	"time"

	"github.com/EdlinOrg/prominentcolor"
)

type PinApp struct {
	grpcClient grpcPins.PinsClient
	boardApp   BoardAppInterface
}

type imageInfo struct {
	height       int
	width        int
	averageColor string
}

func NewPinApp(grpcClient grpcPins.PinsClient, boardApp BoardAppInterface) *PinApp {
	return &PinApp{grpcClient, boardApp}
}

type PinAppInterface interface {
	CreatePin(pin *entity.Pin, file io.Reader, extension string) (int, error)
	SavePin(userID int, pinID int) error                             // Add pin to user's initial board
	AddPin(boardID int, pinID int) error                             // Add pin to specified board
	GetPin(pinID int) (*entity.Pin, error)                           // Get pin by pinID
	GetPins(boardID int) ([]entity.Pin, error)                       // Get pins by boardID
	GetLastPinID(userID int) (int, error)                            // Get user's last pin's ID
	SavePicture(pin *entity.Pin) error                               // Update pin's picture properties
	RemovePin(boardID int, pinID int) error                          // Delete pin from board
	DeletePin(pinID int) error                                       // Delete pin entirely
	UploadPicture(pinID int, file io.Reader, extension string) error // Upload pin's image
	GetNumOfPins(numOfPins int) ([]entity.Pin, error)                // Get specified amount of pins
	SearchPins(keywords string) ([]entity.Pin, error)
}

// CreatePin creates passed pin and adds it to native user's board
// It returns pin's assigned ID and nil on success, any number and error on failure
func (pinApp *PinApp) CreatePin(pin *entity.Pin, file io.Reader, extension string) (int, error) {
	initBoardID, err := pinApp.boardApp.GetInitUserBoard(pin.UserID)
	if err != nil {
		return -1, err
	}
	grpcPin := grpcPins.Pin{}
	ConvertToGrpcPin(&grpcPin, pin)
	pinID, err := pinApp.grpcClient.CreatePin(context.Background(), &grpcPin)
	if err != nil {
		return -1, err
	}

	_, err = pinApp.grpcClient.AddPin(context.Background(), &grpcPins.PinInBoard{
		BoardID: int64(initBoardID), PinID: pinID.PinID})
	if err != nil {
		pinApp.grpcClient.DeletePin(context.Background(), pinID)
		return -1, err
	}

	if grpcPin.BoardID != int64(initBoardID) && grpcPin.BoardID != 0 {
		err = pinApp.AddPin(int(grpcPin.BoardID), int(pinID.PinID))
		if err != nil {
			pinApp.grpcClient.DeletePin(context.Background(), pinID)
			return -1, err
		}
	}

	err = pinApp.UploadPicture(int(pinID.PinID), file, extension)
	if err != nil {
		pinApp.grpcClient.DeletePin(context.Background(), pinID)
		return -1, err
	}

	return int(pinID.PinID), nil
}

// SavePin adds any pin to native user's board
// It returns nil on success, error on failure
func (pinApp *PinApp) SavePin(userID int, pinID int) error {
	initBoardID, err := pinApp.boardApp.GetInitUserBoard(userID)
	if err != nil {
		return err
	}

	err = pinApp.AddPin(initBoardID, pinID)
	if err != nil {
		return err
	}

	return nil
}

// AddPin adds pin to chosen board
// It returns nil on success, error on failure
func (pinApp *PinApp) AddPin(boardID int, pinID int) error {
	_, err := pinApp.grpcClient.AddPin(context.Background(), &grpcPins.PinInBoard{
		BoardID: int64(boardID), PinID: int64(pinID),
	})
	return err
}

// GetPin returns pin with passed pinID
// It returns that pin and nil on success, nil and error on failure
func (pinApp *PinApp) GetPin(pinID int) (*entity.Pin, error) {
	grpcPin, err := pinApp.grpcClient.GetPin(context.Background(), &grpcPins.PinID{PinID: int64(pinID)})
	if err != nil {
		return nil, err
	}
	pin := entity.Pin{}
	ConvertFromGrpcPin(&pin, grpcPin)
	return &pin, nil
}

// GetPins returns all the pins with passed boardID
// It returns slice of pins and nil on success, nil and error on failure
func (pinApp *PinApp) GetPins(boardID int) ([]entity.Pin, error) {
	grpcPinsList, err := pinApp.grpcClient.GetPins(context.Background(), &grpcPins.BoardID{BoardID: int64(boardID)})
	if err != nil {
		return nil, err
	}
	return ConvertGrpcPins(grpcPinsList), nil
}

// DeletePin deletes pin with passed pinID, deleting associated comments and board relations
// It returns nil on success and error on failure
func (pinApp *PinApp) DeletePin(pinID int) error {
	pin, err := pinApp.GetPin(pinID)
	if err != nil {
		return err
	}

	_, err = pinApp.grpcClient.DeletePin(context.Background(), &grpcPins.PinID{PinID: int64(pinID)})
	if err != nil {
		return err
	}
	_, err = pinApp.grpcClient.DeleteFile(context.Background(), &grpcPins.FilePath{ImagePath: pin.ImageLink})
	return err
}

// RemovePin deletes pin from user's passed board
// It returns nil on success and error on failure
func (pinApp *PinApp) RemovePin(boardID int, pinID int) error {
	pin, err := pinApp.GetPin(pinID)
	if err != nil {
		return err
	}

	_, err = pinApp.grpcClient.RemovePin(context.Background(), &grpcPins.PinInBoard{
		BoardID: int64(boardID), PinID: int64(pinID),
	})
	if err != nil {
		return err
	}

	refCount, err := pinApp.grpcClient.PinRefCount(context.Background(), &grpcPins.PinID{PinID: int64(pinID)})
	if err != nil {
		return err
	}

	if refCount.Number == 0 {
		_, err = pinApp.grpcClient.DeletePin(context.Background(), &grpcPins.PinID{PinID: int64(pinID)})
		if err != nil {
			return err
		}
		_, err = pinApp.grpcClient.DeleteFile(context.Background(), &grpcPins.FilePath{ImagePath: pin.ImageLink})
		return err
	}

	return nil
}

// SavePicture saves path to image of current pin in database
// It returns nil on success and error on failure
func (pinApp *PinApp) SavePicture(pin *entity.Pin) error {
	grpcPin := grpcPins.Pin{}
	ConvertToGrpcPin(&grpcPin, pin)
	_, err := pinApp.grpcClient.SavePicture(context.Background(), &grpcPin)
	return err
}

// GetLastPinID returns path to image of current pin in database
// It returns nil on success and error on failure
func (pinApp *PinApp) GetLastPinID(userID int) (int, error) {
	grpcPinID, err := pinApp.grpcClient.GetLastPinID(context.Background(), &grpcPins.UserID{Uid: int64(userID)})
	if err != nil {
		return 0, err
	}
	return int(grpcPinID.PinID), err
}

//UploadPicture uploads picture to pin and saves new picture path in S3
// It returns nil on success and error on failure
func (pinApp *PinApp) UploadPicture(pinID int, file io.Reader, extension string) error {
	pin, err := pinApp.GetPin(pinID)
	if err != nil {
		return entity.PinNotFoundError
	}

	fileAsBytes := make([]byte, 0)
	imageStruct := new(imageInfo)
	switch extension {
	case ".pinAppg":
	case ".jpg":
	case ".gif":
		fileAsBytes, _ = io.ReadAll(file) // TODO: this may be too slow, rework somehow? Maybe restore file after reading height/width?
		err = imageStruct.fillFromImage(bytes.NewReader(fileAsBytes))
		if err != nil {
			return fmt.Errorf("Image parsing failed")
		}
	default:
		return fmt.Errorf("File extension not supported")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, err := pinApp.grpcClient.UploadPicture(ctx)
	if err != nil {
		return entity.FileUploadError
	}
	req := &grpcPins.UploadImage{
		Data: &grpcPins.UploadImage_Extension{
			Extension: extension,
		},
	}
	err = stream.Send(req)
	if err != nil {
		log.Fatal("cannot send image info to server: ", err, stream.RecvMsg(nil))
	}
	reader := bufio.NewReader(file)
	buffer := make([]byte, 8*1024*1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}

		req = &grpcPins.UploadImage{
			Data: &grpcPins.UploadImage_ChunkData{
				ChunkData: buffer[:n],
			},
		}
		err = stream.Send(req)
		if err != nil {
			log.Fatal("cannot send chunk to server: ", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot receive response: ", err)
	}

	pin.ImageLink = res.Path
	pin.ImageHeight = imageStruct.height
	pin.ImageWidth = imageStruct.width
	pin.ImageAvgColor = imageStruct.averageColor

	err = pinApp.SavePicture(pin)
	if err != nil {
		pinApp.grpcClient.DeleteFile(context.Background(), &grpcPins.FilePath{ImagePath: res.Path})
		return fmt.Errorf("Pin saving failed")
	}

	return nil
}

// GetNumOfPins generates the main feed
// It returns numOfPins pins and nil on success, nil and error on failure
func (pinApp *PinApp) GetNumOfPins(numOfPins int) ([]entity.Pin, error) {
	grpcPinsList, err := pinApp.grpcClient.GetNumOfPins(context.Background(), &grpcPins.Number{Number: int64(numOfPins)})
	if err != nil {
		return nil, err
	}

	return ConvertGrpcPins(grpcPinsList), nil
}

// SearchPins returns pins by keywords
// It returns suitable pins and nil on success, nil and error on failure
func (pinApp *PinApp) SearchPins(keyWords string) ([]entity.Pin, error) {
	grpcPinsList, err := pinApp.grpcClient.SearchPins(context.Background(), &grpcPins.SearchInput{KeyWords: keyWords})
	if err != nil {
		return nil, err
	}

	return ConvertGrpcPins(grpcPinsList), nil
}

func (imageStruct *imageInfo) fillFromImage(imageFile io.Reader) error {
	image, _, err := image.Decode(imageFile)
	if err != nil {
		return fmt.Errorf("Image decoding failed")
	}

	imageStruct.height, imageStruct.width = image.Bounds().Dy(), image.Bounds().Dx()

	colors, err := prominentcolor.Kmeans(image)
	if err != nil {
		return fmt.Errorf("Could not determine image's most prominent color")
	}
	imageStruct.averageColor = colors[0].AsString()

	return nil
}

func ConvertToGrpcPin(grpcPin *grpcPins.Pin, pin *entity.Pin) {
	grpcPin.UserID = int64(pin.UserID)
	grpcPin.PinID = int64(pin.PinID)
	grpcPin.BoardID = int64(pin.BoardID)
	grpcPin.Title = pin.Title
	grpcPin.Description = pin.Description
	grpcPin.ImageAvgColor = pin.ImageAvgColor
	grpcPin.ImageWidth = int32(pin.ImageWidth)
	grpcPin.ImageHeight = int32(pin.ImageHeight)
	grpcPin.ImageLink = pin.ImageLink
}

func ConvertFromGrpcPin(pin *entity.Pin, grpcPin *grpcPins.Pin) {
	pin.UserID = int(grpcPin.UserID)
	pin.PinID = int(grpcPin.PinID)
	pin.BoardID = int(grpcPin.BoardID)
	pin.Title = grpcPin.Title
	pin.Description = grpcPin.Description
	pin.ImageAvgColor = grpcPin.ImageAvgColor
	pin.ImageWidth = int(grpcPin.ImageWidth)
	pin.ImageHeight = int(grpcPin.ImageHeight)
	pin.ImageLink = grpcPin.ImageLink
}

func ConvertGrpcPins(grpcPins *grpcPins.PinsList) []entity.Pin {
	pins := make([]entity.Pin, 0)
	for _, grpcPin := range grpcPins.Pins {
		pin := entity.Pin{}
		ConvertFromGrpcPin(&pin, grpcPin)
		pins = append(pins, pin)
	}
	return pins
}