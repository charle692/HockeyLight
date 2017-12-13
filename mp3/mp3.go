package mp3

import "os/exec"

// Play plays the specified mp3 file given it's path.
// Returns any error that has occurred along the way.
func Play(filePath string) error {
	cmdArgs := []string{filePath}
	_, err := exec.Command("mpg123", cmdArgs...).Output()

	if err != nil {
		return err
	}

	return nil
}
