    package config

    import (
    "fmt"
    "log"
    "os"
    "strconv"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    )


    func InitDatabase() *gorm.DB {
    // Configure your PostgreSQL database details here
    port, _ := strconv.Atoi(os.Getenv("DB_PORT")) // Convert port to int
    dsn := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"),         // Host ของฐานข้อมูล
        port,                         // Port ของฐานข้อมูล
        os.Getenv("DB_USER"),         // ชื่อผู้ใช้ของฐานข้อมูล
        os.Getenv("DB_PASSWORD"),     // รหัสผ่านของฐานข้อมูล
        os.Getenv("DB_NAME"),         // ชื่อฐานข้อมูล
    )

    // New logger for detailed SQL logging
    newLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
        logger.Config{
            SlowThreshold: time.Second, // Slow SQL threshold
            LogLevel:      logger.Info, // Log level
            Colorful:      true,        // Enable color
        },
    )

    // connect database
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: newLogger, // Add Logger
    })
    if err != nil {
        log.Fatalf("failed to connect to database: %v", err)  // ข้อผิดพลาดจะถูกแสดง
    }

    fmt.Println("Connect Database successful")
    return db
    }

    func CloseDatabase(db *gorm.DB) {
    // Ensure the database connection is closed when the function exits
    sqlDB, err := db.DB()
    if err != nil {
        log.Fatalf("failed to get sql.DB from gorm: %v", err)
    }
    if err := sqlDB.Close(); err != nil {
        log.Fatalf("Failed to close database connection: %v", err)
    }
    }

