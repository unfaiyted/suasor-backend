package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	// URL of the avatar upload endpoint
	url := "http://localhost:8081/api/v1/users/avatar"
	
	// Create a buffer to store our request body as a multipart form
	var requestBody bytes.Buffer
	
	// Create a multipart writer
	multiPartWriter := multipart.NewWriter(&requestBody)
	
	// Create a form file field
	fileWriter, err := multiPartWriter.CreateFormFile("avatar", "test_avatar.jpg")
	if err != nil {
		log.Fatalf("Error creating form file: %v", err)
	}
	
	// Open the file
	file, err := os.Open("./tests/test_avatar.jpg")
	if err != nil {
		log.Fatalf("Error opening test file: %v", err)
	}
	defer file.Close()
	
	// Report the file size
	fileInfo, _ := file.Stat()
	fmt.Printf("File size: %d bytes\n", fileInfo.Size())
	
	// Copy the file content to the form field
	bytesWritten, err := io.Copy(fileWriter, file)
	if err != nil {
		log.Fatalf("Error copying file content: %v", err)
	}
	fmt.Printf("Copied %d bytes to form field\n", bytesWritten)
	
	// Close the multipart writer to finalize the form
	err = multiPartWriter.Close()
	if err != nil {
		log.Fatalf("Error closing multipart writer: %v", err)
	}
	
	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	
	// Set the content type header
	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())
	
	// Add your auth token (you'd need to get this from login first)
	// Manually set it for now with a placeholder
	authToken := "YOUR_AUTH_TOKEN_HERE"
	req.Header.Set("Authorization", "Bearer "+authToken)
	
	// Send the request
	fmt.Println("Sending request with Content-Type:", req.Header.Get("Content-Type"))
	
	// Uncomment to actually send the request
	// resp, err := client.Do(req)
	// if err != nil {
	//     log.Fatalf("Error sending request: %v", err)
	// }
	// defer resp.Body.Close()
	
	// Read the response
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	//     log.Fatalf("Error reading response: %v", err)
	// }
	
	// Print the response
	// fmt.Printf("Status: %s\nResponse: %s\n", resp.Status, string(body))
	
	// Output request details for debugging
	fmt.Println("\nRequest body size:", requestBody.Len(), "bytes")
	fmt.Println("Boundary:", multiPartWriter.Boundary())
	
	// Save the request body to a file for inspection
	err = os.WriteFile("request_debug.txt", requestBody.Bytes(), 0644)
	if err != nil {
		log.Fatalf("Error writing debug file: %v", err)
	}
	fmt.Println("Wrote request to request_debug.txt for inspection")
}