package config

import (
	"fmt"
	"io/ioutil"

	"github.com/PulseDevelopmentGroup/0x626f74/util"
	"github.com/tidwall/gjson"
)

type (
	// BotConfig defines the configuration container for the bot
	BotConfig struct {
		SimpleCommands map[string]string
		Permissions    *BotPermissions
	}

	// BotPermissions contains the permission maps for roles, channels, and
	// users based on the config file.
	BotPermissions struct {
		RoleIDs map[string][]string
		ChanIDs map[string][]string
		UserIDs map[string][]string
	}
)

// Get loads the config from the json file at the path specified
func Get(path string) (*BotConfig, error) {
	json, err := getJSON(path)
	if err != nil {
		return &BotConfig{}, err
	}

	simpleCommands, err := getSimpleCommands(json)
	if err != nil {
		return &BotConfig{}, err
	}

	roles, err := getRoles(json)
	if err != nil {
		return &BotConfig{}, err
	}

	return &BotConfig{
		SimpleCommands: simpleCommands,
		Permissions: &BotPermissions{
			RoleIDs: roles,
		},
	}, nil
}

// TODO: Support URLs
func getJSON(path string) (string, error) {
	file, err := util.InitFile(path)
	if err != nil {
		return "", err
	}

	json, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return "", err
	}

	return string(json), nil
}

func getSimpleCommands(json string) (map[string]string, error) {
	out := make(map[string]string)

	m, ok := gjson.Parse(
		gjson.Get(json, "simpleCommands").String(),
	).Value().(map[string]interface{})
	if !ok {
		return out,
			fmt.Errorf("unable to get list of simple commands from config file")
	}

	/* TODO: Better error handling in the event of a non-string value */
	for k, v := range m {
		out[k] = v.(string)
	}

	return out, nil
}

/* TODO: I feel like this function is really a disaster and needs a rework sometime */
func getRoles(json string) (map[string][]string, error) {
	out := make(map[string][]string)

	r := gjson.Get(json, "permissions")
	r.ForEach(func(key, value gjson.Result) bool {
		var roles []string

		if value.IsArray() {
			for _, v := range value.Array() {
				roles = append(roles, v.String())
			}
		}

		out[key.String()] = roles

		return true
	})
	return out, nil
}
