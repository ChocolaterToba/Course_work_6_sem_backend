package application

import (
	"encoding/json"
	"fmt"
	"pinterest/domain/entity"
	"sync"

	"github.com/gorilla/websocket"
)

type connectionInfo struct {
	csrfToken string
	client    *websocket.Conn
}

type NotificationApp struct {
	notifications      map[int]map[int]entity.Notification
	lastNotificationID int
	connections        map[int]connectionInfo
	mu                 sync.Mutex
	userApp            UserAppInterface
}

func NewNotificationApp(userApp UserAppInterface) *NotificationApp {
	return &NotificationApp{
		notifications: make(map[int]map[int]entity.Notification),
		connections:   make(map[int]connectionInfo),
		userApp:       userApp,
	}
}

type NotificationAppInterface interface {
	AddNotification(notification *entity.Notification) (int, error)               // Add notification to list of user's notifications
	RemoveNotification(userID int, notificationID int) error                      // Remove notification from list of user's notifications
	EditNotification(notification *entity.Notification) error                     // Change fields of notification with same user and notification ID
	GetNotification(userID int, notificationID int) (*entity.Notification, error) // Get notification from db using user's and notification's IDs
	SendAllNotifications(userID int) error                                        // Send all of the notifications that this user has
	SendNotification(userID int, notificationID int) error                        // Send specified  notification to specified user
	ReadNotification(userID int, notificationID int) error                        // Changes notification's status to "Read"
	ChangeClient(userID int, client *websocket.Conn) error                        // Switches client  that was assigned to user
	ChangeToken(userID int, csrfToken string) error                               // Change user's CRSF token
	CheckToken(userID int, csrfToken string) error                                // Check if passed token is correct (nil on success)
}

func (notificationApp *NotificationApp) AddNotification(notification *entity.Notification) (int, error) {
	notificationApp.mu.Lock()
	defer notificationApp.mu.Unlock()

	notificationsMap, found := notificationApp.notifications[notification.UserID]
	if !found {
		_, err := notificationApp.userApp.GetUser(notification.UserID)
		if err != nil {
			return 0, entity.UserNotFoundError
		}
		notificationsMap = make(map[int]entity.Notification)
	}

	notification.NotificationID = notificationApp.lastNotificationID
	notificationApp.lastNotificationID++

	notificationsMap[notification.NotificationID] = *notification
	notificationApp.notifications[notification.UserID] = notificationsMap
	return notification.NotificationID, nil
}

func (notificationApp *NotificationApp) RemoveNotification(userID int, notificationID int) error {
	notificationApp.mu.Lock()
	defer notificationApp.mu.Unlock()

	notificationsMap, found := notificationApp.notifications[userID]
	if !found {
		return entity.UserNotFoundError
	}
	if notificationsMap == nil {
		return entity.NoNotificationsError
	}

	_, found = notificationsMap[notificationID]
	if !found {
		return entity.NotificationNotFoundError
	}

	delete(notificationsMap, notificationID)
	notificationApp.notifications[userID] = notificationsMap
	return nil
}

func (notificationApp *NotificationApp) EditNotification(notification *entity.Notification) error {
	notificationApp.mu.Lock()
	defer notificationApp.mu.Unlock()

	notificationsMap, found := notificationApp.notifications[notification.UserID]
	if !found {
		return entity.UserNotFoundError
	}
	if notificationsMap == nil {
		return entity.NoNotificationsError
	}

	_, found = notificationsMap[notification.NotificationID]
	if !found {
		return entity.NotificationNotFoundError
	}

	notificationApp.notifications[notification.UserID][notification.NotificationID] = *notification
	return nil
}

func (notificationApp *NotificationApp) GetNotification(userID int, notificationID int) (*entity.Notification, error) {
	notificationApp.mu.Lock()
	defer notificationApp.mu.Unlock()

	notificationsMap, found := notificationApp.notifications[userID]
	if !found {
		return nil, entity.UserNotFoundError
	}
	if notificationsMap == nil {
		return nil, entity.NoNotificationsError
	}

	notification, found := notificationsMap[notificationID]
	if !found {
		return nil, entity.NotificationNotFoundError
	}
	return &notification, nil
}

