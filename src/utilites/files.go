package utilites

import "os"
import "bytes"
import "io"
import "fmt"

func CountLines(name string)(uint32) {
	file, err := os.Open(name)
	if err != nil {
		return 0
	}
	
	defer file.Close()
	
	buf := make([]byte, 32*1024)
	count := 0
	
	for {
		c, err := file.Read(buf)
		count += bytes.Count(buf[:c], []byte{'\n'})
		
		switch {
			case err == io.EOF:
				return uint32(count)
				
			case err != nil:
				return 0
		}
	}
}


// Copy file
func CopyFile(sourceName string, destinationName string) {

	// First, let's check if the source exists
	if _, err := os.Stat(sourceName); err != nil {	
		return
	}

	// Now, let's check if the destination exists
	if _, err := os.Stat(destinationName); err == nil {
		
		// If it does, let's remove it 
		os.Remove(destinationName)
	}


	// OK, now let's create an empty file
	destinationFile, err := os.Create(destinationName)
	if err != nil {
		fmt.Println("Warning: I can't create a file")
		// TODO
		return
	}
	
	defer destinationFile.Close()
	
	// Let's open the source file
	sourceFile, err := os.Open(sourceName)
	if err != nil {
		// TODO
		return
	}
	defer sourceFile.Close()


	// And copy!
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		fmt.Println(err)
		// TODO
	}


}
