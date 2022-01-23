package main

import (
  "fmt"
  "os/exec"
  "github.com/spf13/cobra"
  "regexp"
)

func main() {
  var cmdValidator = &cobra.Command{
    Use:   "key : status, set or unset",
    Short: "Switch priv_validator_key.json symlink between keys/validator_key.json and keys/fullnode_key.json",
    Long: `If symlink is set to keys/validator_key.json, then the node is set to sign and validate blocks
make sure that you only have ONE symlink set to the keys/validator_key.json`,
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
      action := string(args[0])
      switch action {
      case "set":
	fmt.Print(switch_mode("validator_key.json"))
      case "unset":
	fmt.Print(switch_mode("fullnode_key.json"))
      case "status":
        fmt.Print(status())
      default:
        fmt.Println("Only 'status', 'set' and 'unset' actions are available")
      }
    },
  }

  var rootCmd = &cobra.Command{Use: "val_cli"}
  rootCmd.AddCommand(cmdValidator)
  rootCmd.Execute()
}

func switch_mode(key string) (string) {
  var runCmd string
  runCmd = "runuser -u cosmos -- ln -f -s /home/cosmos/config/keys/" + key + " /home/cosmos/config/priv_validator_key.json"
  tmp1, _ := exec.Command("bash", "-euxc", runCmd).CombinedOutput()
  runCmd = "runuser -u cosmos -- ls -l /home/cosmos/config/priv_validator_key.json"
  tmp2, _ := exec.Command("bash", "-euxc", runCmd).CombinedOutput()
  stdout := string(tmp1) + string(tmp2)

  match, _ := regexp.MatchString("may not be used by non-root users", stdout)
  if match {
    return "Should be run with sudo\n"
  } else {
    return stdout
  }
}

func status() (string) {
  var stdout string
  var runCmd string
  runCmd = "ls -l /home/cosmos/config | grep -i priv_validator_key.json"
  tmp1, _ := exec.Command("bash", "-c", runCmd).Output()
  runCmd = "hostname | tr -d '\n'"
  tmp2, _ := exec.Command("bash", "-c", runCmd).Output()
  match, _ := regexp.MatchString("fullnode_key.json", string(tmp1))

  if len(tmp1) > 0 {
    if match {
	stdout = "** " + string(tmp2) + " : in BACKUP mode **\n\n" + string(tmp1)
    } else {
        stdout = "** " + string(tmp2) + " : in VALIDATOR mode **\n\n" + string(tmp1)
    }
  } else {
      stdout = "** " + string(tmp2) + " : tendermint key is missing **\nPlease run 'val_cli key set' or 'val_cli key unset'\n"
  }
  return string(stdout)
}
