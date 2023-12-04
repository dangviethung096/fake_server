# Hướng dẫn sử dụng

### Requirements:
- Golang 1.21.3 trở lên
- Docker

### Chạy chương trình trên local
- B1: Để thực hiện chạy chương trình tại local, chúng ta cần cài đặt Golang 1.21.3 trở lên
- B2: Sau khi cài đặt xong, chúng ta cần clone project về máy
- B3: Truy cập vào folder account_service/init, thực hiện chạy file: `./init.sh`
- B4: Truy cập vào folder core, thực hiện chạy lệnh: `go mod tidy`
- B5: Tương tự truy cập vào folder account_service, thực hiện chạy lệnh: `go mod tidy`
- B6: Sau đó truy cập vào folder account_service, chạy file ./run_localhost.sh
- B7: Server sẽ thực hiện mặc định chạy tại port 8080