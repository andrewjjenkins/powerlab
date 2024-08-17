package megarac

import "fmt"

type powerCommandBody struct {
	PowerCommand int `json:"power_command"`
}

func (api *Api) PowerCommand(command int) error {
	body := powerCommandBody{PowerCommand: command}
	res, err := api.Post("/api/actions/power", body)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("failed power command %d: %v", command, err)
	}
	return nil
}
