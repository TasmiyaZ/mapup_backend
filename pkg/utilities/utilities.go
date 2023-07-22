/*
***

package for common utility functions

***
*/
package utilities

import "io/ioutil"

const Token = "cGFyc2VnZW9qc29u" //hardcoded token for api validation

func ReadDataFromFile(directory string, fileName string) ([]byte, error) {

	//creating full file path
	var storagePath = directory + "/" + fileName

	//getting file data as bytes
	bs, Readerr := ioutil.ReadFile(storagePath)
	if Readerr != nil {
		return nil, Readerr
	}
	// returing file contents as bytes
	return bs, nil
}
