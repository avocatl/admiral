package prompter

import "github.com/manifoldco/promptui"

// Confirm will ask a comfirmation prompt and return an error if it
// was not accepted or nil if it was.
func Confirm(q string, p *promptui.Prompt) error {
	if p == nil {
		p = &promptui.Prompt{
			IsConfirm: true,
			Label:     q,
		}
	} else {
		p.Label = q
		p.IsConfirm = true
	}

	_, err := p.Run()
	if err != nil {
		return err
	}

	return nil
}
