package prompter

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/manifoldco/promptui"
)

// OptionsThreshold limits the number of items
// displayed by default on a select prompter.
const OptionsThreshold = 15

// SelectString prompts a select screen for the given options
// and using the speficied label.
//
// It returns an interface and an error (if exists).
func Select(q string, opts []string) (interface{}, error) {
	s := newSelector(q, opts)
	_, v, err := s.Run()
	if err != nil {
		return nil, err
	}

	return v, nil
}

// SelectRemoteOptions will fetch the slice of strings to
// present as options from the given request, it will use the
// default http.Client, so your request must contain all the
// information required for it to work, like authorization headers
// or any other parameter that is required by the API receiving the
// request.
func SelectRemoteOptions(q string, req *http.Request) (interface{}, error) {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var opts []string
	{
		c, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(c, &opts); err != nil {
			return nil, err
		}
	}

	s := newSelector(q, opts)

	_, v, err := s.Run()
	if err != nil {
		return nil, err
	}

	return v, nil
}

func newSelector(q string, opts []string) *promptui.Select {
	s := promptui.Select{
		Label: q,
		Items: opts,
		Size:  len(opts),
	}

	if len(opts) > OptionsThreshold {
		s.Size = OptionsThreshold
	}

	return &s
}