func sendMessage(client *websocket.Conn, msg []byte) error {
	w, err := client.NextWriter(websocket.TextMessage)
	if err != nil {
		return fmt.Errorf("Could not start writing")
	}

	w.Write(msg)
	w.Close()
	return nil
}

func (notificationApp *NotificationApp) SendAllNotifications(userID int) error {
	notificationApp.mu.Lock()
	defer notificationApp.mu.Unlock()

	notificationsMap, found := notificationApp.notifications[userID]
	if !found {
		_, err := notificationApp.userApp.GetUser(userID)
		if err != nil {
			return entity.UserNotFoundError
		}

		notificationsMap = make(map[int]entity.Notification)
	}

	connection, found := notificationApp.connections[userID]
	if !found {
		return entity.NotificationsClientNotSetError
	}

	allNotifications := entity.MessageManyNotifications{Type: entity.AllNotificationsTypeKey, Notifications: make([]entity.Notification, 0)}

	for _, notification := range notificationsMap {
		allNotifications.Notifications = append(allNotifications.Notifications, notification)
	}

	// TODO: maybe move sending notifications to handler???
	msg, err := json.Marshal(allNotifications)
	if err != nil {
		return fmt.Errorf("Could not parse messages into JSON")
	}

	err = sendMessage(connection.client, msg)

	return err
}

func (notificationApp *NotificationApp) SendNotification(userID int, notificationID int) error {
	notificationApp.mu.Lock()
	defer notificationApp.mu.Unlock()

	notificationsMap, found := notificationApp.notifications[userID]
	if !found {
		return entity.UserNotFoundError
	}

	notification, found := notificationsMap[notificationID]
	if !found {
		return entity.NotificationNotFoundError
	}

	connection, found := notificationApp.connections[userID]
	if !found {
		return entity.NotificationsClientNotSetError
	}

	notificationMsg := entity.MessageOneNotification{Type: entity.OneNotificationTypeKey, Notification: notification}

	msg, err := json.Marshal(notificationMsg)
	if err != nil {
		return fmt.Errorf("Could not parse message into JSON")
	}

	err = sendMessage(connection.client, msg)

	return err
}

func (notificationApp *NotificationApp) ReadNotification(userID int, notificationID int) error {
	notificationApp.mu.Lock()
	defer notificationApp.mu.Unlock()

	notificationsMap, found := notificationApp.notifications[userID]
	if !found {
		return entity.UserNotFoundError
	}

	notification, found := notificationsMap[notificationID]
	if !found {
		return entity.NotificationNotFoundError
	}

	if notification.IsRead {
		return entity.NotificationAlreadyReadError
	}

	notification.IsRead = true
	notificationsMap[notificationID] = notification
	notificationApp.notifications[userID] = notificationsMap
	return nil
}

func (notificationApp *NotificationApp) ChangeClient(userID int, client *websocket.Conn) error {
	notificationApp.mu.Lock()
	defer notificationApp.mu.Unlock()

	connection, found := notificationApp.connections[userID]
	if !found {
		_, err := notificationApp.userApp.GetUser(userID)
		if err != nil {
			return entity.UserNotFoundError
		}

		connection = connectionInfo{}
	}

	if connection.client != nil {
		connection.client.Close()
	}

	connection.client = client
	notificationApp.connections[userID] = connection
	return nil
}

func (notificationApp *NotificationApp) ChangeToken(userID int, csrfToken string) error {
	notificationApp.mu.Lock()
	defer notificationApp.mu.Unlock()

	connection, found := notificationApp.connections[userID]
	if !found {
		_, err := notificationApp.userApp.GetUser(userID)
		if err != nil {
			return entity.UserNotFoundError
		}

		connection = connectionInfo{}
	}

	connection.csrfToken = csrfToken
	notificationApp.connections[userID] = connection
	return nil
}

func (notificationApp *NotificationApp) CheckToken(userID int, csrfToken string) error {
	notificationApp.mu.Lock()
	defer notificationApp.mu.Unlock()

	connection, found := notificationApp.connections[userID]
	if !found {
		_, err := notificationApp.userApp.GetUser(userID)
		if err != nil {
			return entity.UserNotFoundError
		}

		connection = connectionInfo{}
	}

	if connection.csrfToken != csrfToken {
		return fmt.Errorf("Incorrect CSRF token")
	}

	return nil
}