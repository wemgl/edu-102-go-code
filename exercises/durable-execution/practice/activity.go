package translation

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/activity"
	"io"
	"net/http"
	"net/url"
)

func TranslateTerm(ctx context.Context, input TranslationActivityInput) (TranslationActivityOutput, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("TranslateTerm has been executed", "Term", input.Term, "LanguageCode", input.LanguageCode)

	lang := url.QueryEscape(input.LanguageCode)
	term := url.QueryEscape(input.Term)
	url := fmt.Sprintf("http://localhost:9998/translate?lang=%s&term=%s", lang, term)

	resp, err := http.Get(url)
	if err != nil {
		logger.Error("External translate request failed", "Error", err)
		return TranslationActivityOutput{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Attempt to read response body failed", "Error", err)
		return TranslationActivityOutput{}, err
	}

	// This string will contain either the translated term, if the service could
	// perform the translation, or the error message, if it was unsuccessful
	content := string(body)

	status := resp.StatusCode
	if status >= 400 {
		logger.Error("External request's response indicates a client error", "StatusCode", status)
		// This means that we successfully called the service, but it could not
		// perform the translation for some reason
		return TranslationActivityOutput{},
			fmt.Errorf("HTTP Error %d: %s", status, content)
	}

	logger.Debug("TranslateTerm completed", "Translation", content)
	output := TranslationActivityOutput{
		Translation: content,
	}

	return output, nil
}
