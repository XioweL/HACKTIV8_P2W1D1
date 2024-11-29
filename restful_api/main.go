package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Profile struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func createProfileHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Ensure the body is closed after processing
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	var t Profile
	err := decoder.Decode(&t)

	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		// Return error for invalid JSON body
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorMessage{Message: "Invalid Body Input"})
		return
	}

	// Validate required fields
	if t.Name == "" || t.Age <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorMessage{Message: "Invalid Body Input"})
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile created successfully"})
}

func main() {
	router := httprouter.New()
	router.POST("/profile", createProfileHandler)

	// Tambahkan log untuk memastikan server berjalan
	println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", router)
}

/*
1. go run main.go > Allow
2. Buka Postman pilih POST > masukan Endpointnya : http://localhost:8080/profile
3. Buka tab Body di bagian bawah.
4. Pilih raw dan ubah tipe data ke JSON dari dropdown di sebelah kanan.
5. Masukan data JSON ini di kolom editor :
	{
    "name": "John Doe",
    "age": 30
	}
6. Send
7. Periksa Response
INI DATA VALID
{
  "message": "Profile created successfully"
}
INI DATA !VALID
{
  "message": "Invalid Body Input"
}

8. Tambahkan Header (Opsional)
 8.1 Buka tab Headers.
 8.2 Tambahkan header:
		Key: Content-Type
		Value: application/json
*/
