#!/bin/bash

URL="http://localhost:8080/health"
TOTAL_CALLS=1000

echo "Memulai simulasi $TOTAL_CALLS panggilan API ke $URL..."

for ((i=1; i<=TOTAL_CALLS; i++))
do
   # Menggunakan curl untuk memanggil API
   # -s: mode silent
   # -o /dev/null: membuang output body agar terminal bersih
   # -w: menampilkan HTTP status code
   status=$(curl -s -o /dev/null -w "%{http_code}" "$URL")
   
   echo "Request ke-$i: Status $status"
done

echo "Simulasi selesai."