package main

import (
	"log"
	"os"
	"os/exec"
)

type Message struct {
	cmdMsg  *exec.Cmd
	Display string
}

func show(text string, display string) (*Message, error) {
	cmd := exec.Command("swaynag", "--message", text, "--output", display, "--layer", "overlay")

	err := cmd.Start()
	if err != nil {
		logError("Unable to show swaynag in display '%s'.\n", display)
		return nil, err
	}

	go func() {
		cmd.Wait()
	}()

	return &Message{cmdMsg: cmd, Display: display}, nil
}

func ShowMessage(text string, message Message) (*Message, error) {
	if message.cmdMsg != nil {
		return nil, nil
	}
	return show(text, message.Display)
}

func CloseMessage(message Message) {
	log.Printf("close message process %v", message.cmdMsg)
	if err := message.cmdMsg.Process.Signal(os.Kill); err != nil {
		log.Printf("error during quit message %v", err)
	}

}

func ShowAll(text string, messages []Message) []Message {
	var openMessages []Message
	for _, message := range messages {
		newMessage, _ := ShowMessage(text, message)
		if newMessage == nil {
			openMessages = append(openMessages, message)
		} else {
			openMessages = append(openMessages, *newMessage)
		}
	}

	return openMessages
}

func CloseAll(messages []Message) {
	for _, message := range messages {
		CloseMessage(message)
	}
}
