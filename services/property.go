package services

import (
	
	"fmt"
	"Beego_API/models"
	"encoding/json"
	"errors"
	"os"
	"strings"
)

func loadProperties(path string)([]models.Property,error){
      
	  data,err := os.ReadFile(path)
	  

	  if err != nil {
		return nil,fmt.Errorf("JSON file read failed : %w",err)
	  }

	  var properties []models.Property

	  err = json.Unmarshal(data,&properties)

	  if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Converstion of Propeties data json -> struct failed : %w",err)
	  }
	  

	  return properties,nil
	  
}

type AllProperties struct {
	Properties []models.Property
}

func NewProperties(filePath string)(*AllProperties,error) {
    
	 filePath = strings.TrimSpace(filePath)

	 if(filePath == "") {
             return  nil,errors.New("invalid empty file path")
	 }

	 properties,err := loadProperties(filePath)

	 if err != nil {
		return  nil,err
	 }

	 

	 return  &AllProperties{Properties: properties},nil

}

var PropertyService AllProperties