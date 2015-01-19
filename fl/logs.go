package logs

import "os"

func Wr_file(fl, text string) bool {
	f, err := os.OpenFile(fl, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
		return false
	}

	defer f.Close()

	_, err = f.WriteString(text + "\t")
	if err != nil {
		panic(err)
		return false
	}
	return true

}
