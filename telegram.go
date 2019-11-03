package deploy5p

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func Send(text string) {
	logger := logrus.New()
	api_key := "722098929:AAELO91c9c-UzRdfP3AN0NBk3tDSNGQjymc"
	chat_id := "-1001364821881"

	logger.Info(fmt.Sprintf("Trying to send: %s\n", text))

	output := make(map[string]interface{})

	response, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", api_key, chat_id, text))
	if err != nil {
		panic(err)
	}

	if err := json.NewDecoder(response.Body).Decode(&output); err != nil {
		panic(err)
	}

	logger.Info(fmt.Sprintf("Sending succeeded: %s\n", output))
}
