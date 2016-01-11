// Copyright 2015, Stève Sfartz
// Licensed under the MIT License

// Utility to persist messages among calls
// Note : recordings are stored to Recorder Server or Tropo File Storage, a pointer to the recording audio file is stored with each voice message entry
//
package machine

import (
	"log"
	"time"
	"errors"
	"encoding/json"

	"github.com/golang/glog"
	bolt "github.com/boltdb/bolt"
)

const (
	BOLT_BUCKET= "Answering-Machine"
)

type StorageInterface interface {
	CreateVoiceMessage() *VoiceMessage
	Store(msg *VoiceMessage)
}

type VoiceMessageStorage struct {
	db 			*bolt.DB 		// database
}

type MachineProgress string
const (
	STARTED MachineProgress = "STARTED"
	RECORDED MachineProgress = "RECORDER"
	EMPTY MachineProgress = "EMPTY"
	FAILED MachineProgress = "FAILED"
)

type CheckedStatus string
const (
	NEW CheckedStatus = "NEW"
	CHECKED CheckedStatus = "CHECKED"
	DELETED CheckedStatus = "DELETED"
	UNDEFINED CheckedStatus = "UNDEFINED"
)

type VoiceMessage struct {
	CallID			string
	CreatedAt		time.Time
	CallerNumber    string
	Progress    	MachineProgress // enum of STARTED, RECORDED, EMPTY, FAILED
	Recording   	string // URL of the audio recording
	Duration		int // number of seconds
	Transcript  	string // transcript contents if successful
	Status      	CheckedStatus // enum of NEW, CHECKED, DELETED, UNDEFINED
	CheckedAt		time.Time
}


func NewStorage(dbName string) (*VoiceMessageStorage, error) {
	// Open the datafile in current directory, creates the db if it doesn't pre-exist.
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		glog.Fatalf("Cannot create database: %s", dbName)
		log.Fatal(err)
		return nil, errors.New("Cannot create voice messages database")
	}

	// Let's create the bucket if does not pre-exists
	err = db.Update(func(tx *bolt.Tx) error {
		glog.V(2).Infof("Creating bucket %s\n", BOLT_BUCKET)
		_, err := tx.CreateBucket([]byte(BOLT_BUCKET))
		if err != nil {
			glog.Infof("Cannot create bucket %s", BOLT_BUCKET)
			return errors.New("Cannot create bucket: " + BOLT_BUCKET)
		}

		glog.V(0).Infof("Created bucket %s to persist voice messages", BOLT_BUCKET)
		return nil
	})
	if err != nil {
		glog.Fatalf("Cannot create bucket %s", BOLT_BUCKET)
		log.Fatal(err)
		return nil, errors.New("Cannot create bucket")
	}

	return &VoiceMessageStorage{db}, nil
}


func (storage *VoiceMessageStorage) close() {
	log.Printf("[INFO] STORAGE Closing database")
	storage.db.Close()
}


func (storage *VoiceMessageStorage) CreateVoiceMessage(callID string, callerNumber string) *VoiceMessage {
	return &VoiceMessage{
		CallID:			callID,
		CreatedAt:		time.Now(),
		CallerNumber:   callerNumber,
		Progress:   	STARTED,
		Recording: 		"",
		Duration:		0,
		Transcript:  	"",
		Status:      	UNDEFINED,
		CheckedAt:		time.Time{},
	}
}


func (storage *VoiceMessageStorage) Store(msg *VoiceMessage) error {
	glog.V(2).Infof("Storing voice message for callID: %s", msg.CallID)

	err := storage.db.Update(func(tx *bolt.Tx) error {
		encoded, err1 := json.Marshal(msg)
		if err1 != nil {
			glog.V(0).Infof("Cannot encode message with callID: %s\n", msg.CallID)
			return err1
		}

		b := tx.Bucket([]byte(BOLT_BUCKET))
		err2 := b.Put([]byte(msg.CallID), encoded)
		if err2 != nil {
			glog.V(0).Infof("Error while storing voice message with callID: %s\n", msg.CallID)
		}
		return err2
	})
	if err != nil {
		glog.Warningf("Cannot store voice message for callID: %s", msg.CallID)
		return errors.New("Cannot storage message")
	}

	glog.V(0).Infof("Stored voice message for callID: %s", msg.CallID)
	return nil
}

func (storage *VoiceMessageStorage) FetchNewVoiceMessages() [](*VoiceMessage) {
	glog.V(2).Infof("Fetching all voice messages")

	//TODO

	return nil
}

func (storage *VoiceMessageStorage) FetchAllVoiceMessages() [](*VoiceMessage) {
	glog.V(2).Infof("Fetching all voice messages")

	var total int
	var messages [](*VoiceMessage) = make([](*VoiceMessage), 0, 10)
	err := storage.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BOLT_BUCKET))
		total = b.Stats().KeyN
		glog.V(3).Infof("Total of %d messages in bucket", total)

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			glog.V(3).Infof("key=%s, value=%s\n", k, v)

			var msg VoiceMessage
			err := json.Unmarshal(v, &msg)
			if err != nil {
				glog.V(0).Infof("json decode failed for voice message with key %s\n", k)
			} else {
				messages = append(messages, &msg)
			}
		}
		return nil
	})
	if err != nil {
		// Not a problem, we'll return the messages we ware able to fetch
	}

	return messages
}


func (storage *VoiceMessageStorage) GetVoiceMessageForCallID(callID string) (*VoiceMessage, error) {
	glog.V(2).Infof("Retreiving voice message for callID: %s\n", callID)

	var msg VoiceMessage

	err := storage.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BOLT_BUCKET))
		value := bucket.Get([]byte(callID))
		if value == nil {
			glog.V(2).Infof("No voice message found for callID: %s\n", callID)
			return errors.New("No voice message found")
		}

		err := json.Unmarshal(value, &msg)
		if err != nil {
			glog.V(0).Infof("json decode failed for voice message with callID %s\n", callID)
			return errors.New("Could not decode voice message")
		}

		return nil
	})
	if err != nil {
		glog.V(2).Infof("Retreiving voice message for callID: %s\n", callID)
		return nil, err
	}

	return &msg, nil
}











