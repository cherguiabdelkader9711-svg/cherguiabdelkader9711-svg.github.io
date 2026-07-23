package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/kkdai/youtube/v2"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html lang="ar" dir="rtl">
	<head>
	    <meta charset="UTF-8">
	    <title>محمل فيديوهات يوتيوب</title>
	    <style>
	        body { font-family: Tahoma, sans-serif; background: #0f172a; color: #fff; text-align: center; padding-top: 50px; }
	        input { width: 60%; padding: 12px; font-size: 16px; border-radius: 5px; border: none; outline: none; }
	        button { padding: 12px 25px; font-size: 16px; background: #22c55e; color: #fff; border: none; border-radius: 5px; cursor: pointer; font-weight: bold; }
	        button:hover { background: #16a34a; }
	        .container { background: #1e293b; padding: 40px; border-radius: 10px; display: inline-block; box-shadow: 0 4px 15px rgba(0,0,0,0.5); }
	    </style>
	</head>
	<body>
	    <div class="container">
	        <h2>تحميل فيديوهات يوتيوب</h2>
	        <form action="/download" method="GET">
	            <input type="text" name="url" placeholder="أدخل رابط فيديو يوتيوب هنا..." required>
	            <br><br>
	            <button type="submit">تحميل الفيديو</button>
	        </form>
	    </div>
	</body>
	</html>
	`
	fmt.Fprint(w, html)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	videoURL := r.URL.Query().Get("url")
	if videoURL == "" {
		http.Error(w, "الرجاء إدخال رابط الفيديو", http.StatusBadRequest)
		return
	}

	client := youtube.Client{}
	video, err := client.GetVideo(videoURL)
	if err != nil {
		http.Error(w, "فشل في جلب بيانات الفيديو: "+err.Error(), http.StatusInternalServerError)
		return
	}

	format := &video.Formats[0]
	stream, size, err := client.GetStream(video, format)
	if err != nil {
		http.Error(w, "فشل في بدء التحميل: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stream.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"video.mp4\""))
	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", size))

	_, _ = io.Copy(w, stream)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/download", downloadHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}

	fmt.Println("Server is running on port " + port)
	http.ListenAndServe(":"+port, nil)
}
