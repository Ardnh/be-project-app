# be-project-app

`be-project-app` adalah aplikasi backend REST API untuk mengelola daftar tugas (todolist) dalam sebuah proyek. Aplikasi ini bertujuan untuk mencatat hal-hal yang perlu dilakukan ketika akan melakukan suatu proyek.

## Fitur

- **CRUD Proyek**: Buat, Baca, Perbarui, dan Hapus proyek.
- **CRUD Tugas**: Buat, Baca, Perbarui, dan Hapus tugas yang terkait dengan proyek.
- **Autentikasi**: Pengguna harus masuk untuk mengelola proyek dan tugas mereka.
- **Validasi Input**: Semua endpoint API memvalidasi input untuk memastikan data yang dimasukkan sesuai.

## Teknologi yang Digunakan

- **Go**: Bahasa pemrograman utama untuk pengembangan backend.
- **Fiber**: Framework web untuk Go.
- **GORM**: ORM untuk mengelola database.
- **SQLite/MySQL/PostgreSQL**: Database untuk menyimpan data (pilih sesuai kebutuhan).
- **JWT**: Untuk autentikasi dan otorisasi.
- **Docker**: Untuk containerisasi aplikasi.

## Struktur Proyek

be-project-app/
├── controllers
│ ├── projectController.go
│ ├── taskController.go
│ └── userController.go
├── models
│ ├── project.go
│ ├── task.go
│ └── user.go
├── repositories
│ ├── projectRepository.go
│ ├── taskRepository.go
│ └── userRepository.go
├── routes
│ ├── projectRoutes.go
│ ├── taskRoutes.go
│ └── userRoutes.go
├── services
│ ├── projectService.go
│ ├── taskService.go
│ └── userService.go
├── main.go
├── go.mod
└── go.sum


## Instalasi dan Pengaturan

### Prasyarat

- [Go](https://golang.org/dl/)
- [Docker](https://www.docker.com/get-started) (opsional, jika menggunakan Docker)

### Langkah Instalasi

1. **Clone repositori ini**

    ```sh
    git clone https://github.com/username/be-project-app.git
    cd be-project-app
    ```

2. **Instal dependensi**

    ```sh
    go mod tidy
    ```

3. **Konfigurasi Database**

    Buat file `.env` di root direktori dan tambahkan konfigurasi database Anda:

    ```env
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=yourusername
    DB_PASSWORD=yourpassword
    DB_NAME=yourdbname
    JWT_SECRET=your_jwt_secret
    ```

4. **Migrasi Database**

    Jalankan migrasi untuk membuat tabel di database:

    ```sh
    go run main.go migrate
    ```

5. **Menjalankan Aplikasi**

    Jalankan server:

    ```sh
    go run main.go
    ```

    Server akan berjalan di `http://localhost:3000`.

### Menggunakan Docker

Jika Anda ingin menjalankan aplikasi menggunakan Docker, Anda dapat menggunakan `Dockerfile` dan `docker-compose.yml` yang telah disediakan.

1. **Build Docker image**

    ```sh
    docker build -t be-project-app .
    ```

2. **Jalankan Docker container**

    ```sh
    docker-compose up
    ```

    Server akan berjalan di `http://localhost:3000`.

## Endpoint API

### Autentikasi

- **POST /api/register**: Mendaftar pengguna baru.
- **POST /api/login**: Masuk pengguna.

### Proyek

- **GET /api/projects**: Mendapatkan semua proyek.
- **GET /api/projects/:id**: Mendapatkan proyek berdasarkan ID.
- **POST /api/projects**: Membuat proyek baru.
- **PUT /api/projects/:id**: Memperbarui proyek berdasarkan ID.
- **DELETE /api/projects/:id**: Menghapus proyek berdasarkan ID.

### Tugas

- **GET /api/projects/:projectId/tasks**: Mendapatkan semua tugas dari proyek.
- **GET /api/tasks/:id**: Mendapatkan tugas berdasarkan ID.
- **POST /api/projects/:projectId/tasks**: Membuat tugas baru untuk proyek tertentu.
- **PUT /api/tasks/:id**: Memperbarui tugas berdasarkan ID.
- **DELETE /api/tasks/:id**: Menghapus tugas berdasarkan ID.

## Kontribusi

1. Fork repositori ini
2. Buat branch fitur baru (`git checkout -b fitur-anda`)
3. Commit perubahan Anda (`git commit -am 'Tambah fitur'`)
4. Push ke branch (`git push origin fitur-anda`)
5. Buat Pull Request

## Lisensi

Proyek ini dilisensikan di bawah [MIT License](LICENSE).

