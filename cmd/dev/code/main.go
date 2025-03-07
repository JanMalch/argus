package main

import (
	"github.com/janmalch/argus/internal/tui"
	"github.com/rivo/tview"
)

func main() {
	v := tui.NewCodeView()
	v.SetBorder(true).SetTitle(" JSON Example ")
	v.SetText(minifiedContent, "application/json")

	if err := tview.NewApplication().SetRoot(v, true).Run(); err != nil {
		panic(err)
	}
}

var content = `{
	"id": 101,
	"name": "Sample Data",
	"isActive": true,
	"timestamp": "2025-01-25T12:34:56Z",
	"price": 19.99,
	"tags": ["json", "example", "data"],
	"longDescription": "This is an example of a very long string that is intentionally written to be more than 200 characters. It is used to demonstrate how JSON can accommodate various data types, including long strings like this one, without any issues or constraints.",
	"coordinates": {
	  "latitude": 37.7749,
	  "longitude": -122.4194
	},
	"settings": {
	  "theme": "dark",
	  "notifications": {
		"email": true,
		"sms": false,
		"push": true
	  }
	},
	"items": [
	  {
		"id": 1,
		"name": "Item One",
		"quantity": 3,
		"price": 9.99
	  },
	  {
		"id": 2,
		"name": "Item Two",
		"quantity": 1,
		"price": 14.99
	  }
	],
	"user": {
	  "id": 5001,
	  "username": "user123",
	  "email": "user123@example.com",
	  "address": {
		"street": "123 Example St",
		"city": "Sampleville",
		"postalCode": "12fileviewer
		"country": "Exampleland"
	  }
	},
	"status": null,
	"metrics": [12, 25, 37, 49, 58],
	"metadata": {
	  "createdBy": "admin",
	  "createdOn": "2024-12-01",
	  "lastUpdatedBy": "editor",
	  "lastUpdatedOn": "2025-01-15"
	},
	"permissions": {
	  "read": true,
	  "write": false,
	  "execute": false
	},
	"version": 1.0,
	"notes": [
	  "This is the first note.",
	  "This is the second note, which is slightly longer to show variation in length.",
	  "A final note to complete the list."
	],
	"booleanList": [true, false, true, true, false],
	"nestedArray": [
	  [1, 2, 3],
	  [4, 5, 6],
	  [7, 8, 9]
	],
	"complexObject": {
	  "key1": "value1",
	  "key2": {
		"subKey1": 123,
		"subKey2": [1, 2, 3],
		"subKey3": {
		  "deepKey": "deepValue"
		}
	  }
	},
	"emptyArray": [],
	"emptyObject": {}
  }
  `
var minifiedContent = `{"id":101,"name":"Sample Data","isActive":true,"timestamp":"2025-01-25T12:34:56Z","price":19.99,"tags":["json","example","data"],"longDescription":"This is an example of a very long string that is intentionally written to be more than 200 characters. It is used to demonstrate how JSON can accommodate various data types, including long strings like this one, without any issues or constraints.","coordinates":{"latitude":37.7749,"longitude":-122.4194},"settings":{"theme":"dark","notifications":{"email":true,"sms":false,"push":true}},"items":[{"id":1,"name":"Item One","quantity":3,"price":9.99},{"id":2,"name":"Item Two","quantity":1,"price":14.99}],"user":{"id":5001,"username":"user123","email":"user123@example.com","address":{"street":"123 Example St","city":"Sampleville","postalCode":"12345","country":"Exampleland"}},"status":null,"metrics":[12,25,37,49,58],"metadata":{"createdBy":"admin","createdOn":"2024-12-01","lastUpdatedBy":"editor","lastUpdatedOn":"2025-01-15"},"permissions":{"read":true,"write":false,"execute":false},"version":1.0,"notes":["This is the first note.","This is the second note, which is slightly longer to show variation in length.","A final note to complete the list."],"booleanList":[true,false,true,true,false],"nestedArray":[[1,2,3],[4,5,6],[7,8,9]],"complexObject":{"key1":"value1","key2":{"subKey1":123,"subKey2":[1,2,3],"subKey3":{"deepKey":"deepValue"}}},"emptyArray":[],"emptyObject":{}}`
