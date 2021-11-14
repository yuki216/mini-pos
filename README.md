# mini-pos aplikasi beck end

## Changelog
- **v1**: membuat mini-pos beck end

## Keterangan
Ini adalah contoh project dengan  Clean Architecture dengan bahasa Go (Golang).


Aturan Arsitektur Bersih oleh Paman Bob
  * Independen dari Kerangka. Arsitektur tidak tergantung pada keberadaan beberapa perpustakaan perangkat lunak yang sarat fitur. Ini memungkinkan Anda untuk menggunakan kerangka kerja seperti itu sebagai alat, daripada harus menjejalkan sistem Anda ke dalam batasannya yang terbatas.
  * Dapat diuji. Aturan bisnis dapat diuji tanpa UI, Database, Server Web, atau elemen eksternal lainnya.
  * Independen dari UI. UI dapat berubah dengan mudah, tanpa mengubah sistem lainnya. UI Web dapat diganti dengan UI konsol, misalnya, tanpa mengubah aturan bisnis.
  * Independen dari Database. Anda dapat menukar Oracle atau SQL Server, untuk Mongo, BigTable, CouchDB, atau yang lainnya. Aturan bisnis Anda tidak terikat ke database.
  * Independen dari agensi eksternal mana pun. Sebenarnya aturan bisnis Anda sama sekali tidak tahu apa-apa tentang dunia luar.

More at https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html

Proyek ini mempunyai 4 Domain layer :
 * Models Layer
 * Repository Layer
 * Usecase Layer  
 * Delivery Layer

##### The diagram:

![golang clean architecture](https://github.com/yuki216/mini-pos/arsitektur.png)


#### Fitur Aplikasi
 * Login Menggunakan JWT dan password menggunakan enkripsi bcrypt
 * Register User
 * Master Produk
 * Master Customer
 * Master Supplier
 * Master Purchase
 * Master Produk Per Outlet
 * POS (cart, payment)
 * untuk prinsip-prinsip OOP hanya sedikit di implementasikan karena di golang kebanyakan menggunakan patern 
 * untuk aplikasi ini menggunakan patern adapter
 * test menggunakan postman file terlampir 
 
### Tabel DB 
 * users
 * roles
 * user_outlets
 * outlets
 * suppliers
 * products
 * purchases
 * outlet_products
 * customers
 * categories
 * orders
 * orderItem
 * payments
 * untuk setiap tabel sudah di berikan index sebagai optimizing 

### Cara menjalankan aplikasi
> Pastikan import dml dari file mini-pos.sql di mysql
> Pastikan telah install make file
> isi terlebih dahulu file .env.yml di folder config rename terlebih dabulu di tambahkan ekstensi *.yml

### 


#### Run the Applications
Here is the steps to run it with `docker-compose`

```bash
#move to directory
$ cd workspace

# Clone into YOUR $GOPATH/src
$ git clone https://github.com/yuki216/mini_pos.git

#move to project
$ cd mini_pos

# Build the docker image first
$ make docker

# Run the application
$ make run  / go run /app/main.go

# check if the containers are running
$ docker ps


# Stop
$ make stop
```


### Referensi: https://medium.com/golangid/mencoba-golang-clean-architecture-c2462f355f41

