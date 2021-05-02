package pin

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"pinterest/application"
	"pinterest/domain/entity"
	"strconv"

	"github.com/gorilla/mux"
)

type PinInfo struct {
	pinApp   application.PinAppInterface
	boardApp application.BoardAppInterface
	s3App    application.S3AppInterface
	logger   *zap.Logger
}

func NewPinInfo(pinApp application.PinAppInterface,
	s3App application.S3AppInterface,
	boardApp application.BoardAppInterface,
	logger *zap.Logger) *PinInfo {
	return &PinInfo{
		pinApp:   pinApp,
		boardApp: boardApp,
		s3App:    s3App,
		logger:   logger,
	}
}

const maxPostPictureBodySize int = 8 * 1024 * 1024 // 8 mB

func (pinInfo *PinInfo) HandleAddPin(w http.ResponseWriter, r *http.Request) {
	bodySize := r.ContentLength
	userID := r.Context().Value(entity.CookieInfoKey).(*entity.CookieInfo).UserID

	if bodySize < 0 { // No picture was passed
		pinInfo.logger.Info(entity.NoPicturePassed.Error(), zap.String("url", r.RequestURI),
			zap.Int("for user", userID), zap.String("method", r.Method))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if bodySize > int64(maxPostPictureBodySize) { // Picture is too large
		pinInfo.logger.Info(entity.TooLargePicture.Error(), zap.String("url", r.RequestURI),
			zap.Int("for user", userID), zap.String("method", r.Method))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	r.ParseMultipartForm(bodySize)
	jsonData := r.FormValue("pinInfo") // TODO: replace string constants with keys
	currPin := entity.Pin{}
	err := json.Unmarshal([]byte(jsonData), &currPin)
	if err != nil {
		pinInfo.logger.Info(err.Error(), zap.String("url", r.RequestURI),
			zap.Int("for user", userID), zap.String("method", r.Method))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	currPin.UserID = userID
	if currPin.BoardID != 0 {
		err = pinInfo.boardApp.CheckBoard(userID, currPin.BoardID)
		if err != nil {
			pinInfo.logger.Info(
				err.Error(), zap.String("url", r.RequestURI),
				zap.Int("for user", userID), zap.String("method", r.Method))
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	currPin.PinId, err = pinInfo.pinApp.CreatePin(&currPin)
	if err != nil {
		pinInfo.logger.Info(
			err.Error(), zap.String("url", r.RequestURI),
			zap.Int("for user", userID), zap.String("method", r.Method))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	file, _, err := r.FormFile("pinImage")
	if err != nil {
		pinInfo.logger.Info(err.Error(), zap.String("url", r.RequestURI),
			zap.Int("for user", userID), zap.String("method", r.Method))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = pinInfo.pinApp.UploadPicture(currPin.PinId, file)
	if err != nil {
		pinInfo.logger.Info(err.Error(), zap.String("url", r.RequestURI),
			zap.Int("for user", userID), zap.String("method", r.Method))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pinID := entity.PinID{currPin.PinId}
	body, err := json.Marshal(pinID)
	if err != nil {
		pinInfo.logger.Info(err.Error(), zap.String("url", r.RequestURI), zap.String("method", r.Method))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func (pinInfo *PinInfo) HandleAddPinToBoard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardID, err := strconv.Atoi(vars[string(entity.IDKey)])
	if err != nil {
		pinInfo.logger.Info(err.Error(), zap.String("url", r.RequestURI), zap.String("method", r.Method))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(entity.CookieInfoKey).(*entity.CookieInfo).UserID

	err = pinInfo.boardApp.CheckBoard(userID, boardID)
	if err != nil {
		pinInfo.logger.Info(
			err.Error(), zap.String("url", r.RequestURI),
			zap.Int("for user", userID), zap.String("method", r.Method))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	pinID, err := strconv.Atoi(vars["pinID"])
	if err != nil {
		pinInfo.logger.Info(err.Error(), zap.String("url", r.RequestURI), zap.String("method", r.Method))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = pinInfo.pinApp.AddPin(boardID, pinID)
	if err != nil {
		pinInfo.logger.Info(
			err.Error(), zap.String("url", r.RequestURI),
			zap.Int("for user", userID), zap.String("method", r.Method))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (pinInfo *PinInfo) HandleSavePin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pinId, err := strconv.Atoi(vars[string(entity.IDKey)])
	if err != nil {
		pinInfo.logger.Info(err.Error(), zap.String("url", r.RequestURI), zap.String("method", r.Method))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(entity.CookieInfoKey).(*entity.CookieInfo).UserID

	err = pinInfo.pinApp.SavePin(userID, pinId)
	if err != nil {
		pinInfo.logger.Info(
			err.Error(), zap.String("url", r.RequestURI),
			zap.Int("for user", userID), zap.String("method", r.Method))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (pinInfo *PinInfo) HandleDelPinByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardID, err := strconv.Atoi(vars[string(entity.IDKey)])
	if err != nil {
		pinInfo.logger.Info(err.Error(), zap.String("url", r.RequestURI), zap.String("method", r.Method))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(entity.CookieInfoKey).(*entity.CookieInfo).UserID
	err = pinInfo.boardApp.CheckBoard(userID, boardID)
	if err != nil {
		pinInfo.logger.Info(
			err.Error(), zap.String("url", r.RequestURI),
			zap.Int("for user", userID), zap.String("method", r.Method))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	pinID, err := strconv.Atoi(vars["pinID"])
	if err != nil {
		pinInfo.logger.Info(err.Error(), zap.String("url", r.RequestURI), zap.String("method", r.Method))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = pinInfo.pinApp.DeletePin(boardID, pinID)
	if err != nil {
		pinInfo.logger.Info(
			err.Error(), zap.String("url", r.RequestURI),
			zap.Int("for user", userID), zap.String("method", r.Method))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (pinInfo *PinInfo) HandleGetPinByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	pinId, err := strconv.Atoi(vars[string(entity.IDKey)])
	if err != nil {
		pinInfo.logger.Info(err.Error(), zap.String("url", r.RequestURI), zap.String("method", r.Method))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resultPin, err := pinInfo.pinApp.GetPin(pinId)
	if err != nil {
		pinInfo.logger.Info(err.Error(), zap.String("url", r.RequestURI), zap.String("method", r.Method))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	body, err := json.Marshal(resultPin)
	if err != nil {
		pinInfo.logger.Info(err.Error(), zap.String("url", r.RequestURI), zap.String("method", r.Method))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (pinInfo *PinInfo) HandleGetPinsByBoardID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardID, err := strconv.Atoi(vars[string(entity.IDKey)])
	if err != nil {
		pinInfo.logger.Info(err.Error(), zap.String("url", r.RequestURI), zap.String("method", r.Method))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	boardPins, err := pinInfo.pinApp.GetPins(boardID)
	if err != nil {
		pinInfo.logger.Info(err.Error(), zap.String("url", r.RequestURI), zap.String("method", r.Method))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Pins := entity.PinsOutput{boardPins}

	pinsBody, err := json.Marshal(Pins)
	if err != nil {
		pinInfo.logger.Info(err.Error(), zap.String("url", r.RequestURI), zap.String("method", r.Method))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(pinsBody)
}

//// HandleUploadPicture takes picture from request and assigns it to current pin
//func (pinInfo *PinInfo) HandleUploadPicture(w http.ResponseWriter, r *http.Request) {
//	defer r.Body.Close()
//
//	bodySize := r.ContentLength
//	if bodySize < 0 { // No picture was passed
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	if bodySize > int64(maxPostPictureBodySize) { // Picture is too large
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	r.ParseMultipartForm(bodySize)
//	file, _, err := r.FormFile("pinImage")
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	defer file.Close()
//
//	userID := r.Context().Value(entity.CookieInfoKey).(*entity.CookieInfo).UserID
//	err = pinInfo.pinApp.UploadPicture(userID, file) // TODO: change userID to pinID
//
//	if err != nil {
//		log.Println(err)
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusNoContent)
//}

func (pinInfo *PinInfo) HandlePinsFeed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	numOfPins, err := strconv.Atoi(vars["num"])
	if err != nil {
		pinInfo.logger.Info(err.Error(), zap.String("url", r.RequestURI), zap.String("method", r.Method))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	feedPins, err := pinInfo.pinApp.GetNumOfPins(numOfPins)
	if err != nil {
		pinInfo.logger.Info(err.Error(), zap.String("url", r.RequestURI), zap.String("method", r.Method))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Pins := entity.PinsOutput{feedPins}

	pinsBody, err := json.Marshal(Pins)
	if err != nil {
		pinInfo.logger.Info(err.Error(), zap.String("url", r.RequestURI), zap.String("method", r.Method))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(pinsBody)
}