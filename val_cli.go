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
	r := switch_mode("validator_key.json")
	fmt.Print(r)
      case "unset":
	r := switch_mode("fullnode_key.json")
	fmt.Print(r)
      case "status":
        r := status()
        fmt.Print(r)
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
  var stdout string
  runCmd = "runuser -u cosmos -- rm /home/cosmos/config/priv_validator_key.json || true"
  tmp1, _ := exec.Command("bash", "-euxc", runCmd).CombinedOutput()
  runCmd = "runuser -u cosmos -- ln -s /home/cosmos/config/keys/" + key + " /home/cosmos/config/priv_validator_key.json"
  tmp2, _ := exec.Command("bash", "-euxc", runCmd).CombinedOutput()
  runCmd = "runuser -u cosmos -- ls -l /home/cosmos/config/priv_validator_key.json"
  tmp3, _ := exec.Command("bash", "-euxc", runCmd).CombinedOutput()
  stdout = string(tmp1) + string(tmp2) + string(tmp3)

  match, _ := regexp.MatchString("may not be used by non-root users", stdout)
  if match {
    return "Should be run with sudo\n"
  } else {
    return stdout
  }
}

func status() (string) {
  var stdout string
  runCmd := "ls -l /home/cosmos/config | grep -i priv_validator_key.json"
  tmp1, _ := exec.Command("bash", "-c", runCmd).Output()
  match, _ := regexp.MatchString("fullnode_key.json", string(tmp1))

  if len(tmp1) > 0 {
    if match {
        stdout = "** IN BACKUP MODE **\n\n" + string(tmp1)
    } else {
        stdout = "** IN VALIDATOR MODE **\n\n" + string(tmp1)
    }
  } else {
      stdout = "** MISSING priv_validator_key.json **\nPlease run 'val_cli key set' or 'val_cli key set'\n"
  }
  return string(stdout)
}
